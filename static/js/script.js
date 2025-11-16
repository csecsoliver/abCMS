
function starryNight() {
    // Create a canvas element and append it to the body
    const canvas = document.createElement('canvas');
    document.body.appendChild(canvas);
    const ctx = canvas.getContext('2d');

    // Set the canvas size to the full window size
    canvas.width = window.innerWidth;
    canvas.height = window.innerHeight;

    // Style the canvas to be in the background
    canvas.style.position = 'fixed';
    canvas.style.top = '0';
    canvas.style.left = '0';
    canvas.style.zIndex = '-1';

    // Array to hold the stars
    const stars = [];

    // Number of stars to create
    const numStars = 300;

    // Create the stars
    for (let i = 0; i < numStars; i++) {
        stars.push({
            x: Math.random() * canvas.width,
            y: Math.random() * canvas.height,
            radius: Math.random() * 1.5,
            alpha: Math.random(),
            twinkle: Math.random() * 0.04 + 0.01 // Twinkle speed
        });
    }

    // Function to draw a single star
    function drawStar(star) {
        ctx.beginPath();
        ctx.arc(star.x, star.y, star.radius, 0, Math.PI * 2);
        ctx.fillStyle = `rgba(255, 255, 255, ${star.alpha})`;
        ctx.fill();
    }

    // Function to update and redraw the stars
    function update() {
        // Clear the canvas
        ctx.clearRect(0, 0, canvas.width, canvas.height);

        // Update and draw each star
        stars.forEach(star => {
            // Make the stars twinkle
            star.alpha += star.twinkle;
            if (star.alpha > 1 || star.alpha < 0) {
                star.twinkle = -star.twinkle;
            }
            drawStar(star);
        });

        // Request the next animation frame
        requestAnimationFrame(update);
    }

    // Start the animation
    update();

    // Handle window resize
    window.addEventListener('resize', () => {
        canvas.width = window.innerWidth;
        canvas.height = window.innerHeight;
    });
}

// Run the starry night animation
starryNight();
