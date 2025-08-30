package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"mojo-autotech/config"
	a "mojo-autotech/handler/attedance"
	h "mojo-autotech/handler/user_authentication"

	"mojo-autotech/model/attedance"
	"mojo-autotech/model/user_authentication"
)

func main() {
	cfg := config.NewConfig()
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	db := config.NewDB()
	_ = db.AutoMigrate(&user_authentication.User{})
	_ = db.AutoMigrate(&attedance.Attendance{})

	h.HttpHandler(router)
	a.HttpAttendanceHandler(router)

	addr := cfg.Srv.Host + ":" + cfg.Srv.Port
	log.Println("Server running at http://" + addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
