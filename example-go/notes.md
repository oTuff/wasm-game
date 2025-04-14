# Golang wasm example test

Based on the guide from
[Ebitengine's documentation](https://ebitengine.org/en/documents/webassembly.html)
compiling to Wasm is really simple and easy. It relies on Golang's build in
support for Wasm: https://go.dev/wiki/WebAssembly.

In two simple copy-paste commands you have a working Wasm file of your go
project.

To run the server locally to serve the website DenoJS, an alternative runtime to
the popular NodeJS (and also made by the same author), is used for all the
examples to keep the comparison fair.

Since one of the focus points for this project is developer experience Deno also
falls into that category?

Deno has all the tooling in one place and is trying to become the de facto
JavaScript tool chain, similar to how both Rust and Golang has a complete tool
chain.

Deno also has build-in `benchmark` tool that might come in handy for the
comparison.
