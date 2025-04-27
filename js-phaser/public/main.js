//Slightly modified version of: https://github.com/phaserjs/examples/blob/17d166da9ce77185b75f6f603dd91791137c4fcd/public/3.55/src/game%20objects/blitter/benchmark%20test%203.js#L4
let blitter;
let gravity = 0.75;
let idx = 1;
const screenWidth = 640;
const screenHeight = 480;
const bunnyScale = 0.2;

class Example extends Phaser.Scene {
  constructor() {
    super();
  }

  preload() {
    // Set the base URL to load assets from the local server
    this.load.setBaseURL("/public/assets/");

    // Load the single bunny image
    this.load.image("bunny", "bunny.png");
    // bunny.setScale(0.2); // Adjust this value to scale the bunny as needed

    this.numbers = [];
    this.iter = 0;
  }

  launch() {
    const bunny = blitter.create(0, 0, "bunny");

    bunny.data.vx = Math.random() * 5;
    bunny.data.vy = Math.random() * 5;
    bunny.data.bounce = 0.8 + (Math.random() * 0.2);
  }

  create() {
    // Create a blitter object using the loaded bunny image
    blitter = this.add.blitter(0, 0, "bunny");

    // Create several bunny objects using the bunny image
    for (var i = 0; i < 10; ++i) {
      this.launch();
    }

    this.updateDigits();
  }

  update() {
    if (this.input.activePointer.isDown) {
      for (var i = 0; i < 250; ++i) {
        this.launch();
      }

      this.updateDigits();
    }

    for (
      var index = 0, length = blitter.children.list.length;
      index < length;
      ++index
    ) {
      var bunny = blitter.children.list[index];

      bunny.data.vy += gravity;

      bunny.y += bunny.data.vy;
      bunny.x += bunny.data.vx;

      // Handling the bouncing and boundary logic
      if (bunny.x > 640) {
        bunny.x = 640;
        bunny.data.vx *= -bunny.data.bounce;
      } else if (bunny.x < 0) {
        bunny.x = 0;
        bunny.data.vx *= -bunny.data.bounce;
      }

      if (bunny.y > 480) {
        bunny.y = 480;
        bunny.data.vy *= -bunny.data.bounce;
      }
    }
  }

  updateDigits() {
    const len = Phaser.Utils.String.Pad(
      blitter.children.list.length.toString(),
      7,
      "0",
      1,
    );
  }
}

// Phaser game configuration
const config = {
  type: Phaser.AUTO, //Phaser.WEBGL,
  parent: "phaser-example", // Make sure this matches the ID of the container in your HTML
  scene: [Example],
  width: screenWidth,
  height: screenHeight,
  // pixelArt: true, // Use pixel art mode if you're working with pixel-based art (optional)
};

// Create the Phaser game instance
const game = new Phaser.Game(config);
