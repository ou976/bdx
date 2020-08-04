package render

import (
	"net/http"
)

// Render インターフェースはJSON、XML、YAMLなどで実装します。
type Render interface {
	// Render 与えられたインターフェースオブジェクトをマーシャルし、カスタムContentTypeでデータを書き込みます
	Render(http.ResponseWriter) error
	// WriteContentType レスポンスにContentTypeを書き込みます
	WriteContentType(w http.ResponseWriter)
}

var (
	_ Render = JSON{}
	_ Render = JSONAscii{}
	_ Render = YAML{}
)

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}
