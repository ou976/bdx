package middleware

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"

	"github.com/belldata-dx/bdx/interfaces"
	"github.com/go-playground/locales"
	"github.com/go-playground/locales/ja"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	jaTranslations "github.com/go-playground/validator/v10/translations/ja"
	"gopkg.in/yaml.v2"
)

type Trans struct {
	Lang                        locales.Translator
	RegisterDefaultTranslations func(*validator.Validate, ut.Translator) error
}

var Translang = Trans{
	Lang:                        ja.New(),
	RegisterDefaultTranslations: jaTranslations.RegisterDefaultTranslations,
}

var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
	trans    ut.Translator
)

func ChangeTranslate() {
	uni = ut.New(Translang.Lang, Translang.Lang)
	trans, _ = uni.GetTranslator(Translang.Lang.Locale())
	validate = validator.New()
	Translang.RegisterDefaultTranslations(validate, trans)
}

func init() {
	ChangeTranslate()
}

// Content-Type MIME of the most common data formats.
const (
	MIMEJSON              = "application/json"
	MIMEHTML              = "text/html"
	MIMEXML               = "application/xml"
	MIMEXML2              = "text/xml"
	MIMEPlain             = "text/plain"
	MIMEPOSTForm          = "application/x-www-form-urlencoded"
	MIMEMultipartPOSTForm = "multipart/form-data"
	MIMEPROTOBUF          = "application/x-protobuf"
	MIMEYAML              = "application/x-yaml"
)

type errorResponse struct {
	Code          int               `json:"code" xml:"code" yaml:"code"`
	Error         string            `json:"error,omitempty" xml:"error,omitempty" yaml:"error,omitempty"`
	ErrorDescript string            `json:"error_descript,omitempty" xml:"error_descript,omitempty" yaml:"error_descript,omitempty"`
	ErrorDetail   map[string]string `json:"error_detail,omitempty" xml:"error_detail,omitempty" yaml:"error_detail,omitempty"`
}

type MIMEType int

const (
	JSON MIMEType = iota
	XML
	YAML
	HTML
)

func checkAccept(r *http.Request) MIMEType {
	accept := r.Header.Get("Accept")
	m, _ := checkContent(accept)
	return m
}

func checkContentType(r *http.Request) (MIMEType, bool) {
	contentType := r.Header.Get("Content-Type")
	return checkContent(contentType)
}

func checkContent(content string) (MIMEType, bool) {
	if strings.HasPrefix(content, MIMEJSON) {
		return JSON, true
	} else if strings.HasPrefix(content, MIMEXML) {
		return XML, true
	} else if strings.HasPrefix(content, MIMEXML2) {
		return XML, true
	} else if strings.HasPrefix(content, MIMEYAML) {
		return YAML, true
	} else {
		return JSON, false
	}
}

func convResBody(m MIMEType, data interface{}) []byte {
	switch m {
	case JSON:
		buf, _ := json.Marshal(&data)
		return buf
	case XML:
		buf, _ := xml.Marshal(&data)
		return buf
	case YAML:
		buf, _ := yaml.Marshal(&data)
		return buf
	default:
		buf, _ := json.Marshal(&data)
		return buf
	}
}

func convReqBody(m MIMEType, data []byte, instance interface{}) error {
	switch m {
	case JSON:
		return json.Unmarshal(data, instance)
	case XML:
		return xml.Unmarshal(data, instance)
	case YAML:
		return yaml.Unmarshal(data, instance)
	default:
		return errors.New("no content")
	}
}

// Validator は`Request.Body`にセットされている`JSON文字列`を`x`の引数の型に変換し、
//
// `go-playground.Validator`でバリデーション処理を行うミドルウェア
//
// エラーが発生した際はエラーレスポンスへ変換して後続のミドルウェア、ハンドルを実行しない
func Validator(x interface{}) interfaces.BdxHandlerFunc {
	var t reflect.Type
	t = reflect.TypeOf(x)
	return func(ctx interfaces.Context) {
		r := ctx.Request()
		// Get Mthod以外
		if r.Method != http.MethodGet {
			var mimeType MIMEType
			var ok bool
			if mimeType, ok = checkContentType(r); !ok {
				errorBody := errorResponse{
					Code:          http.StatusBadRequest,
					Error:         "Invalid Content-Type",
					ErrorDescript: "Invalid Content-Type",
					ErrorDetail:   nil,
				}
				buf := convResBody(checkAccept(r), errorBody)
				ctx.AbortWithStatusAndMessage(http.StatusBadRequest, buf)
				ctx.Next()
			} else {
				body, err := ioutil.ReadAll(r.Body)
				if err != nil {
					errorBody := errorResponse{
						Code:          http.StatusBadRequest,
						Error:         "Invalid body parser",
						ErrorDescript: "Invalid body parser",
						ErrorDetail:   nil,
					}
					buf := convResBody(checkAccept(r), errorBody)
					ctx.AbortWithStatusAndMessage(http.StatusBadRequest, buf)
				} else {
					r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
					instance := reflect.New(t).Interface()
					convReqBody(mimeType, body, &instance)
					err = validate.Struct(instance)
					if err != nil {
						messages := err.(validator.ValidationErrors).Translate(trans)
						errMsgMap := make(map[string]string)
						for key, val := range messages {
							field := strings.Split(key, ".")
							errMsgMap[field[len(field)-1]] = val
						}
						errorBody := errorResponse{
							Code:          http.StatusBadRequest,
							Error:         "Invalid body parser",
							ErrorDescript: "Invalid body parser",
							ErrorDetail:   errMsgMap,
						}
						buf := convResBody(checkAccept(r), errorBody)
						ctx.AbortWithStatusAndMessage(http.StatusBadRequest, buf)
					}
				}
			}
		}
		ctx.Next()
	}
}
