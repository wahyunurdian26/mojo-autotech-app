package main

import (
	"log"
	"mojo-autotech/config"
	a "mojo-autotech/handler/attedance"
	h "mojo-autotech/handler/user_authentication"
	"mojo-autotech/model/user_authentication"

	"mojo-autotech/model/attedance"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.NewConfig()
	router := gin.Default()
	db := config.NewDB()
	db.AutoMigrate(&user_authentication.User{})
	db.AutoMigrate(&attedance.Attendance{})

	h.HttpHandler(router)
	a.HttpAttendanceHandler(router)

	if err := router.Run(cfg.Srv.Host + ":" + cfg.Srv.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
