package utils

import (
	"bookweb/config"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var (
	Db *sql.DB
)

// InitDB 初始化数据库连接
// InitDB 初始化数据库连接
func InitDB(cfg *config.DbConfig) {
	if Db != nil {
		Db.Close()
	}
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DbName)
	Db, err = sql.Open(cfg.Driver, dsn)
	if err != nil {
		fmt.Printf("DB Connection Error: %v\n", err)
	} else {
		// Test connection
		if err := Db.Ping(); err != nil {
			fmt.Printf("DB Ping Error: %v\n", err)
		} else {
			fmt.Println("Database connection established.")
		}
	}
}
