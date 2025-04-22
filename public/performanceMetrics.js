export class PerformanceMetrics {
  static markWasmStart() {
    performance.mark("start");
  }

  static markWasmEndAndLog() {
    performance.mark("end");
    performance.measure("WASM execution", "start", "end");
    const entry = performance.getEntriesByName("WASM execution")[0];
    if (entry) {
      const { name, duration, startTime } = entry;
      console.log(
        `${name} time: ${duration.toFixed(2)} ms (started at ${
          startTime.toFixed(2)
        } ms after page load)`,
      );
    }
  }

  constructor() {
    this.frames = 0;
    this.lastSecond = performance.now();
    this.frameTimes = [];
    this.lastInput = 0;
    this.inputLatencies = [];
    this.lastFrame = undefined;

    globalThis.addEventListener("mousedown", () => {
      this.lastInput = performance.now();
    });

    this.tick = this.tick.bind(this);
  }

  start() {
    requestAnimationFrame(this.tick);
  }

  tick() {
    this.frames++;
    const now = performance.now();

    // Frame time
    if (this.lastFrame !== undefined) {
      this.frameTimes.push(now - this.lastFrame);
    }
    this.lastFrame = now;

    // Input latency
    if (this.lastInput > 0) {
      this.inputLatencies.push(now - this.lastInput);
      this.lastInput = 0;
    }

    // Log once per second
    if (now - this.lastSecond >= 1000) {
      this._logStats();
      this.frames = 0;
      this.lastSecond = now;
      this.frameTimes = [];
      this.inputLatencies = [];
    }

    requestAnimationFrame(this.tick);
  }

  logWasmSize(wasmPath) {
    fetch(wasmPath)
      .then((r) => r.arrayBuffer())
      .then((buffer) => {
        const mb = buffer.byteLength / 1024 / 1024;
        console.log(
          `WASM size: ${mb.toFixed(2)} MB (${buffer.byteLength} bytes)`,
        );
      });
  }

  _formatMB(bytes) {
    return (bytes / 1024 / 1024).toFixed(2);
  }

  _average(arr) {
    if (!arr.length) return 0;
    return arr.reduce((a, b) => a + b, 0) / arr.length;
  }

  _min(arr) {
    return arr.length ? Math.min(...arr) : 0;
  }

  _max(arr) {
    return arr.length ? Math.max(...arr) : 0;
  }

  _logStats() {
    const fps = this.frames;
    const avgFrame = this._average(this.frameTimes);
    const minFrame = this._min(this.frameTimes);
    const maxFrame = this._max(this.frameTimes);
    const avgInput = this._average(this.inputLatencies);
    const minInput = this._min(this.inputLatencies);
    const maxInput = this._max(this.inputLatencies);

    let memoryMsg = "";
    if (performance.memory) {
      memoryMsg = ` | JS Heap Used: ${
        this._formatMB(performance.memory.usedJSHeapSize)
      } MB`;
    }

    console.log(
      `FPS: ${fps}` +
        ` | Frame: avg=${avgFrame.toFixed(2)}ms min=${
          minFrame.toFixed(2)
        }ms max=${maxFrame.toFixed(2)}ms` +
        ` | Input latency: avg=${avgInput.toFixed(2)}ms min=${
          minInput.toFixed(2)
        }ms max=${maxInput.toFixed(2)}ms` +
        memoryMsg,
    );
  }
}
