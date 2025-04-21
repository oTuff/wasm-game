package main

import (
	"embed"
	"image/png"
	"log"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 640
	screenHeight = 480
	upperBound   = 40 // Amount of pixels from the top
)

//go:embed assets
var content embed.FS

type Bunny struct {
	PosX, PosY      float64
	ScaleX, ScaleY  float64
	SpeedX, SpeedY  float64
	BounceVariation float64
	WobbleFactor    float64
}

func NewBunny() *Bunny {
	return &Bunny{
		PosX:   0.0,
		PosY:   0.0,
		ScaleX: 0.5,
		ScaleY: 0.5,
		SpeedX: (rand.Float64() * 2) + 2, // 2 to 4
		SpeedY: (rand.Float64() * 2) + 2, // 2 to 4
	}
}

type Game struct {
	Sprite  *ebiten.Image
	Bunnies []Bunny
	Gravity float64
}

func (g *Game) edgeDetection(b *Bunny) {
	scaledWidth := float64(g.Sprite.Bounds().Dx()) * b.ScaleX
	scaledHeight := float64(g.Sprite.Bounds().Dy()) * b.ScaleY

	if b.PosX < 0 || b.PosX > screenWidth-scaledWidth {
		b.SpeedX = -b.SpeedX
	}

	if b.PosY > screenHeight-scaledHeight {
		b.PosY = screenHeight - scaledHeight
		b.SpeedY = -b.SpeedY
	} else if b.PosY < upperBound && b.SpeedY < 0 {
		b.SpeedY *= 0.7
	}

}

func (g *Game) AddBunnies(count int) {
	newBunnies := make([]Bunny, count)
	for i := range count {
		newBunnies[i] = *NewBunny()
	}
	g.Bunnies = append(g.Bunnies, newBunnies...)
}

func (g *Game) Update() error {
	for i := range g.Bunnies {
		bunny := &g.Bunnies[i]

		bunny.PosX += bunny.SpeedX
		bunny.PosY += bunny.SpeedY
		bunny.SpeedY += g.Gravity
		g.edgeDetection(bunny)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, b := range g.Bunnies {
		op := &ebiten.DrawImageOptions{}

		op.GeoM.Scale(b.ScaleX, b.ScaleY)

		op.GeoM.Translate(b.PosX, b.PosY)

		screen.DrawImage(g.Sprite, op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	// Load image
	imageFile, err := content.Open("assets/bunny.png")
	if err != nil {
		log.Fatal("error opening image:", err)
	}
	defer imageFile.Close()

	imageData, err := png.Decode(imageFile)
	if err != nil {
		log.Fatal("error decoding image data:", err)
	}

	sprite := ebiten.NewImageFromImage(imageData)

	// Create game instance
	game := &Game{
		Sprite:  sprite,
		Bunnies: make([]Bunny, 0), // 1000), // Pre-allocate capacity for performance
		Gravity: 0.75,
	}

	// Add initial bunnies
	game.AddBunnies(5)

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("BunnyMark Go Ebitengine")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
