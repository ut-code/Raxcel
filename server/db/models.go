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
}

type Token struct {
	Id        string    `json:"id" gorm:"primaryKey"`
	UserId    string    `json:"userId" gorm:"not null;index"`
	Token     string    `json:"token" gorm:"unique;not null"`
	ExpiresAt time.Time `json:"expiresAt" gorm:"not null"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	User      User      `json:"user,omitempty" gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
}

func Migrate() {
	db, err := ConnectDB()
	if err != nil {
		log.Fatal("failed to connect db")
	}
	db.AutoMigrate(&User{}, &Token{})
}
