export function initFormHandling() {
  const validateForm = (form) => {    
    // Reset previous validation messages
      const errorMessages = form.querySelectorAll('.validation-message');
      errorMessages.forEach(msg => msg.remove());

      // Get form fields
      const name = form.querySelector('input[name="name"]');
      const email = form.querySelector('input[name="email"]');
      let isValid = true;

      // Name validation
      if (name && name.value.trim() === '') {
          showError(name, 'Моля, въведете вашето име');
          isValid = false;
      } else if (name && name.value.length < 2) {
          showError(name, 'Името трябва да е поне 2 символа');
          isValid = false;
      }

      // Email validation
      if (email && email.value.trim() === '') {
          showError(email, 'Моля, въведете вашия имейл');
          isValid = false;
      } else if (email && !isValidEmail(email.value)) {
          showError(email, 'Моля, въведете валиден имейл адрес');
          isValid = false;
      }

      if (!isValid) {
          evt.preventDefault(); // Stop HTMX from sending the request
      }
  };

  // Helper functions
  const showError = (input, message) => {
    const errorDiv = document.createElement('div');
    errorDiv.className = 'validation-message';
    errorDiv.textContent = message;
    input.parentNode.appendChild(errorDiv);
    input.classList.add('error');
  };

  const isValidEmail = (email) => {
      return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);
  };

  // Real-time validation
  const setupRealTimeValidation = () => {
      const forms = document.querySelectorAll('form');
      forms.forEach(form => {
          const inputs = form.querySelectorAll('input');
          inputs.forEach(input => {
              input.addEventListener('input', function() {
                  this.classList.remove('error');
                  const errorMessage = this.parentNode.querySelector('.validation-message');
                  if (errorMessage) {
                      errorMessage.remove();
                  }
              });
          });
      });
  };

  document.addEventListener('htmx:beforeRequest', (evt) => {
      const form = evt.detail.elt;
      if (form.tagName === 'FORM' && !validateForm(form)) {
          evt.preventDefault();
      }
  });

  document.addEventListener('htmx:afterSettle', setupRealTimeValidation);
};
