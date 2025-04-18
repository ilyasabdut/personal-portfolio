package handlers

import (
	"os"
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
		AppEnv string

	}{
		Projects: projects,
		AppEnv: os.Getenv("APP_ENV"),
	}

	utils.RenderTemplate(w, "index", data)
}
