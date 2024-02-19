// yourproject/delivery/handlers/user_handler.go
package handlers

import (
	"cleancode/pkg/entity"
	"cleancode/pkg/helpers"
	"cleancode/pkg/usecase/interfaces"
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
		// ConfirmPassword: c.Request.FormValue("confirmPassword"),
	}

	errors, isValid := helpers.ValidateSignup(user)
	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	log.Printf("User registered: %s", user.Email)

	c.JSON(http.StatusOK, "User Details Validated.. Proceed to verification")
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
func (h *UserHandler) LogoutHandler(c *gin.Context) {
	fmt.Println("user logut")

	c.SetCookie("auth", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logout successful",
	})

}
func (h *UserHandler) LoginPost(c *gin.Context) {
	Newmail := c.Request.FormValue("email")
	Newpassword := c.Request.FormValue("password")

	compare, data, err := h.UserUseCase.LoginUser(Newmail, Newpassword)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": data.PasswordError, // You might want to customize this based on your error handling.
		})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating access token"})
		return
	}

	refreshToken, err := helpers.GenerateRefreshToken(claims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating refresh token"})
		return
	}

	UserLoginDetails := &entity.TokenUser{
		Users:        claims,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	c.SetCookie("auth", string(helpers.CreateJson(UserLoginDetails)), 0, "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user":    claims,
		"tokens":  UserLoginDetails,
	})
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

func (h *UserHandler) HomeHandler(c *gin.Context) {
	var data entity.Invalid
	data.LoginStatus = true
	_, err := c.Cookie("auth")
	if err != nil {
		data.LoginStatus = false
		fmt.Println("err")
	}
	fmt.Println("islogin", data.LoginStatus)
	c.HTML(http.StatusOK, "home.html", gin.H{
		// "Productvariants": ProductVariants,
		"IsLogin":         data.LoginStatus,
		"ProductVariants": nil,
	})
}

func (h *UserHandler) AddAddress(c *gin.Context) {
	userID, _ := helpers.GetID(c)
	var address entity.UserAddress
	if err := c.ShouldBindJSON(&address); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": "fields provided are in wrong format"})
		return
	}
	// err := validator.New().Struct(address)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "constraints does not match"})
	// 	return
	// }
	err := h.UserUseCase.AddAddress(int(userID), address)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed adding address"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Success": "Address added successfully"})

}
func (ur *UserHandler) GetAllAddress(c *gin.Context) {
	userID, _ := helpers.GetID(c)
	addressInfo, err := ur.UserUseCase.GetAllAddress(int(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve details"})
		return
	}
	c.JSON(http.StatusOK, addressInfo)

}

func (ur *UserHandler) UserProfileHandler(c *gin.Context) {
	userID, err := helpers.GetID(c)
	fmt.Println("user", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve userid"})
		return
	}
	userDetails, err := ur.UserUseCase.GetUserDetails(int(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve details"})
		return
	}
	c.JSON(http.StatusOK, userDetails)

}
