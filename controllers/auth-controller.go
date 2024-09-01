package controllers

import (
	"learn-golang/initializers"
	"learn-golang/models"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *gin.Context) {
	var authInput models.AuthInput

	//assign payload to authInput
	if err := c.ShouldBindJSON(&authInput); err != nil {
		errorResponse := models.ResponseError{
			Error:   "Invalid request body",
			Success: false,
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	//check if user exists
	var userFound models.User
	initializers.DB.Where("email = ?", authInput.Email).Find(&userFound)

	if userFound.ID != 0 {
		errorResponse := models.ResponseError{
			Error:   "User already exists",
			Success: false,
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	// generate password hash
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(authInput.Password), bcrypt.DefaultCost)

	if err != nil {
		errorResponse := models.ResponseError{
			Error:   "Error generating password hash",
			Success: false,
		}
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	user := models.User{
		Email:    authInput.Email,
		Password: string(passwordHash),
		Name:     authInput.Name,
	}

	initializers.DB.Create(&user)

	userResponse := models.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	successResponse := models.ResponseSuccess{
		Data:    userResponse,
		Success: true,
	}

	c.JSON(http.StatusOK, successResponse)
}

func Login(c *gin.Context) {
	var authInput models.AuthLoginInput

	if err := c.ShouldBindJSON(&authInput); err != nil {
		errorResponse := models.ResponseError{
			Error:   "Invalid request body",
			Success: false,
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	var userFound models.User
	initializers.DB.Where("email = ?", authInput.Email).Find(&userFound)

	if userFound.ID == 0 {
		errorResponse := models.ResponseError{
			Error:   "User does not exist",
			Success: false,
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	// check if password matches
	err := bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(authInput.Password))

	if err != nil {
		errorResponse := models.ResponseError{
			Error:   "Invalid password",
			Success: false,
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	generatedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userFound.ID,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})

	token, err := generatedToken.SignedString([]byte(os.Getenv(("SECRET"))))

	if err != nil {
		errorResponse := models.ResponseError{
			Error:   "Error generating token",
			Success: false,
		}
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	loginResponse := models.LoginResponse{
		Token: token,
		User: models.UserResponse{
			ID:    userFound.ID,
			Name:  userFound.Name,
			Email: userFound.Email,
		},
	}

	successResponse := models.ResponseSuccess{
		Data:    loginResponse,
		Success: true,
	}

	c.JSON(http.StatusOK, successResponse)
}

func GetUserById(c *gin.Context) {
	var userFound models.User

	idParam := c.Param("id")

	if idParam == "" {
		errorResponse := models.ResponseError{
			Error:   "ID parameter is required",
			Success: false,
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	// Convert the ID parameter to an integer
	id, err := strconv.Atoi(idParam)
	if err != nil {
		errorResponse := models.ResponseError{
			Error:   "Invalid ID parameter",
			Success: false,
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	initializers.DB.Where("id = ?", id).Find(&userFound)

	if userFound.ID == 0 {
		errorResponse := models.ResponseError{
			Error:   "User does not exist",
			Success: false,
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	userResponse := models.UserResponse{
		ID:    userFound.ID,
		Name:  userFound.Name,
		Email: userFound.Email,
	}

	successResponse := models.ResponseSuccess{
		Data:    userResponse,
		Success: true,
	}

	c.JSON(http.StatusOK, successResponse)
}
