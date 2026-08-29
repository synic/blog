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
    if (window.errorBoxTimeout) clearTimeout(window.errorBoxTimeout);
    window.errorBoxTimeout = undefined;
  }

  function handleHtmxError(event) {
    const ctx = event?.detail?.ctx;
    if (ctx?.response?.status == 404) {
      window.location = ctx.response.raw?.url || ctx.request?.action;
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
      if (!lb._albumImgs || lb._albumImgs.length === 0 || lb._currentIndex === -1) return;
      lb._currentIndex = (lb._currentIndex + dir + lb._albumImgs.length) % lb._albumImgs.length;
      const nextEl = lb._albumImgs[lb._currentIndex];
      updateLightboxImage(lb, nextEl);
    }

    close.addEventListener("click", hide);
    prev.addEventListener("click", (e) => { e.stopPropagation(); navigate(-1); });
    next.addEventListener("click", (e) => { e.stopPropagation(); navigate(1); });

    let lbTouchStartX = 0;
    let lbTouchStartY = 0;
    let lbTouchStartTime = 0;
    let lbSwiped = false;

    lb.addEventListener("touchstart", (e) => {
      if (e.touches.length !== 1) return;
      lbTouchStartX = e.touches[0].clientX;
      lbTouchStartY = e.touches[0].clientY;
      lbTouchStartTime = Date.now();
      lbSwiped = false;
    }, { passive: true });

    lb.addEventListener("touchend", (e) => {
      if (!lbTouchStartTime) return;
      const touchEndX = e.changedTouches[0].clientX;
      const touchEndY = e.changedTouches[0].clientY;
      const dx = touchEndX - lbTouchStartX;
      const dy = touchEndY - lbTouchStartY;
      lbTouchStartTime = 0;

      if (Math.abs(dx) > 30 && Math.abs(dx) > Math.abs(dy)) {
        lbSwiped = true;
        if (dx < 0) {
          navigate(1);
        } else {
          navigate(-1);
        }
        setTimeout(() => { lbSwiped = false; }, 300);
      }
    }, { passive: true });

    let lbPointerStartX = 0;
    let lbPointerStartY = 0;
    let lbIsPointerDown = false;

    lb.addEventListener("pointerdown", (e) => {
      if (e.pointerType === "touch") return;
      if (e.button !== 0) return;
      if (e.target.closest(".album-lightbox-nav") || e.target.closest(".album-lightbox-close")) return;
      lbPointerStartX = e.clientX;
      lbPointerStartY = e.clientY;
      lbIsPointerDown = true;
      lbSwiped = false;
    });

    lb.addEventListener("pointerup", (e) => {
      if (e.pointerType === "touch" || !lbIsPointerDown) return;
      lbIsPointerDown = false;
      const dx = e.clientX - lbPointerStartX;
      const dy = e.clientY - lbPointerStartY;
      if (Math.abs(dx) > 30 && Math.abs(dx) > Math.abs(dy)) {
        lbSwiped = true;
        if (dx < 0) {
          navigate(1);
        } else {
          navigate(-1);
        }
        setTimeout(() => { lbSwiped = false; }, 300);
      }
    });

    lb.addEventListener("pointercancel", () => {
      lbIsPointerDown = false;
    });

    lb.addEventListener("dragstart", (e) => {
      e.preventDefault();
    });

    lb.addEventListener("click", (e) => {
      if (lbSwiped) {
        lbSwiped = false;
        return;
      }
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

    let currentIndex = 0;
    const activeDot = album.querySelector(".album-dot.is-active");
    if (activeDot && activeDot.dataset.index) {
      currentIndex = parseInt(activeDot.dataset.index, 10) || 0;
    }

    function setActiveState(index) {
      currentIndex = (index + items.length) % items.length;
      dots.forEach((d, i) => {
        d.classList.toggle("is-active", i === currentIndex);
      });
      const currentImg = items[currentIndex]?.querySelector("img");
      if (currentImg && caption) {
        caption.textContent = currentImg.alt || "";
      }
    }

    function goToIndex(index, smooth = true) {
      setActiveState(index);
      scroller.scrollTo({
        left: currentIndex * scroller.clientWidth,
        behavior: smooth ? "smooth" : "instant",
      });
    }

    function updateActiveState() {
      if (scroller.clientWidth > 0) {
        const index = Math.round(scroller.scrollLeft / scroller.clientWidth);
        setActiveState(index);
      }
    }

    scroller.addEventListener("scroll", () => {
      clearTimeout(scroller._scrollTimeout);
      scroller._scrollTimeout = setTimeout(updateActiveState, 50);
    }, { passive: true });

    prev?.addEventListener("click", (e) => {
      e.stopPropagation();
      goToIndex(currentIndex - 1);
    });

    next?.addEventListener("click", (e) => {
      e.stopPropagation();
      goToIndex(currentIndex + 1);
    });

    dots.forEach((dot, i) => {
      dot.addEventListener("click", (e) => {
        e.stopPropagation();
        goToIndex(i);
      });
    });

    let touchStartX = 0;
    let touchStartY = 0;
    let touchStartTime = 0;
    let swiped = false;

    album.addEventListener("touchstart", (e) => {
      if (e.touches.length !== 1) return;
      touchStartX = e.touches[0].clientX;
      touchStartY = e.touches[0].clientY;
      touchStartTime = Date.now();
      swiped = false;
    }, { passive: true });

    album.addEventListener("touchend", (e) => {
      if (!touchStartTime) return;
      const touchEndX = e.changedTouches[0].clientX;
      const touchEndY = e.changedTouches[0].clientY;
      const dx = touchEndX - touchStartX;
      const dy = touchEndY - touchStartY;
      touchStartTime = 0;

      if (Math.abs(dx) > 30 && Math.abs(dx) > Math.abs(dy)) {
        swiped = true;
        if (dx < 0) {
          goToIndex(currentIndex + 1);
        } else {
          goToIndex(currentIndex - 1);
        }
        setTimeout(() => { swiped = false; }, 300);
      }
    }, { passive: true });

    let pointerStartX = 0;
    let pointerStartY = 0;
    let isPointerDown = false;

    album.addEventListener("pointerdown", (e) => {
      if (e.pointerType === "touch") return;
      if (e.button !== 0) return;
      if (e.target.closest(".album-nav") || e.target.closest(".album-dots")) return;
      pointerStartX = e.clientX;
      pointerStartY = e.clientY;
      isPointerDown = true;
      swiped = false;
    });

    album.addEventListener("pointerup", (e) => {
      if (e.pointerType === "touch" || !isPointerDown) return;
      isPointerDown = false;
      const dx = e.clientX - pointerStartX;
      const dy = e.clientY - pointerStartY;
      if (Math.abs(dx) > 30 && Math.abs(dx) > Math.abs(dy)) {
        swiped = true;
        if (dx < 0) {
          goToIndex(currentIndex + 1);
        } else {
          goToIndex(currentIndex - 1);
        }
        setTimeout(() => { swiped = false; }, 300);
      }
    });

    album.addEventListener("pointercancel", () => {
      isPointerDown = false;
    });

    album.addEventListener("dragstart", (e) => {
      e.preventDefault();
    });

    const albumImgs = Array.from(album.querySelectorAll(".lightbox-img"));
    albumImgs.forEach((img, i) => {
      img.addEventListener("click", (e) => {
        if (swiped) {
          swiped = false;
          return;
        }
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

  function formatCommentDates() {
    document.querySelectorAll("time.comment-time[datetime]").forEach((el) => {
      const dt = el.getAttribute("datetime");
      if (!dt) return;
      const d = new Date(dt);
      if (isNaN(d.getTime())) return;
      el.textContent = new Intl.DateTimeFormat(undefined, {
        year: "numeric",
        month: "short",
        day: "numeric",
        hour: "numeric",
        minute: "2-digit",
      }).format(d);
    });
  }

  function syncSearchInputs() {
    const search = new URLSearchParams(location.search).get("search") || "";
    document.querySelectorAll("#search-nav, #search").forEach((el) => {
      if (el === document.activeElement) return;
      if (el.value !== search) el.value = search;
    });
  }

  function init() {
    showScrollToTopButton();
    initAlbums();
    initLightboxImages();
    formatCommentDates();
    syncSearchInputs();
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
    
    document.getElementById("errorbox")?.addEventListener("click", hideErrorBox);

    window.addEventListener("htmx:after:swap", init);

    document.body.addEventListener("htmx:config:request", (event) => {
      const csrfToken = document.cookie
        .split("; ")
        .find((row) => row.startsWith("csrf_token="))
        ?.split("=")[1];
      if (csrfToken) {
        event.detail.ctx.request.headers["X-CSRF-Token"] = csrfToken;
      }
    });

    document.body.addEventListener("htmx:response:error", handleHtmxError);
    document.body.addEventListener("htmx:error", handleHtmxError);
  });
})();
