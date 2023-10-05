package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestLogOut(t *testing.T) {

	// создаем тестовый контекст Gin с помощью httptest
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	LogOut(c)

	// проверяем статус-код ответа
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d but got %d", http.StatusOK, w.Code)
	}

	// проверяем значение куки
	cookies := w.Header().Get("Set-Cookie")
	if cookies != "LoginToken=0; Path=/; HttpOnly; SameSite=Lax" {
		t.Errorf("Expected cookie value %q but got %q", "LoginToken=0; Path=/; HttpOnly; SameSite=Lax", cookies)
	}
}

func TestLogin(t *testing.T) {
	// Создание фейкового контекста
	//router := gin.New()
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	// Создание тестовых данных
	//user := models.User{Login: "testuser", Password: "4535234shzfgszhggzb"}
	//initializers.DB.Create(&user)

	// Подготовка запроса
	request := gin.H{
		"Login":    "testuser",
		"Password": "password123",
	}
	body, _ := json.Marshal(request)
	c.Request, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(body))

	// Выполнение запроса
	Login(c)

	// Проверка статуса ответа
	if c.Writer.Status() != http.StatusOK {
		t.Errorf("Expected status code %d but got %d", http.StatusOK, c.Writer.Status())
	}

	// Проверка установки куки
	cookies := c.Writer.Header()["Set-Cookie"]
	if len(cookies) != 1 || !strings.Contains(cookies[0], "LoginToken=") {
		t.Errorf("Expected Set-Cookie header with LoginToken, but got %v", cookies)
	}
}
