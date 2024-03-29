package middleware

import (
	"cleancode/pkg/entity"
	"cleancode/pkg/helpers"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the auth cookie
		authCookie, err := c.Cookie("adminAuth")
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/admin/login")
			c.AbortWithStatus(http.StatusSeeOther)
			return
		}

		// Decode the JSON content of the auth cookie
		var token entity.TokenUser
		err = json.NewDecoder(strings.NewReader(authCookie)).Decode(&token)
		if err != nil {
			// Redirect to login if there's an error decoding the auth cookie
			c.Redirect(http.StatusSeeOther, "/admin/login")
			c.AbortWithStatus(http.StatusSeeOther)
			return
		}

		claims, err := helpers.ParseToken(token.AccessToken)
		if err != nil {
			log.Println("Error parsing token:", err)
			c.Redirect(http.StatusSeeOther, "/admin/login")
			c.AbortWithStatus(http.StatusSeeOther)
			return
		}
		refreshTokenClaims, err := helpers.ParseToken(token.RefreshToken)
		if err != nil {
			log.Println("Error parsing token:", err)
			c.Redirect(http.StatusSeeOther, "/admin/login")
			c.AbortWithStatus(http.StatusSeeOther)
			return
		}
		if time.Now().Unix() > claims.StandardClaims.ExpiresAt {
			if time.Now().Unix() > refreshTokenClaims.StandardClaims.ExpiresAt {
				log.Println("Token has expired")
				c.Redirect(http.StatusSeeOther, "/admin/login")
				c.AbortWithStatus(http.StatusSeeOther)
				return

			} else {

				c.SetCookie("adminAuth", "", -1, "/admin", "", true, true)
				accessToken, err := helpers.GenerateAccessToken(*claims)
				if err != nil {
					log.Println("Error creating access token token:", err)
					c.Redirect(http.StatusSeeOther, "/admin/login")
					c.AbortWithStatus(http.StatusSeeOther)
					return
				}
				UserAuth := &entity.TokenUser{
					// Users:        claims,
					AccessToken:  accessToken,
					RefreshToken: token.RefreshToken,
				}
				userDetailsJSON := helpers.CreateJson(UserAuth)

				c.SetCookie("adminAuth", string(userDetailsJSON), 0, "/admin", "", true, true)
			}

		}

		if err != nil {
			log.Println("Token not present in cookie:", err)
			c.Redirect(http.StatusSeeOther, "/login")
			c.AbortWithStatus(http.StatusSeeOther)
			return
		}
		if claims.Role != "admin" && claims.Role != "" {
			log.Println("User not logined In")
			c.Redirect(http.StatusSeeOther, "/admin/login")
			c.AbortWithStatus(http.StatusSeeOther)
			return
		}
		if claims.Status != "active" {
			log.Println("User is blocked")
			c.Redirect(http.StatusSeeOther, "/admin/login")
			c.AbortWithStatus(http.StatusSeeOther)
			return
		}
		c.Next()
	}
}
