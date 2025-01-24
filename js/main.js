document.addEventListener("DOMContentLoaded", () => {
  const logoLink = document.getElementById("logo-link");
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