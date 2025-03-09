export function initModalHandling() {
  document.body.addEventListener('click', (event) => {
      const button = event.target.closest("[data-form]");
      // add a variable to get the form type
      if (button) {
          const formUrl = button.getAttribute("data-form");
          loadFormIntoModal(formUrl);
      }
  });

  document.addEventListener("keydown", function(event) {
      if (event.key === "Escape") {
          closeModal();
      }
  });

  document.body.addEventListener('click', (event) => {
      if (event.target.matches("#close-large-modal", "#close-large-modal-btn","#close-small-modal", "#close-small-modal-btn", ".modal-background")) {
          closeModal();
      }
  });

  document.body.addEventListener('click', (event) => {
      if (event.target.id === "close-large-modal" || 
        event. target.closest === "close-large-modal-btn") {
          closeModal();
          event.stopPropagation();
      }
  });

  // add the form type to the parameters of this function 
  function loadFormIntoModal(formUrl) {
    // First, make sure the modal exists or create it if needed
    let modal = document.getElementById("large-modal");
    if (!modal) {
      // If modal doesn't exist yet, fetch it first
      htmx.ajax('GET', '/templates/components/modal-large.html', {
        target: "body",
        swap: "beforeend",
        handler: function() {
          // After loading the modal structure, fetch the form
          modal = document.getElementById("large-modal");
          setTimeout(() => {
            modal.classList.add("is-active");
          }, 100);

          // Ensure the close button works after the modal is added
          document.getElementById('close-large-modal-btn').addEventListener('click', function() {
            modal.classList.remove('is-active');
          });

          htmx.ajax('GET', formUrl, {
            target: "#large-modal-content",
            swap: "innerHTML"
          });
        }
      });
    } else {
      // If modal already exists, just show it and load the form
      setTimeout(() => {
        modal.classList.add("is-active");
      }, 100);

      // Ensure the close button works after the modal is added
      document.getElementById('close-large-modal-btn').addEventListener('click', function() {
        modal.classList.remove('is-active');
      });
      htmx.ajax('GET', formUrl, {
        target: "#large-modal-content",
        swap: "innerHTML"
      });
    }
    
    // Add event listeners for closing
    document.querySelectorAll('#close-large-modal, #close-large-modal-btn').forEach(element => {
      element.addEventListener('click', closeModal);
    });
    
    // Prevent body scrolling when modal is active
    document.body.style.overflow = "hidden";
  }
   
  function closeModal() {
      const modal = document.getElementById("large-modal");
      if (modal) {
          modal.classList.remove("is-active");
          setTimeout(() => modal.remove(), 300);
          document.body.style.overflow = "";
      }
  }
}
