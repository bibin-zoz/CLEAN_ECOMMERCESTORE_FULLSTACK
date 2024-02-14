// yourproject/usecase/interfaces/user_usecase.go
package interfaces

import (
	"cleancode/entity"
)

type UserUseCase interface {
	RegisterUser(user *entity.User) error
	GenerateOTP() (string, error)
	SendOTP(otp, email, femail, epassword string) error
	VerifyOTP(otp string, email string) error
	HashPassword(password string) (string, error)
}
