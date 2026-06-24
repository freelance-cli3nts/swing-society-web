# DeepSeek - Setting Up Firebase Realtime Database Integration - DeepSeek

## Setting Up Firebase Realtime Database Integration - DeepSeek

Date exported: 3/28/2026, 12:51:01 AM

**You said:**

setting up google realtime database in a website build with HTML, CSS, JS, HTMX, and Go(back-end) that is run on Google cloud run.

I want a detailed step by step instructions for setup

**DeepSeek said:**

Setting up Google Firebase Realtime Database in a website built with HTML, CSS, JS, HTMX, and Go (back-end) that is hosted on Google Cloud Run involves several steps. Below is a detailed step-by-step guide to help you set it up:

Step 1: Set Up Firebase Project

Create a Firebase Project:Go to the Firebase Console.Click on "Add Project".Enter a project name and follow the prompts to create the project.Enable Realtime Database:In the Firebase Console, navigate to Build > Realtime Database.Click "Create Database".Choose a location for your database and select "Start in test mode" (for development purposes). You can change the rules later for production.Get Firebase Configuration:In the Firebase Console, go to Project Settings (gear icon).Scroll down to the "Your apps" section and click "Add app".Choose Web (</>).Register the app and note down the Firebase configuration object (API key, auth domain, database URL, etc.). You’ll need this for the front-end.

Step 2: Set Up Google Cloud Run

Deploy Your Go Backend to Cloud Run:Ensure your Go backend is containerized using Docker.Push your Docker image to Google Container Registry (GCR).Deploy the container to Cloud Run:bashgcloud run deploy YOUR_SERVICE_NAME --image gcr.io/YOUR_PROJECT_ID/YOUR_IMAGE_NAME --platform managedNote the Cloud Run service URL.Set Up Firebase Admin SDK in Go:Install the Firebase Admin SDK for Go:bashgo get firebase.google.com/goInitialize the Firebase Admin SDK in your Go backend:goimport (
    "context"
    "log"
    "firebase.google.com/go"
    "google.golang.org/api/option"
)

func main() {
    opt := option.WithCredentialsFile("path/to/serviceAccountKey.json")
    app, err := firebase.NewApp(context.Background(), nil, opt)
    if err != nil {
        log.Fatalf("error initializing app: %v\n", err)
    }

    client, err := app.Database(context.Background())
    if err != nil {
        log.Fatalf("error initializing database client: %v\n", err)
    }

    // Use the client to interact with the Realtime Database
}Download the Firebase service account key from the Firebase Console (Project Settings > Service Accounts) and place it in your Go project.

Step 3: Integrate Firebase Realtime Database with Front-End

Add Firebase SDK to Your HTML:Include the Firebase SDK in your HTML file:html<script src="https://www.gstatic.com/firebasejs/10.0.0/firebase-app.js"></script>
<script src="https://www.gstatic.com/firebasejs/10.0.0/firebase-database.js"></script>Initialize Firebase in Your JS:Use the Firebase configuration object from Step 1 to initialize Firebase:javascriptconst firebaseConfig = {
    apiKey: "YOUR_API_KEY",
    authDomain: "YOUR_AUTH_DOMAIN",
    databaseURL: "YOUR_DATABASE_URL",
    projectId: "YOUR_PROJECT_ID",
    storageBucket: "YOUR_STORAGE_BUCKET",
    messagingSenderId: "YOUR_MESSAGING_SENDER_ID",
    appId: "YOUR_APP_ID"
};

const app = firebase.initializeApp(firebaseConfig);
const database = firebase.database();Interact with the Realtime Database:Example: Write data to the database:javascriptfunction writeData() {
    database.ref('users/').set({
        username: "exampleUser",
        email: "user@example.com"
    });
}Example: Read data from the database:javascriptfunction readData() {
    database.ref('users/').once('value').then((snapshot) => {
        const data = snapshot.val();
        console.log(data);
    });
}

Step 4: Connect Front-End with Back-End

HTMX for Dynamic Updates:Use HTMX to send requests to your Go backend. For example:html<button hx-post="/update-data" hx-trigger="click">
    Update Data
