// id_transform.go
// ID 转换工具
// 支持基于算术规则的 ID 转换，用于多站共享数据库场景
package utils

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

// Op 定义单个算术操作
type Op struct {
	Type  string // +, -, *, /
	Value int
}

var (
	transOps  []Op
	transLock sync.RWMutex
)

// ParseIdTransRule 解析转换规则字符串
// 格式: "*2,+100" (顺序执行:先乘2,后加100)
func ParseIdTransRule(rule string) error {
	transLock.Lock()
	defer transLock.Unlock()
	transOps = []Op{}

	if rule == "" {
		return nil
	}

	parts := strings.Split(rule, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if len(part) < 2 {
			continue
		}
		op := part[:1]
		valStr := part[1:]
		val, err := strconv.Atoi(valStr)
		if err != nil {
			return fmt.Errorf("invalid value in rule %s: %v", part, err)
		}
		switch op {
		case "+", "-", "*", "/":
			transOps = append(transOps, Op{Type: op, Value: val})
		default:
			return fmt.Errorf("invalid operator in rule %s", part)
		}
	}
	return nil
}

// EncodeID 将真实 ID 转换为展示 ID
func EncodeID(id int) int {
	transLock.RLock()
	defer transLock.RUnlock()
	res := id
	for _, op := range transOps {
		switch op.Type {
		case "+":
			res += op.Value
		case "-":
			res -= op.Value
		case "*":
			res *= op.Value
		case "/":
			if op.Value != 0 {
				res /= op.Value
			}
		}
	}
	return res
}

// DecodeID 将展示 ID 还原为真实 ID
func DecodeID(displayID int) int {
	transLock.RLock()
	defer transLock.RUnlock()
	res := displayID
	// 逆序执行，操作取反
	for i := len(transOps) - 1; i >= 0; i-- {
		op := transOps[i]
		switch op.Type {
		case "+":
			res -= op.Value
		case "-":
			res += op.Value
		case "*":
			if op.Value != 0 {
				res /= op.Value
			}
		case "/":
			res *= op.Value
		}
	}
	return res
}
