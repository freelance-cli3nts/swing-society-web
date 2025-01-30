package internal

// import (
//     "log"
//     "net/http"
//     "os"
//     "path/filepath"
// )

// validation of forms & email notifications API	
// type FormData struct {
// 	Name           string
// 	Email          string
// 	Phone          string
// 	Message        string
// 	PreferredContact string
// }


// Success response HTML
// const successHTML = `
// <div class="success-message" style="text-align: center; padding: 2rem;">
// <h3 style="color: #28a745;">Благодарим за съобщението!</h3>
// <p>Ще се свържем с вас възможно най-скоро.</p>
// </div>
// `




// // contact form handler function
// func serveContactForm(w http.ResponseWriter, r *http.Request) {
// 	http.ServeFile(w, r, "templates/contact-form.html")
// }


// // form response handling
// func handleFormSubmission(w http.ResponseWriter, r *http.Request) {
// 	// Only allow POST method
// 	if r.Method != http.MethodPost {
// 			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 			return
// 	}

// 	// Parse form data
// 	err := r.ParseForm()
// 	if err != nil {
// 			http.Error(w, "Error parsing form", http.StatusBadRequest)
// 			return
// 	}

// 	// Get form data
// 	formData := FormData{
// 			Name:            r.FormValue("name"),
// 			Email:           r.FormValue("email"),
// 			Phone:           r.FormValue("phone"),
// 			Message:         r.FormValue("message"),
// 			PreferredContact: r.FormValue("menu-item"),
// 	}

// 	// Validate required fields
// 	if formData.Name == "" || formData.Message == "" {
// 			http.Error(w, "Required fields missing", http.StatusBadRequest)
// 			return
// 	}

// 	// Validate at least one contact method
// 	if formData.Email == "" && formData.Phone == "" {
// 			http.Error(w, "Please provide either email or phone", http.StatusBadRequest)
// 			return
// 	}

// 	// Send email notification
// 	err = sendEmailNotification(formData)
// 	if err != nil {
// 			log.Printf("Error sending email: %v", err)
// 			// Continue anyway - we don't want to show an error to the user if email fails
// 	}

// 	// Return success message
// 	w.Header().Set("Content-Type", "text/html; charset=utf-8")
// 	fmt.Fprint(w, successHTML)
// }

// func handleFormClose(w http.ResponseWriter, r *http.Request) {
// 	// Only allow GET method
// 	if r.Method != http.MethodGet {
// 			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 			return
// 	}

// 	// Return empty response
// 	w.Header().Set("Content-Type", "text/html; charset=utf-8")
// 	fmt.Fprint(w, "")
// }

// func sendEmailNotification(data FormData) error {
// 	// Get email configuration from environment variables
// 	smtpHost := os.Getenv("SMTP_HOST")
// 	smtpPort := os.Getenv("SMTP_PORT")
// 	smtpUser := os.Getenv("SMTP_USER")
// 	smtpPass := os.Getenv("SMTP_PASS")
// 	recipientEmail := os.Getenv("RECIPIENT_EMAIL")

// 	if smtpHost == "" || smtpPort == "" || smtpUser == "" || smtpPass == "" || recipientEmail == "" {
// 			return fmt.Errorf("missing email configuration")
// 	}

// 	// Create email message
// 	m := gomail.NewMessage()
// 	m.SetHeader("From", smtpUser)
// 	m.SetHeader("To", recipientEmail)
// 	m.SetHeader("Subject", "New Contact Form Submission")
	
// 	// Create email body
// 	body := fmt.Sprintf(`
// 	New message from Swing Society website:
	
// 	Name: %s
// 	Email: %s
// 	Phone: %s
// 	Preferred Contact: %s
// 	Message: %s
// 	`, data.Name, data.Email, data.Phone, data.PreferredContact, data.Message)
	
// 	m.SetBody("text/plain", body)

// 	// Create email dialer
// 	d := gomail.NewDialer(smtpHost, 587, smtpUser, smtpPass)

// 	// Send email
// 	if err := d.DialAndSend(m); err != nil {
// 			return fmt.Errorf("failed to send email: %v", err)
// 	}

// 	return nil
// }
