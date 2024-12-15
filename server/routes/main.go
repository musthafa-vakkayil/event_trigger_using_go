package routes

import (
	"context"
	"database/sql"
	config "event_trigger/config/db"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
)

type Api struct {
	Engine *gin.Engine
	srv    *http.Server
}

func (a *Api) Init(version string) (*gin.Engine, error) {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	a.Engine = r
	a.AddReotes(version)
	a.WithServer()
	a.MakeMigrations()
	return r, nil
}

func (a *Api) AddReotes(version string) {
	HealthRoutes(a.Engine, version)
	UserRoutes(a.Engine, version)
	TriggerRoutes(a.Engine, version)
}

func (a *Api) MakeMigrations() {
	db, err := sql.Open("postgres", "postgres://myuser:mypassword@postgres_db:5432/eventdb?sslmode=disable")
	if err != nil {
		log.Fatalf("Error connecting to DB: %v", err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging DB: %v", err)
	}

	exprs := []string{config.CREATE_TABLE}

	for _, expr := range exprs {
		log.Printf("Executing SQL: %s\n", expr)
		if _, err := db.ExecContext(context.Background(), expr); err != nil {
			log.Printf("Error executing query: %v", err)
		} else {
			log.Println("Query executed successfully.")
		}
	}
}

func (a *Api) WithServer() Api {
	ch := make(chan Api)

	go func() {
		a.srv = &http.Server{
			Addr:         fmt.Sprintf(":%s", "80"),
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
