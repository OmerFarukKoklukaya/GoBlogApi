package models

import (
	"github.com/uptrace/bun"
	"time"
)

type Comment struct {
	bun.BaseModel
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	DeletedAt time.Time `bun:",soft_delete"`

	ID       uint   `bun:"id,pk,autoincrement" json:"id"`
	Content  string `bun:"content" json:"content"`
	IsParent bool   `bun:"is_parent" json:"is_parent"`

	UserID uint  `bun:"user_id" json:"userID"`
	User   *User `bun:"rel:belongs-to,join:user_id=id" json:"user"`

	BlogID uint  `bun:"blog_id" json:"blogID"`
	Blog   *Blog `bun:"rel:belongs-to,join:blog_id=id" json:"blog"`

	ParentCommentID uint     `bun:"parent_comment_id" json:"parent_comment_id"`
	ParentComment   *Comment `bun:"rel:belongs-to,join:parent_comment_id=id" json:"parent_comment"`

	ChildComments []*Comment `bun:"rel:has-many,join:id=parent_comment_id" json:"child_comments"`
}
