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
