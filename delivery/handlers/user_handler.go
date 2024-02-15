// yourproject/delivery/handlers/user_handler.go
package handlers

import (
	"cleancode/entity"
	"cleancode/helpers"
	"cleancode/usecase/interfaces"
	"fmt"
	"log"
	"net/http"
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
var User entity.User

func (h *UserHandler) RegisterUser(c *gin.Context) {

	user := entity.User{
		Username: c.Request.FormValue("username"),
		Email:    c.Request.FormValue("email"),
		Number:   c.Request.FormValue("number"),
		Password: c.Request.FormValue("password"),
	}
	User = user
	c.SetCookie("usermail", user.Email, 360, "/", "localhost", false, true)

	// // Bind request data to User struct
	// if err := c.ShouldBindJSON(&user); err != nil {
	// 	fmt.Println("error", user)
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }
	fmt.Println("hii")
	// Call the use case to register the user

	fmt.Println("hii")
	c.Redirect(http.StatusFound, "/verify")
}

func (h *UserHandler) Signup(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Expires", "0")

	c.HTML(http.StatusOK, "signup.html", nil)
}
func (h *UserHandler) LoginHandler(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Expires", "0")
	var data entity.Invalid
	data.LoginStatus = false
	c.HTML(http.StatusOK, "login.html", data)

}

func (h *UserHandler) LoginPost(c *gin.Context) {
	Newmail := c.Request.FormValue("email")
	Newpassword := c.Request.FormValue("password")
	// var compare entity.Compare
	// var data entity.Invalid

	compare, data, error := h.UserUseCase.LoginUser(Newmail, Newpassword)
	if error != nil {
		c.HTML(http.StatusOK, "login.html", data)
		fmt.Println("hi")
		return
	}

	claims := entity.Claims{
		ID:       compare.ID,
		Username: compare.Username,
		Email:    compare.Email,
		Role:     compare.Role,
		Status:   compare.Status,
	}

	accessToken, err := helpers.GenerateAccessToken(claims)
	if err != nil {
		fmt.Println("Error generating access token:", err)
		return
	}

	refreshToken, err := helpers.GenerateRefreshToken(claims)
	if err != nil {
		fmt.Println("Error generating refresh token:", err)

		return
	}

	UserLoginDetails := &entity.TokenUser{
		// Users:        claims,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	userDetailsJSON := helpers.CreateJson(UserLoginDetails)

	c.SetCookie("auth", string(userDetailsJSON), 0, "/", "", true, true)

	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Expires", "0")

	c.Redirect(http.StatusFound, "/home")
}

var lastOTPSendTime time.Time

func (h *UserHandler) VerifyHandler(c *gin.Context) {
	fmt.Println("asa ")
	//for avoiding req in 60secondss
	if time.Since(lastOTPSendTime) < 60*time.Second {
		c.HTML(http.StatusOK, "verify.html", gin.H{"Message": "Please wait before requesting a new OTP"})
		return
	}
	otp := helpers.GenerateOTP()

	useremail, err := c.Cookie("usermail")
	if err != nil {
		c.String(http.StatusNotFound, "Cookie not found")
		return
	}
	helpers.SendOTP(otp, useremail)
	lastOTPSendTime = time.Now()
	if err != nil {
		c.HTML(http.StatusOK, "verify.html", gin.H{"Message": "Please wait before requesting a new OTP"})
		return

	}

	c.HTML(http.StatusOK, "verify.html", gin.H{"Message": "OTP sented"})

	// Update the last OTP send time
	lastOTPSendTime = time.Now()
}

func (h *UserHandler) VerifyPost(c *gin.Context) {
	var verifyData entity.VerifyData

	if err := c.ShouldBindJSON(&verifyData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	otp := verifyData.OTP
	email, err := c.Cookie("usermail")
	if err != nil {
		c.String(http.StatusNotFound, "Cookie not found")
		return
	}
	status := helpers.VerifyOTP(otp, email)
	log.Println("verifypost", otp, status)

	if !status {
		// Handle the case when OTP verification fails
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OTP"})
		return
	}

	hasedPassword, _ := helpers.HashPassword(User.Password)
	newUser := entity.User{
		Username: User.Username,
		Email:    User.Email,
		Number:   User.Number,
		Password: hasedPassword,
	}

	// err := db.DB.Create(&newUser).Error
	if err := h.UserUseCase.RegisterUser(&newUser); err != nil {
		fmt.Println("", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	if err != nil {
		// Check for duplicate key violation or other errors
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error creating user"})
		return
	}
	// helpers.UpdateReferalCount(User.ReferalDetails.ReferalCode)
	// referalDetails := entity.ReferalDetails{
	// 	UserID:      1, // Replace with the actual user ID
	// 	Count:       0,
	// 	ReferalCode: helpers.GenerateRandomReferalCode(),
	// }

	// // Save to the database
	// err = db.DB.Create(&referalDetails).Error
	// if err != nil {
	// 	// Check for duplicate key violation or other errors
	// 	fmt.Println(err)
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Error creating Referal id"})
	// 	return
	// }

	// Redirect to /login with a success message
	c.JSON(http.StatusOK, gin.H{"message": "User created successfully. Please log in."})
}
