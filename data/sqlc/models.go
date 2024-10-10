// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package data

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID         int32       `json:"id"`
	Email      string      `json:"email"`
	FirstName  pgtype.Text `json:"first_name"`
	LastName   pgtype.Text `json:"last_name"`
	Password   string      `json:"password"`
	UserActive int32       `json:"user_active"`
	UpdatedAt  time.Time   `json:"updated_at"`
	CreatedAt  time.Time   `json:"created_at"`
}