</button>In your Go backend, handle the request and update the Firebase Realtime Database:gohttp.HandleFunc("/update-data", func(w http.ResponseWriter, r *http.Request) {
    ref := client.NewRef("users/")
    err := ref.Set(context.Background(), map[string]interface{}{
        "username": "newUser",
        "email":    "newuser@example.com",
    })
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
})CORS Configuration:Ensure your Go backend allows requests from your front-end domain by setting appropriate CORS headers:gofunc enableCORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        next.ServeHTTP(w, r)
    })
}

Step 5: Deploy and Test

Deploy Your Front-End:Host your HTML, CSS, and JS files on a static hosting service (e.g., Firebase Hosting, Google Cloud Storage, or any CDN).Test the Integration:Verify that data is being written to and read from the Firebase Realtime Database.Ensure your Go backend is correctly updating the database via HTMX requests.

Step 6: Secure Your Database

Update Firebase Realtime Database Rules:Go to the Firebase Console and navigate to Realtime Database > Rules.Update the rules to restrict access as needed. For example:json{
  "rules": {
    ".read": "auth != null",
    ".write": "auth != null"
  }
}Implement Authentication (Optional):Use Firebase Authentication to secure access to your database.

By following these steps, you should have a fully functional website with Firebase Realtime Database integrated into your HTML, CSS, JS, HTMX, and Go stack, hosted on Google Cloud Run. Let me know if you need further clarification!

**You said:**

I am a bit confused why do I have to set it up on both the front and back-end?

**DeepSeek said:**

Great question! The confusion is understandable, and it’s important to clarify why Firebase Realtime Database might need to be set up on both the front-end and back-end, depending on your use case. Let me break it down for you:

Why Set Up Firebase on the Front-End?

Direct Client-Side Interaction:Firebase Realtime Database is designed to allow direct client-side interaction. This means your front-end (HTML/JS) can read from and write to the database directly without needing to go through your back-end.This is useful for real-time updates, such as chat applications, live notifications, or collaborative tools, where the front-end needs to listen for changes in the database instantly.Real-Time Updates:Firebase SDK provides real-time listeners (e.g., on('value')) that allow your front-end to react immediately when data changes in the database. This is a key feature of Firebase Realtime Database.Reduced Back-End Load:By allowing the front-end to interact directly with Firebase, you reduce the load on your back-end server (Go in this case), as it doesn’t need to handle every database read/write operation.

Why Set Up Firebase on the Back-End?

Secure Operations:While Firebase allows direct client-side interaction, it’s often unsafe to let the front-end perform sensitive operations (e.g., updating critical data or deleting records). This is because Firebase security rules alone might not be sufficient to enforce complex business logic.By using the Firebase Admin SDK on your back-end, you can securely perform these operations without exposing sensitive logic to the client.Server-Side Logic:Your back-end (Go) might need to perform operations that require server-side logic, such as:Validating data before saving it to the database.Integrating with other services (e.g., sending emails, processing payments).Performing batch updates or complex queries.Authentication and Authorization:If you’re using Firebase Authentication, your back-end might need to verify user tokens or enforce custom authorization rules before allowing certain operations.Offloading Work from the Client:For computationally expensive tasks or tasks that require access to sensitive credentials (e.g., API keys), it’s better to handle them on the back-end.

When to Use Front-End vs. Back-End Firebase Setup

| Use Case                            | Front-End | Back-End |
| ----------------------------------- | --------- | -------- |
| Real-time updates (e.g., live chat) | ✅         | ❌        |
| Simple read/write operations        | ✅         | ❌        |
| Secure or sensitive operations      | ❌         | ✅        |
| Complex business logic              | ❌         | ✅        |
| Integration with other services     | ❌         | ✅        |

Do You Need Both?

It depends on your application’s requirements:

If your app only needs real-time updates and simple read/write operations, you can use Firebase entirely on the front-end and skip the back-end setup.If your app requires secure operations, server-side logic, or integration with other services, you’ll need to set up Firebase on both the front-end and back-end.

Simplified Setup for Your Use Case

If you’re just starting out and want to keep things simple:

Use Firebase only on the front-end for real-time updates and basic read/write operations.Use your Go back-end for handling sensitive operations or integrating with other services.

For example:

Front-End: Use Firebase SDK to read/write data and listen for real-time updates.Back-End: Use Firebase Admin SDK to handle secure operations or complex logic.

Example Workflow

