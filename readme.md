# 📜 Swing Society Website

A modern, high-performance website for Swing Society dance school, Focusing on simplicity and speed, while providing dynamic content and mobile friendly for better UX. 

## 🎯 Current Status

- ✅ Basic infrastructure setup
- ✅ HTMX integration for dynamic content
- ✅ Basic routing and template serving
- ✅ Static file handling
- ✅ Firebase Realtime Database integration
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
* ✅ Firebase Realtime Database integration
* ✅ Configurable CORS settings

### In Development
* 🏗️ Form submissions
* 🏗️ Email notifications
* 🏗️ Analytics tracking
* 🏗️ Event scheduling system

## 🛠️ Technologies

- **Backend**: Go 1.21+
- **Frontend**: HTML5, CSS3, HTMX, Bulma
- **Database**: Firebase Realtime Database
- **Deployment**: Google Cloud Run
- **Containerization**: Docker
- **CDN**: Cloudflare
- **CI/CD**: Google Cloud Build (planned)

## 📌 Requirements

- Go 1.21 or higher
- Google Cloud SDK (for deployment)
- Firebase Admin SDK credentials
- `.env` file for environment variables
- `config.json` for application configuration

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
# For development with live reloading
air -c air.toml

# Or using Go directly
go run server/main.go
```
The site will be available at `http://localhost:8080`

  ## Deployment

  GCP Infrastructure
  Resource                Value
  Backend                 projectswingsociety-backend
  Firebase                projectswing-society-realtime-data
  Cloud Run               servicess-go
  Region                  europe-north1
  Service account         ss-go1@swingsociety-backend.iam.gserviceaccount.com
  Artifact Registry       europe-north1-docker.pkg.dev/swingsociety-backend/ss-go/ss-go

  ### First-time Setup (new collaborator)
  1. Request GCP access from the project owner (mr.anastasov@gmail.com). You need roles/editor on swingsociety-backend.
  2. Authenticate:
  
  ```bash  
  gcloud auth login
  gcloud auth configure-docker europe-north1-docker.pkg.dev
  ```
  3. Deploy
  ```
  bash
  ./deploy.sh
  ```

  The script builds the Docker image, pushes to Artifact Registry, and deploys to Cloud Run with the correct service account. No credential files needed.

  ***Firebase Authentication***
  Firebase access uses GCP Application Default Credentials (ADC) — no key files. On Cloud Run, ss-go1 service account authenticates automatically. For local development, set GOOGLE_CREDENTIALS in .env with the JSON content of a service account key (never a file path, never committed).
  
  ***Security***

  No credential files in the repository — enforced via .gitignore and gitleaks pre-commit hook
  Firebase auth via Workload Identity (ADC) — zero service account keys in use
  Cloud Run configured with dedicated service account ss-go1 (not the default compute SA)
  Unauthenticated Cloud Run access disabled

## 📝 Development Tasks

### MVP
- [x] validate concept
- [x] custom domain setup
- [x] implement back-end logic

### High Priority
- [ ] Design 
   [ ] Redesign Hero Page that demands focus
     [ ] Implement the carousel feature 
     [ ] Fix positioning on mobile/tablet view of button and text
   [ ] Add main pages & route them
    [ ] Landing page
      [ ] header
        [ ] nav - mobile desktop tablet
        [ ] link all pages
        [ ] logo functionality
        [*] social icons proper links
        [*] sign up button
      [ ] main 
        [ ] hero section - make it swappable container
          [ ] mobile - desktop - tablet view of hero
        [ ] dynamic calendar
        [ ] add two column container 4 classes & blog
      [ ] footer
        [ ] mobile - desktop - tablet view of footer
      [ ] Dynamic calendar: 60 - 120 min
  [ ] Class view - MVP: 30 min -> if time -> fancy version: 60 min
    [ ] board 4 all classes
      [ ] coupled
      [ ] solo
      [ ] old clips
    [*] Forms
    [ ] About US
    [ ] Fix colors
    [ ] Layouts
  [x] Coupled Classes
    [o] Carousel view - fix image & video loading
      [ ] Solo Classes
      [ ] Old Clips
      [ ] Teachers section
      [ ] Founders section
      [ ] Organization Concept
      [ ] Parties
      [ ] Festival Page
      [ ] Intense classes Offer
      [ ] Podcast 
   [ ] Content
      [ ] Coupled Classes
      [ ] Solo Classes
      [ ] Old Clips
- [ ] Functionality
    [x] Forms submission -> Writing to the database -> Front-end confirmation message & redirect -> Receiving welcome confirmation email
      [x] Contact form
      [x] Registration form
      [x] Newsletter subscription
      [x] Implement form validation
      [ ] Email validation and sending logic
      [ ] Complete email notification system
      [ ] front-end logic
        [ ]form validation
        [ ]form submission
        [ ]form confirmation
        [ ]form error handling
        [ ]form success message
        [ ]form reset
        [ ]form submission button
        [ ]form submission button loading state
        [ ]form submission button disabled state
        [ ]form submission button success state
        [ ]form submission button error state
        [ ]form submission button disabled state
- [x] Set up Firebase Realtime Database
    [ ] email validation and email confirmation logic
    [x] Configure CORS settings
    [ ] Add proper error handling
    [ ] Calendar
  
### Medium Priority
- [ ] Add event scheduling system
- [ ] Add sub-pages & route them
    [ ] Prices
    [ ] Terms of Use
    [ ] GDPR
    [ ] Privacy policy
    [ ] Events
    [ ] Projects
- [ ] Add analytics tracking
- [ ] CMS or Admin dashboard
  - [ ] alternative to CMS using google drive
- [ ] Implement user authentication
- [ ] Social calendar integration
  - [ ] Google Calendar & iCalendar integration
  - [ ] Facebook events integration
  - [ ] Layout with today, upcoming & past events




### Low Priority
- [ ] Implement logging system
- [ ] Implement caching
- [ ] Add performance monitoring
- [ ] Create backup system
- [x] Implement rate limiting
- [ ] Create automated tests for functionality

### Testing phase:
- [ ] Make an alpha test of the website in the Swing Troupe group 
- [ ] Make a beta test of the website inside the Swing Society Classes group 
  - [ ] ensure testers with various devices: 
    - [ ] desktop: 24 inch, 15 inch, 14 inch, 13 inch
    - [ ] tablet: 10 inch, 8 inch
    - [ ] mobile: 6 inch, 5 inch 4 inch 

Simple metrics logging (without Google Cloud integration)
Implement proper Google Cloud monitoring
* Automated test suit 4 all the security features:
  * Rate limiting
  * Security headers
  * Static file serving
  * Template routing
  * Metrics logging


Form submission 30 - 90 min
-> Writing to the database
-> Receiving welcome confirmation email
-> Front-end message 4 form submission


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


