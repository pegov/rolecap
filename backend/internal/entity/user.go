package entity

import "time"

type UserRow struct {
	Id       int    `db:"id"`
	Email    string `db:"email"`
	Username string `db:"username"`
	Password string `db:"password"`

	Active   bool `db:"active"`
	Verified bool `db:"verified"`

	CreatedAt time.Time `db:"created_at"`
	LastLogin time.Time `db:"last_login"`

	OAuthProvider string `db:"oauth_provider"`
	OAuthSid      string `db:"oauth_sid"`
}

type OAuth struct {
	UserId   int    `db:"user_id" json:"userId"`
	Provider string `db:"provider" json:"provider"`
	Sid      string `db:"sid" json:"sid"`
}

type User struct {
	Id       int    `db:"id" json:"id"`
	Email    string `db:"email" json:"email"`
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"password"`

	Active   bool `db:"active" json:"active"`
	Verified bool `db:"verified" json:"verified"`

	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	LastLogin time.Time `db:"last_login" json:"lastLogin"`

	OAuth *OAuth
}

func UserFromUserRow(row *UserRow) *User {
	var oauth *OAuth
	if row.OAuthProvider != "" {
		oauth = &OAuth{
			UserId:   row.Id,
			Provider: row.OAuthProvider,
			Sid:      row.OAuthSid,
		}
	}

	return &User{
		Id:       row.Id,
		Email:    row.Email,
		Username: row.Username,
		Password: row.Password,

		Active:   row.Active,
		Verified: row.Verified,

		CreatedAt: row.CreatedAt,
		LastLogin: row.LastLogin,

		OAuth: oauth,
	}
}
