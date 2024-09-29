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

	ID      uint   `bun:"id,pk,autoincrement" json:"id"`
	Title   string `bun:"title" json:"title"`
	Body    string `bun:"body" json:"body"`
	Summary string `bun:"summary" json:"summary"`
	Image   string `bun:"image" json:"image"`

	UserID uint  `bun:"user_id" json:"user_id"`
	User   *User `bun:"rel:belongs-to,join:user_id=id" json:"user"`
}

type Comment struct {
	bun.BaseModel
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	DeletedAt time.Time `bun:",soft_delete"`

	ID      uint   `bun:"id,pk,autoincrement" json:"id"`
	Content string `bun:"content" json:"content"`

	UserID uint  `bun:"user_id" json:"user_id"`
	User   *User `bun:"rel:belongs-to,join:user_id=id" json:"user"`
	BlogID uint  `bun:"blog_id" json:"blog_id"`
	Blog   *Blog `bun:"rel:belongs-to,join:blog_id=id" json:"blog"`
}
