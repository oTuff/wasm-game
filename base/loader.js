document.getElementById("load-engine").addEventListener("click", async () => {
  const engine = document.getElementById("engine-select").value;
  const canvas = document.getElementById("game-canvas");
  canvas.getContext("2d").clearRect(0, 0, canvas.width, canvas.height); // Clear canvas

  try {
    // if (engine === "rust-bevy") {
    //   const wasm = await import(`./rust-bevy/bevy_bunnymark.js`);
    //   await wasm.default(); // Initialize WASM module
    //   console.log("Rust-Bevy engine loaded.");
    // } else
    if (engine === "go-ebitengine") {
      const go = new Go();
      const result = await WebAssembly.instantiateStreaming(
        fetch("./go-ebitengine/main.wasm"),
        go.importObject,
      );
      go.run(result.instance); // Start Go runtime
      console.log("Go-Ebitengine engine loaded.");
    } else if (engine === "js-phaser") {
      await import("./js-phaser/bundle.js");
      console.log("JS-Phaser engine loaded.");
    } else {
      console.error("Unknown engine selected.");
    }
  } catch (error) {
    console.error("Error loading engine:", error);
  }
});
