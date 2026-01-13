document.addEventListener('DOMContentLoaded', function () {
    const images = document.querySelectorAll('img[data-fullsrc]');
    images.forEach(img => {
        console.log("added event listener")
        img.addEventListener('mouseover', function (event) {
            this.src = this.dataset.fullsrc;
            console.log("changed src");
        });
    });
});