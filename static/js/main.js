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


// Registration form 
// document.addEventListener('DOMContentLoaded', function() {
//   const form = document.getElementById('contact-form');
//   const emailInput = document.getElementById('email');
//   const phoneInput = document.getElementById('phone');
//   const validationMessage = document.getElementById('contact-validation');

//   form.addEventListener('submit', function(e) {
//     // Check if at least one contact method is provided
//     if (!emailInput.value && !phoneInput.value) {
//       e.preventDefault();
//       validationMessage.style.display = 'block';
//       return false;
//     }
//     validationMessage.style.display = 'none';
//   });

//   // Hide validation message when user starts typing in either field
//   [emailInput, phoneInput].forEach(input => {
//     input.addEventListener('input', function() {
//       if (emailInput.value || phoneInput.value) {
//         validationMessage.style.display = 'none';
//       }
//     });
//   });
// });

