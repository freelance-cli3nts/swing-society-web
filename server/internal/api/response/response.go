// internal/api/response/response.go

package response

import (
    "encoding/json"
    "errors" // Added built-in errors package
    "log"
    "net/http"
    customerrors "swing-society-website/server/internal/errors" // Aliased our custom errors package
)

// Response represents a standardized API response
type Response struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   *ErrorInfo  `json:"error,omitempty"`
}

// ErrorInfo represents error information in responses
type ErrorInfo struct {
    Type    string            `json:"type"`
    Message string            `json:"message"`
    Details map[string]string `json:"details,omitempty"`
}

// JSON sends a JSON response with the given status code
func JSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    
    response := Response{
        Success: status >= 200 && status < 300,
        Data:    data,
    }
    
    // Ensure empty slices are serialized as [] instead of null
    if data == nil {
        response.Data = make([]interface{}, 0)
    }
    
    if err := json.NewEncoder(w).Encode(response); err != nil {
        log.Printf("Error encoding response: %v", err)
    }
}

// HTMLFragment sends an HTML fragment response (for HTMX requests)
func HTMLFragment(w http.ResponseWriter, status int, html string) {
    w.Header().Set("Content-Type", "text/html")
    w.WriteHeader(status)
    w.Write([]byte(html))
}

// Error handles error responses
func Error(w http.ResponseWriter, err error) {
    var appErr *customerrors.AppError
    if !errors.As(err, &appErr) {
        // If it's not an AppError, wrap it as an internal error
        appErr = customerrors.NewInternalError("An unexpected error occurred", err)
    }

    response := Response{
        Success: false,
        Error: &ErrorInfo{
            Type:    string(appErr.Type),
            Message: appErr.Message,
            Details: appErr.Details,
        },
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(appErr.HTTPStatus)
    
    if err := json.NewEncoder(w).Encode(response); err != nil {
        log.Printf("Error encoding error response: %v", err)
    }
}

// Example usage in a handler:
func ExampleHandler(w http.ResponseWriter, r *http.Request) {
    // Success case
    data := map[string]string{"message": "Success"}
    JSON(w, http.StatusOK, data)
    
    // Error case
    err := customerrors.NewValidationError("Invalid input", map[string]string{
        "email": "Invalid email format",
    })
    Error(w, err)
}