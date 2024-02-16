// user_usecase_impl.go
package usecase

import (
	"cleancode/pkg/entity"
	"cleancode/pkg/helpers"
	repository "cleancode/pkg/respository/interfaces"
	interfaceUseCase "cleancode/pkg/usecase/interfaces"
	"errors"
	"fmt"
	"regexp"
)

type UserUseCase struct {
	UserRepository repository.UserRepository
}

func NewUserUseCase(userRepository repository.UserRepository) interfaceUseCase.UserUseCase {
	return &UserUseCase{UserRepository: userRepository}

}

func (uc *UserUseCase) LoginUser(email, password string) (entity.Compare, entity.Invalid, error) {
	var data entity.Invalid
	var compare entity.Compare
	if email == "" {
		data.EmailError = "Email should not be empty"
		return compare, data, fmt.Errorf("email should not be empty")
	}
	if password == "" {
		data.PasswordError = "Password should not be empty"

		return compare, data, fmt.Errorf(data.PasswordError)
	}
	compare, err := uc.UserRepository.FetchUser(email)
	if err != nil {
		fmt.Println("Error querying the database:", err)
		data.EmailError = "user not found"
		return compare, data, fmt.Errorf(data.EmailError)
	}

	err = helpers.VerifyPassword(password, compare.Password)
	if err != nil {
		fmt.Println("helopass")
		fmt.Println("compare", compare)
		data.PasswordError = "Check password again"

		return compare, data, fmt.Errorf(data.PasswordError)
	}
	if compare.Role != "user" {
		data.RoleError = "Click here for admin login -->"
		return compare, data, fmt.Errorf(data.RoleError)
	}
	if compare.Status != "active" {
		fmt.Println("user blocked")
		data.StatusError = "User is blocked"
		return compare, data, fmt.Errorf(data.StatusError)
	}
	return compare, data, nil
}

func (uc *UserUseCase) RegisterUser(user *entity.User) error {
	fmt.Println("usecase")
	if user.Username == "" {
		return fmt.Errorf("name should not be empty")
	}

	// Check if username already exists
	usernameExists, err := uc.UserRepository.CheckExistingUsername(user.Username)
	if err != nil {
		return fmt.Errorf("error checking existing username: %v", err)
	}
	if usernameExists {
		return fmt.Errorf("username already exists")
	}

	if user.Email == "" {
		return fmt.Errorf("email should not be empty")
	}

	// Check if email already exists
	emailExists, err := uc.UserRepository.CheckExistingEmail(user.Email)
	if err != nil {
		return fmt.Errorf("error checking existing email: %v", err)
	}
	if emailExists {
		return fmt.Errorf("email already exists")
	}

	// Validate email format
	emailPattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	emailRegex := regexp.MustCompile(emailPattern)
	if !emailRegex.MatchString(user.Email) {
		return fmt.Errorf("email not in the correct format")
	}

	if user.Number == "" {
		return fmt.Errorf("number should not be empty")
	}

	// Check if number already exists
	numberExists, err := uc.UserRepository.CheckExistingNumber(user.Number)
	if err != nil {
		return fmt.Errorf("error checking existing number: %v", err)
	}
	if numberExists {
		return fmt.Errorf("number already exists")
	}

	// Validate mobile number format
	numberPattern := `^[0-9]{10}$`
	numberRegex := regexp.MustCompile(numberPattern)
	if !numberRegex.MatchString(user.Number) {
		return fmt.Errorf("invalid Mobile Number")
	}

	if user.Password == "" {
		return fmt.Errorf("password should not be empty")
	}

	if len(user.Password) < 6 {
		return fmt.Errorf("password length should be 6")
	}

	// Check if the user already exists with the provided email
	emailCount, err := uc.UserRepository.CheckExistingEmail(user.Email)
	if err != nil {
		return fmt.Errorf("error checking existing user: %v", err)
	}
	if emailCount {
		return fmt.Errorf("user already exists")
	}

	// Call repository method to save the user
	if err := uc.UserRepository.SaveUser(user); err != nil {
		return fmt.Errorf("error saving user: %v", err)
	}

	return nil

}

func (uc *UserUseCase) AddAddress(userID int, address entity.UserAddress) error {
	fmt.Println("adddddd")
	err := uc.UserRepository.AddAddress(userID, address)
	if err != nil {
		return errors.New("could not add the address")
	}
	return nil
}

func (uc *UserUseCase) GetAllAddress(userId int) ([]entity.AddressInfoResponse, error) {
	addressInfo, err := uc.UserRepository.GetAllAddress(userId)
	if err != nil {
		return []entity.AddressInfoResponse{}, err
	}
	return addressInfo, nil

}

// func (uc *UserUseCase) GenerateOTP() (string, error) {
// 	source := rand.NewSource(time.Now().UnixNano())
// 	randGen := rand.New(source)
// 	return fmt.Sprintf("%06d", randGen.Intn(1000000)), nil
// }

// var otpMap = make(map[string]string)

// func (uc *UserUseCase) SendOTP(otp, email, femail, epassword string) error {
// 	from := femail
// 	password := epassword
// 	to := email
// 	log.Println("email", email, otp)
// 	smtpServer := "smtp.gmail.com"
// 	smtpPort := "587"
// 	otpMap[email] = otp

// 	auth := smtp.PlainAuth("", from, password, smtpServer)

// 	message := fmt.Sprintf("Subject: Your OTP\n\nYour OTP is: %s", otp)

// 	err := smtp.SendMail(smtpServer+":"+smtpPort, auth, from, []string{to}, []byte(message))
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
// func (uc *UserUseCase) VerifyOTP(otp string, email string) error {
// 	userEmail := email
// 	enteredOTP := otp

// 	storedOTP, ok := otpMap[userEmail]
// 	if !ok {
// 		return errors.New("OTP not found for the given Email")
// 	}

// 	if enteredOTP == storedOTP {
// 		// Clear the OTP from the map after successful verification
// 		delete(otpMap, userEmail)
// 		// Render HTML page with a success message
// 		// c.HTML(http.StatusOK, "verify.html", gin.H{"message": "OTP verified successfully"})
// 		// // Send JSON response with the same success message
// 		// c.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully"})
// 		return nil
// 	} else {
// 		return errors.New("invalid otp")
// 	}
// }
