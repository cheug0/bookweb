package model

import "github.com/julienschmidt/httprouter"

// ContextKey 用于context键的类型
type ContextKey string

const (
	// ParamsKey 用于存储 httprouter 参数的 Key
	ParamsKey ContextKey = "params"
	// StartTimeKey 用于存储请求开始时间的 Key
	StartTimeKey ContextKey = "startTime"
)

// Params 路由参数类型别名
type Params = httprouter.Params
