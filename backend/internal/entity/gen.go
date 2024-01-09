package entity

import (
	"encoding/json"
	"time"
)

type GenRow struct {
	Id          int             `db:"id"`
	UserId      int             `db:"user_id"`
	Username    string          `db:"username"`
	Title       string          `db:"title"`
	Description string          `db:"description"`
	Access      int             `db:"access"`
	Body        json.RawMessage `db:"body"`
	BodyMeta    json.RawMessage `db:"body_meta"`
	Views       int             `db:"views"`
	DateAdded   time.Time       `db:"date_added"`
	DateUpdated time.Time       `db:"date_updated"`
	Active      bool            `db:"active"`
}

type Gen struct {
	Id          int
	User        Owner
	Title       string
	Description string
	Access      int
	Body        json.RawMessage
	BodyMeta    json.RawMessage
	Views       int
	DateAdded   time.Time
	DateUpdated time.Time
	Active      bool
}

func GenFromGenRow(row *GenRow) *Gen {
	return &Gen{
		Id: row.Id,
		User: Owner{
			Id:       row.UserId,
			Username: row.Username,
		},
		Title:       row.Title,
		Description: row.Description,
		Access:      row.Access,
		Body:        row.Body,
		BodyMeta:    row.BodyMeta,
		Views:       row.Views,
		DateAdded:   row.DateAdded,
		DateUpdated: row.DateUpdated,
		Active:      row.Active,
	}
}
