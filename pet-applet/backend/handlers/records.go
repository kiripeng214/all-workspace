package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"pet-applet-backend/database"
	"pet-applet-backend/models"
)

func GetRecords(c *gin.Context) {
	petID := c.Param("petId")

	var exists bool
	database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM pets WHERE id=?)", petID).Scan(&exists)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "宠物不存在"})
		return
	}

	rows, err := database.DB.Query(
		"SELECT id, pet_id, schedule_id, time, food_type, amount, notes, created_at FROM feeding_records WHERE pet_id = ? ORDER BY created_at DESC",
		petID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	records := []models.FeedingRecord{}
	for rows.Next() {
		var r models.FeedingRecord
		var scheduleID sql.NullString
		if err := rows.Scan(&r.ID, &r.PetID, &scheduleID, &r.Time, &r.FoodType, &r.Amount, &r.Notes, &r.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if scheduleID.Valid {
			r.ScheduleID = scheduleID.String
		}
		records = append(records, r)
	}
	c.JSON(http.StatusOK, records)
}

func GetTodayRecords(c *gin.Context) {
	petID := c.Param("petId")

	var exists bool
	database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM pets WHERE id=?)", petID).Scan(&exists)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "宠物不存在"})
		return
	}

	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).UnixMilli()
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999, now.Location()).UnixMilli()

	rows, err := database.DB.Query(
		"SELECT id, pet_id, schedule_id, time, food_type, amount, notes, created_at FROM feeding_records WHERE pet_id = ? AND created_at >= ? AND created_at <= ? ORDER BY created_at DESC",
		petID, startOfDay, endOfDay,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	records := []models.FeedingRecord{}
	for rows.Next() {
		var r models.FeedingRecord
		var scheduleID sql.NullString
		if err := rows.Scan(&r.ID, &r.PetID, &scheduleID, &r.Time, &r.FoodType, &r.Amount, &r.Notes, &r.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if scheduleID.Valid {
			r.ScheduleID = scheduleID.String
		}
		records = append(records, r)
	}
	c.JSON(http.StatusOK, records)
}

type createRecordInput struct {
	ScheduleID string `json:"scheduleId"`
	Time       string `json:"time"`
	FoodType   string `json:"foodType"`
	Amount     string `json:"amount"`
	Notes      string `json:"notes"`
}

func CreateRecord(c *gin.Context) {
	petID := c.Param("petId")

	var exists bool
	database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM pets WHERE id=?)", petID).Scan(&exists)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "宠物不存在"})
		return
	}

	var input createRecordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据无效"})
		return
	}

	id := generateID()
	now := time.Now().UnixMilli()
	foodType := input.FoodType
	if foodType == "" {
		foodType = "粮食"
	}
	amount := input.Amount
	if amount == "" {
		amount = "一份"
	}

	_, err := database.DB.Exec(
		"INSERT INTO feeding_records (id, pet_id, schedule_id, time, food_type, amount, notes, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		id, petID, input.ScheduleID, input.Time, foodType, amount, input.Notes, now,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.FeedingRecord{
		ID:         id,
		PetID:      petID,
		ScheduleID: input.ScheduleID,
		Time:       input.Time,
		FoodType:   foodType,
		Amount:     amount,
		Notes:      input.Notes,
		CreatedAt:  now,
	})
}

func DeleteRecord(c *gin.Context) {
	id := c.Param("id")
	if _, err := database.DB.Exec("DELETE FROM feeding_records WHERE id = ?", id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已删除"})
}
