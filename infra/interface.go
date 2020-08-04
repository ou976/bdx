package infra

import (
	"github.com/jinzhu/gorm"
)

type RepositoryImple struct {
	Master *gorm.DB
	Slave  *gorm.DB
}
