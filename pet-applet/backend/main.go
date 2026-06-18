package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"pet-applet-backend/config"
	"pet-applet-backend/database"
	"pet-applet-backend/handlers"
)

func main() {
	cfg := config.Load()
	database.Init(cfg)

	r := gin.Default()

	api := r.Group("/api")
	{
		// Pets
		api.GET("/pets", handlers.GetPets)
		api.GET("/pets/:id", handlers.GetPet)
		api.POST("/pets", handlers.CreatePet)
		api.PUT("/pets/:id", handlers.UpdatePet)
		api.DELETE("/pets/:id", handlers.DeletePet)

		// Feeding Schedules
		api.GET("/pets/schedules/:petId", handlers.GetSchedules)
		api.POST("/pets/schedules/:petId", handlers.CreateSchedule)
		api.PUT("/schedules/:id", handlers.UpdateSchedule)
		api.DELETE("/schedules/:id", handlers.DeleteSchedule)

		// Meta
		api.GET("/meta/breeds", handlers.GetBreeds)

		// Feeding Records
		api.GET("/pets/records/:petId", handlers.GetRecords)
		api.GET("/pets/records/today/:petId", handlers.GetTodayRecords)
		api.POST("/pets/records/:petId", handlers.CreateRecord)
		api.DELETE("/records/:id", handlers.DeleteRecord)
	}

	log.Printf("🐾 服务启动于 http://localhost:%s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("启动失败: %v", err)
	}
}
