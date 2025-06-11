//go:build js && wasm

package main

import (
	"embed"
	"fmt"
	"image/png"
	"log"
	"math/rand/v2"
	"syscall/js"

	"codeberg.org/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 640
	screenHeight = 480
	upperBound   = 40 // Amount of pixels from the top
	bunnyScale   = 0.2
)

//go:embed assets
var content embed.FS

var sprite *ebiten.Image

// Components
type PositionData struct {
	X, Y float64
}

type VelocityData struct {
	X, Y float64
}

func NewBunny() gohan.Entity {
	bunny := gohan.NewEntity()

	bunny.AddComponent(&PositionData{
		X: rand.Float64() * 5,
		Y: rand.Float64() * 5,
	})

	bunny.AddComponent(&VelocityData{
		X: (rand.Float64() * 2) + 2, // 2 to 4
		Y: (rand.Float64() * 2) + 2, // 2 to 4
	})

	return bunny

}

type Game struct {
}

func AddBunnies(count int) {
	for range count {
		NewBunny()
	}
}

func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		AddBunnies(10)
	}

	// Right mouse button round up to nearest 100
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton2) {
		current := len(gohan.AllEntities())
		toAdd := 100 - (current % 100)
		if toAdd == 0 {
			toAdd = 100
		}
		AddBunnies(toAdd)
	}

	// Middle mouse button round up to nearest 1000
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton1) {
		current := len(gohan.AllEntities())
		toAdd := 1000 - (current % 1000)
		if toAdd == 0 {
			toAdd = 1000
		}
		AddBunnies(toAdd)
	}
	return gohan.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	err := gohan.Draw(screen)
	if err != nil {
		panic(err)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("fps: %.0f\ntps: %.0f\nbunnies: %v", ebiten.ActualFPS(), ebiten.ActualTPS(), gohan.CurrentEntities()))

}

func (g *Game) exposeMetrics() {
	js.Global().Set("getGoMetrics", js.FuncOf(func(this js.Value, args []js.Value) any {
		return map[string]any{
			"fps":     ebiten.ActualFPS(),
			"tps":     ebiten.ActualTPS(),
			"bunnies": gohan.CurrentEntities(),
		}
	}))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

type movementSystem struct {
	Position *PositionData
	Speed    *VelocityData
}

func NewMovementSystem() *movementSystem {
	return &movementSystem{}
}

func (s *movementSystem) Update(_ gohan.Entity) error {

	s.Position.X += s.Speed.X
	s.Position.Y += s.Speed.Y
	s.Speed.Y += 0.7

	scaledWidth := float64(sprite.Bounds().Dx()) * bunnyScale
	scaledHeight := float64(sprite.Bounds().Dy()) * bunnyScale
	if s.Position.X < 0 || s.Position.X > screenWidth-scaledWidth {
		s.Speed.X = -s.Speed.X
	}

	if s.Position.Y > screenHeight-scaledHeight {
		s.Position.Y = screenHeight - scaledHeight
		s.Speed.Y = -s.Speed.Y
	} else if s.Position.Y < upperBound && s.Speed.Y < 0 {
		s.Speed.Y *= 0.7
	}
	return nil
}

func (s *movementSystem) Draw(_ gohan.Entity, screen *ebiten.Image) error {
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Scale(bunnyScale, bunnyScale)

	op.GeoM.Translate(s.Position.X, s.Position.Y)

	screen.DrawImage(sprite, op)

	return nil
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

	// set global sprite variable
	sprite = ebiten.NewImageFromImage(imageData)

	// Create game instance
	game := &Game{}

	gohan.AddSystem(NewMovementSystem())

	// expose JavaScript function to get metrics
	game.exposeMetrics()

	// Add initial bunnies
	AddBunnies(10)

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("BunnyMark Go Ebitengine")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
