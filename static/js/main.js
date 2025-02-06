document.addEventListener("DOMContentLoaded", handleLogoLogic);
document.body.addEventListener("htmx:afterSwap", handleLogoLogic);

function handleLogoLogic() {
    const logoLink = document.getElementById("logo-link");
    console.log("I got the logo link:", logoLink);
    const currentPath = window.location.pathname;
    console.log("Current path:", currentPath);

    if (currentPath === "/" || currentPath === "") {
        console.log("I am on the main page");
        // Add click event listener for main page
        logoLink.addEventListener("click", (e) => {
            e.preventDefault();
            console.log("Logo clicked on main page");
            htmx.ajax('GET', '/templates/about.html', {
                target: '#main-view',
                swap: 'innerHTML'
            }).then(() => {
                console.log("About content loaded");
            }).catch(error => {
                console.error("Error loading about content:", error);
            });
        });
    } else {
        console.log("I am on another page");
        // Reset to home page navigation
        logoLink.href = "/";
        // Remove any existing click listeners
        logoLink.replaceWith(logoLink.cloneNode(true));
    }
}

// Hamburger menu for mobile
document.addEventListener('DOMContentLoaded', () => {
  const hamburger = document.getElementById('hamburger');
  const navList = document.getElementById('nav-list');
  const menuItems = document.querySelectorAll('.has-submenu');

  // Toggle hamburger menu
  hamburger.addEventListener('click', () => {
    navList.classList.toggle('open');
    hamburger.textContent = navList.classList.contains('open') ? '✕' : '☰';
  });

  // Handle submenu clicks on mobile
  menuItems.forEach(item => {
    const link = item.querySelector('a');
    link.addEventListener('click', (e) => {
      if (window.innerWidth <= 768) {
        e.preventDefault();
        const subMenu = item.querySelector('.submenu');
        subMenu.classList.toggle('open');
      }
    });
  });

  // Close menu when clicking outside
  document.addEventListener('click', (e) => {
    if (!navList.contains(e.target) && !hamburger.contains(e.target)) {
      navList.classList.remove('open');
      hamburger.textContent = '☰';
      menuItems.forEach(item => item.classList.remove('active'));
    }
  });

  // Reset menu state on window resize
  window.addEventListener('resize', () => {
    if (window.innerWidth > 768) {
      navList.classList.remove('open');
      hamburger.textContent = '☰';
      menuItems.forEach(item => item.classList.remove('active'));
    }
  });
});

// Wait for both DOM and HTMX to be ready
document.addEventListener('DOMContentLoaded', function() {
  // HTMX configuration
  const configureHtmx = function() {
        console.log('Checking for HTMX...'); // Debug log
        if (window.htmx) {
        console.log('HTMX found!'); // Debug log
        htmx.config.ignoreFeatureWarnings = true;
          
          const path = window.location.pathname;
          console.log('Current path:', path); // Debug log
            
          // Get current path and handle pretty URLs
          const urlMapping = {
              '/classes/solo-jazz': '/templates/classes/solojazz.html',
              '/classes/lindy-hop': '/templates/classes/lindyhop.html',
              '/classes/old-clips': '/templates/classes/oldclips.html',
              '/events/parties': '/templates/events/parties.html',
              '/events/festivals': '/templates/events/festivals.html',
              '/events/workshops': '/templates/events/workshops.html'
          };
          console.log('URL found in mapping:', path in urlMapping); // Debug log

          if (path in urlMapping) {
              console.log('Loading template:', urlMapping[path]); // Debug log

              htmx.ajax('GET', urlMapping[path], {target: '#main-view'});
          }
      } else {
          // If HTMX isn't loaded yet, wait a bit and try again
          setTimeout(configureHtmx, 50);
      }
  };

  configureHtmx();
});


// Form validation
document.addEventListener('htmx:beforeRequest', function(evt) {
  const form = evt.detail.elt;
  if (form.tagName === 'FORM') {
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
  }
});

// Helper functions
function showError(input, message) {
  const errorDiv = document.createElement('div');
  errorDiv.className = 'validation-message';
  errorDiv.textContent = message;
  input.parentNode.appendChild(errorDiv);
  input.classList.add('error');
}

function isValidEmail(email) {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);
}

// Real-time validation
document.addEventListener('htmx:afterSettle', function() {
  const forms = document.querySelectorAll('form');
  forms.forEach(form => {
      const inputs = form.querySelectorAll('input');
      inputs.forEach(input => {
          input.addEventListener('input', function() {
              // Remove error styling when user starts typing
              this.classList.remove('error');
              const errorMessage = this.parentNode.querySelector('.validation-message');
              if (errorMessage) {
                  errorMessage.remove();
              }
          });
      });
  });
});


// Registration form show/hide partner name field
document.addEventListener('DOMContentLoaded', function() {
  const aloneRadios = document.querySelectorAll('input[name="registerAlone"]');
  const partnerNameGroup = document.getElementById('partner-name-group');

  aloneRadios.forEach(radio => {
      radio.addEventListener('change', function() {
          if (this.value === 'no') {
              partnerNameGroup.style.display = 'block';
              setTimeout(() => partnerNameGroup.classList.add('visible'), 10);
          } else {
              partnerNameGroup.classList.remove('visible');
              setTimeout(() => partnerNameGroup.style.display = 'none', 300);
          }
      });
  });
});


// // Modal window functionality
// document.addEventListener('DOMContentLoaded', () => {
//   const modal = document.getElementById('registrationModal');
//   const openBtn = document.getElementById('openModal');
//   const closeBtn = modal.querySelector('.close-btn');
  
//   function openModal() {
//       modal.classList.add('active');
//       document.body.style.overflow = 'hidden'; // Prevent background scrolling
//   }
  
//   function closeModal() {
//       modal.classList.remove('active');
//       document.body.style.overflow = ''; // Restore scrolling
//   }
  
//   openBtn.addEventListener('click', openModal);
//   closeBtn.addEventListener('click', closeModal);
  
//   // Close on outside click
//   modal.addEventListener('click', (e) => {
//       if (e.target === modal) {
//           closeModal();
//       }
//   });
  
//   // Close on escape key
//   document.addEventListener('keydown', (e) => {
//       if (e.key === 'Escape' && modal.classList.contains('active')) {
//           closeModal();
//       }
//   });
// });