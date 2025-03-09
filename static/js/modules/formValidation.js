/**
 * Swing Dance Registration Form Validation
 * This module handles form validation and interactions for the Swing dance registration form.
 */

// Main validation module
export function initFormValidation() {
  // DOM elements cache
  let form, emailInput, phoneInput, emailError, phoneError;

  /**
   * Initialize form validation
   * @param {Object} config - Configuration object (optional)
   */
  function init(config = {}) {
    // Get form elements
    form = document.getElementById('registration-form');
    emailInput = document.getElementById('email');
    phoneInput = document.getElementById('phone');
    emailError = document.getElementById('email-error');
    phoneError = document.getElementById('phone-error');

    if (!form) {
      console.error('Registration form not found');
      return;
    }

    // Bind event listeners
    bindEvents();
  }

  /**
   * Bind all event listeners
   */
  function bindEvents() {
    // Live email validation
    emailInput?.addEventListener('blur', validateEmailField);
    
    // Live phone validation
    phoneInput?.addEventListener('blur', validatePhoneField);
    
    // Form submission validation
    form.addEventListener('submit', handleFormSubmit);
  }

  /**
   * Validate email field on blur
   */
  function validateEmailField() {
    if (emailInput.value && !validateEmail(emailInput.value)) {
      showError(emailInput, emailError, 'Моля, въведете валиден имейл адрес.');
    } else {
      clearError(emailInput, emailError);
    }
  }

  /**
   * Validate phone field on blur
   */
  function validatePhoneField() {
    if (phoneInput.value && !validatePhone(phoneInput.value)) {
      showError(phoneInput, phoneError, 'Моля, въведете валиден телефонен номер (например: 0888123456 или +359888123456).');
    } else {
      clearError(phoneInput, phoneError);
    }
  }

  /**
   * Handle form submission
   * @param {Event} event - Submit event
   */
  function handleFormSubmit(event) {
    let isValid = true;
    
    // Validate email
    if (!validateEmail(emailInput.value)) {
      showError(emailInput, emailError, 'Моля, въведете валиден имейл адрес.');
      isValid = false;
    }
    
    // Validate phone if provided
    if (phoneInput.value && !validatePhone(phoneInput.value)) {
      showError(phoneInput, phoneError, 'Моля, въведете валиден телефонен номер (например: 0888123456 или +359888123456).');
      isValid = false;
    }
    
    // Prevent form submission if validation fails
    if (!isValid) {
      event.preventDefault();
      // Scroll to the first error
      document.querySelector('.input-error')?.scrollIntoView({ 
        behavior: 'smooth', 
        block: 'center' 
      });
    }
  }

  /**
   * Show error message for a field
   * @param {HTMLElement} input - Input element
   * @param {HTMLElement} errorElement - Error message element
   * @param {string} message - Error message
   */
  function showError(input, errorElement, message) {
    input.classList.add('input-error');
    errorElement.textContent = message;
  }

  /**
   * Clear error message for a field
   * @param {HTMLElement} input - Input element
   * @param {HTMLElement} errorElement - Error message element
   */
  function clearError(input, errorElement) {
    input.classList.remove('input-error');
    errorElement.textContent = '';
  }

  /**
   * Toggle partner field visibility
   * @param {boolean} isAlone - Whether registering alone
   */
  function togglePartnerField(isAlone) {
    const partnerField = document.getElementById('partner-name-group');
    if (partnerField) {
      partnerField.style.display = isAlone ? 'none' : 'block';
    }
  }

  /**
   * Validate email format
   * @param {string} email - Email to validate
   * @returns {boolean} - Whether email is valid
   */
  function validateEmail(email) {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
  }
  
  /**
   * Validate phone format (Bulgarian)
   * @param {string} phone - Phone to validate
   * @returns {boolean} - Whether phone is valid
   */
  function validatePhone(phone) {
    // Skip validation if phone is empty (since it's optional)
    if (!phone) return true;
    
    // Bulgarian phone format: +359xxxxxxxxx or 0xxxxxxxxx
    const phoneRegex = /^(\+359|0)[0-9]{9}$/;
    return phoneRegex.test(phone.replace(/\s/g, ''));
  }

  // Public API
  return {
    init,
    togglePartnerField,
    validateEmail,
    validatePhone
  };
};

