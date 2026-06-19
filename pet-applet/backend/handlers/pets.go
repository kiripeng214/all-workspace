package handlers

import (
	"database/sql"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"pet-applet-backend/database"
	"pet-applet-backend/models"
)

func generateID() string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 8)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

func GetPets(c *gin.Context) {
	rows, err := database.DB.Query("SELECT id, avatar, name, breed, birthday, weight, notes, created_at FROM pets ORDER BY created_at DESC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	pets := []models.Pet{}
	for rows.Next() {
		var p models.Pet
		if err := rows.Scan(&p.ID, &p.Avatar, &p.Name, &p.Breed, &p.Birthday, &p.Weight, &p.Notes, &p.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		pets = append(pets, p)
	}
	c.JSON(http.StatusOK, pets)
}

func GetPet(c *gin.Context) {
	id := c.Param("id")
	var p models.Pet
	err := database.DB.QueryRow("SELECT id, avatar, name, breed, birthday, weight, notes, created_at FROM pets WHERE id = ?", id).
		Scan(&p.ID, &p.Avatar, &p.Name, &p.Breed, &p.Birthday, &p.Weight, &p.Notes, &p.CreatedAt)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "宠物不存在"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}

type createPetInput struct {
	Avatar   string `json:"avatar"`
	Name     string `json:"name" binding:"required"`
	Breed    string `json:"breed"`
	Birthday string `json:"birthday"`
	Weight   string `json:"weight"`
	Notes    string `json:"notes"`
}

func CreatePet(c *gin.Context) {
	var input createPetInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少必填字段 name"})
		return
	}

	id := generateID()
	now := time.Now().UnixMilli()
	avatar := input.Avatar
	if avatar == "" {
		avatar = "🐾"
	}

	_, err := database.DB.Exec(
		"INSERT INTO pets (id, avatar, name, breed, birthday, weight, notes, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		id, avatar, input.Name, input.Breed, input.Birthday, input.Weight, input.Notes, now,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.Pet{
		ID:        id,
		Avatar:    avatar,
		Name:      input.Name,
		Breed:     input.Breed,
		Birthday:  input.Birthday,
		Weight:    input.Weight,
		Notes:     input.Notes,
		CreatedAt: now,
	})
}

type updatePetInput struct {
	Avatar   string `json:"avatar"`
	Name     string `json:"name"`
	Breed    string `json:"breed"`
	Birthday string `json:"birthday"`
	Weight   string `json:"weight"`
	Notes    string `json:"notes"`
}

func UpdatePet(c *gin.Context) {
	id := c.Param("id")

	var existing models.Pet
	err := database.DB.QueryRow("SELECT id, avatar, name, breed, birthday, weight, notes FROM pets WHERE id = ?", id).
		Scan(&existing.ID, &existing.Avatar, &existing.Name, &existing.Breed, &existing.Birthday, &existing.Weight, &existing.Notes)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "宠物不存在"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var input updatePetInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据无效"})
		return
	}

	avatar := input.Avatar
	if avatar == "" {
		avatar = existing.Avatar
	}
	name := input.Name
	if name == "" {
		name = existing.Name
	}
	breed := input.Breed
	if breed == "" {
		breed = existing.Breed
	}
	birthday := input.Birthday
	if birthday == "" {
		birthday = existing.Birthday
	}
	weight := input.Weight
	if weight == "" {
		weight = existing.Weight
	}
	notes := input.Notes
	if notes == "" {
		notes = existing.Notes
	}

	_, err = database.DB.Exec(
		"UPDATE pets SET avatar=?, name=?, breed=?, birthday=?, weight=?, notes=? WHERE id=?",
		avatar, name, breed, birthday, weight, notes, id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.Pet{
		ID:        id,
		Avatar:    avatar,
		Name:      name,
		Breed:     breed,
		Birthday:  birthday,
		Weight:    weight,
		Notes:     notes,
		CreatedAt: existing.CreatedAt,
	})
}

func DeletePet(c *gin.Context) {
	id := c.Param("id")

	var exists bool
	if err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM pets WHERE id=?)", id).Scan(&exists); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "宠物不存在"})
		return
	}

	// 外键 ON DELETE CASCADE 会自动删除 feeding_schedules 和 feeding_records
	if _, err := database.DB.Exec("DELETE FROM pets WHERE id = ?", id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "已删除"})
}
