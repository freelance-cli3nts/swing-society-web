// internal/api/handlers/templates.go
package handlers

import (
    "net/http"
    "strings"
    "swing-society-website/server/internal/api/response"
    customerrors "swing-society-website/server/internal/errors"
    "swing-society-website/server/internal/storage"
)

type TemplateHandler struct {
    storage storage.TemplateStorage
}

func NewTemplateHandler(storage storage.TemplateStorage) *TemplateHandler {
    return &TemplateHandler{
        storage: storage,
    }
}

func (h *TemplateHandler) HandleTemplate(w http.ResponseWriter, r *http.Request, templatePath string) {
    cleanPath := strings.TrimPrefix(templatePath, "/templates")
    
    w.Header().Set("Cache-Control", "private, max-age=10")
    
    if !h.storage.TemplateExists(cleanPath) {
        response.Error(w, customerrors.NewAppError(
            customerrors.ErrorTypeFileSystem,
            "Template not found",
            nil,
            map[string]string{"path": cleanPath},
        ))
        return
    }
    
    if r.Header.Get("HX-Request") == "true" {
        content, err := h.storage.GetTemplate(cleanPath)
        if err != nil {
            response.Error(w, customerrors.NewInternalError(
                "Failed to read template",
                err,
            ))
            return
        }
        response.HTMLFragment(w, http.StatusOK, content)
        return
    }

    content, err := h.storage.GetIndexTemplate()
    if err != nil {
        response.Error(w, customerrors.NewInternalError(
            "Failed to read index template",
            err,
        ))
        return
    }
    response.HTMLFragment(w, http.StatusOK, content)
}
