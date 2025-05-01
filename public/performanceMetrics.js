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

  constructor(metrics, name) {
    this.frames = 0;
    this.lastSecond = performance.now();
    this.frameTimes = [];
    this.lastFrame = undefined;
    this.metrics = metrics;
    this.name = name;
    this.lastBunnyCount = 0;
    this.pendingInputTime = 0;
    this.lastMetrics = null;
    this.pollInterval = 1;

    this.tick = this.tick.bind(this);
  }

  start() {
    requestAnimationFrame(this.tick);
    this.startMetricsPolling();
  }

  tick() {
    this.frames++;

    const now = performance.now();

    // Frame time
    if (this.lastFrame !== undefined) {
      this.frameTimes.push(now - this.lastFrame);
    }
    this.lastFrame = now;

    // Log once per second
    if (now - this.lastSecond >= 1000) {
      const element = document.querySelector("canvas");

      const eventOptions = {
        bubbles: true,
        cancelable: true,
        clientX: globalThis.innerWidth / 2,
        clientY: globalThis.innerHeight / 2,
        button: 0,
      };

      if (element && this.frames >= 60) {
        const mouseDown = new MouseEvent("mousedown", eventOptions);
        const pointerDown = new PointerEvent("pointerdown", eventOptions);
        element.dispatchEvent(mouseDown);
        element.dispatchEvent(pointerDown);

        this.pendingInputTime = performance.now();
      } else {
        console.log("stoped as fps was too low", this.frames);
        const mouseUp = new MouseEvent("mouseup", eventOptions);
        element.dispatchEvent(mouseUp);

        const pointerUp = new PointerEvent("pointerup", eventOptions);
        element.dispatchEvent(pointerUp);
      }

      this._logStats();
      this.frames = 0;
      this.lastSecond = now;
      this.frameTimes = [];
    }

    requestAnimationFrame(this.tick);
  }

  startMetricsPolling() {
    const poll = () => {
      const now = performance.now();
      const metrics = this.metrics();

      if (this.pendingInputTime && this.lastMetrics) {
        if (metrics.bunnies > this.lastMetrics.bunnies) {
          const delay = now - this.pendingInputTime;
          console.log(`bunny spawn delay: ${delay.toFixed(2)} ms`);
          this.pendingInputTime = 0;
        }
      }

      this.lastMetrics = metrics;
      setTimeout(poll, this.pollInterval);
    };
    poll();
  }

  _logStats() {
    const avgFrame = this.frameTimes.length
      ? this.frameTimes.reduce((a, b) => a + b, 0) / this.frameTimes.length
      : 0;
    const minFrame = this.frameTimes.length ? Math.min(...this.frameTimes) : 0;
    const maxFrame = this.frameTimes.length ? Math.max(...this.frameTimes) : 0;

    let memoryMsg = "";
    if (performance.memory) {
      memoryMsg = ` | JS Heap Used: ${
        (performance.memory.usedJSHeapSize / 1024 / 1024).toFixed(2)
      } MB`;
    }

    const browser = navigator.userAgentData
      ? `${
        navigator.userAgentData.brands.map((b) => b.brand + " " + b.version)
          .join(", ")
      }`
      : navigator.userAgent;

    console.log(this.name);
    console.log(browser);
    console.log(this.metrics());

    console.log(
      `FPS: ${this.frames}` +
        ` | Frame: avg=${avgFrame.toFixed(2)}ms min=${
          minFrame.toFixed(2)
        }ms max=${maxFrame.toFixed(2)}ms` +
        memoryMsg,
    );
  }
}
