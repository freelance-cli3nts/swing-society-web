// Modal window functionality
document.addEventListener('DOMContentLoaded', () => {
  const modal = document.getElementById('registrationModal');
  const openBtn = document.getElementById('openModal');
  const closeBtn = modal.querySelector('.close-btn');
  
  function openModal() {
      modal.classList.add('active');
      document.body.style.overflow = 'hidden'; // Prevent background scrolling
  }
  
  function closeModal() {
      modal.classList.remove('active');
      document.body.style.overflow = ''; // Restore scrolling
  }
  
  openBtn.addEventListener('click', openModal);
  closeBtn.addEventListener('click', closeModal);
  
  // Close on outside click
  modal.addEventListener('click', (e) => {
      if (e.target === modal) {
          closeModal();
      }
  });
  
  // Close on escape key
  document.addEventListener('keydown', (e) => {
      if (e.key === 'Escape' && modal.classList.contains('active')) {
          closeModal();
      }
  });
});