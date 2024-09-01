package controllers

import (
	"learn-golang/initializers"
	"learn-golang/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func BuyBook(c *gin.Context) {
	var buyInput models.BuyInput

	if err := c.ShouldBindJSON(&buyInput); err != nil {
		errorResponse := models.ResponseError{
			Error:   "Invalid request body",
			Success: false,
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}
	currentUser, exists := c.Get("curentUser")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get user ID"})
		return
	}

	userId := currentUser.(models.User).ID

	bookId := buyInput.BookID

	var book models.Book
	initializers.DB.Where("id = ?", bookId).Find(&book)

	if book.ID == 0 {
		errorResponse := models.ResponseError{
			Error:   "Book does not exist",
			Success: false,
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	var buyer models.Buyer
	initializers.DB.Where("book_id = ? AND user_id = ?", bookId, userId).Find(&buyer)

	if buyer.ID != 0 {
		errorResponse := models.ResponseError{
			Error:   "User already bought this book",
			Success: false,
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	buyer = models.Buyer{
		BookID: bookId,
		UserID: userId,
	}

	if err := initializers.DB.Create(&buyer).Error; err != nil {
		errorResponse := models.ResponseError{
			Error:   err.Error(),
			Success: false,
		}
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	successResponse := models.ResponseSuccess{
		Data:    buyer.ID,
		Success: true,
	}

	c.JSON(http.StatusOK, successResponse)
}
