import { serveDir } from "https://deno.land/std/http/file_server.ts";
import { serve } from "https://deno.land/std/http/server.ts";

const PORT = 8080;

// Serve the current directory
const handler = (req) => {
  return serveDir(req, {
    fsRoot: ".",
    urlRoot: "",
    showDirListing: true,
    enableCors: true,
  });
};

console.log(`Server is running on http://localhost:${PORT}`);
serve(handler, { port: PORT });