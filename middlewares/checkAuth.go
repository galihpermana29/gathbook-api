package middlewares

import (
	"fmt"
	"learn-golang/initializers"
	"learn-golang/models"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func CheckAuth(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		errorResponse := models.ResponseError{
			Error:   "Unauthorized",
			Success: false,
		}
		c.JSON(http.StatusUnauthorized, errorResponse)
		c.Abort()
		return
	}

	authToken := strings.Split(authHeader, " ")

	if len(authToken) != 2 || authToken[0] != "Bearer" {

		errorResponse := models.ResponseError{
			Error:   "Unauthorized",
			Success: false,
		}
		c.JSON(http.StatusUnauthorized, errorResponse)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	tokenString := authToken[1]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil || !token.Valid {
		errorResponse := models.ResponseError{
			Error:   "Invalid token or expired token",
			Success: false,
		}
		c.JSON(http.StatusUnauthorized, errorResponse)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		errorResponse := models.ResponseError{
			Error:   "Invalid token",
			Success: false,
		}
		c.JSON(http.StatusUnauthorized, errorResponse)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		errorResponse := models.ResponseError{
			Error:   "Expired token",
			Success: false,
		}
		c.JSON(http.StatusUnauthorized, errorResponse)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var user models.User
	initializers.DB.Where("id = ?", claims["id"]).Find(&user)

	if user.ID == 0 {
		errorResponse := models.ResponseError{
			Error:   "Invalid token",
			Success: false,
		}
		c.JSON(http.StatusUnauthorized, errorResponse)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("curentUser", user)
	c.Next()
}
