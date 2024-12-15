package middleware

import (
	"database/sql"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func PostgresMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dbString := os.Getenv("DB_CONNECTION")
		db, err := sql.Open("postgres", dbString)
		if err != nil {
			panic(err)
		}
		if err := db.Ping(); err != nil {
			panic(err)
		}
		db.SetMaxIdleConns(100)
		db.SetConnMaxLifetime(time.Hour)
		ctx.Set("postgresClient", db)
	}
}
