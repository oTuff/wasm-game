export function Layout({
  title,
  children,
}: {
  title: string;
  children: preact.ComponentChildren;
}) {
  return (
    <html>
      <head>
        <title>{title}</title>
        <meta charSet="UTF-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
        <link
          rel="stylesheet"
          href="https://unpkg.com/@picocss/pico@latest/css/pico.min.css"
        />
      </head>
      <body class="container">
        <header>
          <nav>
            <ul>
              <li>
                <a href="/">Home</a>
              </li>
              <li>
                <a href="/go-ebitengine">Go Ebitengine</a>
              </li>
              <li>
                <a href="/rust-bevy">Rust Bevy</a>
              </li>
              <li>
                <a href="/js-phaser">JS Phaser</a>
              </li>
              <li>
                <a href="/lua-love2d">Lua Love2d</a>
              </li>
            </ul>
          </nav>
        </header>
        <main style="display: flex; justify-content: center;">{children}</main>{" "}
        <p style="font-weight: bold;">Controls:</p>
        <ul>
          <li>'left click' to rapidly add bunnies</li>
          <li>'right click' to round up 100 bunnies</li>
          <li>'middle click' to round up 1000 bunnies</li>
        </ul>
      </body>
    </html>
  );
}

export function Home() {
  return (
    <div>
      <h1>Game Engine Wasm Benchmark</h1>
      <p>
        Navigate to the different benchmarks by clicking on the links on the
        navbar at the top of the page
      </p>
    </div>
  );
}

export function ExampleGo() {
  return (
    <div>
      <h1>Go Wasm Example</h1>
      <iframe
        src="/go-ebitengine/public/main.html"
        width="640"
        height="480"
        style={{ border: "0" }}
      ></iframe>
    </div>
  );
}

export function ExampleRust() {
  return (
    <div>
      <h1>Rust WASM Example</h1>
      <iframe
        src="/rust-bevy/public/main.html"
        width="640"
        height="480"
      ></iframe>
    </div>
  );
}

export function ExampleJS() {
  return (
    <div>
      <h1>JavaScript Example</h1>
      <iframe
        src="/js-phaser/public/main.html"
        width="640"
        height="480"
        scrolling="no"
      ></iframe>
      <p>This example is incomplete - all the controls won't work</p>
    </div>
  );
}

export function ExampleLua() {
  return (
    <div>
      <h1>Lua Example</h1>
      <iframe
        src="/lua-love2d/public/main.html"
        width="640"
        height="480"
        scrolling="no"
      ></iframe>
      <p>This example is incomplete</p>
    </div>
  );
}
