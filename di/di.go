package di

import (
	"errors"
	"reflect"

	"github.com/belldata-dx/bdx/infra"
)

var cache = make(map[interface{}]reflect.Value)

type (
	// Definition モジュール名と生成の方法の構造体
	Definition struct {
		Name    interface{}
		Builder interface{}
		DiName  []interface{}
	}
	// Container DI Container
	Container struct {
		definitions map[interface{}]*Definition
	}
)

type DiType int

const (
	DB DiType = 0
)

var diType = 0

func Increment() DiType {
	diType++
	return DiType(diType)
}

func (d *Definition) get(values ...reflect.Value) reflect.Value {
	if val, ok := cache[d.Name]; ok {
		return val
	}
	fv := reflect.ValueOf(d.Builder)
	result := fv.Call(values)
	val := result[0]
	cache[d.Name] = val
	return val
}

// New DI Containerコンストラクタ
func New() *Container {
	return &Container{
		definitions: map[interface{}]*Definition{
			DB: &Definition{
				Name:    DB,
				Builder: infra.NewDBInit,
			},
		},
	}
}

// Set DI Containerにモジュールを登録する
func (c *Container) Set(d *Definition) interface{} {
	if d.Name == nil {
		d.Name = Increment()
	}
	c.definitions[d.Name] = d
	return d.Name
}

// Get DI Containerからモジュールを取り出す
func (c *Container) get(key interface{}) reflect.Value {
	if d, ok := c.definitions[key]; ok {
		values := []reflect.Value{}
		for _, name := range d.DiName {
			val := c.get(name)
			values = append(values, val)
		}
		result := d.get(values...)
		return result
	}
	return reflect.Value{}
}

// Get DI Containerからモジュールを取り出す
//
// この時に依存関係は全て解決される。
func (c *Container) Get(key interface{}) interface{} {
	if d, ok := c.definitions[key]; ok {
		values := []reflect.Value{}
		for _, name := range d.DiName {
			if key == name {
				panic(errors.New("自身の名前が依存関係に設定されています。"))
			}
			val := c.get(name)
			values = append(values, val)
		}
		result := d.get(values...)
		return result.Interface()
	}
	return nil
}
