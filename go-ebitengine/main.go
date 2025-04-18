package main

import (
	"embed"
	"fmt"
	"image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets
var content embed.FS

type Game struct {
	Sprite *ebiten.Image
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Scale(0.5, 0.5)

	screen.DrawImage(g.Sprite, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	// Load image
	imageFile, err := content.Open("assets/bunny.png")
	if err != nil {
		log.Fatal("error opening image:", err)
	}
	defer func() {
		if imageFile != nil {
			imageFile.Close()
		}
	}()

	imageData, err := png.Decode(imageFile)
	if err != nil {
		log.Fatal("error decoding image data:", err)
	}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("BunnyMark Go Ebitengine")
	if err := ebiten.RunGame(&Game{Sprite: ebiten.NewImageFromImage(imageData)}); err != nil {
		log.Fatal(err)
	}
}
