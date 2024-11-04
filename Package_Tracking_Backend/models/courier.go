package models

type Courier struct {
	ID        int64  `db:"id" json:"id"`
	Name      string `db:"name" json:"name"`
	Email     string `db:"email" json:"email"`         // Optional
	Available bool   `db:"available" json:"available"` // Optional: Courier availability
}
