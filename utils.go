package bdx

import (
	"os"
	"path"
	"reflect"
	"runtime"
)

type (
	// B response body
	B map[string]interface{}
)

func lastChar(str string) uint8 {
	if str == "" {
		panic("文字列の長さを 0 にすることはできません。")
	}
	return str[len(str)-1]
}

func nameOfFunction(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func joinPaths(absolutePath, relativePath string) string {
	if relativePath == "" {
		return absolutePath
	}

	finalPath := path.Join(absolutePath, relativePath)
	if lastChar(relativePath) == '/' && lastChar(finalPath) != '/' {
		return finalPath + "/"
	}
	return finalPath
}

func assert1(guard bool, text string) {
	if !guard {
		panic(text)
	}
}

func resolveAddress(addr []string) string {
	switch len(addr) {
	case 0:
		if port := os.Getenv("HTTP_PORT"); port != "" {
			DefaultLogger.Infof("環境変数のPORT=\"%s\"", port)
			return ":" + port
		}
		DefaultLogger.Info("環境変数`HTTP_PORT`が未定義です。 デフォルトでは`:8080`を使用します。")
		return ":8080"
	case 1:
		return addr[0]
	default:
		panic("パラメータが多すぎます。")
	}
}
