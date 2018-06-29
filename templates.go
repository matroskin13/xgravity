package xgravity

import (
	"bytes"
	"go/format"
	"strings"
	"text/template"
)

type HTTPTemplateData struct {
	InterfaceName string
}

func GetEndpointTemplate(iName string, endpoint *Endpoint) []byte {
	str := `e.{{.method}}("{{.path}}", func(c echo.Context) error {
    {{.iName}}.{{.methodName}}(1)
	return c.String(404, "not found")
})`
	t := template.Must(template.New("http_endpoint").Parse(str))

	var buf bytes.Buffer

	t.Execute(&buf, map[string]interface{}{
		"method":     endpoint.Method,
		"path":       endpoint.Path,
		"iName":      iName,
		"methodName": endpoint.Name,
	})

	return buf.Bytes()
}

func GetHTTPTemplate(parentPackage string, entities []Entity) ([]byte, error) {
	var buf bytes.Buffer

	funcs := template.FuncMap{
		"toLower": strings.ToLower,
	}

	var sEndpoints []string

	for _, entity := range entities {
		for _, endpoint := range entity.Endpoints {
			sEndpoints = append(
				sEndpoints,
				string(GetEndpointTemplate("entity"+entity.Name, &endpoint)),
			)
		}
	}

	t := template.Must(template.New("http_template").Funcs(funcs).Parse(`package api

import (
    "errors"
	"net/http"


    // parent "{{.parentPackage}}"
	"github.com/labstack/echo"
    "github.com/matroskin13/xgravity"
)

func Start({{range .entities}}entity{{.Name}} {{.Name}}Api,{{end}} port string) error {
	e := echo.New()

	{{.sEndpoints}}

	return e.Start(":"+port)
}	`))

	t.Execute(&buf, map[string]interface{}{
		"entities":      entities,
		"parentPackage": parentPackage,
		"sEndpoints":    strings.Join(sEndpoints, ""),
	})

	return format.Source(buf.Bytes())
}
