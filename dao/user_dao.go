package dao

import (
	"bookweb/model"
	"bookweb/utils"
	"fmt"
)

// CheckUserNameAndPassword 验证用户名和密码
func CheckUserNameAndPassword(username string, password string) (*model.User, error) {
	sqlStr := "select id,username,password,email, IFNULL(last_login_time, ''), IFNULL(current_login_time, '') from users where username = ? and password = ?"
	row := utils.Db.QueryRow(sqlStr, username, password)
	user := &model.User{}
	err := row.Scan(&user.Id, &user.Username, &user.Password, &user.Email, &user.LastLoginTime, &user.CurrentLoginTime)
	if err != nil {
		return nil, err
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

// SaveUser 保存用户
func SaveUser(username string, password string, email string) error {
	sqlStr := "insert into users(username,password,email) values(?,?,?)"
	_, err := utils.Db.Exec(sqlStr, username, password, email)
	if err != nil {
		fmt.Println("执行错误：", err)
		return err
	}
	return nil
}

// UpdateUser 更新用户信息
func UpdateUser(user *model.User) error {
	sqlStr := "update users set password = ?, email = ? where id = ?"
	_, err := utils.Db.Exec(sqlStr, user.Password, user.Email, user.Id)
	if err != nil {
		fmt.Println("执行错误：", err)
		return err
	}
	return nil
}

// GetUserByID 根据ID获取用户
func GetUserByID(id int) (*model.User, error) {
	sqlStr := "select id,username,password,email, IFNULL(last_login_time, ''), IFNULL(current_login_time, '') from users where id = ?"
	row := utils.Db.QueryRow(sqlStr, id)
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
