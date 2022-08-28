package models

import "github.com/kataras/iris/v12"

type Responser struct {
	StatusCode int
	Message    string
	Err        error
	Data       iris.Map
}