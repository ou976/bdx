package examples

import (
	bdx "github.com/belldata-dx/bdx"
	logger ""
	"github.com/belldata-dx/bdx/di"
	"github.com/belldata-dx/bdx/infra"
	"github.com/belldata-dx/bdx/interfaces"
)

type (
	// StudentModel モデル
	StudentModel struct {
		ID   int
		Name string
	}

	// StudentRepo リポジトリインターフェイス
	StudentRepo interface {
		Create(model StudentModel) (result *StudentModel, err error)
		Find() (result []*StudentModel)
	}

	// studentInfra infrastructure層のインターフェイス
	studentInfra struct {
		infra.RepositoryImple
	}

	// IStudentUseCase UseCase層のインターフェイス
	IStudentUseCase interface {
		Create(model StudentModel) (result *StudentModel, err error)
		Find() (result []*StudentModel)
	}

	// studentUseCaseImpl UseCase層の実装を持つ構造体
	studentUseCaseImpl struct {
		repo StudentRepo
	}

	// StudentHandle application層のインターフェイス
	StudentHandle interface {
		Post(c interfaces.Context)
		Get(c interfaces.Context)
	}

	// studentHandler application層の実装を持つ構造体
	studentHandler struct {
		u IStudentUseCase
	}
)

var (
	students = []StudentModel{}
)

// NewStudentInfra infrastructure層の実装を持つインスタンス
func NewStudentInfra(db *infra.DB) StudentRepo {
	return &studentInfra{
		RepositoryImple: infra.RepositoryImple{
			Master: db.Master,
			Slave:  db.Slave,
		},
	}
}

// Create infrastructure層の実装
func (u *studentInfra) Create(model StudentModel) (resultModel *StudentModel, err error) {
	model.ID = len(students) + 1
	students = append(students, model)
	return &model, nil
}

// Find infrastructure層の実装
func (u *studentInfra) Find() (result []*StudentModel) {
	for _, u := range students {
		result = append(result, &u)
	}
	return result
}

// NewStudentUseCase UseCase層の実装を持つインスタンス
func NewStudentUseCase(repo StudentRepo) IStudentUseCase {
	return &studentUseCaseImpl{repo}
}

// Create UseCase層の実装
func (u *studentUseCaseImpl) Create(model StudentModel) (resultModel *StudentModel, err error) {
	return u.repo.Create(model)
}

// Find UseCase層の実装
func (u *studentUseCaseImpl) Find() (result []*StudentModel) {
	return u.repo.Find()
}

// NewStudentHandler application層の実装を持つインスタンス
func NewStudentHandler(u IStudentUseCase) StudentHandle {
	return &studentHandler{u}
}

// Post applicatoin層の実装
func (u *studentHandler) Post(c interfaces.Context) {
	s, _ := u.u.Create(StudentModel{Name: "test"})
	c.JSON(200, s)
}

// Get application層の実装
func (u *studentHandler) Get(c interfaces.Context) {
	s := u.u.Find()
	c.JSON(200, s)
}

var (
	studentInf    interface{}
	studentUse    interface{}
	studentHandle interface{}
)

func main() {

	// DIコンテナの生成
	container := di.New()

	// コンテナへセット
	studentInf = container.Set(&di.Definition{
		// DI
		DiName:  []interface{}{di.DB},
		Builder: NewStudentInfra,
	})

	// コンテナへセット
	studentUse = container.Set(&di.Definition{
		// DI
		DiName:  []interface{}{studentInf},
		Builder: NewStudentUseCase,
	})

	// コンテナへセット
	studentHandle = container.Set(&di.Definition{
		// DI
		DiName:  []interface{}{studentUse},
		Builder: NewStudentHandler,
	})

	// コンテナから取得(この時に依存関係が全て解決される。)
	sHandler := container.Get(studentHandle).(StudentHandle)

	// ルーターのインスタンス生成
	router := bdx.Default()
	// ログレベルの設定
	router.SetLogLevel(logger.Debug)

	router.POST("/student", sHandler.Post)
	router.GET("/student", sHandler.Get)

	router.Run()
}
