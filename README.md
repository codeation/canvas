# WebAssembly impress terminal

This is a part of cross-platform GUI Library for Go. See https://github.com/codeation/impress

The WebAssembly terminal is a tool for using the browser as a GUI canvas.

Reasons to have a browser version for a GUI app:

- Access from anywhere to view some of the data created on the desktop.
- Demonstration of desktop GUI application.
- Periodically launched application with a graphical interface without installation on computer.

## Hello world

<img src="https://codeation.github.io/images/canvas_hello.png" width="749" height="685" />

[The hello world source](https://github.com/codeation/canvas/blob/master/examples/wasm/wasm.go) 
located in [examples folder](https://github.com/codeation/canvas/tree/master/examples).

To run :

```
git clone https://github.com/codeation/canvas.git
cd canvas
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .
GOOS=js GOARCH=wasm go build -o main.wasm ./cmd/
go run ./examples/wasm
```

To see results, open `http://localhost:8080/` in you browser. Please, try with Brave or Chrome browser.

## Project State

### Notes

- The performance of the browser will not be enough for the regular use.
- This project is still in the early stages of development and is not yet in a usable state.
- The project tested on Debian 11.7 with Brave browser (Version 1.51.114 Chromium: 113.0.5672.92).
