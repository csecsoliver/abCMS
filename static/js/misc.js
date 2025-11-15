document.addEventListener('DOMContentLoaded', function () {
    const images = document.querySelectorAll('img[data-fullsrc]');
    images.forEach(img => {
        console.log("added event listener")
        img.addEventListener('mouseover', function (event) {
            // if (event.cancelable){

            // event.preventDefault();
            this.src = this.dataset.fullsrc;
            console.log("changed src");
            // menu = new MouseEvent("contextmenu", {
            //     bubbles: true,
            //     cancelable: false,
            //     view: window,
            //     button: 2,
            // })
            // event.target.dispatchEvent(menu);
            // } else {
            // console.log("not cancellable, continuing")
            // }

        });
    });
});