package models

import (
	"time"

	"github.com/lib/pq"
)

type Book struct {
	ID        uint           `gorm:"primaryKey"`
	Title     string         `gorm:"not null"`
	Author    string         `gorm:"not null"`
	Synopsis  string         `json:"synopsis"`
	Price     string         `gorm:"not null"`
	Cover     pq.StringArray `gorm:"type:text[]"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Topics    []Topic `gorm:"foreignKey:BookID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Buyers    []Buyer `gorm:"foreignKey:BookID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type Topic struct {
	ID        uint       `gorm:"primaryKey"`
	Title     string     `gorm:"not null"`
	BookID    uint       `gorm:"not null"`
	SubTopics []SubTopic `gorm:"foreignKey:TopicID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type SubTopic struct {
	ID       uint      `gorm:"primaryKey"`
	Title    string    `gorm:"not null"`
	TopicID  uint      `gorm:"not null"`
	Contents []Content `gorm:"foreignKey:SubTopicID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type Content struct {
	ID          uint   `gorm:"primaryKey"`
	Description string `gorm:"not null"`
	SubTopicID  uint   `gorm:"not null"`
}

type Buyer struct {
	ID     uint `gorm:"primaryKey"`
	BookID uint `gorm:"not null"`
	UserID uint `gorm:"not null"`
}
