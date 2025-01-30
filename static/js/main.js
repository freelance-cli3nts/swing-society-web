// logo link to about us section on main page, and to main page on other pages
document.addEventListener("DOMContentLoaded", () => {
  const logoLink = document.getElementById("logo-icon");
  const currentPage = window.location.pathname;

  // Check if we're on the main page
  if (currentPage.includes("index.html") || currentPage === "/") {
    // If on the main page, set link to scroll to About Us section
    logoLink.setAttribute("href", "#about");
  } else {
    // On other pages, set link to go back to the main page
    logoLink.setAttribute("href", "index.html");
  }
});

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

