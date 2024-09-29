package models

import (
	"github.com/uptrace/bun"
	"time"
)

type Blog struct {
	bun.BaseModel
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	DeletedAt time.Time `bun:",soft_delete"`
	ID        int       `bun:"id,pk,autoincrement" json:"id"`
	Title     string    `bun:"title" json:"title"`
	Body      string    `bun:"body" json:"body"`
	Summary   string    `bun:"summary" json:"summary"`
	Image     string    `bun:"image" json:"image"`
	UserID    int       `bun:"user_id" json:"user_id"`
	User      *User     `bun:"rel:belongs-to,join:user_id=id" json:"user"`
}
