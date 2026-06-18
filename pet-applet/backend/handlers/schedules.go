package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"pet-applet-backend/database"
	"pet-applet-backend/models"
)

func GetSchedules(c *gin.Context) {
	petID := c.Param("petId")

	var exists bool
	database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM pets WHERE id=?)", petID).Scan(&exists)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "宠物不存在"})
		return
	}

	rows, err := database.DB.Query("SELECT id, pet_id, time, food_type, amount FROM feeding_schedules WHERE pet_id = ? ORDER BY time", petID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	schedules := []models.FeedingSchedule{}
	for rows.Next() {
		var s models.FeedingSchedule
		if err := rows.Scan(&s.ID, &s.PetID, &s.Time, &s.FoodType, &s.Amount); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		schedules = append(schedules, s)
	}
	c.JSON(http.StatusOK, schedules)
}

type createScheduleInput struct {
	Time     string `json:"time" binding:"required"`
	FoodType string `json:"foodType"`
	Amount   string `json:"amount"`
}

func CreateSchedule(c *gin.Context) {
	petID := c.Param("petId")

	var exists bool
	database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM pets WHERE id=?)", petID).Scan(&exists)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "宠物不存在"})
		return
	}

	var input createScheduleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少必填字段 time"})
		return
	}

	id := generateID()
	foodType := input.FoodType
	if foodType == "" {
		foodType = "粮食"
	}
	amount := input.Amount
	if amount == "" {
		amount = "一份"
	}

	_, err := database.DB.Exec(
		"INSERT INTO feeding_schedules (id, pet_id, time, food_type, amount) VALUES (?, ?, ?, ?, ?)",
		id, petID, input.Time, foodType, amount,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.FeedingSchedule{
		ID:       id,
		PetID:    petID,
		Time:     input.Time,
		FoodType: foodType,
		Amount:   amount,
	})
}

type updateScheduleInput struct {
	Time     string `json:"time"`
	FoodType string `json:"foodType"`
	Amount   string `json:"amount"`
}

func UpdateSchedule(c *gin.Context) {
	id := c.Param("id")

	var existing models.FeedingSchedule
	err := database.DB.QueryRow("SELECT id, pet_id, time, food_type, amount FROM feeding_schedules WHERE id = ?", id).
		Scan(&existing.ID, &existing.PetID, &existing.Time, &existing.FoodType, &existing.Amount)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "喂养时间不存在"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var input updateScheduleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据无效"})
		return
	}

	time := input.Time
	if time == "" {
		time = existing.Time
	}
	foodType := input.FoodType
	if foodType == "" {
		foodType = existing.FoodType
	}
	amount := input.Amount
	if amount == "" {
		amount = existing.Amount
	}

	_, err = database.DB.Exec("UPDATE feeding_schedules SET time=?, food_type=?, amount=? WHERE id=?",
		time, foodType, amount, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.FeedingSchedule{
		ID:       id,
		PetID:    existing.PetID,
		Time:     time,
		FoodType: foodType,
		Amount:   amount,
	})
}

func DeleteSchedule(c *gin.Context) {
	id := c.Param("id")
	database.DB.Exec("DELETE FROM feeding_schedules WHERE id = ?", id)
	c.JSON(http.StatusOK, gin.H{"message": "已删除"})
}
