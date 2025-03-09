// Function to send the email data to the backend for validation
export function initValidation() {
  const email = document.getElementById('email').value;
  console.log("Email being sent: ", email);  // Check if the email is populated

  // Check if the email field is not empty before sending
  if (email.trim() !== "") {
      const data = { email: email };

      // Make the AJAX request to validate the email
      fetch('/api/validate-email', {
          method: 'POST',
          headers: {
              'Content-Type': 'application/json',  // Set content type to JSON
          },
          body: JSON.stringify(data)  // Send the email as a JSON payload
      })
      .then(response => response.json())  // Parse the JSON response
      .then(data => {
          // Handle the success or failure response
          const emailErrorElement = document.getElementById('email-error');
          console.log("Response data: ", data);  // Log the response data
          if (data.message) {
              // Display success message
              emailErrorElement.textContent = "Valid email!";
              emailErrorElement.style.color = 'green';
          } else if (data.error) {
              // Display error message
              emailErrorElement.textContent = data.error.message;
              emailErrorElement.style.color = 'red';
          }
      })
      .catch(error => {
          console.error('Error:', error);
          const emailErrorElement = document.getElementById('email-error');
          emailErrorElement.textContent = "An unexpected error occurred.";
          emailErrorElement.style.color = 'red';
      });
  }
}
