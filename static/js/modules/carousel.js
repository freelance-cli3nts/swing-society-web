export function initCarousels() {
  console.log("Initializing carousels");

  // document.querySelectorAll('.carousel').forEach(carousel => {
    const carousels = document.querySelectorAll('.carousel');
    console.log(`Found ${carousels.length} carousels`);

    carousels.forEach(carousel => {
      const type = carousel.dataset.carousel;
      console.log(`Setting up carousel type: ${type}`);
      
      let items = [];
      // Your existing item setup...
  
      setupCarousel(carousel);
    });

    // const type = carousel.dataset.carousel;
    
    // let items = [];

    // if (type === 'hero') {
    //   items = [
    //     { type: 'image', src: '/static/assets/images/hero1.jpg', alt: 'Swing Dance' },
    //     { type: 'video', src: '/static/assets/videos/hero.mp4' },
    //     { type: 'text', text: 'Join our Swing Society today!' }
    //   ];
    // } else if (type === 'classes') {
    //   items = [
    //     { type: 'image', src: '/static/assets/images/class1.jpg', alt: 'Class Image 1' },
    //     { type: 'image', src: '/static/assets/images/class2.jpg', alt: 'Class Image 2' },
    //     { type: 'text', text: 'Learn Lindy Hop from the best instructors!' }
    //   ];
    // }

    // setupCarousel(carousel, items);
  // });
}

function setupCarousel(carousel, items) {
  if (!carousel) {
    console.error("Carousel element not found");
    return;
  }

  const carouselInner = carousel.querySelector('.carousel-inner');
  const prevButton = carousel.querySelector('.carousel-control.prev');
  const nextButton = carousel.querySelector('.carousel-control.next');
  
  let currentIndex = 0;

  carouselInner.innerHTML = items.map(item => {
    if (item.type === 'image') {
      return `<div class="carousel-item"><img src="${item.src}" alt="${item.alt}"></div>`;
    } else if (item.type === 'video') {
      return `<div class="carousel-item"><video controls muted playsinline>
                <source src="${item.src}" type="video/mp4">
              </video></div>`;
    } else if (item.type === 'text') {
      return `<div class="carousel-item"><p>${item.text}</p></div>`;
    }
  }).join('');


  if (!carouselInner) {
    console.error("Carousel inner container not found");
    return;
  }

  const carouselItems = carousel.querySelectorAll('.carousel-item');
  console.log(`Carousel has ${carouselItems.length} items`);

  if (carouselItems.length === 0) {
    console.warn("No items found in carousel");
    return;
  }

  function updateCarousel() {
    carouselInner.style.transform = `translateX(-${currentIndex * 100}%)`;
  }

    // Setup event listeners
    if (prevButton) {
      prevButton.addEventListener('click', () => {
        currentIndex = (currentIndex > 0) ? currentIndex - 1 : carouselItems.length - 1;
        updateCarousel();
      });
    }
  
    if (nextButton) {
      nextButton.addEventListener('click', () => {
        currentIndex = (currentIndex < carouselItems.length - 1) ? currentIndex + 1 : 0;
        updateCarousel();
      });
    }

  // prevButton.addEventListener('click', () => {
  //   currentIndex = (currentIndex > 0) ? currentIndex - 1 : carouselItems.length - 1;
  //   updateCarousel();
  // });

  // nextButton.addEventListener('click', () => {
  //   currentIndex = (currentIndex < carouselItems.length - 1) ? currentIndex + 1 : 0;
  //   updateCarousel();
  // });

  updateCarousel();

  // Optional: Auto-rotate
  if (carousel.dataset.autorotate === 'true') {
    const interval = parseInt(carousel.dataset.interval || '5000');
    const autoRotate = setInterval(() => {
      currentIndex = (currentIndex < carouselItems.length - 1) ? currentIndex + 1 : 0;
      updateCarousel();
    }, interval);
    
    // Cleanup function
    return function cleanup() {
      clearInterval(autoRotate);
    };
  }
}

