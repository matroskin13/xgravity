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

func GetHTTPTemplate(parentPackage string, entities []Entity) []byte {
	var buf bytes.Buffer

	funcs := template.FuncMap{
		"toLower": strings.ToLower,
	}

	template.Must(template.New("http_template").Funcs(funcs).Parse(`package api

import (
    "errors"
	"net/http"
    // parent "{{.parentPackage}}"

    "github.com/matroskin13/xgravity"
)

func Start({{range .entities}}entity{{.Name}} {{.Name}}Api,{{end}} port string) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		id := xgravity.GetParam(r.URL.Path, 1)

		{{range $entity := .entities}}
			{{range $entity.Methods}}
				{{ if (eq . "get")}}
					if r.URL.Path == "/{{$entity.Name | toLower}}" && r.Method == "GET" && id > 0 {
						advert, err := entity{{$entity.Name}}.get{{$entity.Name}}ById(id)
						if err != nil {
							xgravity.ErrorResponse(w, err)
							return
						}

						xgravity.SuccessResponse(w, advert)

						return
					}
				{{ end }}
			{{end}}
		{{end}}

		xgravity.ErrorResponse(w, errors.New("not found"))
	})
	http.ListenAndServe(":"+port, nil)
}	`)).Execute(&buf, map[string]interface{}{
		"entities":      entities,
		"parentPackage": parentPackage,
	})

	b, _ := format.Source(buf.Bytes())

	return b
}
