package controller

import (
	"bookweb/utils"
	"net/http"
	"strconv"
)

// NotFound 处理 404 错误并渲染自定义模版
func NotFound(w http.ResponseWriter, r *http.Request) {
	// 使用通用数据获取方法，确保 SEO 兜底生效
	data := GetCommonData(r).Add("CurrentTitle", "页面未找到")

	// 设置状态码为 404
	w.WriteHeader(http.StatusNotFound)

	// 使用预编译模板（自动根据PC/移动端选择）
	t := GetRenderTemplate(w, r, "error.html")
	if t == nil {
		// 如果模版获取失败，回退到默认 404
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}
	t.Execute(w, data)
}

// GetIDOr404 尝试获取整数 ID，如果失败则直接渲染 404 页面
// 返回 (id, 是否成功)
func GetIDOr404(w http.ResponseWriter, r *http.Request, name string) (int, bool) {
	val, err := strconv.Atoi(utils.GetRouteParam(r, name))
	if err != nil {
		NotFound(w, r)
		return 0, false
	}
	// 如果是小说 ID，进行解码
	if name == "aid" || name == "articleid" {
		original := val
		val = utils.DecodeID(val)
		// Debug Log
		if original != val {
			// fmt.Printf("DEBUG: GetIDOr404 Decode %s: %d -> %d\n", name, original, val)
		}
	}
	return val, true
}

// GetID 获取整数 ID，如果失败返回 0, false (不渲染 404)
func GetID(w http.ResponseWriter, r *http.Request, name string) (int, bool) {
	val, err := strconv.Atoi(utils.GetRouteParam(r, name))
	if err != nil {
		return 0, false
	}
	// 如果是小说 ID，进行解码
	if name == "aid" || name == "articleid" {
		val = utils.DecodeID(val)
	}
	return val, true
}
