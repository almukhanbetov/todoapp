package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Todo struct {
	ID   int    `json:"id"`
	Task string `json:"task"`
}

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	// Подключение к БД
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbpool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("unable to connect to db: %v", err)
	}
	defer dbpool.Close()

	r := gin.Default()

	// healthcheck
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// получить все задачи
	r.GET("/todos", func(c *gin.Context) {
		rows, err := dbpool.Query(ctx, "SELECT id, task FROM todos")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var todos []Todo
		for rows.Next() {
			var t Todo
			if err := rows.Scan(&t.ID, &t.Task); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			todos = append(todos, t)
		}
		c.JSON(http.StatusOK, todos)
	})

	// добавить задачу
	r.POST("/todos", func(c *gin.Context) {
		var t Todo
		if err := c.ShouldBindJSON(&t); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
			return
		}

		err := dbpool.QueryRow(ctx,
			"INSERT INTO todos (task) VALUES ($1) RETURNING id", t.Task).Scan(&t.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, t)
	})

	log.Println("todoapp started on :8081")
	r.Run(":8081")
}
