export function initNavigation() {
  const mobileNavToggle = document.querySelector('.mobile-nav-toggle');
  const navList = document.querySelector('.nav-list');
  const dropdownToggles = document.querySelectorAll('.dropdown-toggle');

  if (!mobileNavToggle || !navList) return;

  // Toggle mobile menu
  mobileNavToggle.addEventListener('click', () => {
    const isVisible = navList.getAttribute('data-visible') === 'true';
    navList.setAttribute('data-visible', !isVisible);
    mobileNavToggle.setAttribute('aria-expanded', !isVisible);
    mobileNavToggle.querySelector('.hamburger').textContent = isVisible ? '☰' : '✕';
  });

  // Handle dropdown toggles
  dropdownToggles.forEach(toggle => {
    toggle.addEventListener('click', (e) => {
      if (window.innerWidth <= 768) {
        e.preventDefault();
        const parent = toggle.closest('.has-submenu');
        
        // Close other dropdowns
        document.querySelectorAll('.has-submenu').forEach(item => {
          if (item !== parent && item.classList.contains('active')) {
            item.classList.remove('active');
          }
        });
        
        // Toggle current dropdown
        parent.classList.toggle('active');
      }
    });
  });

  // Handle arrow clicks separately if needed
  document.querySelectorAll('.dropdown-icon').forEach(arrow => {
    arrow.addEventListener('click', (e) => {
      if (window.innerWidth <= 768) {
        e.preventDefault();
        e.stopPropagation(); // Prevent triggering the button click
        const parent = arrow.closest('.has-submenu');
        parent.classList.toggle('active');
      }
    });
  });

  // Close menu when clicking outside
  document.addEventListener('click', (e) => {
    if (!navList.contains(e.target) && !mobileNavToggle.contains(e.target)) {
      navList.setAttribute('data-visible', 'false');
      mobileNavToggle.setAttribute('aria-expanded', 'false');
      mobileNavToggle.querySelector('.hamburger').textContent = '☰';
      
      // Close all dropdowns
      document.querySelectorAll('.has-submenu').forEach(item => {
        item.classList.remove('active');
      });
    }
  });
}

document.addEventListener('DOMContentLoaded', initNavigation);
document.body.addEventListener('htmx:afterSwap', initNavigation);
