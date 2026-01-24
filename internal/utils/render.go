package utils

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func RenderTemplate(w http.ResponseWriter, page string, data interface{}) {
	partials, _ := filepath.Glob("templates/partials/*.html")

	files := append([]string{
		filepath.Join("templates", "layout.html"),
		filepath.Join("templates", page+".html"),
	}, partials...)

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		log.Println("Error parsing templates:", err)
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Execution error", http.StatusInternalServerError)
	}
}