Front-End:User interacts with the website (e.g., submits a form).Firebase SDK writes data directly to the Realtime Database.Firebase SDK listens for real-time updates and reflects them on the UI.Back-End:User triggers an action that requires server-side logic (e.g., payment processing).HTMX sends a request to your Go back-end.Go back-end uses Firebase Admin SDK to securely update the database or perform other tasks.

Conclusion

You don’t necessarily have to set up Firebase on both the front-end and back-end. It depends on your application’s needs:

Use the front-end Firebase SDK for real-time updates and simple operations.Use the back-end Firebase Admin SDK for secure or complex operations.

If you’re still unsure, start with the front-end setup and add the back-end integration only when needed. Let me know if you’d like further clarification or help with a specific part of the setup!

**You said:**

Ok, for now it's going to be only user submited forms. But for some reason I thought because they are sent to the server -> database should be also configured and set-up there

**DeepSeek said:**

Ah, I see where the confusion is coming from! You're absolutely correct that traditionally, form submissions are sent to a server, which then processes the data and stores it in a database. However, Firebase Realtime Database is a bit different because it allows direct client-to-database communication, bypassing the need for a traditional server in some cases.

Let’s clarify this and provide a clear path forward for your use case (user-submitted forms).

Traditional Workflow (Server-Centric)

In a traditional setup:

The user submits a form on the front-end.The form data is sent to a server (e.g., your Go back-end).The server processes the data (e.g., validates it, applies business logic).The server saves the data to a database (e.g., MySQL, PostgreSQL, etc.).

In this setup, the database is only accessible from the server, and the front-end has no direct access to it.

Firebase Realtime Database Workflow (Client-Centric)

With Firebase Realtime Database:

The user submits a form on the front-end.The form data is sent directly to Firebase Realtime Database using the Firebase SDK.Firebase handles saving the data and can notify other clients in real-time if needed.

In this setup, the front-end interacts directly with the database, and the server (Go back-end) is not involved unless you need it for additional processing.

Why You Might Still Want to Use the Server

Even though Firebase allows direct client-to-database communication, there are good reasons to involve your Go back-end for form submissions:

Data Validation:You might want to validate the form data on the server to ensure it meets certain criteria before saving it to the database.Security:If you allow the front-end to write directly to the database, you need to rely heavily on Firebase Security Rules to prevent malicious or invalid data from being written.By sending data to your Go back-end first, you can enforce stricter validation and security checks.Business Logic:If your form submissions require additional processing (e.g., sending an email, integrating with a third-party API), you’ll need to handle this on the server.Data Transformation:You might need to transform or enrich the data before saving it to the database (e.g., adding timestamps, generating IDs, etc.).

How to Set Up Form Submissions with Firebase and Go

Here’s how you can set up your form submissions to use both Firebase and your Go back-end:

Option 1: Front-End Directly Writes to Firebase

Front-End:Use the Firebase SDK to write form data directly to the Realtime Database.Example:javascriptfunction submitForm() {
    const formData = {
        name: document.getElementById('name').value,
        email: document.getElementById('email').value,
        message: document.getElementById('message').value
    };

    const database = firebase.database();
    database.ref('formSubmissions/').push(formData)
        .then(() => {
            alert('Form submitted successfully!');
        })
        .catch((error) => {
            console.error('Error submitting form: ', error);
        });
}Firebase Security Rules:Configure Firebase Security Rules to restrict who can read/write to the database.Example:json{
  "rules": {
    "formSubmissions": {
      ".read": "auth != null",  // Only authenticated users can read
      ".write": "auth != null" // Only authenticated users can write
    }
  }
}

Option 2: Front-End Sends Data to Go Back-End, Which Writes to Firebase

Front-End:Use HTMX or Fetch API to send form data to your Go back-end.Example with HTMX:html<form hx-post="/submit-form" hx-trigger="submit">
    <input type="text" id="name" name="name" placeholder="Name" required>
    <input type="email" id="email" name="email" placeholder="Email" required>
    <textarea id="message" name="message" placeholder="Message" required></textarea>
    <button type="submit">Submit</button>
