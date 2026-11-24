function initImageHover(root) {
    const scope = root && root.querySelectorAll ? root : document;
    const images = scope.querySelectorAll('img[data-full-src]');
    images.forEach(img => {
        if (img.dataset.hoverBound === 'true') {
            return;
        }
        img.dataset.hoverBound = 'true';
        img.addEventListener('mouseover', function () {
            if (this.dataset.fullSrc && this.src !== this.dataset.fullSrc) {
                this.src = this.dataset.fullSrc;
            }
        });
    });
}

if (typeof htmx !== 'undefined' && typeof htmx.onLoad === 'function') {
    htmx.onLoad(function (content) {
        initImageHover(content);
    });
}

document.addEventListener('DOMContentLoaded', function () {
    initImageHover(document);
});
