package templates

import (
	"github.com/gin-contrib/multitemplate"
)

var (
	UserTemplate = userTemplate{}
)

type userTemplate struct{}

func (template userTemplate) CreateUserTemplate() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("index", "templates/tpl.html")
	return r
}
