// user_service.go
// 用户服务
// 处理用户相关的业务逻辑
package service

import (
	"bookweb/dao"
	"bookweb/model"
	"errors"
	"regexp"
	"time"
)

// Login 用户登录
func Login(username string, password string) (*model.User, error) {
	user, err := dao.CheckUserNameAndPassword(username, password)
	if err != nil {
		return nil, err
	}

	// 更新登录时间
	// 上次登录时间 = 数据库中当前的 CurrentLoginTime
	newLastLoginTime := user.CurrentLoginTime
	// 本次登录时间 = 当前时间
	newCurrentLoginTime := time.Now().Format("2006-01-02 15:04:05")

	// 更新数据库
	err = dao.UpdateLoginTime(user.Id, newLastLoginTime, newCurrentLoginTime)
	if err != nil {

	}

	user.LastLoginTime = newLastLoginTime
	user.CurrentLoginTime = newCurrentLoginTime

	return user, nil
}

// Register 用户注册
func Register(username string, password string, email string) error {
	// 校验密码长度
	if len(password) < 6 {
		return errors.New("密码长度不能少于6位")
	}

	// 校验邮箱格式
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("邮箱格式不正确")
	}
	// 验证用户名是否已存在
	exist, err := dao.CheckUserName(username)
	if err != nil {
		return err
	}
	if exist {
		return errors.New("用户名已存在")
	}
	return dao.SaveUser(username, password, email)
}

// UpdateUserInfo 更新用户信息
func UpdateUserInfo(userID int, password string, email string) error {
	// Check for empty input first
	if password == "" && email == "" {
		return errors.New("没有需要修改的内容")
	}

	// 获取用户信息
	user, err := dao.GetUserByID(userID)
	if err != nil {
		return err
	}

	// 如果密码不为空，则更新密码
	if password != "" {
		// 校验密码长度
		if len(password) < 6 {
			return errors.New("密码长度不能少于6位")
		}
		user.Password = password
	}

	// 如果邮箱不为空，则更新邮箱
	if email != "" {
		// 校验邮箱格式
		emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
		if !emailRegex.MatchString(email) {
			return errors.New("邮箱格式不正确")
		}
		user.Email = email
	}

	return dao.UpdateUser(user)
}
