function initImageHover() {
    const images = document.querySelectorAll('img[data-full-src]');
    images.forEach(img => {
        console.log("added event listener");
        img.addEventListener('mouseover', function (event) {
            this.src = this.dataset.fullSrc;
            console.log("changed src");
        });
    });
}

document.addEventListener('DOMContentLoaded', initImageHover);