package helpers

import (
	"cleancode/pkg/entity"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func CreateToken(user entity.Claims, expireTime time.Time) (string, error) {
	expirationTime := expireTime // Adjust as needed
	claims := &entity.Claims{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
		Email:    user.Email,
		Status:   user.Status,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	// fmt.Println("us", user.Username)
	// fmt.Printf("Claims: %+v\n", claims)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtKey := []byte(os.Getenv("jwtKey"))
	fmt.Println("JWT Key:", jwtKey)
	signedToken, err := token.SignedString(jwtKey)
	if err != nil {
		fmt.Println("Error signing token:", err)
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return "", err
	}
	return signedToken, nil

	// c.SetCookie("token", signedToken, int(expirationTime.Unix()), "/", "exclusivestore.xyz", false, true)

	// c.Status(http.StatusOK)

}

func GenerateAccessToken(user entity.Claims) (string, error) {

	expirationTime := time.Now().Add(15 * time.Minute)
	tokenString, err := CreateToken(user, expirationTime)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GenerateRefreshToken(user entity.Claims) (string, error) {

	expirationTime := time.Now().Add(24 * 90 * time.Hour)
	tokenString, err := CreateToken(user, expirationTime)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseToken(tokenString string) (*entity.Claims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &entity.Claims{}, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("jwtKey")), nil
	})
	if err != nil {
		fmt.Println("Access token expired", err)
	}
	claims, ok := token.Claims.(*entity.Claims)
	if !ok {
		return nil, errors.New("failed to extract claims from token")
	}

	return claims, nil
}

func CreateJson(token *entity.TokenUser) (userDetailsJSON []byte) {
	userDetailsJSON, err := json.Marshal(token)
	if err != nil {
		fmt.Println("Error converting UserDetails to JSON:", err)

		return
	}
	return userDetailsJSON

}
func GetID(c *gin.Context) (uint, error) {
	usercookie, err := c.Cookie("auth")
	if err != nil {
		fmt.Println("Error retrieving auth cookie:", err)
		return 0, err
	}

	fmt.Println("Cookie Content:", usercookie)

	var token entity.TokenUser
	err = json.NewDecoder(strings.NewReader(usercookie)).Decode(&token)
	if err != nil {
		fmt.Println("Error decoding UserDetails:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user details"})
		return 0, err
	}

	Claims, err := ParseToken(token.AccessToken)
	if err != nil {
		fmt.Println("Error fetching UserDetails from token:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user details from token"})
		return 0, err
	}

	return Claims.ID, nil
}
func ValidateSignup(user entity.User) (errorMsg string, isValid bool) {
	if user.Username == "" {
		return "Name should not be empty", false
	}

	if user.Email == "" {
		return "Email should not be empty", false
	}

	// Validate email format
	emailPattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	emailRegex := regexp.MustCompile(emailPattern)
	if !emailRegex.MatchString(user.Email) {
		return "Email not in the correct format", false
	}

	if user.Number == "" {
		return "Number should not be empty", false
	}

	// Validate mobile number format
	numberPattern := `^[0-9]{10}$`
	numberRegex := regexp.MustCompile(numberPattern)
	if !numberRegex.MatchString(user.Number) {
		return "Invalid Mobile Number", false
	}

	if user.Password == "" {
		return "Password should not be empty", false
	}

	if len(user.Password) < 6 {
		return "Password length should be at least 6 characters", false
	}

	return "", true
}
