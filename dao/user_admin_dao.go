package dao

import (
	"bookweb/model"
	"bookweb/utils"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// GetUserListAdmin 后台获取用户列表
func GetUserListAdmin(page, pageSize int) ([]*model.User, int, error) {
	offset := (page - 1) * pageSize

	// 统计总数
	var total int
	err := utils.Db.QueryRow("SELECT COUNT(*) FROM users").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 查询列表
	rows, err := utils.Db.Query("SELECT id, username FROM users ORDER BY id DESC LIMIT ?, ?", offset, pageSize)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		u := &model.User{}
		err := rows.Scan(&u.Id, &u.Username)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}
	return users, total, nil
}

// DeleteUserAdmin 删除用户
func DeleteUserAdmin(id int) error {
	// 删除书架
	utils.Db.Exec("DELETE FROM bookcase WHERE userid = ?", id)
	// 删除书签
	utils.Db.Exec("DELETE FROM bookmark WHERE userid = ?", id)
	// 删除用户
	_, err := utils.Db.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}

// GetUserByIDAdmin 根据ID获取用户（后台用）
func GetUserByIDAdmin(id int) (*model.User, error) {
	sqlStr := "SELECT id, username, password, email FROM users WHERE id = ?"
	row := utils.Db.QueryRow(sqlStr, id)
	u := &model.User{}
	err := row.Scan(&u.Id, &u.Username, &u.Password, &u.Email)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// UpdateUserAdmin 更新用户信息（后台用，密码使用 bcrypt 加密）
func UpdateUserAdmin(id int, username, password, email string) error {
	var err error
	if password != "" {
		// 使用 bcrypt 加密密码
		hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if hashErr != nil {
			return fmt.Errorf("密码加密失败: %v", hashErr)
		}
		sqlStr := "UPDATE users SET username = ?, password = ?, email = ? WHERE id = ?"
		_, err = utils.Db.Exec(sqlStr, username, string(hashedPassword), email, id)
	} else {
		// 不更新密码
		sqlStr := "UPDATE users SET username = ?, email = ? WHERE id = ?"
		_, err = utils.Db.Exec(sqlStr, username, email, id)
	}
	return err
}
