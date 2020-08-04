package render

import (
	"net/http"

	"gopkg.in/yaml.v2"
)

type YAML struct {
	Data interface{} `yaml:"data"`
}

var yamlContentType = []string{"application/x-yaml; charset=utf-8"}

// Render 与えられたインターフェースオブジェクトをマーシャルし、カスタムContentTypeでデータを書き込みます。(YAML)
func (r YAML) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	bytes, err := yaml.Marshal(r.Data)
	if err != nil {
		return err
	}
	_, err = w.Write(bytes)
	return err
}

// WriteContentType レスポンスにContentTypeを書き込みます。
func (r YAML) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, yamlContentType)
}
