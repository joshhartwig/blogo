document.querySelectorAll('nav a').forEach(a => {
  if (a.getAttribute('href') === location.pathname) {
    a.classList.add('active');
  }
});
