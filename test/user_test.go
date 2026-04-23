package test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"temukan-api/internal/middleware"
	"testing"

	"temukan-api/internal/handler"
	"temukan-api/internal/helper"
	"temukan-api/internal/model"
	"temukan-api/internal/repository"
	"temukan-api/internal/usecase"

	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ── Globals ────────────────────────────────────────────────────────────────

var (
	testDB     *gorm.DB
	testRouter *gin.Engine
)

// ── Setup ──────────────────────────────────────────────────────────────────

func setupTestDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME_TEST"), // pastikan ada DB_NAME_TEST di .env
		os.Getenv("DB_PORT"),
		os.Getenv("DB_TIMEZONE"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "pgcrypto";`).Error; err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&model.User{}); err != nil {
		return nil, err
	}

	return db, nil
}

func setupRouter(db *gorm.DB) *gin.Engine {
	validate := validator.New()

	userRepo := repository.NewUserRepository(db)
	uc := usecase.NewUserUsecase(userRepo, validate)
	h := handler.NewUserHandlerImpl(uc)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(middleware.ErrorRecovery())

	api := r.Group("/api/v1")
	api.POST("/users", h.Create)

	return r
}

func truncateTables(db *gorm.DB) {
	db.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")
}

func TestMain(m *testing.M) {
	if err := godotenv.Load("../.env.test"); err != nil {
		panic("Failed to load .env.test: " + err.Error())
	}

	db, err := setupTestDB()
	if err != nil {
		panic(err)
	}
	testDB = db
	testRouter = setupRouter(testDB)

	code := m.Run()
	os.Exit(code)
}

// ── Helpers ────────────────────────────────────────────────────────────────

func doPost(router *gin.Engine, url, body string) *http.Response {
	req := httptest.NewRequest(http.MethodPost, url, strings.NewReader(body))
	req.Header.Add("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec.Result()
}

func parseBody(r *http.Response) map[string]any {
	b, _ := io.ReadAll(r.Body)
	var m map[string]any
	json.Unmarshal(b, &m)
	return m
}

// ── Register (Create User) ─────────────────────────────────────────────────

// 1. Happy path — semua field valid termasuk phone (opsional)
func TestRegisterSuccess(t *testing.T) {
	truncateTables(testDB)

	resp := doPost(testRouter, "/api/v1/users", `{
		"name":     "Budi Santoso",
		"email":    "budi@mail.com",
		"password": "rahasia123",
		"role":     "finder"
	}`)
	body := parseBody(resp)

	assert.Equal(t, 201, resp.StatusCode)
	assert.Equal(t, "OK", body["status"])
	assert.Equal(t, "User created successfully", body["message"])

	data := body["data"].(map[string]any)
	assert.Equal(t, "Budi Santoso", data["name"])
	assert.Equal(t, "budi@mail.com", data["email"])
	assert.Equal(t, "finder", data["role"])
	assert.NotEmpty(t, data["id"])
	assert.Nil(t, data["password"]) // password tidak boleh muncul di response
}

// 2. Happy path — dengan phone opsional
func TestRegisterSuccessWithPhone(t *testing.T) {
	truncateTables(testDB)

	resp := doPost(testRouter, "/api/v1/users", `{
		"name":     "Ani Rahayu",
		"email":    "ani@mail.com",
		"password": "rahasia123",
		"role":     "seeker",
		"phone":    "08123456789"
	}`)
	body := parseBody(resp)

	assert.Equal(t, 201, resp.StatusCode)
	assert.Equal(t, "OK", body["status"])

	data := body["data"].(map[string]any)
	assert.Equal(t, "Ani Rahayu", data["name"])
	assert.Equal(t, "seeker", data["role"])
}

// 3. Happy path — role volunteer
func TestRegisterSuccessRoleVolunteer(t *testing.T) {
	truncateTables(testDB)

	resp := doPost(testRouter, "/api/v1/users", `{
		"name":     "Candra",
		"email":    "candra@mail.com",
		"password": "rahasia123",
		"role":     "volunteer"
	}`)
	body := parseBody(resp)

	assert.Equal(t, 201, resp.StatusCode)
	assert.Equal(t, "OK", body["status"])

	data := body["data"].(map[string]any)
	assert.Equal(t, "volunteer", data["role"])
}

// 4. Password tidak boleh ada di response
func TestRegisterPasswordNotExposed(t *testing.T) {
	truncateTables(testDB)

	resp := doPost(testRouter, "/api/v1/users", `{
		"name":     "Dian",
		"email":    "dian@mail.com",
		"password": "rahasia123",
		"role":     "finder"
	}`)
	body := parseBody(resp)

	assert.Equal(t, 201, resp.StatusCode)
	data := body["data"].(map[string]any)
	assert.Nil(t, data["password"])
}

// 5. Duplicate email → 409 CONFLICT
func TestRegisterDuplicateEmail(t *testing.T) {
	truncateTables(testDB)

	payload := `{
		"name":     "Eko",
		"email":    "eko@mail.com",
		"password": "rahasia123",
		"role":     "finder"
	}`

	first := doPost(testRouter, "/api/v1/users", payload)
	assert.Equal(t, 201, first.StatusCode)

	second := doPost(testRouter, "/api/v1/users", payload)
	body := parseBody(second)

	assert.Equal(t, 409, second.StatusCode)
	assert.Equal(t, "CONFLICT", body["status"])
	assert.NotEmpty(t, body["error"])
}

// 6. Field name kosong → 400
func TestRegisterMissingName(t *testing.T) {
	truncateTables(testDB)

	resp := doPost(testRouter, "/api/v1/users", `{
		"email":    "noname@mail.com",
		"password": "rahasia123",
		"role":     "finder"
	}`)
	body := parseBody(resp)

	assert.Equal(t, 400, resp.StatusCode)
	assert.Equal(t, "BAD REQUEST", body["status"])
}

// 7. Field email kosong → 400
func TestRegisterMissingEmail(t *testing.T) {
	truncateTables(testDB)

	resp := doPost(testRouter, "/api/v1/users", `{
		"name":     "Fajar",
		"password": "rahasia123",
		"role":     "finder"
	}`)
	body := parseBody(resp)

	assert.Equal(t, 400, resp.StatusCode)
	assert.Equal(t, "BAD REQUEST", body["status"])
}

// 8. Field password kosong → 400
func TestRegisterMissingPassword(t *testing.T) {
	truncateTables(testDB)

	resp := doPost(testRouter, "/api/v1/users", `{
		"name":  "Gilang",
		"email": "gilang@mail.com",
		"role":  "finder"
	}`)
	body := parseBody(resp)

	assert.Equal(t, 400, resp.StatusCode)
	assert.Equal(t, "BAD REQUEST", body["status"])
}

// 9. Field role kosong → 400
func TestRegisterMissingRole(t *testing.T) {
	truncateTables(testDB)

	resp := doPost(testRouter, "/api/v1/users", `{
		"name":     "Hendra",
		"email":    "hendra@mail.com",
		"password": "rahasia123"
	}`)
	body := parseBody(resp)

	assert.Equal(t, 400, resp.StatusCode)
	assert.Equal(t, "BAD REQUEST", body["status"])
}

// 10. Format email tidak valid → 400
func TestRegisterInvalidEmailFormat(t *testing.T) {
	truncateTables(testDB)

	resp := doPost(testRouter, "/api/v1/users", `{
		"name":     "Indra",
		"email":    "bukan-email",
		"password": "rahasia123",
		"role":     "finder"
	}`)
	body := parseBody(resp)

	assert.Equal(t, 400, resp.StatusCode)
	assert.Equal(t, "BAD REQUEST", body["status"])
}

// 11. Role tidak valid (bukan finder/seeker/volunteer) → 400
func TestRegisterInvalidRole(t *testing.T) {
	truncateTables(testDB)

	resp := doPost(testRouter, "/api/v1/users", `{
		"name":     "Joko",
		"email":    "joko@mail.com",
		"password": "rahasia123",
		"role":     "admin"
	}`)
	body := parseBody(resp)

	assert.Equal(t, 400, resp.StatusCode)
	assert.Equal(t, "BAD REQUEST", body["status"])
}

// 12. Body kosong {} → 400
func TestRegisterEmptyBody(t *testing.T) {
	truncateTables(testDB)

	resp := doPost(testRouter, "/api/v1/users", `{}`)
	body := parseBody(resp)

	assert.Equal(t, 400, resp.StatusCode)
	assert.Equal(t, "BAD REQUEST", body["status"])
}

// 13. JSON tidak valid → 400
func TestRegisterInvalidJSON(t *testing.T) {
	truncateTables(testDB)

	resp := doPost(testRouter, "/api/v1/users", `{ invalid json }`)
	body := parseBody(resp)

	assert.Equal(t, 400, resp.StatusCode)
	assert.Equal(t, "BAD REQUEST", body["status"])
}

// 14. Verifikasi data tersimpan di DB dengan benar
func TestRegisterDataSavedToDB(t *testing.T) {
	truncateTables(testDB)

	resp := doPost(testRouter, "/api/v1/users", `{
		"name":     "Kiki",
		"email":    "kiki@mail.com",
		"password": "rahasia123",
		"role":     "seeker"
	}`)
	assert.Equal(t, 201, resp.StatusCode)

	var user model.User
	testDB.Where("email = ?", "kiki@mail.com").First(&user)

	assert.Equal(t, "Kiki", user.Name)
	assert.Equal(t, "kiki@mail.com", user.Email)
	assert.Equal(t, model.RoleSeeker, user.Role)
	assert.NotEmpty(t, user.ID)
	// Password harus di-hash, bukan plain text
	assert.NotEqual(t, "rahasia123", user.Password)
	assert.True(t, helper.CheckPasswordHash("rahasia123", user.Password))
}

// 15. Verifikasi phone tersimpan di DB jika dikirim
func TestRegisterPhoneSavedToDB(t *testing.T) {
	truncateTables(testDB)

	resp := doPost(testRouter, "/api/v1/users", `{
		"name":     "Lina",
		"email":    "lina@mail.com",
		"password": "rahasia123",
		"role":     "volunteer",
		"phone":    "08129999888"
	}`)
	assert.Equal(t, 201, resp.StatusCode)

	var user model.User
	testDB.Where("email = ?", "lina@mail.com").First(&user)

	assert.NotNil(t, user.Phone)
	assert.Equal(t, "08129999888", *user.Phone)
}