</form>Go Back-End:Handle the form submission, validate the data, and write it to Firebase Realtime Database using the Firebase Admin SDK.Example:gohttp.HandleFunc("/submit-form", func(w http.ResponseWriter, r *http.Request) {
    // Parse form data
    name := r.FormValue("name")
    email := r.FormValue("email")
    message := r.FormValue("message")

    // Validate data (example)
    if name == "" || email == "" || message == "" {
        http.Error(w, "All fields are required", http.StatusBadRequest)
        return
    }

    // Save data to Firebase
    ref := client.NewRef("formSubmissions/")
    err := ref.Push(context.Background(), map[string]interface{}{
        "name":    name,
        "email":   email,
        "message": message,
    })
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Form submitted successfully!"))
})Firebase Security Rules:Since the Go back-end uses the Firebase Admin SDK, it bypasses Firebase Security Rules. You can restrict access to the database to only your back-end.

Which Option Should You Choose?

Option 1 (Front-End Directly Writes to Firebase):Use this if your form submissions are simple and don’t require server-side validation or processing.Ensure you configure Firebase Security Rules properly to prevent abuse.Option 2 (Front-End Sends Data to Go Back-End):Use this if you need server-side validation, security, or additional processing.This is the more traditional and secure approach.

Recommendation for Your Use Case

Since you’re dealing with user-submitted forms, I recommend Option 2 (sending data to your Go back-end first). This gives you more control over the data and ensures it’s validated and processed securely before being saved to Firebase.

Let me know if you’d like help implementing either option!

**You said:**

I would perfer to go with option 2. Can we restart the setup process with this in mind

**DeepSeek said:**

Absolutely! Let’s restart the setup process with Option 2 in mind, where the front-end sends form data to your Go back-end, and the Go back-end writes the data to Firebase Realtime Database. This approach ensures better control, validation, and security for your form submissions.

Step-by-Step Setup for Option 2

Step 1: Set Up Firebase Project

Create a Firebase Project:Go to the Firebase Console.Click on "Add Project" and follow the prompts to create a new project.Enable Realtime Database:In the Firebase Console, navigate to Build > Realtime Database.Click "Create Database".Choose a location and start in test mode (for development). You can tighten security rules later.Get Firebase Admin SDK Credentials:Go to Project Settings > Service Accounts.Click "Generate New Private Key" to download a JSON file containing your service account credentials. This file will be used in your Go back-end to authenticate with Firebase.

Step 2: Set Up Go Back-End

Initialize Your Go Project:Create a new Go project or use an existing one.Initialize a Go module (if not already done):bashgo mod init your-project-nameInstall Firebase Admin SDK for Go:Install the Firebase Admin SDK:bashgo get firebase.google.com/goSet Up Firebase in Your Go Back-End:Create a file (e.g., main.go) and initialize the Firebase Admin SDK:gopackage main

import (
    "context"
    "log"
    "net/http"
    "firebase.google.com/go"
    "google.golang.org/api/option"
)

var client *firebase.Database

func main() {
    // Initialize Firebase Admin SDK
    opt := option.WithCredentialsFile("path/to/serviceAccountKey.json") // Path to your Firebase service account key
    app, err := firebase.NewApp(context.Background(), nil, opt)
    if err != nil {
        log.Fatalf("error initializing Firebase app: %v\n", err)
    }

    // Initialize Realtime Database client
    client, err = app.Database(context.Background())
    if err != nil {
        log.Fatalf("error initializing database client: %v\n", err)
    }

    // Set up HTTP routes
    http.HandleFunc("/submit-form", handleFormSubmission)

    // Start the server
    log.Println("Server started on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleFormSubmission(w http.ResponseWriter, r *http.Request) {
    // Parse form data
    name := r.FormValue("name")
    email := r.FormValue("email")
    message := r.FormValue("message")

    // Validate data (example)
    if name == "" || email == "" || message == "" {
        http.Error(w, "All fields are required", http.StatusBadRequest)
        return
    }

    // Save data to Firebase Realtime Database
    ref := client.NewRef("formSubmissions/")
    err := ref.Push(context.Background(), map[string]interface{}{
        "name":    name,
        "email":   email,
        "message": message,
    })
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Respond to the client
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Form submitted successfully!"))
}Run Your Go Back-End:Start your Go server:bashgo run main.go

Step 3: Set Up Front-End

Create an HTML Form:Create an HTML file (e.g., index.html) with a form:html<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Form Submission</title>
    <script src="https://unpkg.com/htmx.org"></script>
</head>
<body>
    <h1>Submit Form</h1>
    <form hx-post="/submit-form" hx-trigger="submit" hx-target="#response">
        <label for="name">Name:</label>
        <input type="text" id="name" name="name" required><br><br>

        <label for="email">Email:</label>
        <input type="email" id="email" name="email" required><br><br>

        <label for="message">Message:</label>
        <textarea id="message" name="message" required></textarea><br><br>

        <button type="submit">Submit</button>
    </form>

    <div id="response"></div>
</body>
</html>Serve the HTML File:Use your Go back-end to serve the HTML file:gohttp.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "index.html")
})

