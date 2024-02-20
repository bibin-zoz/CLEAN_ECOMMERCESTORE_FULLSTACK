package handlers

import (
	"bytes"
	"cleancode/pkg/entity"
	"cleancode/pkg/usecase/mock"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
	userHandler := &UserHandler{}

	router := gin.Default()
	router.POST("/register", userHandler.RegisterUser)

	requestBody := bytes.NewBufferString("username=test&email=test@example.com&number=1234506789&password=secretpassword")
	request, err := http.NewRequest("POST", "/register", requestBody)
	assert.NoError(t, err)

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder := httptest.NewRecorder()

	router.ServeHTTP(responseRecorder, request)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	assert.Equal(t, "\"User Details Validated.. Proceed to verification\"", responseRecorder.Body.String())
}

func Test_LoginHandler(t *testing.T) {
	testCase := map[string]struct {
		Newmail, Newpassword string
		buildStub            func(useCaseMock *mock.MockUserUseCase, newmail, newpassword string)
		checkResponse        func(t *testing.T, responseRecorder *httptest.ResponseRecorder)
	}{
		"Success": {
			Newmail:     "bibin@gmail.com",
			Newpassword: "asdasd",
			buildStub: func(useCaseMock *mock.MockUserUseCase, newmail, newpassword string) {
				compare := entity.Compare{
					ID:       1,
					Password: "asdasd",
					Username: "shah&007",
					Email:    "sha@123.com",
					Role:     "user",
					Status:   "active",
				}

				invalid := entity.Invalid{}

				useCaseMock.EXPECT().LoginUser(newmail, newpassword).Times(1).Return(compare, invalid, nil)
			},
			checkResponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, responseRecorder.Code)
			},
		},

		"User couldn't login": {
			Newmail:     "bibin@gmail.com",
			Newpassword: "no password",
			buildStub: func(useCaseMock *mock.MockUserUseCase, newmail, newpassword string) {
				compare := entity.Compare{}
				invalid := entity.Invalid{
					PasswordError: "Check password again",
				}

				useCaseMock.EXPECT().LoginUser("bibin@gmail.com", "no password").Times(1).Return(compare, invalid, fmt.Errorf("cannot login up"))
			},
			checkResponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, responseRecorder.Code)

				var responseBody map[string]string
				err := json.Unmarshal(responseRecorder.Body.Bytes(), &responseBody)
				assert.NoError(t, err)
				assert.Equal(t, "Check password again", responseBody["error"])
			},
		},
	}

	for testName, test := range testCase {
		test := test
		t.Run(testName, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			mockUseCase := mock.NewMockUserUseCase(ctrl)
			test.buildStub(mockUseCase, test.Newmail, test.Newpassword)
			UserHandler := NewUserHandler(mockUseCase)
			server := gin.Default()

			server.POST("/login", UserHandler.LoginPost)
			jsonData, err := json.Marshal(map[string]string{
				"email":    test.Newmail,
				"password": test.Newpassword,
			})

			assert.NoError(t, err)
			body := bytes.NewBuffer(jsonData)
			mockRequest, err := http.NewRequest(http.MethodPost, "/login", body)
			assert.NoError(t, err)
			responseRecorder := httptest.NewRecorder()
			server.ServeHTTP(responseRecorder, mockRequest)
			test.checkResponse(t, responseRecorder)
		})
	}
}
