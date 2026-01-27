// main.go (genpwd)
// 密码生成工具
// 用于生成 bcrypt 加密的密码哈希，可用于重置管理员或用户密码
package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "admin123" // 默认密码，可通过命令行参数修改
	if len(os.Args) > 1 {
		password = os.Args[1]
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Password:", password)
	fmt.Println("Bcrypt Hash:", string(hash))
	fmt.Println()
	fmt.Printf("SQL: UPDATE admin SET password = '%s' WHERE username = 'admin';\n", string(hash))
}
