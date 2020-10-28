package main

import (
	"authengine/controllers"
	"authengine/persistence/orm/authdb"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	echoPrometheus "github.com/globocom/echo-prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	host := os.Getenv("DB_HOST")
	if host == "" {
		fmt.Println("Host cannot be empty")
		panic("DB_HOST cannot be empty, application intialization failed...")
	}
	port := 8669
	dbport := os.Getenv("DB_PORT")
	if dbport == "" {
		fmt.Println("Port cannot be empty")
		panic("DB_PORT cannot be empty, application intialization failed...")
	} else {
		portnumber, err := strconv.Atoi(dbport)
		if err != nil {
			fmt.Println("Port is not a valid number")
			panic("Port is not a valid number, please enter a valid number for DB_PORT. Application initialization failed...")
		} else {
			port = portnumber
		}
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		fmt.Println("User cannot be empty")
		panic("DB_USER cannot be empty, application intialization failed...")
	}
	password := os.Getenv("DB_PASSWORD")
	if user == "" {
		fmt.Println("User cannot be empty")
		panic("DB_USER cannot be empty, application intialization failed...")
	}
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		fmt.Println("Database name cannot be empty")
		panic("DB_NAME cannot be empty, application intialization failed...")
	}
	sslmode := os.Getenv("DB_SSL_MODE")
	if sslmode == "" {
		fmt.Println("SSL mode cannot be empty")
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
	fmt.Println("Successfully connected to database!")
	// // Create Server and Route Handlers
	// r := mux.NewRouter()
	// r.Use(controllers.TrackResponseTime)

	srv := &http.Server{

		Addr:         ":8083",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Configure Logging
	logFileLocation := os.Getenv("LOG_FILE_LOCATION")
	if logFileLocation != "" {
		log.SetOutput(&lumberjack.Logger{
			Filename:   logFileLocation,
			MaxSize:    50, // megabytes
			MaxBackups: 3,
			MaxAge:     28,   //days
			Compress:   true, // disabled by default
		})
		log.Println("Successfully initialized log file...")
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
			Compress:   true, // disabled by default
		},
	}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, "Role", echo.HeaderAuthorization, "Refresh-Token", echo.HeaderXRealIP},
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

	// Methods that check application themselves and use the applicaiton information
	auth.POST("/:application/login", env.Login)
	auth.POST("/:application/otp/verify", env.VerifyOtp)

	// User operations
	apiAdminAuth.GET("/:application/user", env.GetUsers)
	apiNoAuth.GET("/:application/user/:username", env.CheckAvailability)
	apiAuth.GET("/:application/user/:username", env.GetUser)
	auth.POST("/:application/user", env.Register)
	apiAuth.PUT("/:application/user", env.UpdateUser)
	apiAuth.DELETE("/:application/user/:username", env.DeleteUser)

	// User Language operations
	apiAuth.GET("/:application/user/language/:username", env.GetUserLanguages)
	apiAuth.POST("/:application/user/language/:username/:language/:proficiency", env.AddUserLanguage)
	apiAuth.DELETE("/:application/user/language/:username/:language", env.DeleteUserLanguages)

	// User Role operations
	apiAuth.PUT("/:application/user/role/:newRole/:oldRole/:username", env.UpdateUserRole)
	apiAuth.POST("/:application/user/role/:role/:username", env.AddUserToRole)

	// Language operations
	apiAuth.GET("/:application/language/:language", env.GetLanguage)
	apiAuth.GET("/:application/language", env.GetLanguages)
	apiAdminAuth.POST("/:application/language/:language", env.AddLanguage)
	apiAdminAuth.PUT("/:application/language/:language/:newLanguage", env.UpdateLanguage)
	apiAdminAuth.DELETE("/:application/language/:language", env.DeleteLanguage)

	// Language proficiency operations
	apiAuth.GET("/:application/proficiency/:proficiency", env.GetLanguageProficiency)
	apiAuth.GET("/:application/proficiency", env.GetLanguageProficiencies)
	apiAdminAuth.POST("/:application/proficiency/:proficiency", env.AddLanguageProficiency)
	apiAdminAuth.PUT("/:application/proficiency/:proficiency/:newProficiency", env.UpdateLanguageProficiency)
	apiAdminAuth.DELETE("/:application/proficiency/:proficiency", env.DeleteLanguageProficiency)

	// Timezone Operations
	apiAuth.GET("/:application/timezone/:timezone", env.GetTimezone)
	apiAuth.GET("/:application/timezone", env.GetTimezones)
	apiAdminAuth.POST("/:application/timezone", env.AddTimezone)
	apiAdminAuth.PUT("/:application/timezone/:timezone", env.UpdateTimezone)
	apiAdminAuth.DELETE("/:application/timezone/:timezone", env.DeleteTimezone)

	// Application Operations
	apiAuth.GET("/:application/application/:application", env.GetApplication)
	apiAdminAuth.GET("/:application/application", env.GetApplications)
	apiAdminAuth.POST("/:application/application", env.AddApplication)
	apiAdminAuth.PUT("/:application/application/:application", env.UpdateApplication)
	apiAdminAuth.DELETE("/:application/application/:application", env.DeleteApplication)

	// Countries Operations
	apiAuth.GET("/:application/country/:country", env.GetCountry)
	apiAuth.GET("/:application/country", env.GetCountries)
	apiAdminAuth.POST("/:application/country", env.AddCountry)
	apiAdminAuth.PUT("/:application/country/:country", env.UpdateCountry)
	apiAdminAuth.DELETE("/:application/country/:country", env.DeleteCountry)

	// States Operations
	apiAuth.GET("/:application/state/:state", env.GetState)
	apiAuth.GET("/:application/state", env.GetStates)
	apiAdminAuth.POST("/:application/state/:state/:country", env.AddState)
	apiAdminAuth.PUT("/:application/state/:state/:newState", env.UpdateState)
	apiAdminAuth.DELETE("/:application/state/:state", env.DeleteState)
	apiAuth.GET("/:application/state/country/:country", env.GetStatesByCountry)

	// Roles Operations
	apiAdminAuth.GET("/:application/role/:role", env.GetRole)
	apiAdminAuth.GET("/:application/role", env.GetRoles)
	apiAdminAuth.POST("/:application/role", env.AddRole)
	apiAdminAuth.PUT("/:application/role/:role", env.UpdateRole)
	apiAdminAuth.DELETE("/:application/role/:role", env.DeleteRole)
	apiAdminAuth.GET("/:application/role/application/:application", env.GetRolesByApplication)
	apiAdminAuth.POST("/:application/role/:role/:application", env.AddApplicationRole)

	// go func() {
	// 	log.Println("Starting Server...")
	// 	if err := srv.ListenAndServe(); err != nil {
	// 		log.Fatal(err)
	// 	}
	// }()
	go func() {
		log.Println("Starting Server...")
		e.Logger.Fatal(e.StartServer(srv))
	}()

	// Graceful Shutdown
	waitForShutdown(srv)
}

func waitForShutdown(srv *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan
	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	srv.Shutdown(ctx)

	log.Println("Shutting down...")
	os.Exit(0)
}
