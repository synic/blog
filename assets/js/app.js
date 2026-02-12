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

function handleHtmxError(event) {
  if (event?.detail?.xhr?.status == 404) {
    window.location =
      event.detail.pathInfo.finalRequestPath ||
      event.detail.pathInfo.requestPath;
    return;
  }

  const errorBox = document.getElementById("errorbox");
  errorBox.style.display = "flex";
  clearTimeout(window.errorBoxTimeout);
  window.errorBoxTimeout = setTimeout(hideErrorBox, 3500);
}

window.addEventListener("load", () => {
  showScrollToTopButton();
  window.addEventListener("scroll", showScrollToTopButton);
  window.addEventListener("resize", showScrollToTopButton);
  window.addEventListener("htmx:afterSwap", showScrollToTopButton);

  document.body.addEventListener("htmx:configRequest", (event) => {
    const csrfToken = document.cookie
      .split("; ")
      .find((row) => row.startsWith("csrf_token="))
      ?.split("=")[1];
    if (csrfToken) {
      event.detail.headers["X-CSRF-Token"] = csrfToken;
    }
  });

  document.body.addEventListener("htmx:onLoadError", handleHtmxError);
  document.body.addEventListener("htmx:responseError", handleHtmxError);
  document.body.addEventListener("htmx:sendError", handleHtmxError);
});
