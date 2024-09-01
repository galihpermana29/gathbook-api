package models

type BookResponse struct {
	ID       uint            `json:"id"`
	Title    string          `json:"title" binding:"required"`
	Author   string          `json:"author" binding:"required"`
	Synopsis string          `json:"synopsis"`
	Cover    []string        `json:"cover" binding:"required"`
	Price    string          `json:"price" binding:"required"`
	Topics   []TopicResponse `json:"topics" binding:"required"`
	Buyers   []uint          `json:"buyers" binding:"required"`
}

type TopicResponse struct {
	Title     string             `json:"title" binding:"required"`
	SubTopics []SubTopicResponse `json:"subTopics" binding:"required"`
}

type SubTopicResponse struct {
	Title    string            `json:"title" binding:"required"`
	Contents []ContentResponse `json:"contents" binding:"required"`
}

type ContentResponse struct {
	Description string `json:"description" binding:"required"`
}
