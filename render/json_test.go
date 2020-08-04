package render_test

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/belldata-dx/bdx/render"
	"github.com/stretchr/testify/assert"
)

func TestJson(t *testing.T) {
	w := httptest.NewRecorder()
	js := render.JSON{
		Data: map[string]string{
			"data": "test",
		},
	}
	js.Render(w)
	reader, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, string(reader), `{"data":"test"}`)
	content := "application/json; charset=utf-8"
	contentRes := w.Header().Get("Content-Type")
	assert.Equal(t, content, contentRes)
}
