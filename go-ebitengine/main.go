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
)

const (
	screenWidth  = 640
	screenHeight = 480
	upperBound   = 40 // Amount of pixels from the top
	bunnyScale   = 0.2
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
		PosX:   rand.Float64() * 5,
		PosY:   rand.Float64() * 5,
		ScaleX: bunnyScale,
		ScaleY: bunnyScale,
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
	newBunnies := make([]Bunny, count, count)
	// fmt.Printf("count: %v\n", count)
	// fmt.Printf("before: cap %v, len %v, %p\n", cap(newBunnies), len(newBunnies), newBunnies)
	for i := range count {
		newBunnies[i] = *NewBunny()
		// newBunnies = append(newBunnies, *NewBunny())
		// g.Bunnies = append(g.Bunnies, *NewBunny())
		// fmt.Printf("cap %v, len %v, %p\n", cap(newBunnies), len(newBunnies), newBunnies)
	}
	g.Bunnies = append(g.Bunnies, newBunnies...)
}

func (g *Game) Update() error {

	// fmt.Printf("cap %v, len %v, %p\n", cap(g.Bunnies), len(g.Bunnies), g.Bunnies)

	// Left mouse button rapid fire add 10 bunnies
	if ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
		g.AddBunnies(10)
	}

	// Right mouse button round up to nearest 100
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton2) {
		current := len(g.Bunnies)
		toAdd := 100 - (current % 100)
		if toAdd == 0 {
			toAdd = 100
		}
		g.AddBunnies(toAdd)
	}

	// Middle mouse button round up to nearest 1000
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton1) {
		current := len(g.Bunnies)
		toAdd := 1000 - (current % 1000)
		if toAdd == 0 {
			toAdd = 1000
		}
		g.AddBunnies(toAdd)
	}

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

	ebitenutil.DebugPrint(screen, fmt.Sprintf("fps: %.0f\ntps: %.0f\nbunnies: %v", ebiten.ActualFPS(), ebiten.ActualTPS(), len(g.Bunnies)))
}

func (g *Game) exposeMetrics() {
	js.Global().Set("getGoMetrics", js.FuncOf(func(this js.Value, args []js.Value) any {
		return map[string]any{
			"fps":     ebiten.ActualFPS(),
			"tps":     ebiten.ActualTPS(),
			"bunnies": len(g.Bunnies),
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

	sprite := ebiten.NewImageFromImage(imageData)

	// Create game instance
	game := &Game{
		Sprite:  sprite,
		Bunnies: make([]Bunny, 0, 1000), // Pre-allocate capacity for performance
		Gravity: 0.75,
	}
	// expose JavaScript function to get metrics
	game.exposeMetrics()

	// Add initial bunnies
	game.AddBunnies(10)

	// ebiten.SetScreenClearedEveryFrame(false)
	// ebiten.SetVsyncEnabled(true)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("BunnyMark Go Ebitengine")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
