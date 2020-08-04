package di_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/belldata-dx/bdx"
	logger "github.com/belldata-dx/bdx-logger"
	"github.com/belldata-dx/bdx/di"
	"github.com/belldata-dx/bdx/infra"
	"github.com/belldata-dx/bdx/interfaces"
	"github.com/go-playground/assert/v2"
	_ "github.com/lib/pq"
)

type (
	ParentRepo interface {
		Create(model ParentModel) (result *ParentModel, err error)
		Find() (result []*ParentModel)
	}
	IParentUseCase interface {
		Create(model ParentModel) (result *ParentModel, err error)
		Find() (result []*ParentModel)
	}
	ParentModel struct {
		ID   int
		Name string
	}
	parentInfra struct {
		infra.RepositoryImple
	}
	parentUseCaseImpl struct {
		repo ParentRepo
	}

	IStudentParentUseCase interface {
		Create(s StudentModel, p ParentModel) (sresult *StudentModel, presult *ParentModel, err error)
		Find() (sresult []*StudentModel, presult []*ParentModel, err error)
	}
	studentParentUseCaseImpl struct {
		srepo StudentRepo
		prepo ParentRepo
	}
	StudentParentHandle interface {
		Post(c interfaces.Context)
		Get(c interfaces.Context)
	}
	studentParentHandler struct {
		u IStudentParentUseCase
	}
)

var (
	parents = []ParentModel{}
)

func NewStudentParentUseCase(repo StudentRepo, prepo ParentRepo) IStudentParentUseCase {
	return &studentParentUseCaseImpl{repo, prepo}
}

func (sp *studentParentUseCaseImpl) Create(s StudentModel, p ParentModel) (sresult *StudentModel, presult *ParentModel, err error) {
	err = nil
	sresult, err = sp.srepo.Create(s)
	presult, err = sp.prepo.Create(p)
	return
}

func (sp *studentParentUseCaseImpl) Find() (sresult []*StudentModel, presult []*ParentModel, err error) {
	err = nil
	sresult = sp.srepo.Find()
	presult = sp.prepo.Find()
	return
}

func NewStudentParentHandle(u IStudentParentUseCase) StudentParentHandle {
	return &studentParentHandler{u}
}

func (h *studentParentHandler) Post(c interfaces.Context) {
	s, p, _ := h.u.Create(StudentModel{Name: "test"}, ParentModel{Name: "test"})
	c.JSON(200, bdx.B{"s": s, "p": p})
}

func (h *studentParentHandler) Get(c interfaces.Context) {
	s, p, _ := h.u.Find()
	c.JSON(200, bdx.B{"s": s, "p": p})
}

func NewParentInfra(db *infra.DB) ParentRepo {
	return &parentInfra{
		RepositoryImple: infra.RepositoryImple{
			Master: db.Master,
			Slave:  db.Slave,
		},
	}
}

func (u *parentInfra) Create(model ParentModel) (resultModel *ParentModel, err error) {
	model.ID = len(parents) + 1
	parents = append(parents, model)
	return &model, nil
}

func (u *parentInfra) Find() (result []*ParentModel) {
	for _, u := range parents {
		result = append(result, &u)
	}
	return result
}

func NewParentUseCase(repo ParentRepo) IParentUseCase {
	return &parentUseCaseImpl{repo}
}

func (u *parentUseCaseImpl) Create(model ParentModel) (resultModel *ParentModel, err error) {
	return u.repo.Create(model)
}

func (u *parentUseCaseImpl) Find() (result []*ParentModel) {
	return u.repo.Find()
}

func TestMultiDI(t *testing.T) {
	container := di.New()
	singleContainer(container)
	parentInf := di.Increment()
	parentUse := di.Increment()

	container.Set(&di.Definition{
		Name:    parentInf,
		DiName:  []interface{}{di.DB},
		Builder: NewParentInfra,
	})
	container.Set(&di.Definition{
		Name:    parentUse,
		DiName:  []interface{}{parentInf},
		Builder: NewParentUseCase,
	})

	spUse := di.Increment()
	spHandle := di.Increment()
	container.Set(&di.Definition{
		Name:    spUse,
		DiName:  []interface{}{studentInf, parentInf},
		Builder: NewStudentParentUseCase,
	})
	container.Set(&di.Definition{
		Name:    spHandle,
		DiName:  []interface{}{spUse},
		Builder: NewStudentParentHandle,
	})

	router := bdx.Default()
	router.SetLogLevel(logger.Debug)

	spHandler := container.Get(spHandle).(StudentParentHandle)
	router.POST("/user", spHandler.Post)
	router.GET("/user", spHandler.Get)

	request(router, http.MethodPost, "/user", "")
	w := request(router, http.MethodGet, "/user", "")
	read, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, `{"p":[{"ID":1,"Name":"test"}],"s":[{"ID":1,"Name":"test"}]}`, string(read))

}

type header struct {
	Key   string
	Value string
}

func request(r http.Handler, method, path, body string, headers ...header) *httptest.ResponseRecorder {
	headers = append(headers, header{
		Key:   "Content-Type",
		Value: "application/json",
	})
	reader := strings.NewReader(body)
	req := httptest.NewRequest(method, path, reader)
	for _, h := range headers {
		req.Header.Add(h.Key, h.Value)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
