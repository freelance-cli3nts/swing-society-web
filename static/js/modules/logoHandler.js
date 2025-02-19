export function initLogoHandling() {
  const handleLogoLogic = () => {
      const logoLink = document.getElementById("logo-link");
      const currentPath = window.location.pathname;

      if (currentPath === "/" || currentPath === "") {
        // Add click event listener for main page          
        logoLink.addEventListener("click", (e) => {
              e.preventDefault();
              htmx.ajax('GET', '/templates/about.html', {
                  target: '#main-view',
                  swap: 'innerHTML'
              }).catch(error => {
                  console.error("Error loading about content:", error);
              });
          });
      } else {
        // Reset to home page navigation
        logoLink.href = "/";
        // Remove any existing click listeners          
          logoLink.replaceWith(logoLink.cloneNode(true));
      }
  };

}

document.addEventListener("DOMContentLoaded", handleLogoLogic);
document.body.addEventListener("htmx:afterSwap", handleLogoLogic);