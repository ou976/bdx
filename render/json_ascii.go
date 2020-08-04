package render

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/belldata-dx/bdx/util/conv/bytesconv"
)

type JSONAscii struct {
	Data interface{} `json:"data"`
}

var jsonAsciiContentType = []string{"application/json"}

// Render 与えられたインターフェースオブジェクトをマーシャルし、カスタムContentTypeでデータを書き込みます(JSONAscii)
func (r JSONAscii) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	ret, err := json.Marshal(r.Data)
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	for _, r := range bytesconv.BytesToString(ret) {
		cvt := string(r)
		if r >= 128 {
			cvt = fmt.Sprintf("\\u%04x", int64(r))
		}
		buffer.WriteString(cvt)
	}

	_, err = w.Write(buffer.Bytes())
	return err
}

// WriteContentType レスポンスにContentTypeを書き込みます
func (r JSONAscii) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonAsciiContentType)
}
