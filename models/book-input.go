package models

type BookInput struct {
	ID       uint         `json:"id"`
	Title    string       `json:"title" binding:"required"`
	Author   string       `json:"author" binding:"required"`
	Synopsis string       `json:"synopsis"`
	Cover    []string     `json:"cover" binding:"required"`
	Price    string       `json:"price" binding:"required"`
	Topics   []TopicInput `json:"topics" binding:"required"`
}

type TopicInput struct {
	Title     string          `json:"title" binding:"required"`
	SubTopics []SubTopicInput `json:"subTopics" binding:"required"`
}

type SubTopicInput struct {
	Title    string         `json:"title" binding:"required"`
	Contents []ContentInput `json:"contents" binding:"required"`
}

type ContentInput struct {
	Description string `json:"description" binding:"required"`
}
