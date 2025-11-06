async function jumps() {
    setTimeout(()=>{


        let sound = sfxr.generate("synth");
        let image = new Image();
        image.src = "../../1.png"
        image.onload = (e)=>{
            sfxr.play(sound)
            console.log("played jumps")
        }
        console.log("set up jumps")
    }, 2000);
}