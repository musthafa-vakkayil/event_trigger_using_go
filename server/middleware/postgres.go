package middleware

import (
	"database/sql"
	"time"

	"github.com/gin-gonic/gin"
)

func PostgresMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		db, err := sql.Open("postgres", "postgres://myuser:mypassword@postgres_db:5432/eventdb?sslmode=disable")
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
