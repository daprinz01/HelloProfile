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

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"gopkg.in/natefinch/lumberjack.v2"
)

type name struct {
	Name int `json:"name"`
}

// Env is used to export things that can be reused
// type Env struct {
// 	AuthDb *authdb.Queries
// }

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
	// Create Server and Route Handlers
	r := mux.NewRouter()
	apiNoAuth := r.PathPrefix("/api/v1").Subrouter()
	auth := r.PathPrefix("/api/v1/auth").Subrouter()
	// Methods that don't require the application information is it's being verified in a middleware
	apiNoAuth.Use(env.CheckApplication)
	apiNoAuth.HandleFunc("/{application}/refresh", env.RefreshToken).Methods(http.MethodGet)
	apiNoAuth.HandleFunc("/{application}/otp/send", env.SendOtp).Methods(http.MethodPost)
	apiNoAuth.HandleFunc("/{application}/password/reset", env.ResetPassword).Methods(http.MethodPost)

	// Methods that check application themselves and use the applicaiton information
	auth.HandleFunc("/{application}/login", env.Login).Methods(http.MethodPost)
	auth.HandleFunc("/{application}/otp/verify", env.VerifyOtp).Methods(http.MethodPost)

	// Methods that require authentication but don't need the information from the applications or authorization to function
	apiAuth := apiNoAuth.PathPrefix("/app").Subrouter()
	apiAuth.Use(controllers.Authorize)

	// Admin operations authorization
	apiAdminAuth := apiNoAuth.PathPrefix("/admin").Subrouter()
	apiAdminAuth.Use(controllers.AuthorizeAdmin)

	// User operations
	apiAdminAuth.HandleFunc("/{application}/user", env.GetUsers).Methods(http.MethodGet)
	apiNoAuth.HandleFunc("/{application}/user/{username}", env.CheckAvailability).Methods(http.MethodGet)
	apiAuth.HandleFunc("/{application}/user/{username}", env.GetUser).Methods(http.MethodGet)
	auth.HandleFunc("/{application}/user", env.Register).Methods(http.MethodPost)
	apiAuth.HandleFunc("/{application}/user", env.UpdateUser).Methods(http.MethodPut)
	apiAuth.HandleFunc("/{application}/user/{username}", env.DeleteUser).Methods(http.MethodDelete)

	// User Language operations
	apiAuth.HandleFunc("/{application}/user/language/{username}", env.GetUserLanguages).Methods(http.MethodGet)
	apiAuth.HandleFunc("/{application}/user/language/{username}/{language}/{proficiency}", env.AddUserLanguage).Methods(http.MethodPost)
	apiAuth.HandleFunc("/{application}/user/language/{username}/{language}", env.DeleteUserLanguages).Methods(http.MethodDelete)

	// Language operations
	apiAuth.HandleFunc("/{application}/language/{language}", env.GetLanguage).Methods(http.MethodGet)
	apiAuth.HandleFunc("/{application}/language", env.GetLanguages).Methods(http.MethodGet)
	apiAdminAuth.HandleFunc("/{application}/language/{language}", env.AddLanguage).Methods(http.MethodPost)
	apiAdminAuth.HandleFunc("/{application}/language/{language}/{newLanguage}", env.UpdateLanguage).Methods(http.MethodPut)
	apiAdminAuth.HandleFunc("/{application}/language/{language}", env.DeleteLanguage).Methods(http.MethodDelete)

	// Language proficiency operations
	apiAuth.HandleFunc("/{application}/proficiency/{proficiency}", env.GetLanguageProficiency).Methods(http.MethodGet)
	apiAuth.HandleFunc("/{application}/proficiency", env.GetLanguageProficiencies).Methods(http.MethodGet)
	apiAdminAuth.HandleFunc("/{application}/proficiency/{proficiency}", env.AddLanguageProficiency).Methods(http.MethodPost)
	apiAdminAuth.HandleFunc("/{application}/proficiency/{proficiency}/{newProficiency}", env.UpdateLanguageProficiency).Methods(http.MethodPut)
	apiAdminAuth.HandleFunc("/{application}/proficiency/{proficiency}", env.DeleteLanguageProficiency).Methods(http.MethodDelete)

	// Timezone Operations
	apiAuth.HandleFunc("/{application}/timezone/{timezone}", env.GetTimezone).Methods(http.MethodGet)
	apiAuth.HandleFunc("/{application}/timezone", env.GetTimezones).Methods(http.MethodGet)
	apiAdminAuth.HandleFunc("/{application}/timezone", env.AddTimezone).Methods(http.MethodPost)
	apiAdminAuth.HandleFunc("/{application}/timezone/{timezone}", env.UpdateTimezone).Methods(http.MethodPut)
	apiAdminAuth.HandleFunc("/{application}/timezone/{timezone}", env.DeleteTimezone).Methods(http.MethodDelete)

	// Application Operations
	apiAdminAuth.HandleFunc("/{application}/application/{application}", env.GetApplication).Methods(http.MethodGet)
	apiAdminAuth.HandleFunc("/{application}/application", env.GetApplications).Methods(http.MethodGet)
	apiAdminAuth.HandleFunc("/{application}/application", env.AddApplication).Methods(http.MethodPost)
	apiAdminAuth.HandleFunc("/{application}/application/{application}", env.UpdateApplication).Methods(http.MethodPut)
	apiAdminAuth.HandleFunc("/{application}/application/{application}", env.DeleteApplication).Methods(http.MethodDelete)

	// Countries Operations
	apiAuth.HandleFunc("/{application}/country/{country}", env.GetCountry).Methods(http.MethodGet)
	apiAuth.HandleFunc("/{application}/country", env.GetCountries).Methods(http.MethodGet)
	apiAdminAuth.HandleFunc("/{application}/country", env.AddCountry).Methods(http.MethodPost)
	apiAdminAuth.HandleFunc("/{application}/country/{country}", env.UpdateCountry).Methods(http.MethodPut)
	apiAdminAuth.HandleFunc("/{application}/country/{country}", env.DeleteCountry).Methods(http.MethodDelete)

	// States Operations
	apiAuth.HandleFunc("/{application}/state/{state}", env.GetState).Methods(http.MethodGet)
	apiAuth.HandleFunc("/{application}/state", env.GetStates).Methods(http.MethodGet)
	apiAdminAuth.HandleFunc("/{application}/state/{state}/{country}", env.AddState).Methods(http.MethodPost)
	apiAdminAuth.HandleFunc("/{application}/state/{state}/{newState}", env.UpdateState).Methods(http.MethodPut)
	apiAdminAuth.HandleFunc("/{application}/state/{state}", env.DeleteState).Methods(http.MethodDelete)
	apiAuth.HandleFunc("/{application}/state/country/{country}", env.GetStatesByCountry).Methods(http.MethodGet)

	// Roles Operations
	apiAdminAuth.HandleFunc("/{application}/role/{role}", env.GetRole).Methods(http.MethodGet)
	apiAdminAuth.HandleFunc("/{application}/role", env.GetRoles).Methods(http.MethodGet)
	apiAdminAuth.HandleFunc("/{application}/role", env.AddRole).Methods(http.MethodPost)
	apiAdminAuth.HandleFunc("/{application}/role/{role}", env.UpdateRole).Methods(http.MethodPut)
	apiAdminAuth.HandleFunc("/{application}/role/{role}", env.DeleteRole).Methods(http.MethodDelete)
	apiAdminAuth.HandleFunc("/{application}/role/application/{application}", env.GetRolesByApplication).Methods(http.MethodGet)
	apiAdminAuth.HandleFunc("/{application}/role/{role}/{application}", env.AddApplicationRole).Methods(http.MethodPost)

	srv := &http.Server{
		Handler:      r,
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
	go func() {
		log.Println("Starting Server...")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
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
