
package handlers

import (
    "net/http"
    "swing-society-website/server/internal/api/models"
    "swing-society-website/server/internal/api/response"
    customerrors "swing-society-website/server/internal/errors"
    "swing-society-website/server/internal/storage"
)

type ClassHandler struct {
	storage storage.ClassStorage
}

func NewClassHandler(storage storage.ClassStorage) *ClassHandler {
	return &ClassHandler{
			storage: storage,
	}
}

func (h *ClassHandler) HandleClass(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
			response.Error(w, customerrors.NewValidationError(
					"Method not allowed",
					map[string]string{"method": "Only POST method is allowed"},
			))
			return
	}

	var class models.ClassInquiry
	if err := r.ParseForm(); err != nil {
			response.Error(w, customerrors.NewValidationError(
					"Error parsing form",
					map[string]string{"form": "Invalid form data"},
			))
			return
	}

	// Parse form values
	class.Name = r.FormValue("name")
	class.Email = r.FormValue("email")
	class.Phone = r.FormValue("phone")
	class.ClassType = r.FormValue("classType")
	class.Level = r.FormValue("level")
	class.Message = r.FormValue("message")

	if errors := class.Validate(); len(errors) > 0 {
			response.Error(w, customerrors.NewValidationError(
					"Validation failed",
					errors,
			))
			return
	}

	// Store class inquiry
	if err := h.storage.StoreInquiry(&class); err != nil {
			response.Error(w, customerrors.NewInternalError(
					"Failed to store class inquiry",
					err,
			))
			return
	}

	response.JSON(w, http.StatusOK, map[string]string{
			"message": "Thank you for your interest! We'll contact you with class details soon.",
	})
}