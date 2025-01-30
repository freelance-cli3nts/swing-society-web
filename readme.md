📜 Swing Society Website
A fast, lightweight, and interactive website for Swing Society, built with Go, HTMX, and HTML/CSS.

📂 Project Structure
swing-society/
│── static/            # CSS, JS, images, fonts
│── templates/         # HTML files (views)
│── data/              # JSON or other static data
│── .env               # Environment variables
│── main.go            # Main Go application
│── go.mod             # Go module file
│── README.md          # Project documentation

Project New Features Board:
https://www.figma.com/board/YXIM0ao7X6mtGbcUFQQOj8/Swing-Society-web?node-id=7-925&t=weUqv9eMgFAnFSnO-0


🚀 Features
✅ Static website with dynamic components using HTMX library (https://htmx.org/)
✅ Forms submission with Go server-side processing
✅ Automatic email notifications for forms submissions
✅ Mobile-responsive navigation with hamburger menu
✅ Lightweight, fast, and simple Go backend

📌 Requirements
Before running the project, make sure you have:

Go installed (go version to check)
.env file with SMTP email credentials (for email notifications)
🛠️ Installation & Setup
1️⃣ Clone the repository
git clone https://github.com/yourusername/swing-society.git
cd swing-society
2️⃣ Initialize Go module (if not done yet)
go mod tidy
3️⃣ Set up environment variables
Create a .env file in the project root:
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USER=your-email@example.com
SMTP_PASS=your-email-password
RECIPIENT_EMAIL=info@swingsociety.bg
4️⃣ Run the server
go run main.go
Your site will now be available at http://localhost:8080 🎉

🔧 Endpoints
Endpoint	Method	Description
/	GET	Serves the homepage
/contact-form	GET	Loads the contact form (HTMX)
/submit-form	POST	Handles form submission (Go backend)
/close-form	GET	Closes the form
/static/	GET	Serves static assets (CSS, JS, Images)

📧 Contact Form Handling
Uses HTMX for dynamic updates without JavaScript.
Sends user messages via email notifications.
Prevents empty form submissions & validates input.

🖥️ Deployment
Option 1: Run on a VPS
Install Go on the server.
Clone the repository.
Set up environment variables (.env).
Run the server:
nohup go run main.go &
Use Nginx as a reverse proxy for production.

Option 2: Deploy with Docker (Coming Soon)

🤝 Contributing
Feel free to contribute to this project!
Fork the repo.
Create a feature branch (git checkout -b feature-name).
Make your changes & commit (git commit -m "Description").
Push and open a pull request.

📜 License
This project is licensed under the MIT License.

🎷 Swing Society – Bringing Jazz to Life! 🎶
Let me know if you want any changes or additions! 🚀

# Backlog of tasks
## add for us, class and contacts page





🎶 Нова група по SWING танци за начинаещи! 🎶🗓️ Дата на старт: 15 януари 2025🕢 Време: всеки понеделник и сряда от 19:30 ч.📍 Локация: ул. Сердика 25: https://maps.app.goo.gl/6FUq5nkWctsXrxMQA💃🕺 Предишен опит не е необходим. Няма нужда и от партньор.🎷 Искаш ли да танцуваш един от най-забавните и социални танци? Ще учим най-популярния SWING танц - Линди Хоп, който се появил в Ню Йорк през 30-те години на миналия век.📅 Курсът се състои от 16 класа по 70 минути.💲 Цена на курса: 100 лв на месец или 180 лв за целия курс🎉 Предизвикай себе си и стани част от световното SWING и JAZZ танцувално общество, в което те очакват ритъм, стил и много нови приятели!


Firefox:
ctr + shift + m - toggle between mobile and desktop view in dev options
ctr + r - refresh the page

VSCode: 
Ctr+W - close window