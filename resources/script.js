let leaves = new Array();
async function falling_leaves() {
    for (let i = 0; i < 200; i++) {
        leaves.push(new Leaf())

    }
    loo()
}
class Leaf {
    constructor() {
        this.element = document.createElement("div")
        this.element.classList.add("leaf")
        let img = Math.ceil(Math.random() * 4)
        this.element.style.backgroundImage = `url("${img}.png")`
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
            this.life = 60*10
        }
    }
}
async function loo() {
    requestAnimationFrame(loo)
    for (const i of leaves) {
        i.tick()
    }
}
falling_leaves()