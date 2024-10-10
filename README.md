# WebAssembly impress terminal

This is a part of cross-platform GUI Library for Go. See https://github.com/codeation/impress

The WebAssembly terminal is a tool for using the browser as a GUI canvas.

Reasons to have a browser version for a GUI app:

- Access from anywhere to view some of the data created on the desktop.
- Demonstration of desktop GUI application.
- Periodically launched application with a graphical interface without installation on computer.

## Hello world

<img src="https://codeation.github.io/images/canvas_small.png" width="675" height="445" />

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
    app := impress.NewApplication(image.Rect(0, 0, 480, 240), "Hello World Application")
    defer app.Close()

    font := app.NewFont(15, map[string]string{"family": "Verdana"})
    defer font.Close()

    w := app.NewWindow(image.Rect(0, 0, 480, 240), color.RGBA{240, 240, 240, 255})
    defer w.Drop()

    w.Text("Hello, world!", font, image.Pt(200, 100), color.RGBA{0, 0, 0, 255})
    w.Line(image.Pt(200, 120), image.Pt(300, 120), color.RGBA{255, 0, 0, 255})
    w.Show()
    app.Sync()

    for {
        e := <-app.Chan()
        if e == event.DestroyEvent || e == event.KeyExit {
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
- The project tested on Debian 12.6 with Brave browser (Version 1.66.113 Chromium: 125.0.6422.76).
