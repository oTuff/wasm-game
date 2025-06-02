// Slightly modified version of: https://github.com/phaserjs/examples/blob/17d166da9ce77185b75f6f603dd91791137c4fcd/public/3.55/src/game%20objects/blitter/benchmark%20test%203.js
//
// Copyright (c) 2024 Richard Davey, Phaser Studio Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

let bunnies;
const gravity = 0.75;
const screenWidth = 640;
const screenHeight = 480;
const upperBound = 40;
const bunnyScale = 0.2;

class Example extends Phaser.Scene {
  constructor() {
    super();
    // Fixed tickrate at 60 updates per second
    this.tickrate = 1000 / 60;
    this.accumulator = 0;
    this.ticks = 0;
    this.frameTime = 0;
    this.frames = 0;
  }

  preload() {
    this.load.setBaseURL("/public/assets/");
    this.load.image("bunny", "bunny.png");
  }

  launch() {
    const bunny = this.add.sprite(0, 0, "bunny").setScale(bunnyScale);

    bunny.vx = (Math.random() * 2) + 2;
    bunny.vy = (Math.random() * 2) + 2;

    bunny.bounce = 1.0;
    bunny.widthScaled = bunny.displayWidth;
    bunny.heightScaled = bunny.displayHeight;

    bunnies.add(bunny);
  }

  create() {
    bunnies = this.add.group();

    for (let i = 0; i < 10; ++i) {
      this.launch();
    }

    // expose metrics
    let metrics = {
      fps: 0,
      tps: 0,
      bunnies: 0,
    };

    // Diagnostics
    this.diagnostics = this.add.text(10, 10, "", {
      font: "16px Courier",
    }).setDepth(10);
  }

  update(time, delta) {
    this.accumulator += delta;

    while (this.accumulator >= this.tickrate) {
      this.accumulator -= this.tickrate;
      this.ticks++;

      // Update logic at fixed interval
      if (this.input.activePointer.isDown) {
        for (let i = 0; i < 10; ++i) {
          this.launch();
        }
      }

      Phaser.Actions.Call(bunnies.getChildren(), (bunny) => {
        bunny.x += bunny.vx;
        bunny.y += bunny.vy;
        bunny.vy += gravity;

        const maxX = screenWidth - bunny.widthScaled;
        const maxY = screenHeight - bunny.heightScaled;

        // Bounce horizontally
        if (bunny.x < 0) {
          bunny.x = 0;
          bunny.vx = -bunny.vx;
        } else if (bunny.x > maxX) {
          bunny.x = maxX;
          bunny.vx = -bunny.vx;
        }

        // Bounce vertically
        if (bunny.y > maxY) {
          bunny.y = maxY;
          bunny.vy = -bunny.vy;
        } else if (bunny.y < upperBound && bunny.vy < 0) {
          bunny.vy *= 0.7;
        }
      });
    }

    this.frames++;
    this.frameTime += delta;

    if (this.frameTime >= 1000) {
      const fps = this.frames;
      const tps = this.ticks;
      const bunnyCount = bunnies.getChildren().length;

      this.frameTime = 0;
      this.frames = 0;

      this.ticks = 0;

      this.metrics = {
        fps: fps,
        tps: tps,
        bunnies: bunnyCount,
      };
      this.diagnostics.setText([
        `FPS: ${fps}`,
        `TPS: ${tps}`,
        `Bunnies: ${bunnyCount}`,
      ]);
    }
  }
}

const config = {
  type: Phaser.AUTO,
  parent: "phaser-example",
  scene: [Example],
  width: screenWidth,
  height: screenHeight,
  pixelArt: true,
};

const game = new Phaser.Game(config);
