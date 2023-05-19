# WebAssembly impress terminal

This is a part of cross-platform GUI Library for Go. See https://github.com/codeation/impress

The WebAssembly terminal is a tool for using the browser as a GUI canvas.

Reasons to have a browser version for a GUI app:

- Access from anywhere to view some of the data created on the desktop.
- Demonstration of desktop GUI application.
- Periodically launched application with a graphical interface without installation on computer.

## Hello world

<img src="https://codeation.github.io/images/canvas_hello.png" width="749" height="685" />

Source:

```
package main

import (
	"image"
	"image/color"

	"github.com/codeation/impress"
	"github.com/codeation/impress/event"

	_ "github.com/codeation/impress/canvas"
)

func main() {
	app := impress.NewApplication(image.Rect(0, 0, 640, 480), "Hello World Application")
	defer app.Close()

	font := impress.NewFont(15, map[string]string{"family": "Verdana"})
	defer font.Close()

	w := app.NewWindow(image.Rect(0, 0, 640, 480), color.RGBA{240, 240, 240, 0})
	defer w.Drop()

	w.Text("Hello, world!", font, image.Pt(280, 210), color.RGBA{0, 0, 0, 0})
	w.Line(image.Pt(270, 230), image.Pt(380, 230), color.RGBA{255, 0, 0, 0})
	w.Show()
	app.Sync()

	for {
		action := <-app.Chan()
		if action == event.DestroyEvent || action == event.KeyExit {
			break
		}
	}
}
```

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
