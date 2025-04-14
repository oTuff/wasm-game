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

---

In comparison to Golang and Ebitengine, Rust and Bevy do not have a simple
solution to wasm.

Bevy does not really have any docs, but there does exist good
examples(https://github.com/bevyengine/bevy/tree/latest/examples#wasm) and an
unofficial "cheatbook": https://bevy-cheatbook.github.io/platforms/wasm.html
(which is severely outdated)

Rust points to the "rust wasm" book on its main website along with some other
useful links [The official wasm website](https://webassembly.org/) and
https://developer.mozilla.org/en-US/docs/WebAssembly

As with much of the Rust ecosystem one of the following third-party tools are
needed for wasm support:

- wasm-server-runner
- wasm-pack
- trunk-rs

wasm-pack even has its own "book": https://rustwasm.github.io/docs/wasm-pack/

Maybe when following the example wasm setup is quite easy (but it is hidden at
the bottom of a markdown file in the GitHub repo - not as accessible as the
Golang and Ebitengine documentation)
