// Base layout
export const Layout = ({
  title,
  children,
}: {
  title: string;
  children: preact.ComponentChildren;
}) => (
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
              <a href="/go-ebitengine">Go WASM</a>
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

// Pages
export const Home = () => (
  <div>
    <h1>Game Engine Wasm Benchmark</h1>
  </div>
);

export const ExampleGo = () => (
  <div>
    <h1>Go Ebitengine Example</h1>
    <p>Below is the embedded game, with automatic screen scaling.</p>
    <iframe
      src="/example-go/public/main.html"
      width="640"
      height="480"
      style={{ border: "0" }}
    ></iframe>
  </div>
);

export const ExampleRust = () => (
  <div>
    <h1>Rust WASM Example</h1>
    <p>Below is the embedded game, with automatic screen scaling.</p>
  </div>
);

export const ExampleJS = () => (
  <div>
    <h1>JavaScript Example</h1>
    <p>Below is the embedded game, with automatic screen scaling.</p>
  </div>
);
