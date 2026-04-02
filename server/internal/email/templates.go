package email

import (
	"bytes"
	"html/template"
)

// base wraps content in a shared HTML email shell.
const baseTemplate = `<!DOCTYPE html>
<html lang="bg">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<style>
  body { margin: 0; padding: 0; background: #f7f0e8; font-family: Georgia, serif; color: #3a3030; }
  .wrapper { max-width: 600px; margin: 2rem auto; background: #fff; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 12px rgba(0,0,0,0.1); }
  .header { background: #1c1410; padding: 2rem; text-align: center; }
  .header h1 { margin: 0; color: #ffce1b; font-size: 1.6rem; letter-spacing: 0.06em; }
  .header p { margin: 0.25rem 0 0; color: #f0e6d3; font-size: 0.85rem; }
  .body { padding: 2rem 2.5rem; }
  .body h2 { color: #1c1410; margin-top: 0; }
  .body p { line-height: 1.7; margin-bottom: 1rem; }
  .divider { border: none; border-top: 2px solid #ffce1b; margin: 1.5rem 0; }
  .footer { background: #f7f0e8; padding: 1rem 2.5rem; text-align: center; font-size: 0.8rem; color: #888; }
  .footer a { color: #1c1410; }
  .detail-row { display: flex; gap: 0.5rem; margin-bottom: 0.5rem; font-size: 0.95rem; }
  .detail-label { font-weight: bold; min-width: 100px; color: #1c1410; }
</style>
</head>
<body>
<div class="wrapper">
  <div class="header">
    <h1>🎷 Swing Society</h1>
    <p>Where Jazz Happens</p>
  </div>
  <div class="body">{{.Body}}</div>
  <div class="footer">
    <p>© 2025 Swing Society · <a href="https://swingsociety.bg">swingsociety.bg</a></p>
    <p>ул. Сердика 25, София · <a href="mailto:info@swingsociety.bg">info@swingsociety.bg</a></p>
  </div>
</div>
</body>
</html>`

func renderBase(body string) (string, error) {
	tmpl, err := template.New("base").Parse(baseTemplate)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, map[string]template.HTML{"Body": template.HTML(body)}); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// welcomeEmail generates a registration welcome email body.
func welcomeEmail(name string) (string, error) {
	tmpl, err := template.New("welcome").Parse(`
<h2>Здравейте, {{.Name}}! 👋</h2>
<p>Радваме се, че се присъединявате към общността на <strong>Swing Society</strong>!</p>
<p>Вашата регистрация беше получена успешно. Скоро ще се свържем с вас с подробности за следващите стъпки.</p>
<hr class="divider">
<p><em>Hi {{.Name}}! We're excited to have you join the Swing Society community. Your registration has been received and we'll be in touch soon with next steps.</em></p>
<hr class="divider">
<p>До скоро на паркета! 🕺💃<br><strong>Отборът на Swing Society</strong></p>
`)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, map[string]string{"Name": name}); err != nil {
		return "", err
	}
	return renderBase(buf.String())
}

// registrationNotificationEmail generates an admin alert for a new registration.
func registrationNotificationEmail(name, email, phone string) (string, error) {
	tmpl, err := template.New("reg-notify").Parse(`
<h2>Нова регистрация 🎉</h2>
<p>Получена е нова регистрация от сайта:</p>
<div class="detail-row"><span class="detail-label">Име:</span><span>{{.Name}}</span></div>
<div class="detail-row"><span class="detail-label">Имейл:</span><span>{{.Email}}</span></div>
<div class="detail-row"><span class="detail-label">Телефон:</span><span>{{.Phone}}</span></div>
<hr class="divider">
<p>Моля, свържете се с тях в рамките на 48 часа.</p>
`)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, map[string]string{"Name": name, "Email": email, "Phone": phone}); err != nil {
		return "", err
	}
	return renderBase(buf.String())
}

// contactNotificationEmail generates an admin alert for a new contact form submission.
func contactNotificationEmail(name, email, message string) (string, error) {
	tmpl, err := template.New("contact-notify").Parse(`
<h2>Ново запитване 📩</h2>
<p>Получено е ново съобщение от формата за контакт:</p>
<div class="detail-row"><span class="detail-label">Подател:</span><span>{{.Name}}</span></div>
<div class="detail-row"><span class="detail-label">Имейл:</span><span>{{.Email}}</span></div>
<hr class="divider">
<p><strong>Съобщение:</strong></p>
<p>{{.Message}}</p>
`)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, map[string]string{"Name": name, "Email": email, "Message": message}); err != nil {
		return "", err
	}
	return renderBase(buf.String())
}

// newsletterConfirmationEmail generates a newsletter subscription confirmation.
func newsletterConfirmationEmail(name string) (string, error) {
	tmpl, err := template.New("newsletter").Parse(`
<h2>Абонирани сте! 🎶</h2>
<p>Здравейте{{if .Name}}, <strong>{{.Name}}</strong>{{end}}!</p>
<p>Успешно се абонирахте за бюлетина на <strong>Swing Society</strong>. Ще получавате известия за предстоящи класове, партита и събития.</p>
<hr class="divider">
<p><em>You're subscribed to the Swing Society newsletter. We'll keep you updated on upcoming classes, parties, and events.</em></p>
<hr class="divider">
<p>С джазови поздрави,<br><strong>Отборът на Swing Society</strong></p>
`)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, map[string]string{"Name": name}); err != nil {
		return "", err
	}
	return renderBase(buf.String())
}
