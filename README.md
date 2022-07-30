## raylib-go-headless
[![Build Status](https://github.com/icodealot/raylib-go-headless/actions/workflows/build.yml/badge.svg)](https://github.com/icodealot/raylib-go-headless/actions)
[![GoDoc](https://godoc.org/github.com/icodealot/raylib-go-headless/raylib?status.svg)](https://godoc.org/github.com/icodealot/raylib-go-headless/raylib)
[![Go Report Card](https://goreportcard.com/badge/github.com/icodealot/raylib-go)](https://goreportcard.com/report/github.com/icodealot/raylib-go-headless)
[![Examples](https://img.shields.io/badge/learn%20by-examples-0077b3.svg?style=flat-square)](https://github.com/icodealot/raylib-go-headless/tree/master/examples)

My goal with this library is to be able to do headless rendering in Docker containers using raylib bindings for Golang. The bindings themselves are heavily based on the excellent work of [raylib-go](https://github.com/gen2brain/raylib-go)

There is nothing stopping you from trying to use this outside of a Docker container, with a Linux server or distro of your choice, but the examples here will be geared towards the former. On Ubuntu for example you might need to add some libraries such as:

```
sudo apt-get update -y; sudo apt-get -y install libgl1-mesa-dev libosmesa6-dev
```

### Docker Requirements

##### Docker with golang + osmesa

###### Example Dockerfile: 
    
	FROM golang:1.18-alpine
  
    RUN apk update

	RUN apk add \
        build-base \
        mesa-dev \
        mesa-osmesa

    COPY your app code etc...

	RUN go mod tidy

	RUN go build -o yourservice

	ENTRYPOINT yourservice


### Module Installation

    go get -v -u github.com/icodealot/raylib-go-headless/raylib

### Documentation

Documentation on [GoDoc](https://godoc.org/github.com/icodealot/raylib-go-headless/raylib). Also check raylib [cheatsheet](http://www.raylib.com/cheatsheet/cheatsheet.html).

### Example

There is an example of rendering images from raylib and sending them as image/png payloads to HTTP reqeusts.

You will need:

- The docker daemon running
- Access to run ```docker build``` and ```docker run```
- cURL or A web browser to test

Clone the Git repo:
```
git clone https://github.com/icodealot/raylib-go-headless
```

```
cd raylib-go-headless/examples/imageserver
```

Build and run the docker container:
```
docker build -t imageserver-example .
```

```
docker run -dp 8080:8080 imageserver-example:latest
```

Open a web browser to http://localhost:8080/

or

```
curl -o example.png http://localhost:8080/
```

If all goes well you should see something similar to this image.

![example](./examples/imageserver/example.png)

### Example source

```go
package main

import (
	"bytes"
	"fmt"
	"image/png"
	"net/http"
	"strconv"

	rl "github.com/icodealot/raylib-go-headless/raylib"
)

func main() {

	// Raylib image rendering and sending results over the http socket

	rl.InitRaylib()

	render := rl.LoadRenderTexture(800, 450)

	rl.BeginTextureMode(render)
	rl.ClearBackground(rl.RayWhite)
	rl.DrawRectangleGradientV(0, 0, 800, 450, rl.RayWhite, rl.Red)
	rl.DrawText("Hello, World!", 20, 20, 20, rl.DarkGray)
	rl.EndTextureMode()

	image := rl.LoadImageFromTexture(render.Texture)

	rl.ImageFlipVertical(*&image) // gl buffers are y flipped

	// This sends a pre-rendered image but it could be setup to render
	// a new imge for each HTTP request.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		buffer := new(bytes.Buffer)
		if err := png.Encode(buffer, image.ToImage()); err != nil {
			http.Error(w, "error encoding image", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "image/png")
		w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
		_, err := w.Write(buffer.Bytes())
		if err != nil {
			return
		}
	})

	fmt.Printf("Starting server")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}

	defer rl.UnloadRenderTexture(render)
	defer rl.UnloadImage(image)
	defer rl.CloseRaylib()
}
```

### License

raylib-go-headless is licensed under an unmodified zlib/libpng license. View [LICENSE](https://github.com/icodealot/raylib-go-headless/blob/master/LICENSE).
