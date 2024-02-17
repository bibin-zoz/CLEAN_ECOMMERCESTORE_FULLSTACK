package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
	userHandler := &UserHandler{}

	router := gin.Default()
	router.POST("/register", userHandler.RegisterUser)

	requestBody := bytes.NewBufferString("username=test&email=test@example.com&number=123456789&password=secretpassword")
	request, err := http.NewRequest("POST", "/register", requestBody)
	assert.NoError(t, err)

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder := httptest.NewRecorder()

	router.ServeHTTP(responseRecorder, request)

	assert.Equal(t, http.StatusCreated, responseRecorder.Code)

	assert.Equal(t, "\"User successfully signed up\"", responseRecorder.Body.String())
}