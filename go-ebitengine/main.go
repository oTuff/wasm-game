//go:build js && wasm

package main

import (
	"embed"
	"fmt"
	"image/png"
	"log"
	"math/rand/v2"
	"syscall/js"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mlange-42/ark/ecs"
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

var world = ecs.NewWorld()
var mapper = ecs.NewMap2[PositionData, VelocityData](&world)
var filter = ecs.NewFilter2[PositionData, VelocityData](&world)

type Game struct {
}

func (g *Game) AddBunnies(count int) {
	for range count {
		_ = mapper.NewEntity(
			&PositionData{
				X: rand.Float64() * 5,
				Y: rand.Float64() * 5,
			},

			&VelocityData{
				X: (rand.Float64() * 2) + 2, // 2 to 4
				Y: (rand.Float64() * 2) + 2, // 2 to 4
			},
		)
	}
}

func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		g.AddBunnies(10)
	}

	// Right mouse button round up to nearest 100
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton2) {
		current := world.Stats().Entities.Total
		toAdd := 100 - (current % 100)
		if toAdd == 0 {
			toAdd = 100
		}
		g.AddBunnies(toAdd)
	}

	// Middle mouse button round up to nearest 1000
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton1) {
		current := world.Stats().Entities.Total
		toAdd := 1000 - (current % 1000)
		if toAdd == 0 {
			toAdd = 1000
		}
		g.AddBunnies(toAdd)
	}

	query := filter.Query()
	for query.Next() {
		position, speed := query.Get()

		position.X += speed.X
		position.Y += speed.Y
		speed.Y += 0.7

		scaledWidth := float64(sprite.Bounds().Dx()) * bunnyScale
		scaledHeight := float64(sprite.Bounds().Dy()) * bunnyScale
		if position.X < 0 || position.X > screenWidth-scaledWidth {
			speed.X = -speed.X
		}

		if position.Y > screenHeight-scaledHeight {
			position.Y = screenHeight - scaledHeight
			speed.Y = -speed.Y
		} else if position.Y < upperBound && speed.Y < 0 {
			speed.Y *= 0.7
		}

	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	query := filter.Query()
	for query.Next() {
		position, _ := query.Get()

		op := &ebiten.DrawImageOptions{}

		op.GeoM.Scale(bunnyScale, bunnyScale)

		op.GeoM.Translate(position.X, position.Y)

		screen.DrawImage(sprite, op)
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("fps: %.0f\ntps: %.0f\nbunnies: %v", ebiten.ActualFPS(), ebiten.ActualTPS(), world.Stats().Entities.Total))

}

func (g *Game) exposeMetrics() {
	js.Global().Set("getGoMetrics", js.FuncOf(func(this js.Value, args []js.Value) any {
		return map[string]any{
			"fps":     ebiten.ActualFPS(),
			"tps":     ebiten.ActualTPS(),
			"bunnies": world.Stats().Entities.Total,
		}
	}))
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

	// set global sprite variable
	sprite = ebiten.NewImageFromImage(imageData)

	// Create game instance
	game := &Game{}

	// expose JavaScript function to get metrics
	game.exposeMetrics()

	// Add initial bunnies
	game.AddBunnies(10)

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("BunnyMark Go Ebitengine")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
