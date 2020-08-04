package infra

import (
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestDB(t *testing.T) {
	os.Setenv("REPLICA", "Y")
	var count1 int
	var count2 int
	sql := `select count(*) from information_schema."tables"`
	db := NewDBInit()
	err := db.Master.Raw(sql).Row().Scan(&count1)
	if err != nil {
		t.Error(err.Error())
	}
	err = db.Slave.Raw(sql).Row().Scan(&count2)
	assert.Equal(t, count1, count2)
}
