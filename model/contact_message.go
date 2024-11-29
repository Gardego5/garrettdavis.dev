//go:generate msgp
package model

type ContactMessage struct {
	ID        int    `db:"id"`
	Name      string `db:"name" validate:"required"`
	Email     string `db:"email" validate:"email"`
	Message   string `db:"message" validate:"required"`
	CreatedAt Time   `db:"created_at" validate:"required"`
}
