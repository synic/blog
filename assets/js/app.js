(function () {
  if (window.__synicAppJsLoaded) return;
  window.__synicAppJsLoaded = true;

  function showScrollToTopButton() {
    const scrollButton = document.getElementById("scroll-to-top");
    const content = document.getElementById("content");
    if (!scrollButton || !content) return;

    if (window.scrollY > 100) {
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
    if (errorBox) errorBox.style.display = "none";
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
    if (errorBox) {
      errorBox.style.display = "flex";
      if (!window.errorBoxTimeout) {
        window.errorBoxTimeout = setTimeout(hideErrorBox, 1500);
      }
    }
  }

  let lastTouchEnd = 0;
  document.addEventListener("touchend", (event) => {
    const now = Date.now();
    if (now - lastTouchEnd <= 300) {
      event.preventDefault();
    }
    lastTouchEnd = now;
  }, { passive: false });

  function ensureLightbox() {
    let lb = document.getElementById("album-lightbox");
    if (lb) return lb;
    lb = document.createElement("div");
    lb.id = "album-lightbox";
    lb.className = "album-lightbox";
    lb.setAttribute("role", "dialog");
    lb.setAttribute("aria-label", "Image viewer");
    lb.setAttribute("aria-modal", "true");

    const close = document.createElement("button");
    close.type = "button";
    close.className = "album-lightbox-close";
    close.setAttribute("aria-label", "Close image viewer");
    close.textContent = "\u2715";

    const prev = document.createElement("button");
    prev.type = "button";
    prev.className = "album-lightbox-nav album-lightbox-prev";
    prev.setAttribute("aria-label", "Previous image");
    prev.textContent = "\u276E";

    const next = document.createElement("button");
    next.type = "button";
    next.className = "album-lightbox-nav album-lightbox-next";
    next.setAttribute("aria-label", "Next image");
    next.textContent = "\u276F";

    const img = document.createElement("img");
    img.className = "album-lightbox-img";
    img.alt = "";

    lb.appendChild(close);
    lb.appendChild(prev);
    lb.appendChild(next);
    lb.appendChild(img);
    document.body.appendChild(lb);

    function hide() {
      lb.classList.remove("is-open");
      img.src = "";
      img.srcset = "";
      lb._albumImgs = null;
      lb._currentIndex = -1;
    }

    function navigate(dir) {
      if (!lb._albumImgs || lb._currentIndex === -1) return;
      lb._currentIndex = (lb._currentIndex + dir + lb._albumImgs.length) % lb._albumImgs.length;
      const nextEl = lb._albumImgs[lb._currentIndex];
      updateLightboxImage(lb, nextEl);
    }

    close.addEventListener("click", hide);
    prev.addEventListener("click", (e) => { e.stopPropagation(); navigate(-1); });
    next.addEventListener("click", (e) => { e.stopPropagation(); navigate(1); });
    lb.addEventListener("click", (e) => {
      if (e.target === lb) hide();
    });

    document.addEventListener("keydown", (e) => {
      if (!lb.classList.contains("is-open")) return;
      if (e.key === "Escape") hide();
      if (e.key === "ArrowLeft") navigate(-1);
      if (e.key === "ArrowRight") navigate(1);
    });
    return lb;
  }

  function updateLightboxImage(lb, el) {
    const lbImg = lb.querySelector(".album-lightbox-img");
    let src, alt, srcset;
    if (el.tagName === "IMG") {
      src = el.src;
      alt = el.alt;
      srcset = el.srcset;
    } else {
      const img = el.querySelector("img");
      if (img) {
        src = img.src;
        alt = img.alt;
        srcset = img.srcset;
      }
    }
    lbImg.src = src || "";
    lbImg.alt = alt || "";
    lbImg.srcset = srcset || "";
  }

  function openLightbox(el, albumImgs = null, currentIndex = -1) {
    const lb = ensureLightbox();
    lb._albumImgs = albumImgs;
    lb._currentIndex = currentIndex;

    const prev = lb.querySelector(".album-lightbox-prev");
    const next = lb.querySelector(".album-lightbox-next");

    if (albumImgs && albumImgs.length > 1) {
      prev.style.display = "flex";
      next.style.display = "flex";
    } else {
      prev.style.display = "none";
      next.style.display = "none";
    }

    updateLightboxImage(lb, el);
    lb.classList.add("is-open");
  }

  function initAlbum(album) {
    if (album.dataset.albumInitialized === "1") return;
    
    const scroller = album.querySelector(".album-scroller");
    const items = album.querySelectorAll(".album-item");
    const caption = album.querySelector(".album-caption");
    const dots = album.querySelectorAll(".album-dot");
    const prev = album.querySelector(".album-nav-prev");
    const next = album.querySelector(".album-nav-next");
    if (!scroller || items.length === 0) return;

    function setActiveState(index) {
      dots.forEach((d, i) => {
        d.classList.toggle("is-active", i === index);
      });
      const currentImg = items[index]?.querySelector("img");
      if (currentImg && caption) {
        caption.textContent = currentImg.alt || "";
      }
    }

    function updateActiveState() {
      const index = Math.round(scroller.scrollLeft / scroller.clientWidth);
      setActiveState(index);
    }

    scroller.addEventListener("scroll", () => {
      clearTimeout(scroller._scrollTimeout);
      scroller._scrollTimeout = setTimeout(updateActiveState, 50);
    }, { passive: true });

    prev?.addEventListener("click", (e) => {
      e.stopPropagation();
      const currentIndex = Math.round(scroller.scrollLeft / scroller.clientWidth);
      const newIndex = (currentIndex - 1 + items.length) % items.length;
      setActiveState(newIndex);
      scroller.scrollTo({ left: newIndex * scroller.clientWidth, behavior: "smooth" });
    });

    next?.addEventListener("click", (e) => {
      e.stopPropagation();
      const currentIndex = Math.round(scroller.scrollLeft / scroller.clientWidth);
      const newIndex = (currentIndex + 1) % items.length;
      setActiveState(newIndex);
      scroller.scrollTo({ left: newIndex * scroller.clientWidth, behavior: "smooth" });
    });

    dots.forEach((dot, i) => {
      dot.addEventListener("click", (e) => {
        e.stopPropagation();
        setActiveState(i);
        scroller.scrollTo({ left: i * scroller.clientWidth, behavior: "smooth" });
      });
    });

    const albumImgs = Array.from(album.querySelectorAll(".lightbox-img"));
    albumImgs.forEach((img, i) => {
      img.addEventListener("click", (e) => {
        e.stopPropagation();
        openLightbox(img, albumImgs, i);
      });
    });

    album.dataset.albumInitialized = "1";
  }

  function initAlbums() {
    document.querySelectorAll(".album").forEach(initAlbum);
  }

  function initLightboxImages() {
    document.querySelectorAll(".lightbox-img").forEach((img) => {
      if (img.dataset.lightboxInitialized === "1") return;
      if (img.closest(".album")) return;
      img.addEventListener("click", (e) => {
        e.stopPropagation();
        openLightbox(img);
      });
      img.dataset.lightboxInitialized = "1";
    });
  }

  function init() {
    showScrollToTopButton();
    initAlbums();
    initLightboxImages();
  }

  window.addEventListener("load", () => {
    init();
    window.addEventListener("scroll", showScrollToTopButton, { passive: true });
    window.addEventListener("resize", () => {
      showScrollToTopButton();
      document.querySelectorAll(".album-scroller").forEach(s => {
        const activeDot = s.closest(".album").querySelector(".album-dot.is-active");
        if (activeDot) {
          const index = parseInt(activeDot.dataset.index);
          s.scrollTo({ left: index * s.clientWidth });
        }
      });
    }, { passive: true });
    
    window.addEventListener("htmx:afterSwap", init);

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
})();
