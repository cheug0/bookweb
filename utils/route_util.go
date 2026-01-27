// route_util.go
// 路由工具
// 路由参数解析及路径生成辅助函数
package utils

import (
	"bookweb/model"
	"fmt"
	"net/http"
	"strconv"
)

// GetRouteParam 从请求中获取路由参数
// 现在的逻辑非常纯粹：直接从 Context 中读取由 Router 包解析好的参数
// 所有的复杂匹配、正则解析、后缀清理都已经在 Router 层处理完毕
func GetRouteParam(r *http.Request, name string) string {
	params := r.Context().Value(model.ParamsKey)
	if params == nil {
		return ""
	}

	ps := params.(model.Params)

	// Router 层已经确保存入的是纯净的逻辑 ID 值（不带后缀）
	return ps.ByName(name)
}

// GetIntParam 获取路由参数并尝试转换为整数
// 仅返回结果，不处理 HTTP 响应
func GetIntParam(r *http.Request, name string) (int, error) {
	val := GetRouteParam(r, name)
	if val == "" {
		return 0, fmt.Errorf("param %s is empty", name)
	}
	return strconv.Atoi(val)
}
