package render_test

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/belldata-dx/bdx/render"
	"github.com/stretchr/testify/assert"
)

func TestYAML(t *testing.T) {
	w := httptest.NewRecorder()
	yml := render.YAML{
		Data: map[string]string{
			"data": "test",
		},
	}
	yml.Render(w)
	reader, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, string(reader), "data: test\n")
	content := "application/x-yaml; charset=utf-8"
	contentRes := w.Header().Get("Content-Type")
	assert.Equal(t, content, contentRes)

}
