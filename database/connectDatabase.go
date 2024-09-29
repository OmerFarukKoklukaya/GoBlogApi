package database

import (
	"blog/models"
	"context"
	"database/sql"
	"fmt"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var ctx = context.Background()
var DB *bun.DB

func ConnectDatabase() {
	dsn := "postgres://postgres:user@localhost:5432/blog?sslmode=disable"
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	DB = bun.NewDB(sqldb, pgdialect.New())
}

func CreateTables() {
	db := DB

	db.NewCreateTable().IfNotExists().Model(&models.User{}).Exec(ctx)
	db.RegisterModel((*models.RoleToPermission)(nil))

	for _, m := range []interface{}{
		(*models.RoleToPermission)(nil),
		(*models.Role)(nil),
		(*models.Permission)(nil),
		(*models.Blog)(nil),
		(*models.Comment)(nil),
	} {
		_, err := db.NewCreateTable().IfNotExists().Model(m).Exec(ctx)
		if err != nil {
			fmt.Println(err)
		}
	}

}
