package render_test

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/belldata-dx/bdx/render"
	"github.com/stretchr/testify/assert"
)

func TestJsonAscii(t *testing.T) {
	w := httptest.NewRecorder()
	js := render.JSONAscii{
		Data: map[string]string{
			"data": "test",
		},
	}
	js.Render(w)
	reader, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, string(reader), `{"data":"test"}`)
	content := "application/json"
	contentRes := w.Header().Get("Content-Type")
	assert.Equal(t, content, contentRes)
}
