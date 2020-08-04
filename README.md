# [BDX Web Framework](https://github.com/belldata-dx/bdx)

<img src="https://raw.githubusercontent.com/belldata-dx/bdx-doc-image/master/bxGolang.png" align="right">

bdxは Go(Golang)で記述されたWebフレームワークです。

## Contents

- [BDX Web Framework](#bdx-web-framework)
  - [Contents](#contents)
  - [Installation](#installation)
  - [Quick Start](#quick-start)
  - [API Example](#api-example)
    - [`GET`, `POST`, `PUT`, `DELETE`, `OPTIONS`を提供しています](#get-post-put-delete-optionsを提供しています)

## Installation

`go get github.com/belldata-dx/bdx`

## Quick Start

```sh
cat _examples/main.go
```

```go
package main

import (
  "github.com/belldata-dx/bdx"
  "github.com/belldata-dx/bdx/interfaces"
)

func main() {
  r := bdx.New()
  r.GET("/health", func(c interfaces.Context) {
    c.JSON(200, bdx.B{"data": "test"})
  })
  r.Run()
}
```

```sh
$ go run _examples/main.go
2020/08/03 21:02:17 [INFO] [engine] utils.go:resolveAddress:51: 環境変数`HTTP_PORT`が未定義です。 デフォルトでは`:8080`を使用します。
```

```sh
$ curl localhost:8080/health
{"data":"test"}
```

## API Example

### `GET`, `POST`, `PUT`, `DELETE`, `OPTIONS`を提供しています

```go
func main() {
  router := bdx.New()
  router.GET("/get")
  router.POST("/post")
  router.PUT("/put")
  router.DELETE("/delete")
  router.OPTIONS("/options")
}
```

詳細なサンプルは[ここを参照](_examples/domain-driven-design/examples.go)
