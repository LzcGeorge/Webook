package web

type Result struct {
	// 业务错误码
	Code int `json:"code"`
	// 错误信息
	Msg string `json:"msg"`
	// 业务数据
	Data any `json:"data"`
}
