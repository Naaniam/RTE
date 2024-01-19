package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `json:"id" gorm:"type:char(190);primaryKey column:id"`
	Role     string    `json:"role" gorm:"column:role"`
	Mail     string    `json:"mail" gorm:"unique;column:mail"`
	Password string    `json:"password" gorm:"column:password"`
}

type Post struct {
	ID           uuid.UUID `json:"id" gorm:"type:char(190);primaryKey;column:id"`
	RoleID       uuid.UUID `json:"role_id" gorm:"type:type:char(190);;column:role_id"`
	Category     string    `json:"category" gorm:"column:category"`
	Title        string    `json:"title" gorm:"unique;column:title"`
	Description  string    `json:"description" gorm:"column:description"`
	PostDate     time.Time `json:"post_date" gorm:"column:post_date"`
	CommentCount uint      `json:"comment_count" gorm:"column:comment_count"`
	User         User      `json:"-" gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Categories struct {
	Category string `json:"category" `
}

type Comments struct {
	ID       uuid.UUID `json:"id" gorm:"type:char(190);primaryKey column:id"`
	PostID   uuid.UUID `json:"post_id" gorm:"type:char(190); column:post_id"`
	RoleID   uuid.UUID `json:"role_id" gorm:"type:char(190); column:role_id"`
	Feedback string    `json:"feedback" gorm:"primaryKey column:feedback"`
	User     User      `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Post     Post      `gorm:"foreignKey:PostID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
