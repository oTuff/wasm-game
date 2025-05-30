export class PerformanceMetrics {
  static wasmExecDuration = null;

  static logWasmTiming(assetList = []) {
    const entries = performance.getEntriesByType("resource");
    let totalDuration = 0;

    assetList.forEach((asset) => {
      const entry = entries.find((e) => e.name.endsWith(asset));
      if (entry) {
        totalDuration += entry.duration;
      } else {
        console.warn(`[Perf] Resource "${asset}" not found.`);
      }
    });

    PerformanceMetrics.wasmExecDuration = totalDuration.toFixed(2);
  }

  constructor(metrics, lang) {
    this.frames = 0;
    this.lastSecond = performance.now();
    this.frameTimes = [];
    this.lastFrame = undefined;
    this.metrics = metrics;
    this.lang = lang;
    this.lastBunnyCount = 0;
    this.pendingInputTime = 0;
    this.lastMetrics = null;
    this.latency = null;
    this.pollInterval = 1;
    this.browser = (() => {
      const ua = navigator.userAgent;
      if (ua.includes("Firefox/")) return "Firefox";
      if (ua.includes("Chrome/") && ua.includes("Safari/")) return "Chromium";
      if (ua.includes("Safari/") && ua.includes("Version/")) return "GnomeWeb";
      return ua;
    })();
    this.loggedHeader = false;

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
        // console.log("stoped as fps was too low", this.frames);
        const mouseUp = new MouseEvent("mouseup", eventOptions);
        const pointerUp = new PointerEvent("pointerup", eventOptions);
        element.dispatchEvent(mouseUp);
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
          this.latency = now - this.pendingInputTime;
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
    const memory = performance.memory
      ? (performance.memory.usedJSHeapSize / 1024 / 1024).toFixed(2)
      : "";

    const game = this.metrics();

    if (!this.loggedHeader) {
      console.log(
        ",lang,browser,wasm_exec_ms,fps_js,fps_game,tps,bunnies,avg_frame,min_frame,max_frame,heap_mb,click_latency_ms",
      );
      this.loggedHeader = true;
    }

    console.log([
      "," + this.lang,
      `"${this.browser}"`,
      PerformanceMetrics.wasmExecDuration || "",
      this.frames,
      game.fps?.toFixed(2),
      game.tps?.toFixed(2),
      game.bunnies,
      avgFrame.toFixed(2),
      minFrame.toFixed(2),
      maxFrame.toFixed(2),
      memory,
      this.latency?.toFixed(2) || "",
    ].join(","));
  }
}
