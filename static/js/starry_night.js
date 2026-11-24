(function () {
    const STORAGE_KEY = 'abcms-reduced-motion';
    const mediaQuery = window.matchMedia
        ? window.matchMedia('(prefers-reduced-motion: reduce)')
        : null;

    function getStoredPreference() {
        try {
            const val = window.localStorage.getItem(STORAGE_KEY);
            if (val === 'true') return true;
            if (val === 'false') return false;
        } catch (_) {
            // Ignore storage access errors (e.g. private browsing restrictions).
        }
        return null;
    }

    function savePreference(reduced) {
        try {
            window.localStorage.setItem(STORAGE_KEY, reduced ? 'true' : 'false');
        } catch (_) {
            // Ignore storage write errors.
        }
    }

    function resolveInitialReducedMotion() {
        const stored = getStoredPreference();
        if (stored !== null) {
            return stored;
        }
        return Boolean(mediaQuery && mediaQuery.matches);
    }

    function starryNight() {
        if (document.getElementById('starry-night-canvas')) {
            return;
        }

        // Create a persistent canvas element outside <main> so HTMX swaps never reset it
        const canvas = document.createElement('canvas');
        canvas.id = 'starry-night-canvas';
        canvas.setAttribute('hx-preserve', 'true');
        canvas.setAttribute('aria-hidden', 'true');
        document.body.appendChild(canvas);
        const context = canvas.getContext('2d');

        canvas.width = window.innerWidth;
        canvas.height = window.innerHeight;

        canvas.style.position = 'fixed';
        canvas.style.top = '0';
        canvas.style.left = '0';
        canvas.style.zIndex = '-1';
        canvas.style.pointerEvents = 'none';

        // Store normalized coordinates so resizing scales cleanly without blank regions
        const stars = [];
        const starCount = 5000;
        for (let i = 0; i < starCount; i++) {
            stars.push({
                nx: Math.random(),
                ny: Math.random(),
                radius: Math.random() * 1.5,
                alpha: Math.random(),
                twinkle: Math.random() * 0.04 - 0.02
            });
        }

        let reducedMotion = resolveInitialReducedMotion();
        let rafId = null;

        function renderFrame(animate) {
            context.clearRect(0, 0, canvas.width, canvas.height);
            const width = canvas.width;
            const height = canvas.height;

            for (let i = 0; i < stars.length; i++) {
                const star = stars[i];
                if (animate) {
                    star.alpha += star.twinkle;
                    if (star.alpha > 1 || star.alpha < 0) {
                        star.twinkle = -star.twinkle;
                        star.alpha = Math.max(0, Math.min(1, star.alpha));
                    }
                }
                context.beginPath();
                context.arc(star.nx * width, star.ny * height, star.radius, 0, Math.PI * 2);
                context.fillStyle = `rgba(255, 255, 255, ${star.alpha})`;
                context.fill();
            }
        }

        function tick() {
            if (reducedMotion) {
                rafId = null;
                return;
            }
            renderFrame(true);
            rafId = requestAnimationFrame(tick);
        }

        function syncToggleUI() {
            document.documentElement.dataset.reducedMotion = reducedMotion ? 'true' : 'false';
            const toggle = document.getElementById('reduced-motion-toggle');
            if (!toggle) {
                return;
            }
            toggle.setAttribute('aria-checked', reducedMotion ? 'false' : 'true');
            toggle.textContent = reducedMotion ? 'Motion: OFF' : 'Motion: ON';
        }

        function applyMotionMode(nextReduced) {
            reducedMotion = Boolean(nextReduced);
            syncToggleUI();

            if (reducedMotion) {
                if (rafId !== null) {
                    cancelAnimationFrame(rafId);
                    rafId = null;
                }
                // Render a single calm static starfield frame.
                renderFrame(false);
            } else if (rafId === null) {
                tick();
            }
        }

        // Use event delegation so the toggle keeps working after HTMX swaps #site-nav.
        document.addEventListener('click', function (evt) {
            const toggle = evt.target && evt.target.closest
                ? evt.target.closest('#reduced-motion-toggle')
                : null;
            if (!toggle) {
                return;
            }
            evt.preventDefault();
            const next = !reducedMotion;
            savePreference(next);
            applyMotionMode(next);
        });

        // Re-sync the nav toggle text/aria-checked whenever HTMX swaps #site-nav.
        document.addEventListener('htmx:after:swap', syncToggleUI);
        if (typeof htmx !== 'undefined' && typeof htmx.onLoad === 'function') {
            htmx.onLoad(syncToggleUI);
        }

        if (mediaQuery && typeof mediaQuery.addEventListener === 'function') {
            mediaQuery.addEventListener('change', function (evt) {
                if (getStoredPreference() === null) {
                    applyMotionMode(evt.matches);
                }
            });
        }

        window.addEventListener('resize', function () {
            canvas.width = window.innerWidth;
            canvas.height = window.innerHeight;
            if (reducedMotion) {
                renderFrame(false);
            }
        });

        applyMotionMode(reducedMotion);
    }

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', starryNight);
    } else {
        starryNight();
    }
})();
