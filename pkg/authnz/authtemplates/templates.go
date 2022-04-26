package authtemplates

import (
	_ "embed"
	"html/template"
	"io"

	"github.com/freifunkMUC/wg-access-server/pkg/authnz/authruntime"
)

var (
	//go:embed base.go.html
	base string
	//go:embed login.go.html
	loginPage string
	//go:embed simpleauth.go.html
	simpleAuthPage string
)

var (
	baseTemplate       = template.Must(template.New("base").Parse(base))
	loginPageTemplate  = template.Must(template.Must(baseTemplate.Clone()).Parse(loginPage))
	simpleAuthTemplate = template.Must(template.Must(baseTemplate.Clone()).Parse(simpleAuthPage))
)

type LoginPage struct {
	Providers []*authruntime.Provider
}

func RenderLoginPage(w io.Writer, data LoginPage) error {
	return loginPageTemplate.Execute(w, data)
}

type SimpleAuthPage struct {
	PostURL      string
	ErrorMessage string
}

func RenderSimpleAuthPage(w io.Writer, data SimpleAuthPage) error {
	return simpleAuthTemplate.Execute(w, data)
}
