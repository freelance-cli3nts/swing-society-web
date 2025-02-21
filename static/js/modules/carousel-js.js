// carousel-js.js
export function initCarousel() {
  const carousel = document.querySelector('.carousel');
  if (!carousel) return;
  
  const carouselInner = carousel.querySelector('.carousel-inner');
  const items = carousel.querySelectorAll('.carousel-item');
  const prevButton = carousel.querySelector('.carousel-control.prev');
  const nextButton = carousel.querySelector('.carousel-control.next');
  
  if (!carouselInner || !items.length || !prevButton || !nextButton) return;
  
  let currentIndex = 0;
  let autoplayInterval;
  let isVideoPlaying = false;

  function updateCarousel() {
    carouselInner.style.transform = `translateX(-${currentIndex * 100}%)`;
    
    // Update active state and handle video
    items.forEach((item, index) => {
      item.classList.toggle('active', index === currentIndex);
      
      // Handle video when slide changes
      const video = item.querySelector('video');
      if (video) {
        if (index === currentIndex) {
          // If this is a video slide and it's now active
          video.currentTime = 0; // Reset video to start
          const playPromise = video.play();
          if (playPromise !== undefined) {
            playPromise.catch(error => console.log("Auto-play prevented"));
          }
          isVideoPlaying = true;
          // Pause autoplay while video is playing
          if (autoplayInterval) {
            clearInterval(autoplayInterval);
          }
        } else {
          // Pause video if slide is not active
          video.pause();
          video.currentTime = 0;
        }
      }
    });
  }

  function nextSlide() {
    const currentItem = items[currentIndex];
    const video = currentItem.querySelector('video');
    
    // If current slide has a video that's playing, don't auto-advance
    if (video && !video.ended && !video.paused) {
      return;
    }
    
    currentIndex = (currentIndex < items.length - 1) ? currentIndex + 1 : 0;
    updateCarousel();
  }

  function prevSlide() {
    currentIndex = (currentIndex > 0) ? currentIndex - 1 : items.length - 1;
    updateCarousel();
  }

  // Event Listeners
  prevButton.addEventListener('click', () => {
    prevSlide();
    resetAutoplay();
  });

  nextButton.addEventListener('click', () => {
    nextSlide();
    resetAutoplay();
  });

  // Handle video endings
  items.forEach(item => {
    const video = item.querySelector('video');
    if (video) {
      video.addEventListener('ended', () => {
        isVideoPlaying = false;
        startAutoplay();
        nextSlide();
      });
    }
  });

  function startAutoplay() {
    if (!isVideoPlaying) {
      autoplayInterval = setInterval(nextSlide, 5000);
    }
  }

  function resetAutoplay() {
    if (autoplayInterval) {
      clearInterval(autoplayInterval);
    }
    startAutoplay();
  }

  // Initialize carousel
  updateCarousel();
  startAutoplay();

  return function cleanup() {
    if (autoplayInterval) {
      clearInterval(autoplayInterval);
    }
    // Pause all videos on cleanup
    items.forEach(item => {
      const video = item.querySelector('video');
      if (video) {
        video.pause();
      }
    });
  };
}

