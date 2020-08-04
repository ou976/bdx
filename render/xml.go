package render

import (
	"encoding/xml"
	"net/http"
)

type XML struct {
	Data interface{} `xml:"data"`
}

var xmlContentType = []string{"application/xml; charset=utf-8"}

// Render 与えられたインターフェースオブジェクトをマーシャルし、カスタムContentTypeでデータを書き込みます(XML)
func (r XML) Render(w http.ResponseWriter) (err error) {
	r.WriteContentType(w)
	buf, err := xml.MarshalIndent(r.Data, "", "")
	if err != nil {
		panic(err)
	}
	w.Write(buf)
	return
}

// WriteContentType レスポンスにContentTypeを書き込みます
func (r XML) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, xmlContentType)
}
