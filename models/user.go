package models

import (
	"bytes"
	"crypto/sha512"
	"github.com/uptrace/bun"
	"time"
)

type User struct {
	bun.BaseModel
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	DeletedAt time.Time `bun:",soft_delete"`
	ID        uint      `bun:"id,pk,autoincrement" json:"id"`
	Name      string    `bun:"name" json:"name"`
	Password  []byte    `bun:"password" json:"password"`
	Blogs     []*Blog   `bun:"rel:has-many,join:id=user_id" json:"blogs"`
	Role      *Role     `bun:"rel:belongs-to,join:role_id=id" json:"role"`
	RoleID    uint      `bun:"role_id" json:"role_id"`
}

const hashCycle = 3

func (user *User) ChangePassword(password []byte) {
	for i := 0; i < hashCycle; i++ {
		dudukluTencere := sha512.New()
		dudukluTencere.Write(password)
		password = dudukluTencere.Sum(nil)
	}
	user.Password = password
}

func (user *User) CheckPassword(password []byte) bool {
	for i := 0; i < hashCycle; i++ {
		dudukluTencere := sha512.New()
		dudukluTencere.Write(password)
		password = dudukluTencere.Sum(nil)
	}
	if bytes.Equal(user.Password, password) {
		return true
	}
	return false
}
