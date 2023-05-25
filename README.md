# Fiberwebgui

Create cross-platform desktop apps using Fiber and GO! 

This small package just starts the Fiber server and the Chrome* browser in app mode. Doing this allows you to use [Fiber](https://github.com/gofiber/fiber) go webframework to create a desktop application using html/css/js, go html templates anything you would use to create a website.

Fiberwebgui is an adaptation of [flaskwebgui](https://github.com/ClimenteA/flaskwebgui) python package which serves the same purpuse.


## Install

```bash
go get -u github.com/ClimenteA/fiberwebgui
```

## Usage

```go
package main

import (
	"github.com/ClimenteA/fiberwebgui"
	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	fiberwebgui.Run(app)
	// Alternatives:
	// fiberwebgui.RunOnPort(app, 5656)
	// fiberwebgui.RunWithSize(app, 800, 600)
	// fiberwebgui.RunWithSizeOnPort(app, 800, 600, 5656)
	// fiberwebgui.RunBrowser(app)
	// fiberwebgui.RunBrowserOnPort(app, 5656)
}

```

## Distribution

Here are some CLI go commands for cross-platform executables.

Windows 64bit:
```bash
GOOS=windows GOARCH=amd64 go build -ldflags -H=windowsgui -o dist/myapp.exe main.go
```

Linux 64bit:
```bash
GOOS=linux GOARCH=amd64 go build -o dist/myapp main.go
```

Mac 64bit:
```bash
GOOS=darwin GOARCH=amd64 go build -o dist/myapp main.go
```

Of course, modify these commands as needed for your specific hardware architecture.

## Observations

- Parameters `width`, `height` and maybe `fullscreen` may not work on Mac;
- Window control is limited to width, height, fullscreen;
- Remember the GUI is still a browser - pressing F5 will refresh the page + other browser specific things (you can hack it with js though);
- You don't need production level setup - you just have one user to serve;
- If you want to debug/reload features - just run it as you would normally do fiberwebgui does not provide auto-reload;


## Why Fiber and not net/http or Gin or X framework? 
Comming from a Python/JS background I found Fiber the most well documented and easy to use webframework for GO. If you need this to work with other go frameworks you can take a look at the source code and adapt it as needed (nothing to fancy there).
