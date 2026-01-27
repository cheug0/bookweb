// admin_dao.go
// 管理员 DAO
// 处理后台管理员数据的相关操作
package dao

import (
	"bookweb/model"
	"bookweb/utils"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// GetAdminByUsername 根据用户名查询管理员
func GetAdminByUsername(username string) (*model.Admin, error) {
	sqlStr := "SELECT id, username, password FROM admin WHERE username = ?"
	row := utils.Db.QueryRow(sqlStr, username)
	admin := &model.Admin{}
	err := row.Scan(&admin.Id, &admin.Username, &admin.Password)
	if err != nil {
		return nil, err
	}
	return admin, nil
}

// VerifyAdminPassword 验证管理员密码 (bcrypt)
func VerifyAdminPassword(username, password string) (*model.Admin, error) {
	admin, err := GetAdminByUsername(username)
	if err != nil {
		return nil, errors.New("用户名不存在")
	}
	// bcrypt 比对
	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password))
	if err != nil {
		return nil, errors.New("密码错误")
	}
	return admin, nil
}

// CreateAdmin 创建管理员 (bcrypt 加密)
func CreateAdmin(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %v", err)
	}
	sqlStr := "INSERT INTO admin (username, password) VALUES (?, ?)"
	_, err = utils.Db.Exec(sqlStr, username, string(hashedPassword))
	return err
}

// UpdateAdminPassword 更新管理员密码 (bcrypt 加密)
func UpdateAdminPassword(id int, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %v", err)
	}
	sqlStr := "UPDATE admin SET password = ? WHERE id = ?"
	_, err = utils.Db.Exec(sqlStr, string(hashedPassword), id)
	return err
}
