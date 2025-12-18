package db

import (
	"log"
	"time"
)

type User struct {
	Id           string    `json:"id" gorm:"primaryKey"`
	Email        string    `json:"email" gorm:"unique;not null"`
	PasswordHash string    `json:"-" gorm:"not null"`
	CreatedAt    time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
	IsVerified   bool      `json:"isVerified"`
	Tokens       []Token   `json:"tokens,omitempty" gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
	Messages     []Message `json:"messages,omitempty" gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
}

type Token struct {
	Id        string    `json:"id" gorm:"primaryKey"`
	UserId    string    `json:"userId" gorm:"not null;index"`
	Token     string    `json:"token" gorm:"unique;not null"`
	ExpiresAt time.Time `json:"expiresAt" gorm:"not null"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
}

type Message struct {
	Id        string    `json:"id" gorm:"primaryKey"`
	UserId    string    `json:"userId" gorm:"not null;index"`
	Content   string    `json:"content" gorm:"not null"`
	Role      string    `json:"role" gorm:"not null"` // "user" or "assistant"
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
}

func Migrate() {
	db, err := ConnectDB()
	if err != nil {
		log.Fatal("failed to connect db")
	}
	db.AutoMigrate(&User{}, &Token{}, &Message{})
}
