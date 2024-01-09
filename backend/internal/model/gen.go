package model

import (
	"time"

	"github.com/pegov/rolecap/backend/internal/entity"
)

type GenPublic struct {
	Id          int          `json:"id"`
	User        entity.Owner `json:"user"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Access      int          `json:"access"`
	Views       int          `json:"views"`
	DateAdded   time.Time    `json:"dateAdded"`
	DateUpdated time.Time    `json:"dateUpdated"`
	Active      bool         `json:"active"`
}

func GenPublicFromGen(entity *entity.Gen) *GenPublic {
	return &GenPublic{
		Id:          entity.Id,
		User:        entity.User,
		Title:       entity.Title,
		Description: entity.Description,
		Access:      entity.Access,
		Views:       entity.Views,
		DateAdded:   entity.DateAdded,
		DateUpdated: entity.DateUpdated,
		Active:      entity.Active,
	}
}
