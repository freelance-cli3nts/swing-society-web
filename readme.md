# Things to do
## fix the mobile view of the main page
## add for us, class and contacts page


.hamburger {
  display: none;
  font-size: 2rem;
  cursor: pointer;
  z-index: 1001;
  color: var(--black);
}

.nav-list {
  display: flex;
  align-items: center;
  gap: 2rem;
}

.nav-list li {
  position: relative;
}

.nav-list li a {
  font-size: 2.5rem;
  padding: 0.5rem 1rem;
}

.nav-list.open {
  padding: 1rem;
}

.submenu {
  display: none;
  position: static;
  top: 100%;
  left: 0;
  background: var(--shade);
  min-width: 12.5rem;
  box-shadow: 0 0.125rem 0.3125rem rgba(146, 61, 61, 0.2);
  border-radius: 0.25rem;
  z-index: 1000;
}

.submenu li {
  background: var(--shade);
  width: 100%;
  padding: 0.1rem 0;
}

.submenu li a {
  display: block;
  color: var(--black);
  background: var(--shade);
  font-size: 1.5rem;
  transition: all 0.3s ease;

}

.has-submenu.active + li {
  margin-top: auto; /* Push next item down when submenu opens */
}

.has-submenu .submenu {
  max-height: 0;
  overflow: hidden;
  transition: max-height 0.3s ease-out;
}

.has-submenu.active .submenu {
  max-height: 15rem; /* Adjust based on content */
}


.submenu li a:hover {
  background: var(--primary);
  border-radius: 0.25rem;
}
