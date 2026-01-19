package utils

import (
	"bookweb/config"
	"compress/gzip"
	"net/http"
	"strings"
)

type gzipResponseWriter struct {
	http.ResponseWriter
	gzWriter *gzip.Writer
	wrote    bool
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	if !w.wrote {
		w.wrote = true
		// 如果 Content-Encoding 已经被设置（例如 controller 返回了缓存的 gzip 数据），则跳过压缩
		if w.Header().Get("Content-Encoding") != "" {
			return w.ResponseWriter.Write(b)
		}
		// 否则初始化 gzip writer
		w.Header().Set("Content-Encoding", "gzip")
		w.gzWriter = gzip.NewWriter(w.ResponseWriter)
	}

	if w.gzWriter != nil {
		return w.gzWriter.Write(b)
	}
	return w.ResponseWriter.Write(b)
}

func (w *gzipResponseWriter) Close() {
	if w.gzWriter != nil {
		w.gzWriter.Close()
	}
}

// GzipMiddleware GZIP压缩中间件
func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. 检查全局配置是否开启 GZIP
		if !config.GetGlobalConfig().Site.GzipEnabled {
			next.ServeHTTP(w, r)
			return
		}

		// 2. 检查客户端是否支持 gzip
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		// 包装 ResponseWriter
		gzw := &gzipResponseWriter{ResponseWriter: w}
		defer gzw.Close()

		next.ServeHTTP(gzw, r)
	})
}
