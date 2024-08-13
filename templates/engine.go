package templates

import (
	"sync"
	"text/template"
)

var engine Engine

type Engine struct {
	Mutex sync.RWMutex
	Template *template.Template
}

func Init(template *template.Template) {
	engine = Engine{
		Template: template,
	}
}

func E() *Engine {
	return &engine
}
