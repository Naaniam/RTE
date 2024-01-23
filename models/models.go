package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `json:"id" gorm:"type:char(190);primaryKey column:id"`
	Role     string    `json:"role" gorm:"column:role" validate:"required"`
	Mail     string    `json:"mail" gorm:"unique;column:mail" validate:"required,email" `
	Password string    `json:"password" gorm:"column:password" validate:"required"`
}

type Post struct {
	ID           uuid.UUID `json:"id" gorm:"type:char(190);primaryKey;column:id"`
	RoleID       uuid.UUID `json:"role_id" gorm:"type:char(190);;column:role_id"`
	Category     string    `json:"category" gorm:"column:category" validate:"required"`
	Title        string    `json:"title" gorm:"unique;column:title" validate:"required"`
	Description  string    `json:"description" gorm:"column:description"`
	PostDate     time.Time `json:"post_date" gorm:"column:post_date"`
	CommentCount uint      `json:"comment_count" gorm:"column:comment_count"`
	ViewsCount   int       `json:"views_count" gorm:"column:views_count"`
	UserCount    int       `json:"user_count" gorm:"column:user_count"`
	User         User      `json:"-" gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Comments struct {
	ID       uuid.UUID `json:"id" gorm:"type:char(190);primaryKey column:id"`
	PostID   uuid.UUID `json:"post_id" gorm:"type:char(190); column:post_id"`
	RoleID   uuid.UUID `json:"role_id" gorm:"type:char(190); column:role_id"`
	Feedback string    `json:"feedback" gorm:"primaryKey column:feedback" validate:"required"`
	User     User      `json:"-" gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Post     Post      `json:"-" gorm:"foreignKey:PostID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Views struct {
	ID       uuid.UUID `json:"id" gorm:"type:char(190);column:id"`
	Views    int       `json:"views" gorm:"column:views"`
	RoleID   uuid.UUID `json:"role_id" gorm:"type:char(190); column:role_id"`
	PostID   uuid.UUID `json:"post_id" gorm:"type:char(190); column:post_id"`
	IsViewed bool      `json:"-" gorm:"column:is_valid"`
	Post     Post      `json:"-" gorm:"foreignKey:PostID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	User     User      `json:"-" gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
