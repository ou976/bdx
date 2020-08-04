package di_test

import (
	"io/ioutil"
	"net/http"
	"testing"

	bdx "github.com/belldata-dx/bdx"
	logger "github.com/belldata-dx/bdx-logger"
	"github.com/belldata-dx/bdx/di"
	"github.com/belldata-dx/bdx/infra"
	"github.com/belldata-dx/bdx/interfaces"
	"github.com/go-playground/assert/v2"
	_ "github.com/lib/pq"
)

type (
	StudentRepo interface {
		Create(model StudentModel) (result *StudentModel, err error)
		Find() (result []*StudentModel)
	}
	StudentModel struct {
		ID   int
		Name string
	}
	studentInfra struct {
		infra.RepositoryImple
	}
	IStudentUseCase interface {
		Create(model StudentModel) (result *StudentModel, err error)
		Find() (result []*StudentModel)
	}
	studentUseCaseImpl struct {
		repo StudentRepo
	}
	StudentHandle interface {
		Post(c interfaces.Context)
		Get(c interfaces.Context)
	}
	studentHandler struct {
		u IStudentUseCase
	}
)

var (
	students = []StudentModel{}
)

func NewStudentInfra(db *infra.DB) StudentRepo {
	return &studentInfra{
		RepositoryImple: infra.RepositoryImple{
			Master: db.Master,
			Slave:  db.Slave,
		},
	}
}

func (u *studentInfra) Create(model StudentModel) (resultModel *StudentModel, err error) {
	model.ID = len(students) + 1
	students = append(students, model)
	return &model, nil
}

func (u *studentInfra) Find() (result []*StudentModel) {
	for _, u := range students {
		result = append(result, &u)
	}
	return result
}

func NewStudentUseCase(repo StudentRepo) IStudentUseCase {
	return &studentUseCaseImpl{repo}
}

func (u *studentUseCaseImpl) Create(model StudentModel) (resultModel *StudentModel, err error) {
	return u.repo.Create(model)
}

func (u *studentUseCaseImpl) Find() (result []*StudentModel) {
	return u.repo.Find()
}

func NewStudentHandler(u IStudentUseCase) StudentHandle {
	return &studentHandler{u}
}

func (u *studentHandler) Post(c interfaces.Context) {
	s, _ := u.u.Create(StudentModel{Name: "test"})
	c.JSON(200, s)
}

func (u *studentHandler) Get(c interfaces.Context) {
	s := u.u.Find()
	c.JSON(200, s)
}

var (
	studentInf    interface{}
	studentUse    interface{}
	studentHandle interface{}
)

func singleContainer(container *di.Container) {
	studentInf = container.Set(&di.Definition{
		DiName:  []interface{}{di.DB},
		Builder: NewStudentInfra,
	})

	studentUse = container.Set(&di.Definition{
		DiName:  []interface{}{studentInf},
		Builder: NewStudentUseCase,
	})

	studentHandle = container.Set(&di.Definition{
		DiName:  []interface{}{studentUse},
		Builder: NewStudentHandler,
	})
}

func TestSingle(t *testing.T) {
	students = []StudentModel{}
	container := di.New()
	singleContainer(container)
	sHandler := container.Get(studentHandle).(StudentHandle)
	router := bdx.Default()
	router.SetLogLevel(logger.Debug)

	router.POST("/student", sHandler.Post)
	router.GET("/student", sHandler.Get)

	request(router, http.MethodPost, "/student", "")
	w := request(router, http.MethodGet, "/student", "")
	read, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, `[{"ID":1,"Name":"test"}]`, string(read))

}
