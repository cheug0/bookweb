package dao

import (
	"bookweb/model"
	"bookweb/utils"
	"crypto/md5"
	"encoding/hex"
	"errors"
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

// VerifyAdminPassword 验证管理员密码
func VerifyAdminPassword(username, password string) (*model.Admin, error) {
	admin, err := GetAdminByUsername(username)
	if err != nil {
		return nil, errors.New("用户名不存在")
	}
	// MD5 加密比对
	hash := md5.Sum([]byte(password))
	if admin.Password != hex.EncodeToString(hash[:]) {
		return nil, errors.New("密码错误")
	}
	return admin, nil
}

// CreateAdmin 创建管理员（仅供初始化使用）
func CreateAdmin(username, password string) error {
	hash := md5.Sum([]byte(password))
	passwordHash := hex.EncodeToString(hash[:])
	sqlStr := "INSERT INTO admin (username, password) VALUES (?, ?)"
	_, err := utils.Db.Exec(sqlStr, username, passwordHash)
	return err
}

// UpdateAdminPassword 更新管理员密码
func UpdateAdminPassword(id int, password string) error {
	hash := md5.Sum([]byte(password))
	passwordHash := hex.EncodeToString(hash[:])
	sqlStr := "UPDATE admin SET password = ? WHERE id = ?"
	_, err := utils.Db.Exec(sqlStr, passwordHash, id)
	return err
}
