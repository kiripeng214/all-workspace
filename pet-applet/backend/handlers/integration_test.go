package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"pet-applet-backend/config"
	"pet-applet-backend/database"
)

var testRouter *gin.Engine
var testPetID string

func TestMain(m *testing.M) {
	// 初始化测试数据库
	cfg := &config.Config{
		DB: config.DBConfig{
			Host:     "localhost",
			Port:     "3306",
			User:     "root",
			Password: "123456",
			Name:     "pet_applet_test",
		},
		Server: config.Server{
			Port: "3000",
		},
	}

	// 先创建测试数据库
	rootCfg := *cfg
	rootCfg.DB.Name = ""

	tmpDB, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=true&loc=Local",
		rootCfg.DB.User, rootCfg.DB.Password, rootCfg.DB.Host, rootCfg.DB.Port))
	if err != nil {
		log.Printf("跳过集成测试: 数据库不可用 (%v)", err)
		os.Exit(0)
	}
	tmpDB.Exec("CREATE DATABASE IF NOT EXISTS pet_applet_test CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci")
	tmpDB.Close()

	// 初始化应用数据库
	database.Init(cfg)

	// 创建路由
	gin.SetMode(gin.TestMode)
	testRouter = gin.Default()
	api := testRouter.Group("/api")
	{
		api.GET("/pets", GetPets)
		api.GET("/pets/:id", GetPet)
		api.POST("/pets", CreatePet)
		api.PUT("/pets/:id", UpdatePet)
		api.DELETE("/pets/:id", DeletePet)
		api.GET("/pets/schedules/:petId", GetSchedules)
		api.POST("/pets/schedules/:petId", CreateSchedule)
		api.PUT("/schedules/:id", UpdateSchedule)
		api.DELETE("/schedules/:id", DeleteSchedule)
		api.GET("/meta/breeds", GetBreeds)
		api.GET("/pets/records/:petId", GetRecords)
		api.GET("/pets/records/today/:petId", GetTodayRecords)
		api.POST("/pets/records/:petId", CreateRecord)
		api.DELETE("/records/:id", DeleteRecord)
	}

	code := m.Run()

	// 清理测试数据
	cleanup()
	os.Exit(code)
}

func cleanup() {
	database.DB.Exec("DROP TABLE IF EXISTS feeding_records")
	database.DB.Exec("DROP TABLE IF EXISTS feeding_schedules")
	database.DB.Exec("DROP TABLE IF EXISTS pets")
	database.DB.Exec("DROP DATABASE IF EXISTS pet_applet_test")
	database.DB.Close()
}

func jsonRequest(method, path string, body any) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		b, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(b)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req, _ := http.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	return w
}

func TestCreateAndGetPet(t *testing.T) {
	// 创建宠物
	w := jsonRequest("POST", "/api/pets", map[string]string{
		"name":   "旺财",
		"avatar": "🐶",
		"breed":  "金毛",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("创建宠物预期 201，实际 %d: %s", w.Code, w.Body.String())
	}

	var created map[string]any
	json.Unmarshal(w.Body.Bytes(), &created)
	testPetID = created["id"].(string)

	if created["name"] != "旺财" {
		t.Errorf("宠物名预期 '旺财'，实际 %v", created["name"])
	}
}

func TestGetPets_EmptyList(t *testing.T) {
	// 先清理所有
	database.DB.Exec("DELETE FROM pets WHERE id != ?", testPetID)

	w := jsonRequest("GET", "/api/pets", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("获取列表预期 200，实际 %d", w.Code)
	}

	var pets []any
	json.Unmarshal(w.Body.Bytes(), &pets)
	if len(pets) == 0 {
		t.Error("期望列表非空")
	}
}

func TestUpdatePet_Name(t *testing.T) {
	if testPetID == "" {
		t.Skip("需要先创建宠物")
	}

	w := jsonRequest("PUT", "/api/pets/"+testPetID, map[string]string{
		"name": "小强",
	})
	if w.Code != http.StatusOK {
		t.Fatalf("更新宠物预期 200，实际 %d: %s", w.Code, w.Body.String())
	}

	var updated map[string]any
	json.Unmarshal(w.Body.Bytes(), &updated)
	if updated["name"] != "小强" {
		t.Errorf("更新后名称预期 '小强'，实际 %v", updated["name"])
	}
}

func TestUpdatePet_Partial(t *testing.T) {
	if testPetID == "" {
		t.Skip("需要先创建宠物")
	}

	// 只更新 breed，其他字段不应受影响
	w := jsonRequest("PUT", "/api/pets/"+testPetID, map[string]string{
		"breed": "拉布拉多",
	})
	if w.Code != http.StatusOK {
		t.Fatalf("局部更新预期 200，实际 %d", w.Code)
	}

	var updated map[string]any
	json.Unmarshal(w.Body.Bytes(), &updated)
	if updated["breed"] != "拉布拉多" {
		t.Errorf("品种预期 '拉布拉多'，实际 %v", updated["breed"])
	}
	if updated["name"] != "小强" {
		t.Errorf("name 应保持 '小强'，实际 %v", updated["name"])
	}
}

func TestGetPet_NotFound(t *testing.T) {
	w := jsonRequest("GET", "/api/pets/nonexist", nil)
	if w.Code != http.StatusNotFound {
		t.Errorf("不存在的宠物预期 404，实际 %d", w.Code)
	}
}

func TestCreatePet_MissingName(t *testing.T) {
	w := jsonRequest("POST", "/api/pets", map[string]string{
		"avatar": "🐱",
	})
	if w.Code != http.StatusBadRequest {
		t.Errorf("缺少 name 预期 400，实际 %d", w.Code)
	}
}

func TestCreateSchedule(t *testing.T) {
	if testPetID == "" {
		t.Skip("需要先创建宠物")
	}

	w := jsonRequest("POST", "/api/pets/schedules/"+testPetID, map[string]string{
		"time":     "08:00",
		"foodType": "猫粮",
		"amount":   "一份",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("创建计划预期 201，实际 %d: %s", w.Code, w.Body.String())
	}
}

func TestGetSchedules(t *testing.T) {
	if testPetID == "" {
		t.Skip("需要先创建宠物")
	}

	w := jsonRequest("GET", "/api/pets/schedules/"+testPetID, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("获取计划预期 200，实际 %d", w.Code)
	}
}

func TestGetBreeds(t *testing.T) {
	w := jsonRequest("GET", "/api/meta/breeds", nil)
	if w.Code != http.StatusOK {
		t.Errorf("获取品种预期 200，实际 %d", w.Code)
	}

	var result map[string]any
	json.Unmarshal(w.Body.Bytes(), &result)
	if result["petEmojis"] == nil {
		t.Error("品种列表应包含 petEmojis")
	}
}

func TestCreateAndDeletePet(t *testing.T) {
	// 创建临时宠物
	w := jsonRequest("POST", "/api/pets", map[string]string{
		"name": "临时宠物",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("创建宠物预期 201，实际 %d", w.Code)
	}

	var created map[string]any
	json.Unmarshal(w.Body.Bytes(), &created)
	petID := created["id"].(string)

	// 删除它
	w = jsonRequest("DELETE", "/api/pets/"+petID, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("删除宠物预期 200，实际 %d", w.Code)
	}

	// 验证已删除
	w = jsonRequest("GET", "/api/pets/"+petID, nil)
	if w.Code != http.StatusNotFound {
		t.Errorf("删除后查询应 404，实际 %d", w.Code)
	}
}
