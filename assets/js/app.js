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

function hideErrorBox() {
  const errorBox = document.getElementById("errorbox");
  errorBox.style.display = "none";
  window.errorBoxTimeout = undefined;
}

window.addEventListener("load", () => {
  // scroll
  showScrollToTopButton();
  window.addEventListener("scroll", showScrollToTopButton);
  window.addEventListener("resize", showScrollToTopButton);

  // Watch for content changes to update scroll button
  const content = document.getElementById("content");
  if (content) {
    const observer = new MutationObserver(showScrollToTopButton);
    observer.observe(content, { childList: true, subtree: true });
  }
});