Step 4: Deploy to Google Cloud Run

Containerize Your Go Back-End:Create a Dockerfile:DockerfileFROM golang:1.20-alpine
WORKDIR /app
COPY . .
RUN go build -o main .
CMD ["./main"]Build and Push Docker Image:Build the Docker image:bashdocker build -t gcr.io/YOUR_PROJECT_ID/YOUR_IMAGE_NAME .Push the image to Google Container Registry (GCR):bashdocker push gcr.io/YOUR_PROJECT_ID/YOUR_IMAGE_NAMEDeploy to Cloud Run:Deploy the container:bashgcloud run deploy YOUR_SERVICE_NAME --image gcr.io/YOUR_PROJECT_ID/YOUR_IMAGE_NAME --platform managed

Step 5: Test the Setup

Submit the Form:Open your website and submit the form.The form data will be sent to your Go back-end, which will save it to Firebase Realtime Database.Verify Data in Firebase:Go to the Firebase Console and check the Realtime Database section to ensure the data is being saved correctly.

Step 6: Secure Your Setup

Update Firebase Security Rules:Restrict access to your database to only your Go back-end:json{
  "rules": {
    "formSubmissions": {
      ".read": "false",
      ".write": "false"
    }
  }
}Since your Go back-end uses the Firebase Admin SDK, it bypasses these rules.Enable Authentication (Optional):If you want to restrict form submissions to authenticated users, integrate Firebase Authentication.

That’s it! You now have a secure and scalable setup for handling user-submitted forms with Firebase Realtime Database and Go. Let me know if you need further assistance!

**You said:**

I have a question about point 3 of step 1. Where do I save this JSON file

![Image](data:image/svg+xml,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20width%3D%2228%22%20height%3D%2228%22%20fill%3D%22none%22%3E%3Cpath%20d%3D%22M0%2014C0%206.267%206.267%200%2014%200c7.734%200%2014%206.267%2014%2014%200%206.766-4.797%2012.41-11.176%2013.716-.911.187-1.856.284-2.82.284H0V14Z%22%2F%3E%3Cpath%20fill%3D%22%23fff%22%20d%3D%22M23.688%2012.674H17.38l-1.486%202.574h5.202c-.815%204.475-5.793%207.523-10.767%204.945a5.857%205.857%200%200%201-2.503-2.48C4.964%2012.267%208.873%206.786%2014.032%206.786c2.127%200%204.042.931%205.358%202.408l1.333-2.31a9.695%209.695%200%200%200-7.38-2.643c-4.866.335-8.811%204.37-9.048%209.24-.27%205.593%204.2%2010.226%209.733%2010.226%205.534%200%209.748-4.373%209.748-9.747%200-.437-.033-.866-.088-1.285Z%22%2F%3E%3C%2Fsvg%3E)

![Image](data:image/svg+xml,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20width%3D%2228%22%20height%3D%2228%22%20fill%3D%22none%22%3E%3Cpath%20fill%3D%22%23027E6F%22%20d%3D%22M0%2014C0%206.267%206.267%200%2014%200c7.734%200%2014%206.267%2014%2014%200%206.766-4.797%2012.41-11.176%2013.716-.911.187-1.856.284-2.82.284H0V14Z%22%2F%3E%3C%2Fsvg%3E)

**DeepSeek said:**

Great question! The JSON file you download from Firebase (the service account key) contains sensitive credentials that allow your Go back-end to authenticate with Firebase. Where and how you store this file is important for both security and functionality. Here's a detailed explanation:

Where to Save the JSON File

