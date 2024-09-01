package controllers

import (
	"learn-golang/initializers"
	"learn-golang/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetBookByID(c *gin.Context) {
	var book models.Book
	paramId := c.Param("id")

	if err := initializers.DB.Preload("Topics.SubTopics.Contents").Preload("Buyers").First(&book, paramId).Error; err != nil {
		errorResponse := models.ResponseError{
			Error:   "Book not found",
			Success: false,
		}
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	bookResponse := models.BookResponse{
		ID:       book.ID,
		Title:    book.Title,
		Author:   book.Author,
		Synopsis: book.Synopsis,
		Cover:    book.Cover,
		Price:    book.Price,
	}

	for _, buyer := range book.Buyers {
		bookResponse.Buyers = append(bookResponse.Buyers, buyer.UserID)
	}

	for _, topic := range book.Topics {
		topicResponse := models.TopicResponse{
			Title: topic.Title,
		}

		for _, subTopic := range topic.SubTopics {
			subTopicResponse := models.SubTopicResponse{
				Title: subTopic.Title,
			}

			for _, content := range subTopic.Contents {
				contentResponse := models.ContentResponse{
					Description: content.Description,
				}
				subTopicResponse.Contents = append(subTopicResponse.Contents, contentResponse)
			}
			topicResponse.SubTopics = append(topicResponse.SubTopics, subTopicResponse)
		}
		bookResponse.Topics = append(bookResponse.Topics, topicResponse)
	}

	successResponse := models.ResponseSuccess{
		Data:    bookResponse,
		Success: true,
	}

	c.JSON(http.StatusOK, successResponse)
}

func GetBooks(c *gin.Context) {
	var books []models.Book
	if err := initializers.DB.Preload("Topics.SubTopics.Contents").Preload("Buyers").Find(&books).Error; err != nil {
		errorResponse := models.ResponseError{
			Error:   err.Error(),
			Success: false,
		}
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	bookResponses := make([]models.BookResponse, len(books))

	for i, book := range books {
		bookResponse := models.BookResponse{
			ID:       book.ID,
			Title:    book.Title,
			Author:   book.Author,
			Synopsis: book.Synopsis,
			Cover:    book.Cover,
			Price:    book.Price,
		}

		for _, buyer := range book.Buyers {
			bookResponse.Buyers = append(bookResponse.Buyers, buyer.UserID)
		}

		for _, topic := range book.Topics {
			topicResponse := models.TopicResponse{
				Title: topic.Title,
			}

			for _, subTopic := range topic.SubTopics {
				subTopicResponse := models.SubTopicResponse{
					Title: subTopic.Title,
				}

				for _, content := range subTopic.Contents {
					contentResponse := models.ContentResponse{
						Description: content.Description,
					}
					subTopicResponse.Contents = append(subTopicResponse.Contents, contentResponse)
				}
				topicResponse.SubTopics = append(topicResponse.SubTopics, subTopicResponse)
			}
			bookResponse.Topics = append(bookResponse.Topics, topicResponse)
		}

		bookResponses[i] = bookResponse
	}

	successResponse := models.ResponseSuccess{
		Data:    bookResponses,
		Success: true,
	}

	c.JSON(http.StatusOK, successResponse)
}

func UpdateBook(c *gin.Context) {
	var bookInput models.BookInput

	if err := c.ShouldBindJSON(&bookInput); err != nil {
		errorResponse := models.ResponseError{
			Error:   "Invalid request body",
			Success: false,
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	paramId := c.Param("id")

	var bookFromDB models.Book
	if err := initializers.DB.Preload("Topics.SubTopics.Contents").First(&bookFromDB, paramId).Error; err != nil {
		errorResponse := models.ResponseError{
			Error:   "Book not found",
			Success: false,
		}
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	// Update book details
	bookFromDB.Title = bookInput.Title
	bookFromDB.Author = bookInput.Author
	bookFromDB.Synopsis = bookInput.Synopsis
	bookFromDB.Cover = bookInput.Cover
	bookFromDB.Price = bookInput.Price

	if err := initializers.DB.Delete(&bookFromDB.Topics).Error; err != nil {
		errorResponse := models.ResponseError{
			Error:   err.Error(),
			Success: false,
		}
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	bookFromDB.Topics = []models.Topic{}

	// Update or create topics, sub-topics, and contents
	for _, topicInput := range bookInput.Topics {
		topic := models.Topic{
			Title: topicInput.Title,
		}

		for _, subTopicInput := range topicInput.SubTopics {
			subTopic := models.SubTopic{
				Title: subTopicInput.Title,
			}

			for _, contentInput := range subTopicInput.Contents {
				content := models.Content{
					Description: contentInput.Description,
				}
				subTopic.Contents = append(subTopic.Contents, content)
			}

			topic.SubTopics = append(topic.SubTopics, subTopic)
		}

		bookFromDB.Topics = append(bookFromDB.Topics, topic)
	}

	if err := initializers.DB.Save(&bookFromDB).Error; err != nil {
		errorResponse := models.ResponseError{
			Error:   err.Error(),
			Success: false,
		}
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	successResponse := models.ResponseSuccess{
		Data:    bookFromDB.ID,
		Success: true,
	}

	c.JSON(http.StatusOK, successResponse)
}

func CreateBook(c *gin.Context) {
	var bookInput models.BookInput

	if err := c.ShouldBindJSON(&bookInput); err != nil {
		errorResponse := models.ResponseError{
			Error:   "Invalid request body",
			Success: false,
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	var existingBook models.Book
	if err := initializers.DB.Where("title = ? AND author = ?", bookInput.Title, bookInput.Author).First(&existingBook).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create a new book
			book := models.Book{
				Title:    bookInput.Title,
				Author:   bookInput.Author,
				Synopsis: bookInput.Synopsis,
				Cover:    bookInput.Cover,
				Price:    bookInput.Price,
			}

			// Create the topics, sub-topics, and contents
			for _, topicInput := range bookInput.Topics {
				topic := models.Topic{
					Title: topicInput.Title,
				}

				for _, subTopicInput := range topicInput.SubTopics {
					subTopic := models.SubTopic{
						Title: subTopicInput.Title,
					}

					for _, contentInput := range subTopicInput.Contents {
						content := models.Content{
							Description: contentInput.Description,
						}

						subTopic.Contents = append(subTopic.Contents, content)
					}

					topic.SubTopics = append(topic.SubTopics, subTopic)
				}

				book.Topics = append(book.Topics, topic)
			}

			if err := initializers.DB.Create(&book).Error; err != nil {
				errorResponse := models.ResponseError{
					Error:   "Error creating book",
					Success: false,
				}
				c.JSON(http.StatusInternalServerError, errorResponse)
				return
			}

			successResponse := models.ResponseSuccess{
				Data:    book.ID,
				Success: true,
			}

			c.JSON(http.StatusOK, successResponse)
		} else {
			errorResponse := models.ResponseError{
				Error:   err.Error(),
				Success: false,
			}
			c.JSON(http.StatusInternalServerError, errorResponse)
			return
		}
	} else {
		errorResponse := models.ResponseError{
			Error:   "Book with the same title and author already exists",
			Success: false,
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}
}
