// yourproject/entity/user.go
package entity

import "github.com/golang-jwt/jwt"

type VerifyData struct {
	OTP string `json:"otp"`
}
type Invalid struct {
	IDError       string
	NameError     string
	EmailError    string
	NumberError   string
	PasswordError string
	RoleError     string
	CommonError   string
	LoginStatus   bool
	AmountError   string
	DateError     string
	StatusError   string
}

type Claims struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Status   string `json:"status"`
	jwt.StandardClaims
}

type UserSignUp struct {
	Username string `json:"username" validate:"gte=3"`
	Lastname string `json:"lastname" validate:"gte=1"`
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"min=6,max=20"`
	Number   string `json:"number"`
}
type UserDetailsResponse struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Number   string `json:"number"`
}
type UserDetailsAtAdmin struct {
	Id          int    `json:"id"`
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	BlockStatus bool   `json:"block_status"`
}

type LoginDetail struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
	Id        uint   `json:"id"`
	UserId    uint   `json:"user_id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
}
type AddressInfoResponse struct {
	ID         uint   `json:"id"`
	Street     string `gorm:"not null"  json:"street" form:"street" binding:"required"`
	City       string `gorm:"not null" json:"city" form:"city" binding:"required"`
	State      string `gorm:"not null" json:"state" form:"state" binding:"required"`
	PostalCode string `gorm:"not null" json:"postalcode" form:"postalcode" binding:"required"`
	Country    string `gorm:"not null" json:"country" form:"country" binding:"required"`
}
type AddressInfo struct {
}
type UsersProfileDetails struct {
	Firstname string `json:"firstname" `
	Lastname  string `json:"lastname" `
	Email     string `json:"email" `
	Phone     string `json:"phone" `
}

type PaymentDetails struct {
	ID           uint   `json:"id"`
	Payment_Name string `json:"payment_name"`
}

type ChangePassword struct {
	Oldpassword string `json:"old_password"`
	Password    string `json:"password"`
	Repassword  string `json:"re_password"`
}

type ForgotPasswordSend struct {
	Phone string `json:"phone"`
}
type ForgotVerify struct {
	Phone       string `json:"phone" binding:"required" validate:"required"`
	Otp         string `json:"otp" binding:"required"`
	NewPassword string `json:"newpassword" binding:"required" validate:"min=6,max=20"`
}
