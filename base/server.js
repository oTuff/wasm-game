const http = require("http");
const fs = require("fs");
const path = require("path");

const PORT = 8000;

const mimeTypes = {
  ".html": "text/html",
  ".css": "text/css",
  ".js": "application/javascript",
  ".wasm": "application/wasm",
  ".png": "image/png",
  ".jpg": "image/jpeg",
  ".svg": "image/svg+xml",
  ".json": "application/json",
  ".txt": "text/plain",
};

http
  .createServer((req, res) => {
    const urlPath = req.url === "/" ? "/index.html" : req.url;
    const filePath = path.join(__dirname, urlPath);
    const fileExt = path.extname(filePath);

    fs.readFile(filePath, (err, data) => {
      if (err) {
        res.writeHead(404, { "Content-Type": "text/plain" });
        res.end("404 Not Found");
      } else {
        const mimeType = mimeTypes[fileExt] || "application/octet-stream";
        res.writeHead(200, { "Content-Type": mimeType });
        res.end(data);
      }
    });
  })
  .listen(PORT, () => {
    console.log(`Server running at http://localhost:${PORT}`);
  });
