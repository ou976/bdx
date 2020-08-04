package render_test

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/belldata-dx/bdx/render"
	"github.com/stretchr/testify/assert"
)

func TestXML(t *testing.T) {
	w := httptest.NewRecorder()
	xml := render.XML{
		Data: "test",
	}
	err := xml.Render(w)
	assert.Equal(t, err, nil)
	reader, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, string(reader), `<string>test</string>`)
	content := "application/xml; charset=utf-8"
	contentRes := w.Header().Get("Content-Type")
	assert.Equal(t, content, contentRes)
}
