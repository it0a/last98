package index

import (
	"github.com/it0a/last98/page"
	"html/template"
	"net/http"
)

func IndexHandler(response http.ResponseWriter, request *http.Request) {
	data := struct {
		Page page.Page
	}{page.Page{"Index"}}
	tmpl := make(map[string]*template.Template)
	tmpl["index.tmpl"] = template.Must(template.ParseFiles("templates/base.tmpl", "templates/index.tmpl"))
	err := tmpl["index.tmpl"].ExecuteTemplate(response, "base", data)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
}
