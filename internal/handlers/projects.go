package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"

	"portfolio/internal/database"
)

func SearchProjectsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("search")

	projects, err := database.SearchProjects(query)
	if err != nil {
		http.Error(w, "Search failed", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles(filepath.Join("web/templates", "partials", "project_list.html")))
	tmpl.ExecuteTemplate(w, "project_list", struct {
		Projects []database.Project
	}{Projects: projects})
}
