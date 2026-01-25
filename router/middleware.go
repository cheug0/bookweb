package router

import (
	"bookweb/utils"
	"net/http"
	"time"
)

// responseWriter 包装 ResponseWriter 以捕获状态码
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// LoggingMiddleware HTTP请求日志中间件
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 包装 ResponseWriter
		wrapped := &responseWriter{
			ResponseWriter: w,
			status:         200,
		}

		// 执行下一个处理器
		next.ServeHTTP(wrapped, r)

		// 记录日志
		duration := time.Since(start)
		clientIP := getClientIP(r)
		utils.LogHTTP(r.Method, r.URL.Path, wrapped.status, duration, clientIP)
	})
}

// getClientIP 获取客户端真实IP
func getClientIP(r *http.Request) string {
	// 优先从 X-Forwarded-For 获取
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}
	// 其次从 X-Real-IP 获取
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	// 最后使用 RemoteAddr
	return r.RemoteAddr
}
