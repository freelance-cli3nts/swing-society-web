import { initHtmxConfig } from './modules/htmxConfig.js';
import { initLogoHandling } from './modules/logoHandler.js';
import { initNavigation } from './modules/navigation.js';
import { initModalHandling } from './modules/modal.js';
import { initCarousels } from './modules/carousel.js';
import { initCarousel } from './modules/carousel-js.js';


// Initialize all modules
document.addEventListener('DOMContentLoaded', () =>{
  initHtmxConfig();
  initLogoHandling();
  initNavigation();
  initModalHandling();
  initCarousels();
  initCarousel();
})

document.body.addEventListener('htmx:afterSwap', () => {
  // Only call functions that need to work with newly loaded content
  initCarousel();
  initCarousels();
});


let carouselCleanup = null;

document.body.addEventListener('htmx:beforeSwap', () => {
  // Clean up carousel before content changes
  if (typeof carouselCleanup === 'function') {
    carouselCleanup();
    carouselCleanup = null;
  }
});

// check if this is properly written
function togglePartnerField(isAlone) {
  const partnerField = document.getElementById('partner-name-group');
  partnerField.style.display = isAlone ? 'none' : 'block';
}