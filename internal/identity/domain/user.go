package domain

import "time"

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Role      string    `json:"role"` // Owner, Producer, Director, Writer, Editor, Operator, Viewer
	OrgID     string    `json:"org_id"`
	CreatedAt time.Time `json:"created_at"`
}
