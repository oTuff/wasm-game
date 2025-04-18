package main

import (
	"fmt"
	"image/png"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

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
	imageFile, err := os.Open("../assets/bunny.png")
	if err != nil {
		fmt.Println("error opening image")
	}
	defer imageFile.Close()

	imageData, err := png.Decode(imageFile)
	if err != nil {
		fmt.Println("error decoding image data")
	}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("BunnyMark Go Ebitengine")
	if err := ebiten.RunGame(&Game{Sprite: ebiten.NewImageFromImage(imageData)}); err != nil {
		log.Fatal(err)
	}
}
