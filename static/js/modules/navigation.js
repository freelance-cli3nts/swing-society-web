export function initNavigation() {
   // Elements
   const mobileMenuToggle = document.getElementById('mobile-menu-toggle');
   const closeMenuButton = document.getElementById('close-menu');
   const mobileMenu = document.getElementById('mobile-menu');
   const submenuToggles = document.querySelectorAll('.menu-item.has-submenu');
   
   // Toggle mobile menu
   mobileMenuToggle.addEventListener('click', () => {
     mobileMenu.classList.toggle('is-active');
     closeMenuButton.classList.toggle('is-active');
     mobileMenuToggle.classList.toggle('is-active');
   });
   
   // Close mobile menu
   closeMenuButton.addEventListener('click', () => {
     mobileMenu.classList.remove('is-active');
     closeMenuButton.classList.remove('is-active');
     mobileMenuToggle.classList.remove('is-active');
   });
   
   // Toggle submenus
   submenuToggles.forEach(toggle => {
     toggle.addEventListener('click', () => {
       const targetId = toggle.dataset.target;
       const targetSubmenu = document.getElementById(targetId);
       
       // Close other submenus
       document.querySelectorAll('.submenu.is-active').forEach(submenu => {
         if (submenu.id !== targetId) {
           submenu.classList.remove('is-active');
           
           // Remove active class from other toggles
           const otherToggle = document.querySelector(`[data-target="${submenu.id}"]`);
           if (otherToggle) {
             otherToggle.classList.remove('is-active');
           }
         }
       });
       
       // Toggle current submenu
       toggle.classList.toggle('is-active');
       targetSubmenu.classList.toggle('is-active');
     });
   });

  // Close menu when clicking a submenu item (navigation link)
  document.querySelectorAll('.submenu-item').forEach(link => {
    link.addEventListener('click', () => {
      // Close the mobile menu after a short delay to allow HTMX to process
      setTimeout(() => {
        mobileMenu.classList.remove('is-active');
        closeMenuButton.classList.remove('is-active');
        mobileMenuToggle.classList.remove('is-active');
        
        // Also close any open submenus
        document.querySelectorAll('.submenu.is-active').forEach(submenu => {
          submenu.classList.remove('is-active');
        });
        document.querySelectorAll('.menu-item.is-active').forEach(item => {
          item.classList.remove('is-active');
        });
      }, 100);
    });
  });

  // Handle window resize
  window.addEventListener('resize', () => {
    if (window.innerWidth >= 1024) {
      // Hide mobile menu on desktop view
      if (mobileMenu) {
        mobileMenu.classList.remove('is-active');
        if (closeMenuButton) closeMenuButton.classList.remove('is-active');
        if (mobileMenuToggle) mobileMenuToggle.classList.remove('is-active');
      }
    }
  });
}