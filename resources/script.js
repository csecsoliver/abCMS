let leaves = new Array();
async function falling_leaves() {
    // Respect reduced motion preference
    if (window.matchMedia && window.matchMedia('(prefers-reduced-motion: reduce)').matches) {
        console.log(localStorage.getItem("motion") == "optout")
        if ((localStorage.getItem("motion") == null)) {
            localStorage.setItem("motion", confirm("You have reduced motion enabled. Do you want to see the pretty leaves?")?"optin":"optout");
        } else if (localStorage.getItem("motion") == "optout") {
            return;
        }
    }
    for (let i = 0; i < 50; i++) {
        leaves.push(new Leaf())

    }
    loo()
}
class Leaf {
    constructor() {
        this.element = document.createElement("div")
        this.element.classList.add("leaf")
        let imgNum = Math.ceil(Math.random() * 4)
        let imgUrl = `${imgNum}.png`
        this.element.style.backgroundImage = `url("/${imgUrl}")`
        
        // this is to avoid cutting them in half
        let img = new Image()
        img.onload = () => {
            let aspectRatio = img.naturalWidth / img.naturalHeight
            this.element.style.width = (50 * aspectRatio) + "px"
        }
        img.src = "/" + imgUrl
        
        document.querySelector("body").appendChild(this.element)
        this.x = Math.random() * window.innerWidth
        this.y = Math.random() * window.innerHeight
        this.element.style.left = this.x + "px"
        this.element.style.top = this.y + "px"
        this.dx = 0
        this.dy = 0
        this.rotation = Math.random() * 360
        this.drotation = 0
    }
    tick() {
        if(this.disabled){
            this.life -= 1
            if (this.life <= 0){

                leaves.splice(leaves.findIndex(e => e === this), 1)
                this.element.remove()
            }
            return
        }
        this.dx = (this.dx) + (Math.random()*2 - 1)*0.1
        this.dy = (this.dy) + (Math.random()*2 - 1)*0.1
        this.drotation = (this.drotation) + (Math.random()*2 - 1)*0.1
        this.x += this.dx
        this.y += this.dy
        this.rotation += this.drotation
        this.element.style.left = (this.x) + "px"
        this.element.style.top = (this.y) + "px"
        this.element.style.transform = `translate(-50%, -50%) rotate(${this.rotation}deg)`
        if (this.x > window.innerWidth || this.x < 0 || this.y > window.innerHeight || this.y < 0 ){
            leaves.push(new Leaf())
            this.disabled = true
            this.life = 60*40
        }
    }
}
async function loo() {
    requestAnimationFrame(loo)
    // Pause animation when page is hidden
    if (document.hidden) return;
    for (const i of leaves) {
        i.tick()
    }
}
falling_leaves()