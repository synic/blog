window.addEventListener("load", () => {
  function showScrollToTopButton() {
    const scrollButton = document.getElementById("scroll-to-top");
    const content = document.getElementById("content");

    if (window.pageYOffset > 100) {
      scrollButton.style.display = "block";
      const btnRect = scrollButton.getBoundingClientRect();
      const contentRect = content.getBoundingClientRect();
      scrollButton.style.left =
        (contentRect.right - btnRect.width).toString() + "px";
    } else {
      scrollButton.style.display = "none";
    }
  }
  showScrollToTopButton();
  window.addEventListener("scroll", showScrollToTopButton);
  window.addEventListener("resize", showScrollToTopButton);
  window.addEventListener("htmx:afterSwap", showScrollToTopButton);

  function hideErrorBox() {
    const errorBox = document.getElementById("errorbox");
    errorBox.style.display = "none";
    window.errorBoxTimeout = undefined;
  }

  function showErrorBox() {
    const errorBox = document.getElementById("errorbox");
    errorBox.style.display = "flex";
    clearTimeout(window.errorBoxTimeout);
    window.errorBoxTimeout = setTimeout(hideErrorBox, 3500);
  }

  document.body.addEventListener("htmx:onLoadError", showErrorBox);
  document.body.addEventListener("htmx:sendError", showErrorBox);
});
