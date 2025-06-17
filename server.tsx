import { Application, Router, send } from "oak";
import renderToString from "preact-render-to-string";
import {
  ExampleGo,
  ExampleJS,
  ExampleRust,
  ExampleLua,
  Home,
  Layout,
} from "./views.tsx";

const router = new Router();

const pageRoutes = [
  { path: "/", title: "Home", view: <Home /> },
  { path: "/go-ebitengine", title: "Go WASM Example", view: <ExampleGo /> },
  { path: "/rust-bevy", title: "Rust WASM Example", view: <ExampleRust /> },
  { path: "/js-phaser", title: "JavaScript WASM Example", view: <ExampleJS /> },
  { path: "/lua-love2d", title: "Lua WASM Example", view: <ExampleLua /> },
];

router.use(async (ctx, next) => {
  ctx.response.headers.set("Cross-Origin-Embedder-Policy", "require-corp");
  ctx.response.headers.set("Cross-Origin-Opener-Policy", "same-origin");
  await next();
});

pageRoutes.forEach(({ path, title, view: view }) => {
  router.get(path, (ctx) => {
    const body = renderToString(<Layout title={title}>{view}</Layout>);
    ctx.response.body = "<!DOCTYPE html>" + body;
  });
});

// Static files
router.get("/public/(.*)", async (ctx) => {
  const filePath = ctx.request.url.pathname.replace("/public", "");
  await send(ctx, filePath, {
    root: `${Deno.cwd()}/public`,
  });
});

router.get("/:platform/public/(.*)", async (ctx) => {
  const { platform } = ctx.params;
  const filePath = ctx.request.url.pathname.replace(`/${platform}/public`, "");
  await send(ctx, filePath, {
    root: `${Deno.cwd()}/${platform}/public`,
  });
});

const app = new Application();
app.use(router.routes());
app.use(router.allowedMethods());

app.listen({ port: 8080 });
