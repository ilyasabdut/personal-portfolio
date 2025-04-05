package handlers

import (
	"net/http"
	"portfolio/database"
	"portfolio/utils"
)

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	projects, err := database.GetAllProjects()
	if err != nil {
		http.Error(w, "Failed to load projects", http.StatusInternalServerError)
		return
	}

	data := struct {
		Projects []database.Project
	}{
		Projects: projects,
	}

	utils.RenderTemplate(w, "index", data)
}
