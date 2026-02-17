// main.go
package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/JIeeiroSst/hub/handler"
	"github.com/JIeeiroSst/hub/repository"
	"github.com/JIeeiroSst/hub/service"
	ws "github.com/JIeeiroSst/hub/websocket"
)

func main() {
	dbType := repository.DBType(getEnv("DB_TYPE", "postgres"))

	cfg := repository.DBConfig{
		Type: dbType,
		DSN: getEnv("DB_DSN", "host=localhost user=postgres password=postgres dbname=items_db port=5432 sslmode=disable"),

		MongoURI:      getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDBName:   getEnv("MONGO_DB", "items_db"),
		MongoCollName: getEnv("MONGO_COLL", "items"),
	}

	repo, err := repository.NewRepository(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize repository [%s]: %v", dbType, err)
	}
	log.Printf("✅ Database strategy initialized: %s", dbType)

	hub := ws.NewHub()
	go hub.Run() 

	svc := service.NewItemService(repo, hub)
	h := handler.NewItemHandler(svc, hub)

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	v1 := r.Group("/api/v1")
	{
		v1.POST("/items", h.Create)     // Tạo item → tự broadcast real-time
		v1.GET("/items", h.List)        // Lấy list, sort created_at DESC
		v1.GET("/items/:id", h.GetByID) // Lấy theo ID
		v1.GET("/health", h.Health)     // Health check
	}

	r.GET("/ws", h.WebSocket)

	port := getEnv("PORT", "8080")
	log.Printf("🚀 Server running on :%s", port)
	log.Printf("📡 WebSocket endpoint: ws://localhost:%s/ws", port)
	log.Printf("📋 List API: GET http://localhost:%s/api/v1/items", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