Local Development:Save the JSON file in your Go project directory, but do not commit it to version control (e.g., Git).A common practice is to place it in a config or secrets folder:textyour-go-project/
├── main.go
├── go.mod
├── go.sum
└── config/
    └── serviceAccountKey.jsonProduction (Google Cloud Run):Since Cloud Run is a serverless environment, you can't directly place files in the container's filesystem after deployment. Instead, you have two options:Option 1: Include the JSON file in your Docker image (not recommended for production due to security risks).Option 2: Use Google Cloud Secret Manager to securely store and access the JSON file.

How to Securely Handle the JSON File

Option 1: Include the JSON File in Your Docker Image (for Local Development)

Save the JSON file in your project directory (e.g., config/serviceAccountKey.json).Update your Dockerfile to copy the JSON file into the container:DockerfileFROM golang:1.20-alpine
WORKDIR /app
COPY . .
RUN go build -o main .
CMD ["./main"]In your Go code, reference the JSON file path:goopt := option.WithCredentialsFile("config/serviceAccountKey.json")Important: Add the JSON file to your .gitignore to avoid committing it to version control:text# .gitignore
config/serviceAccountKey.json

Option 2: Use Google Cloud Secret Manager (Recommended for Production)

Upload the JSON File to Secret Manager:Go to the Google Cloud Console.Navigate to Security > Secret Manager.Click "Create Secret".Name the secret (e.g., firebase-service-account-key).Upload the contents of your JSON file as the secret value.Grant Access to the Secret:Ensure your Cloud Run service account has permission to access the secret:Go to IAM & Admin > IAM.Find the service account used by Cloud Run (e.g., PROJECT_NUMBER-compute@developer.gserviceaccount.com).Add the role Secret Manager Secret Accessor.Modify Your Go Code to Use Secret Manager:Install the Secret Manager client library for Go:bashgo get cloud.google.com/go/secretmanagerUpdate your Go code to fetch the secret:goimport (
    "context"
    "log"
    "net/http"
    "firebase.google.com/go"
    "google.golang.org/api/option"
    "cloud.google.com/go/secretmanager"
)

func getFirebaseCredentials() ([]byte, error) {
    ctx := context.Background()
    client, err := secretmanager.NewClient(ctx)
    if err != nil {
        return nil, err
    }
    defer client.Close()

    // Replace with your secret name
    name := "projects/YOUR_PROJECT_ID/secrets/firebase-service-account-key/versions/latest"
    result, err := client.AccessSecretVersion(ctx, &secretmanager.AccessSecretVersionRequest{
        Name: name,
    })
    if err != nil {
        return nil, err
    }

    return result.Payload.Data, nil
}

func main() {
    // Fetch Firebase credentials from Secret Manager
    credentials, err := getFirebaseCredentials()
    if err != nil {
        log.Fatalf("error fetching Firebase credentials: %v\n", err)
    }

    // Initialize Firebase Admin SDK
    opt := option.WithCredentialsJSON(credentials)
    app, err := firebase.NewApp(context.Background(), nil, opt)
    if err != nil {
        log.Fatalf("error initializing Firebase app: %v\n", err)
    }

    // Initialize Realtime Database client
    client, err := app.Database(context.Background())
    if err != nil {
        log.Fatalf("error initializing database client: %v\n", err)
    }

    // Set up HTTP routes
    http.HandleFunc("/submit-form", handleFormSubmission)

    // Start the server
    log.Println("Server started on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}Deploy to Cloud Run:When deploying to Cloud Run, ensure the service account has access to the secret.

Best Practices

Never Commit Secrets to Version Control:Always add sensitive files like serviceAccountKey.json to your .gitignore.Use Environment Variables for Local Development:You can also encode the JSON file as a base64 string and store it in an environment variable:bashexport FIREBASE_CREDENTIALS=$(base64 config/serviceAccountKey.json)In your Go code, decode the environment variable:gocredentials, err := base64.StdEncoding.DecodeString(os.Getenv("FIREBASE_CREDENTIALS"))Rotate Secrets Regularly:Periodically update your Firebase service account key and secrets in Secret Manager.

Summary

For local development, save the JSON file in your project directory (e.g., config/serviceAccountKey.json) and reference it in your Go code.For production, use Google Cloud Secret Manager to securely store and access the JSON file.

Let me know if you need further clarification or help with implementing this!


---

[Powered by: chatexport.workpent.com](https://chatexport.workpent.com/)