package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type Movie struct {
	ID            int       `json:"id"`
	Title         string    `json:"title" gorm:"column:title;type:varchar(100)" validate:"required"`
	Description   string    `json:"description" gorm:"column:description;type:text"  validate:"required"`
	Durasi        string    `json:"durasi" gorm:"column:durasi;type:varchar(100)" validate:"required"`
	Artist        string    `json:"artist" gorm:"column:artist;type:varchar(100)" validate:"required"`
	Genre         string    `json:"genre" gorm:"column:genre;type:varchar(100)" validate:"required"`
	Url           string    `json:"url" gorm:"column:url;type:varchar(100)" validate:"required"`
	LikeCount     int       `json:"like_count" gorm:"column:like_count;"`
	DislikeCount  int       `json:"dislike_count" gorm:"column:dislike_count;"`
	CountSelected int       `json:"count_selected" gorm:"column:count_selected;"`
	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
}

func (*Movie) TableName() string {
	return "movies"
}

func (l Movie) Validate() error {
	v := validator.New()
	return v.Struct(l)
}
