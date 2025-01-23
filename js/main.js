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
