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
