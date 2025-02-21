// Wait for both DOM and HTMX to be ready
export function initHtmxConfig() {
    if (window.htmx) {
    htmx.config.ignoreFeatureWarnings = true;
      
      const path = window.location.pathname;
        
      // Get current path and handle pretty URLs
      const urlMapping = {
          '/classes/solo-jazz': '/templates/classes/solojazz.html',
          '/classes/lindy-hop': '/templates/classes/lindyhop.html',
          '/classes/old-clips': '/templates/classes/oldclips.html',
          '/events/parties': '/templates/events/parties.html',
          '/events/festivals': '/templates/events/festivals.html',
          '/events/workshops': '/templates/events/workshops.html'
      };

      if (path in urlMapping) {

          htmx.ajax('GET', urlMapping[path], {target: '#main-view'});
      }
  } else {
      // If HTMX isn't loaded yet, wait a bit and try again
      setTimeout(configureHtmx, 50);
  }
};
