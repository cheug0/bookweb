package dao

import (
	"bookweb/model"
	"bookweb/utils"
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// CheckUserNameAndPassword 验证用户名和密码
func CheckUserNameAndPassword(username string, password string) (*model.User, error) {
	// 先根据用户名查询用户
	sqlStr := "select id,username,password,email, IFNULL(last_login_time, ''), IFNULL(current_login_time, '') from users where username = ?"
	row := utils.Db.QueryRow(sqlStr, username)
	user := &model.User{}
	err := row.Scan(&user.Id, &user.Username, &user.Password, &user.Email, &user.LastLoginTime, &user.CurrentLoginTime)
	if err != nil {
		return nil, err
	}
	// 使用 bcrypt 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("密码错误")
	}
	return user, nil
}

// CheckUserName 验证用户名是否存在
func CheckUserName(username string) (bool, error) {
	sqlStr := "select count(*) from users where username = ?"
	row := utils.Db.QueryRow(sqlStr, username)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// SaveUser 保存用户 (密码使用 bcrypt 加密)
func SaveUser(username string, password string, email string) error {
	// 使用 bcrypt 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %v", err)
	}
	sqlStr := "insert into users(username,password,email) values(?,?,?)"
	_, err = utils.Db.Exec(sqlStr, username, string(hashedPassword), email)
	if err != nil {
		fmt.Println("执行错误：", err)
		return err
	}
	return nil
}

// UpdateUser 更新用户信息 (如果修改密码则使用 bcrypt 加密)
func UpdateUser(user *model.User) error {
	// 检查密码是否是新密码（非 bcrypt hash 格式，长度通常 < 60）
	passwordToSave := user.Password
	if len(user.Password) < 60 {
		// 不是 bcrypt hash，需要加密
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("密码加密失败: %v", err)
		}
		passwordToSave = string(hashedPassword)
	}
	sqlStr := "update users set password = ?, email = ? where id = ?"
	_, err := utils.Db.Exec(sqlStr, passwordToSave, user.Email, user.Id)
	if err != nil {
		fmt.Println("执行错误：", err)
		return err
	}
	return nil
}

// GetUserByID 根据ID获取用户
func GetUserByID(id int) (*model.User, error) {
	var row *sql.Row
	if stmtGetUserByID != nil {
		row = stmtGetUserByID.QueryRow(id)
	} else {
		sqlStr := "select id,username,password,email, IFNULL(last_login_time, ''), IFNULL(current_login_time, '') from users where id = ?"
		row = utils.Db.QueryRow(sqlStr, id)
	}
	user := &model.User{}
	err := row.Scan(&user.Id, &user.Username, &user.Password, &user.Email, &user.LastLoginTime, &user.CurrentLoginTime)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateLoginTime 更新登录时间
func UpdateLoginTime(userID int, lastLogin string, currentLogin string) error {
	sqlStr := "update users set last_login_time = ?, current_login_time = ? where id = ?"
	_, err := utils.Db.Exec(sqlStr, lastLogin, currentLogin, userID)
	if err != nil {
		fmt.Println("执行错误：", err)
		return err
	}
	return nil
}
