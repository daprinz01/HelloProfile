package main

import (
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
	"helloprofile.com/controllers"
	"helloprofile.com/persistence/orm/helloprofiledb"

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
	fields := log.Fields{"microservice": "helloprofile.service"}
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
	helloprofiledatabase := helloprofiledb.New(db)
	env := &controllers.Env{HelloProfileDb: helloprofiledatabase}
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
		Format: "{\"@timestamp\":\"${time_rfc3339}\", \"uri\":\"${uri}\", \"remote_ip\":\"${remote_ip}\", \"host\":\"${host}\", \"id\":\"${id}\", \"method\":\"${method}\", \"user_agent\":\"${user_agent}\", \"status\":\"${status}\", \"error\":\"${error}\", \"latency\":\"${latency}\", \"latency_human\":\"${latency_human}\", \"bytes_in\":\"${bytes_in}\", \"bytes_out\":\"${bytes_out}\", \"message\":\"Echo http request logger\", \"microservice\": \"helloprofile.service\", \"level\":\"info\", \"user_agent\":\"${user_agent}\"}",
	}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		ExposeHeaders: []string{"*"},
		AllowOrigins:  []string{"*"},
		AllowHeaders:  []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, "Role", echo.HeaderAuthorization, "Refresh-Token", echo.HeaderXRealIP},
	}))
	e.Use(controllers.TrackResponseTime)
	e.Use(middleware.Recover())
	// Enable metrics middleware
	e.Use(echoPrometheus.MetricsMiddleware())
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	apiNoAuth := e.Group("/api/v1")
	// apiNoAuth.Use(env.CheckApplication)
	auth := apiNoAuth.Group("/auth")
	// Methods that require authentication but don't need the information from the applications or authorization to function
	apiAuth := apiNoAuth.Group("/app")

	apiAuth.Use(controllers.Authorize)

	// Admin operations authorization
	apiAdminAuth := apiNoAuth.Group("/admin")

	apiAdminAuth.Use(controllers.AuthorizeAdmin)

	// Methods that don't require the application information is it's being verified in a middleware
	apiNoAuth.GET("/refresh", env.RefreshToken)
	apiNoAuth.POST("/otp/send", env.SendOtp)
	apiNoAuth.POST("/password/reset", env.ResetPassword)
	apiNoAuth.POST("/confirm/email", env.DoEmailVerification)
	apiNoAuth.POST("/verify/email", env.VerifyEmailToken)

	// Methods that check application themselves and use the applicaiton information
	auth.POST("/login", env.Login)
	auth.POST("/otp/verify", env.VerifyOtp)
	auth.POST("/login/google", env.GoogleLoginHandler)

	// User operations
	apiAdminAuth.GET("/user", env.GetUsers)
	apiNoAuth.GET("/user/:username/check", env.CheckAvailability)
	apiAuth.GET("/user", env.GetUser)
	auth.POST("/user", env.Register)
	apiAuth.PUT("/user", env.UpdateUser)
	apiAuth.DELETE("/user/:username", env.DeleteUser)

	//Profile operations
	apiNoAuth.GET("/profile", env.GetProfiles)
	apiAuth.POST("/profile/:email", env.AddProfile)
	apiAuth.PUT("/profile", env.UpdateProfile)
	apiAuth.PUT("/profile/:profileId", env.DeleteProfile)

	// User Role operations
	apiAuth.PUT("/user/role/:newRole/:oldRole/:username", env.UpdateUserRole)
	apiAuth.POST("/user/role/:role/:username", env.AddUserToRole)

	// Roles Operations
	apiAdminAuth.GET("/role/:role", env.GetRole)
	apiAdminAuth.GET("/role", env.GetRoles)
	apiAdminAuth.POST("/role", env.AddRole)
	apiAdminAuth.PUT("/role/:role", env.UpdateRole)
	apiAdminAuth.DELETE("/role/:role", env.DeleteRole)

	// Basic block operations
	apiAuth.POST("/blocks/basic/:profileId", env.AddBasicBlock)
	apiAuth.PUT("/blocks/basic", env.UpdateBasicBlock)
	apiAuth.DELETE("/blocks/basic/:id", env.DeleteBasicBlock)

	// Contact block operations
	apiAuth.POST("/blocks/contact/:profileId", env.AddContactBlock)
	apiAuth.PUT("/blocks/contact", env.UpdateContactBlock)
	apiAuth.DELETE("/blocks/contact/:id", env.DeleteContactBlock)

	// Call to action operations
	apiAuth.GET("/callToAction", env.GetCallToActions)

	// Content types operations
	apiAuth.GET("/content/types", env.GetContents)

	// Socials types
	apiAuth.GET("/socials/types", env.GetSocials)

	// Socials operations
	apiAuth.POST("/blocks/socials/:profileId", env.AddSocialsBlock)
	apiAuth.PUT("/blocks/socials", env.UpdateSocialsBlock)
	apiAuth.DELETE("/blocks/socials", env.DeleteSocialBlock)

	// Content operations
	apiAuth.POST("/blocks/content/:profileId", env.AddContentBlock)
	apiAuth.PUT("/blocks/content", env.UpdateContentBlock)
	apiAuth.DELETE("/blocks/content/:contentId", env.DeleteContentBlock)

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
