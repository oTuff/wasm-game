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
            </ul>
          </nav>
        </header>
        <main style="display: flex; justify-content: center;">{children}</main>
      </body>
    </html>
  );
}

export function Home() {
  return (
    <div>
      <h1>Game Engine Wasm Benchmark</h1>
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
    </div>
  );
}

export function ExampleJS() {
  return (
    <div>
      <h1>JavaScript Example</h1>
    </div>
  );
}
