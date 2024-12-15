package routes

import (
	"context"
	"database/sql"
	config "event_trigger/config/db"
	"event_trigger/goroutines"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
)

type Api struct {
	Engine *gin.Engine
	srv    *http.Server
}

func (a *Api) Init(version string) (*gin.Engine, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using system environment variables")
	}
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	a.Engine = r
	a.AddReotes(version)
	a.WithServer()

	// Run migrations and initialize Go routines
	db := a.MakeMigrations()
	a.InitGoRoutines(db)

	return r, nil
}

func (a *Api) AddReotes(version string) {
	HealthRoutes(a.Engine, version)
	UserRoutes(a.Engine, version)
	TriggerRoutes(a.Engine, version)
	EventRoutes(a.Engine, version)

	a.Engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func (a *Api) MakeMigrations() *sql.DB {
	dbString := os.Getenv("DB_CONNECTION")
	db, err := sql.Open("postgres", dbString)
	if err != nil {
		log.Fatalf("Error connecting to DB: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging DB: %v", err)
	}

	exprs := []string{config.CREATE_TABLE, config.CREATE_TRIGGER}

	for _, expr := range exprs {
		log.Printf("Executing SQL: %s\n", expr)
		if _, err := db.ExecContext(context.Background(), expr); err != nil {
			log.Printf("Error executing query: %v", err)
		} else {
			log.Println("Query executed successfully.")
		}
	}

	return db
}

func (a *Api) WithServer() Api {
	ch := make(chan Api)

	go func() {
		a.srv = &http.Server{
			Addr:         fmt.Sprintf(":%s", "8080"),
			Handler:      a.Engine,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		}

		ch <- *a
	}()

	return <-ch
}

func (a *Api) Start(r *gin.Engine) {
	// build and start server
	go func() {
		if err := a.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Panicf("listen: %s\n", err)
		}
	}()

	// Wait for uinterrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)

	// kill (no param) default send syncall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall, SIGKILL but can't be caught, so son't need it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := a.srv.Shutdown(ctx); err != nil {
		log.Panicln("Server Shutfdown: ...", err)
	}

	// catching ctx.Done() timeout of 5 second
	<-ctx.Done()
	log.Println("Server Existing...")
}

func New(apiVersion string) Api {
	api := Api{}
	r, err := api.Init(apiVersion)
	if err != nil {
		panic(err)
	}

	api.Engine = r
	return api
}

func (a *Api) InitGoRoutines(db *sql.DB) {
	// Create channels
	triggerChannel := make(chan string)
	archiveChannel := make(chan string)

	// Start Go routines
	go goroutines.ListenForTriggers(db, triggerChannel)
	go goroutines.ProcessTriggers(triggerChannel, db)
	go goroutines.ArchiveLogs(db, archiveChannel)
	go goroutines.ProcessArchiving(db, archiveChannel)
	go goroutines.DeleteLogs(db)

	log.Println("Go routines initialized and running...")
}
