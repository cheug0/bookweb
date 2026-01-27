// db.go
// 数据库工具
// 处理 MySQL 数据库连接池的初始化与配置
package utils

import (
	"bookweb/config"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	Db *sql.DB
)

// InitDB 初始化数据库连接
func InitDB(cfg *config.DbConfig) {
	if Db != nil {
		Db.Close()
	}
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local&timeout=5s&readTimeout=5s&writeTimeout=5s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DbName)
	Db, err = sql.Open(cfg.Driver, dsn)
	if err != nil {
		LogError("Database", "DB Connection Error: %v", err)
		return
	}

	// 配置连接池
	maxOpen := cfg.MaxOpenConns
	if maxOpen <= 0 {
		maxOpen = 100 // 默认值
	}
	// 强制 MaxIdleConns = MaxOpenConns，保持所有连接常驻，避免频繁重连
	maxIdle := maxOpen

	connMaxLifetime := cfg.ConnMaxLifetime
	if connMaxLifetime <= 0 {
		connMaxLifetime = 300 // 默认5分钟
	}

	Db.SetMaxOpenConns(maxOpen)
	Db.SetMaxIdleConns(maxIdle)
	Db.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Second)
	Db.SetConnMaxIdleTime(1 * time.Hour) // 空闲连接存活1小时，尽量不关闭

	// Test connection
	if err := Db.Ping(); err != nil {
		LogWarn("Database", "DB Ping Error: %v", err)
	} else {
		LogInfo("Database", "Database connection established. Pool: MaxOpen=%d, MaxIdle=%d, MaxLifetime=%ds",
			maxOpen, maxIdle, connMaxLifetime)
	}
}
