import { Application, Router, RouterContext, send } from "oak";
import renderToString from "preact-render-to-string";
import { Layout, Home, ExampleGo, ExampleRust, ExampleJS } from "./views.tsx";

const router = new Router();

// TODO: simplify routing
router.get("/", (ctx) => {
  const body = renderToString(
    <Layout title="Hello oak!">
      <Home />
    </Layout>,
  );
  ctx.response.body = "<!DOCTYPE html>" + body;
});

router.get("/go-ebitengine", (ctx) => {
  const body = renderToString(
    <Layout title="Go WASM Example">
      <ExampleGo />
    </Layout>,
  );
  ctx.response.body = "<!DOCTYPE html>" + body;
});

router.get("/rust-bevy", (ctx) => {
  const body = renderToString(
    <Layout title="Rust WASM Example">
      <ExampleRust />
    </Layout>,
  );
  ctx.response.body = "<!DOCTYPE html>" + body;
});

router.get("/js-phaser", (ctx) => {
  const body = renderToString(
    <Layout title="JavaScript WASM Example">
      <ExampleJS />
    </Layout>,
  );
  ctx.response.body = "<!DOCTYPE html>" + body;
});

// Serve static files from /example-go/public
router.get("/example-go/public/(.*)", async (ctx) => {
  await send(ctx, ctx.request.url.pathname.replace("/example-go/public", ""), {
    root: `${Deno.cwd()}/example-go/public`,
  });
});

// TODO: implement the other static files for the other implementation

const app = new Application();
app.use(router.routes());
app.use(router.allowedMethods());

app.listen({ port: 8080 });

