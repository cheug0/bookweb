// db_test_dao.go
// 数据库测试 DAO
// 提供简单的数据库连接测试功能
package dao

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// TestConnection 测试数据库连接
func TestConnection(driver, host string, port int, user, password, dbname string) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbname)

	db, err := sql.Open(driver, dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	return db.Ping()
}
