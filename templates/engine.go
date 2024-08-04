package templates

import (
	"crow/oraiplayground/utils"
	"sync"
	"text/template"
)

const EngineCtxKey = utils.CtxKey("TemplateEngine")

type Engine struct {
	Mutex sync.RWMutex
	Template *template.Template
}
