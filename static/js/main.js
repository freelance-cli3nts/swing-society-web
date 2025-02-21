import { initHtmxConfig } from './modules/htmxConfig.js';
import { initLogoHandling } from './modules/logoHandler.js';
import { initNavigation } from './modules/navigation.js';
import { initCarousels } from './modules/carousel.js';
import { initCarousel } from './modules/carousel-js.js';
import { initFormHandling } from './modules/formHandler.js';
import { initRegistrationForm } from './modules/registrationForm.js';

// Initialize all modules
document.addEventListener('DOMContentLoaded', () =>{
  initHtmxConfig();
  initLogoHandling();
  initNavigation();
  initCarousels();
  initCarousel();
  initFormHandling();
  initRegistrationForm();
})

document.body.addEventListener('htmx:afterSwap', () => {
  // Only call functions that need to work with newly loaded content
  initCarousel();
  initCarousels();
  initFormHandling();
  initRegistrationForm();
});


let carouselCleanup = null;

document.body.addEventListener('htmx:beforeSwap', () => {
  // Clean up carousel before content changes
  if (typeof carouselCleanup === 'function') {
    carouselCleanup();
    carouselCleanup = null;
  }
});
