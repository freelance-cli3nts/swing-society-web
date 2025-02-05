# 📜 Swing Society Website

A modern, high-performance website for Swing Society dance school, built with Go, HTMX, and HTML/CSS. Focusing on simplicity and speed while providing dynamic content updates without heavy JavaScript.

## 🎯 Current Status

- ✅ Basic infrastructure setup
- ✅ HTMX integration for dynamic content
- ✅ Basic routing and template serving
- ✅ Static file handling
- 🏗️ Form handling (In Progress)
- 🏗️ Email notifications (In Progress)

## 📂 Project Structure
```
swing-society/
├── server/           # Backend server logic
│   ├── internal/         # Core backend logic
│   │   ├── api/             # API handlers
│   │   ├── config/          # Configuration
│   │   ├── middleware/      # HTTP middleware
│   │   └── monitoring/      # Metrics and monitoring
│   └── main.go          # Entry point
├── static/           # Static assets
│   ├── css/             # Stylesheets
│   ├── js/              # JavaScript files
│   ├── data/            # Data files (schedules, etc.)
│   └── assets/          # Media assets
│       ├── images/          # Image files
│       └── videos/          # Video files
└── templates/        # HTML templates
    ├── classes/          # Class-related templates
    ├── events/           # Event templates
    └── forms/            # Form templates
```

## 🚀 Features

### Implemented
* ✅ Server-side routing with Go
* ✅ Dynamic content loading with HTMX
* ✅ Mobile-responsive design
* ✅ Static file serving
* ✅ Template-based pages
* ✅ Basic monitoring metrics

### In Development
* 🏗️ Form submissions
* 🏗️ Email notifications
* 🏗️ Database integration
* 🏗️ Analytics tracking
* 🏗️ Event scheduling system

## 🛠️ Technologies

- **Backend**: Go 1.21+
- **Frontend**: HTML5, CSS3, HTMX
- **Database**: Firestore (planned)
- **Deployment**: Google Cloud Run
- **CI/CD**: Google Cloud Build (planned)

## 📌 Requirements

- Go 1.21 or higher
- Google Cloud SDK (for deployment)
- `.env` file for environment variables

## 🚦 Getting Started

1. **Clone the Repository**
```bash
git clone https://github.com/yourusername/swing-society.git
cd swing-society
```

2. **Set Up Environment Variables**
Create a `.env` file in the root directory:
```env
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USER=your-email@example.com
SMTP_PASS=your-email-password
RECIPIENT_EMAIL=info@swingsociety.bg
GOOGLE_CLOUD_PROJECT=swingsociety-backend
```

3. **Install Dependencies**
```bash
cd server
go mod tidy
```

4. **Run Locally**
```bash
go run main.go
```
The site will be available at `http://localhost:8080`

## 🚀 Deployment

### Local Development
```bash
go run main.go
```

### Google Cloud Run Deployment
```bash
# Build the container
docker build -t gcr.io/swingsociety-backend/ss-go .

# Push to Container Registry
docker push gcr.io/swingsociety-backend/ss-go

# Deploy to Cloud Run
gcloud run deploy ss-go \
  --image gcr.io/swingsociety-backend/ss-go \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated
```

## 📝 Development Tasks

### MVP
- [ ] validate concept
- [ ] finish back-end logic
- [ ] custom domain setup

### High Priority
- [ ] Redesign Hero Page to be focus
- [ ] Complete email notification system
- [ ] Implement form validation
- [ ] Add email validation and actual email sending logic
- [ ] Add all sub-pages & connect them
  - [ ] Add contact form
  - [ ] For us
  - [ ] class
  - [ ] contacts
- [ ] Set up Firestore database
- [ ] Add proper error handling
- [ ] Implement logging system

### Medium Priority
- [ ] Add event scheduling system
- [ ] Implement user authentication
- [ ] Create admin dashboard
- [ ] Add analytics tracking

### Low Priority
- [ ] Implement caching
- [ ] Add performance monitoring
- [ ] Create backup system
- [ ] Implement rate limiting

## 🔗 Resources
- [Project Figma Board](https://www.figma.com/board/YXIM0ao7X6mtGbcUFQQOj8/Swing-Society-web?node-id=7-925&t=weUqv9eMgFAnFSnO-0)
- [HTMX Documentation](https://htmx.org/docs/)
- [Go Documentation](https://golang.org/doc/)

## 🤝 Contributing
Feel free to contribute to this project!
Fork the repo.
Create a feature branch (git checkout -b feature-name).
Make your changes & commit (git commit -m "Description").
Push and open a pull request.


## 📜 License
MIT License - See LICENSE file for details

---
## 🎷 Swing Society – Where Jazz Comes to Life! 🎶
Let me know if you want any changes or additions! 🚀


