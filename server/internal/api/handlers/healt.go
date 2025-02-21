package handlers

import (
	"net/http"
	"swing-society-website/server/internal/api/response"
	"swing-society-website/server/internal/config"
)

func HandleHealth(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{
		"status":      "ok",
		"environment": config.AppConfig.Environment,
	})
}
