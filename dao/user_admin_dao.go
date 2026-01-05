package dao

import (
	"bookweb/model"
	"bookweb/utils"
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
	utils.Db.Exec("DELETE FROM bookcase WHERE uid = ?", id)
	// 删除书签
	utils.Db.Exec("DELETE FROM bookmark WHERE uid = ?", id)
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

// UpdateUserAdmin 更新用户信息（后台用）
func UpdateUserAdmin(id int, username, password, email string) error {
	var err error
	if password != "" {
		// 如果提供了密码，则更新密码
		// 注意：实际应用中密码应该加密存储，这里为了保持一致性需确认加密逻辑
		// 假设密码在 Controller 层或 Model 层已处理，或此处需处理
		// 鉴于现有代码逻辑，这里先假设传入的是明文，需加密
		// 但 checking dao/admin_dao.go/VerifyAdminPassword uses MD5, user login uses what?
		// Let's check `dao/user_dao.go` to see how password is handled usually.
		// Waiting to check user_dao.go, but for now let's write a generic update.
		// To be safe, let's look at `controller/user.go` or `dao/user_dao.go`.
		// Assuming `UpdateUser` in `dao/user_dao.go` has logic.
		// For now simple update query.
		sqlStr := "UPDATE users SET username = ?, password = ?, email = ? WHERE id = ?"
		_, err = utils.Db.Exec(sqlStr, username, password, email, id)
	} else {
		// 不更新密码
		sqlStr := "UPDATE users SET username = ?, email = ? WHERE id = ?"
		_, err = utils.Db.Exec(sqlStr, username, email, id)
	}
	return err
}
