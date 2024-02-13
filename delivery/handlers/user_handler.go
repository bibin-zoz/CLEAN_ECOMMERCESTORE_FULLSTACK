// yourproject/delivery/handlers/user_handler.go
package handlers

import (
	"cleancode/entity"
	"cleancode/usecase/interfaces"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserUseCase interfaces.UserUseCase
}

func NewUserHandler(userUseCase interfaces.UserUseCase) *UserHandler {
	return &UserHandler{UserUseCase: userUseCase}
}

// Rest of the code...

func (h *UserHandler) RegisterUser(c *gin.Context) {

	user := entity.User{
		Username: c.Request.FormValue("username"),
		Email:    c.Request.FormValue("email"),
		Number:   c.Request.FormValue("number"),
		Password: c.Request.FormValue("password"),
	}
	c.SetCookie()

	// // Bind request data to User struct
	// if err := c.ShouldBindJSON(&user); err != nil {
	// 	fmt.Println("error", user)
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }
	fmt.Println("hii")
	// Call the use case to register the user
	if err := h.UserUseCase.RegisterUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	fmt.Println("hii")
	c.Redirect(http.StatusFound, "/verify")
}

func (h *UserHandler) Signup(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Expires", "0")

	c.HTML(http.StatusOK, "signup.html", nil)
}

var lastOTPSendTime time.Time

func (h *UserHandler) VerifyHandler(c *gin.Context) {
	//for avoiding req in 60secondss
	if time.Since(lastOTPSendTime) < 60*time.Second {
		c.HTML(http.StatusOK, "verify.html", gin.H{"Message": "Please wait before requesting a new OTP"})
		return
	}
	otp, err := h.UserUseCase.GenerateOTP()
	from := os.Getenv("email")
	password := os.Getenv("password")
	h.UserUseCase.SendOTP(otp, user.Email, from, password)
	if err != nil {
		c.HTML(http.StatusOK, "verify.html", gin.H{"Message": "Please wait before requesting a new OTP"})
		return

	}

	c.HTML(http.StatusOK, "verify.html", gin.H{"Message": "OTP sented"})

	// Update the last OTP send time
	lastOTPSendTime = time.Now()
}
