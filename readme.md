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
│   ├── js/              # Main JavaScript file
|   |   ├── modules/     # Separate JavaScript functionalities
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
// for development dynamic reloading as you make changes
```bash
go run main.go 
~/go/bin/air -c air.toml 
```
The site will be available at `http://localhost:3001`

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

# Deploy to Cloud Run // make it eu, stop the unauthenticated access
gcloud run deploy ss-go \
  --image gcr.io/swingsociety-backend/ss-go \
  --platform managed \
  --region us-central1 \ 
  --allow-unauthenticated

```

## 📝 Development Tasks

### MVP
- [ ] validate concept
- [ ] custom domain setup
- [ ] implement back-end logic

### High Priority
- [ ] Redesign Hero Page that demanding focus
- [ ] Forms
  - [ ] Contact form
  - [ ] Registration form
  - [ ] Newsletter subscription
  - [ ] Implement form validation
  - [ ] email validation and email sending logic
  - [ ] Complete email notification system
- [ ] Set up Firestore database
- [ ] Add main pages & route them
  - [ ] Classes
  - [ ] Forms
  - [ ] About US
- [ ] Add proper error handling

### Medium Priority
- [ ] Add event scheduling system
- [ ] Add sub-pages & route them
  - [ ] prices
  - [ ] Terms of Use
  - [ ] GDPR
  - [ ] Privacy policy
  - [ ] Events
  - [ ] Projects
- [ ] Add analytics tracking
- [ ] CMS or Admin dashboard
- [ ] Implement user authentication

### Low Priority
- [ ] Implement logging system
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


