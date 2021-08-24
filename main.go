package main

import (
	"authengine/controllers"
	"authengine/persistence/orm/authdb"
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	echoPrometheus "github.com/globocom/echo-prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	// Configure Logging
	log.SetFormatter(&log.JSONFormatter{
		FieldMap: log.FieldMap{
			log.FieldKeyTime: "@timestamp",
			log.FieldKeyMsg:  "message",
			log.FieldKeyFunc: "function",
		},
	})
	fields := log.Fields{"microservice": "persian.black.authengine.service"}
	log.WithFields(fields)
	log.SetLevel(log.TraceLevel)

	logFileLocation := os.Getenv("LOG_FILE_LOCATION")
	if logFileLocation != "" {
		log.SetOutput(&lumberjack.Logger{
			Filename:   logFileLocation,
			MaxSize:    50, // megabytes
			MaxBackups: 3,
			MaxAge:     28,   //days
			Compress:   true, // disabled by default
		})

	}

	log.WithFields(fields).Info("Successfully initialized log file...")

	host := os.Getenv("DB_HOST")
	if host == "" {
		log.WithFields(fields).Warn("Host cannot be empty")
		panic("DB_HOST cannot be empty, application intialization failed...")
	}
	port := 8669
	dbport := os.Getenv("DB_PORT")
	if dbport == "" {
		log.WithFields(fields).Warn("Port cannot be empty")
		panic("DB_PORT cannot be empty, application intialization failed...")
	} else {
		portnumber, err := strconv.Atoi(dbport)
		if err != nil {
			log.WithFields(fields).Warn("Port is not a valid number")
			panic("Port is not a valid number, please enter a valid number for DB_PORT. Application initialization failed...")
		} else {
			port = portnumber
		}
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		log.WithFields(fields).Warn("User cannot be empty")
		panic("DB_USER cannot be empty, application intialization failed...")
	}
	password := os.Getenv("DB_PASSWORD")
	if user == "" {
		log.WithFields(fields).Warn("User cannot be empty")
		panic("DB_USER cannot be empty, application intialization failed...")
	}
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		log.WithFields(fields).Warn("Database name cannot be empty")
		panic("DB_NAME cannot be empty, application intialization failed...")
	}
	sslmode := os.Getenv("DB_SSL_MODE")
	if sslmode == "" {
		log.WithFields(fields).Warn("SSL mode cannot be empty")
		panic("DB_SSL_MODE cannot be empty, application intialization failed...")
	}
	//Connect to database
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	authdatabase := authdb.New(db)
	env := &controllers.Env{AuthDb: authdatabase}
	log.WithFields(fields).Warn("Successfully connected to database!")
	// // Create Server and Route Handlers
	// r := mux.NewRouter()
	// r.Use(controllers.TrackResponseTime)

	srv := &http.Server{

		Addr:         ":8083",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		// Format: "method=${method}, uri=${uri}, status=${status}\n",
		Output: &lumberjack.Logger{
			Filename:   logFileLocation,
			MaxSize:    50, // megabytes
			MaxBackups: 3,
			MaxAge:     28,   //days
			Compress:   true, // disabled by default,
		},
		Format: "{\"@timestamp\":\"${time_rfc3339}\", \"uri\":\"${uri}\", \"remote_ip\":\"${remote_ip}\", \"host\":\"${host}\", \"id\":\"${id}\", \"method\":\"${method}\", \"user_agent\":\"${user_agent}\", \"status\":\"${status}\", \"error\":\"${error}\", \"latency\":\"${latency}\", \"latency_human\":\"${latency_human}\", \"bytes_in\":\"${bytes_in}\", \"bytes_out\":\"${bytes_out}\", \"message\":\"Echo http request logger\", \"microservice\": \"persian.black.authengine.service\", \"level\":\"info\", \"user_agent\":\"${user_agent}\"}",
	}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		ExposeHeaders: []string{"*"},
		AllowOrigins:  []string{"*"},
		AllowHeaders:  []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, "Role", echo.HeaderAuthorization, "Refresh-Token", echo.HeaderXRealIP},
	}))
	e.Use(controllers.TrackResponseTime)
	// e.Use(middleware.CSRF())
	e.Use(middleware.Recover())
	// Enable metrics middleware
	e.Use(echoPrometheus.MetricsMiddleware())
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	apiNoAuth := e.Group("/api/v1")
	apiNoAuth.Use(env.CheckApplication)
	auth := apiNoAuth.Group("/auth")
	// Methods that require authentication but don't need the information from the applications or authorization to function
	apiAuth := apiNoAuth.Group("/app")
	// Configure middleware with the custom claims type
	// config := middleware.JWTConfig{
	// 	Claims:     &models.Claims{},
	// 	SigningKey: []byte(os.Getenv("JWT_SECRET_KEY")),
	// }
	apiAuth.Use(controllers.Authorize)

	// Admin operations authorization
	apiAdminAuth := apiNoAuth.Group("/admin")

	// adminConfig := middleware.JWTConfig{
	// 	Claims: &models.Claims{
	// 		Role: "admin",
	// 	},
	// 	SigningKey: []byte(os.Getenv("JWT_SECRET_KEY")),
	// }
	apiAdminAuth.Use(controllers.AuthorizeAdmin)

	// Methods that don't require the application information is it's being verified in a middleware
	apiNoAuth.GET("/:application/refresh", env.RefreshToken)
	apiNoAuth.POST("/:application/otp/send", env.SendOtp)
	apiNoAuth.POST("/:application/password/reset", env.ResetPassword)
	apiNoAuth.POST("/:application/confirm/email", env.DoEmailVerification)
	apiNoAuth.POST("/:application/verify/email", env.VerifyEmailToken)

	// Methods that check application themselves and use the applicaiton information
	auth.POST("/:application/login", env.Login)
	auth.POST("/:application/otp/verify", env.VerifyOtp)

	// User operations
	apiAdminAuth.GET("/:application/user", env.GetUsers)
	apiNoAuth.GET("/:application/user/:username/check", env.CheckAvailability)
	apiNoAuth.GET("/:application/user/:username", env.GetUser)
	auth.POST("/:application/user", env.Register)
	apiAuth.PUT("/:application/user", env.UpdateUser)
	apiAuth.DELETE("/:application/user/:username", env.DeleteUser)

	// User Language operations
	apiNoAuth.GET("/:application/user/language/:username", env.GetUserLanguages)
	apiAuth.POST("/:application/user/language/:username/:language/:proficiency", env.AddUserLanguage)
	apiAuth.DELETE("/:application/user/language/:username/:language", env.DeleteUserLanguages)

	// User Role operations
	apiAuth.PUT("/:application/user/role/:newRole/:oldRole/:username", env.UpdateUserRole)
	apiAuth.POST("/:application/user/role/:role/:username", env.AddUserToRole)

	// Language operations
	apiNoAuth.GET("/:application/language/:language", env.GetLanguage)
	apiNoAuth.GET("/:application/language", env.GetLanguages)
	apiAdminAuth.POST("/:application/language/:language", env.AddLanguage)
	apiAdminAuth.PUT("/:application/language/:language/:newLanguage", env.UpdateLanguage)
	apiAdminAuth.DELETE("/:application/language/:language", env.DeleteLanguage)

	// Language proficiency operations
	apiNoAuth.GET("/:application/proficiency/:proficiency", env.GetLanguageProficiency)
	apiNoAuth.GET("/:application/proficiency", env.GetLanguageProficiencies)
	apiAdminAuth.POST("/:application/proficiency/:proficiency", env.AddLanguageProficiency)
	apiAdminAuth.PUT("/:application/proficiency/:proficiency/:newProficiency", env.UpdateLanguageProficiency)
	apiAdminAuth.DELETE("/:application/proficiency/:proficiency", env.DeleteLanguageProficiency)

	// Timezone Operations
	apiNoAuth.GET("/:application/timezone/:timezone", env.GetTimezone)
	apiNoAuth.GET("/:application/timezone", env.GetTimezones)
	apiAdminAuth.POST("/:application/timezone", env.AddTimezone)
	apiAdminAuth.PUT("/:application/timezone/:timezone", env.UpdateTimezone)
	apiAdminAuth.DELETE("/:application/timezone/:timezone", env.DeleteTimezone)

	// Application Operations
	apiNoAuth.GET("/:application/application/:application", env.GetApplication)
	apiNoAuth.GET("/:application/application", env.GetApplications)
	apiAdminAuth.POST("/:application/application", env.AddApplication)
	apiAdminAuth.PUT("/:application/application/:application", env.UpdateApplication)
	apiAdminAuth.DELETE("/:application/application/:application", env.DeleteApplication)

	// Countries Operations
	apiNoAuth.GET("/:application/country/:country", env.GetCountry)
	apiNoAuth.GET("/:application/country", env.GetCountries)
	apiAdminAuth.POST("/:application/country", env.AddCountry)
	apiAdminAuth.PUT("/:application/country/:country", env.UpdateCountry)
	apiAdminAuth.DELETE("/:application/country/:country", env.DeleteCountry)

	// States Operations
	apiNoAuth.GET("/:application/state/:state", env.GetState)
	apiNoAuth.GET("/:application/state", env.GetStates)
	apiAdminAuth.POST("/:application/state/:state/:country", env.AddState)
	apiAdminAuth.PUT("/:application/state/:state/:newState", env.UpdateState)
	apiAdminAuth.DELETE("/:application/state/:state", env.DeleteState)
	apiNoAuth.GET("/:application/state/country/:country", env.GetStatesByCountry)

	// Roles Operations
	apiAdminAuth.GET("/:application/role/:role", env.GetRole)
	apiAdminAuth.GET("/:application/role", env.GetRoles)
	apiAdminAuth.POST("/:application/role", env.AddRole)
	apiAdminAuth.PUT("/:application/role/:role", env.UpdateRole)
	apiAdminAuth.DELETE("/:application/role/:role", env.DeleteRole)
	apiAdminAuth.GET("/:application/role/application/:application", env.GetRolesByApplication)
	apiAdminAuth.POST("/:application/role/:role/:application", env.AddApplicationRole)

	go func(fields log.Fields) {
		log.WithFields(fields).Info("Starting Server...")
		e.Logger.Fatal(e.StartServer(srv))
	}(fields)

	// Graceful Shutdown
	waitForShutdown(srv, fields)
}

func waitForShutdown(srv *http.Server, fields log.Fields) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan
	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	srv.Shutdown(ctx)

	log.WithFields(fields).Info("Shutting down...")
	os.Exit(0)
}
