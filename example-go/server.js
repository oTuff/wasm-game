import { serve } from "https://deno.land/std@0.185.0/http/server.ts";

const PORT = 8080;

const handler = async (request) => {
  const url = new URL(request.url);
  const path = url.pathname === "/" ? "/render.html" : url.pathname;

  try {
    const file = await Deno.readFile(`.${path}`);
    const contentType =
      {
        ".html": "text/html",
        ".js": "application/javascript",
        ".wasm": "application/wasm",
        ".mod": "text/plain",
        ".sum": "text/plain",
        ".go": "text/plain",
      }[path.slice(path.lastIndexOf("."))] || "application/octet-stream";

    return new Response(file, {
      headers: { "Content-Type": contentType },
    });
  } catch {
    return new Response("404 Not Found", { status: 404 });
  }
};

console.log(`Server is running on http://localhost:${PORT}`);
serve(handler, { port: PORT });
