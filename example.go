package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/basicfont"
)

const (
	SCREEN_WIDTH    = 640
	SCREEN_HEIGHT   = 480
	TILE_SIZE       = 20
	FRAMES_PER_MOVE = 5 // the lower the faster the snake moves
)

type Pos struct {
	X, Y float64
}

type Player struct {
	Img    *ebiten.Image
	Pieces []Pos
	Dx, Dy float64
}

type Apple struct {
	Img *ebiten.Image
	Pos Pos
}

type Game struct {
	Player       Player
	Apple        Apple
	InputBuffer  []Pos
	FrameCounter int
	Paused       bool
	Score        int
}

func getRandomApplePosition(g *Game) Pos {
	// Create map of snake positions
	snakePositions := make(map[Pos]bool)
	for _, piece := range g.Player.Pieces {
		snakePositions[piece] = true
	}

	// Generate available positions in one pass
	var available []Pos
	for x := 0; x < SCREEN_WIDTH; x += TILE_SIZE {
		for y := 0; y < SCREEN_HEIGHT; y += TILE_SIZE {
			pos := Pos{X: float64(x), Y: float64(y)}
			if !snakePositions[pos] {
				available = append(available, pos)
			}
		}
	}

	if len(available) == 0 {
		fmt.Println("Congratulations! You've won!")
		os.Exit(0)
	}

	return available[rand.Intn(len(available))]
}

func checkCollision(g *Game) {
	head := g.Player.Pieces[0]

	// Check wall collision
	if head.X < 0 || head.X >= SCREEN_WIDTH || head.Y < 0 || head.Y >= SCREEN_HEIGHT {
		fmt.Println("You ran into the wall! Game Over. \n")
		os.Exit(0)
	}

	// Check self-collision
	for i := 1; i < len(g.Player.Pieces); i++ {
		if head.X == g.Player.Pieces[i].X && head.Y == g.Player.Pieces[i].Y {
			fmt.Println("You ran into yourself! Game Over. \n")
			os.Exit(0)
		}
	}
}

func (g *Game) Update() error {
	if !g.Paused {
		// Collect inputs without validation
		if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
			g.InputBuffer = append(g.InputBuffer, Pos{X: 0, Y: -TILE_SIZE})
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
			g.InputBuffer = append(g.InputBuffer, Pos{X: -TILE_SIZE, Y: 0})
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
			g.InputBuffer = append(g.InputBuffer, Pos{X: 0, Y: TILE_SIZE})
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
			g.InputBuffer = append(g.InputBuffer, Pos{X: TILE_SIZE, Y: 0})
		}

		if g.FrameCounter == FRAMES_PER_MOVE {
			g.FrameCounter = 1

			if len(g.InputBuffer) > 0 {
				next := g.InputBuffer[0]
				// Opposite direction check
				if !((g.Player.Dx == -next.X && next.X != 0) ||
					(g.Player.Dy == -next.Y && next.Y != 0)) {
					g.Player.Dx, g.Player.Dy = next.X, next.Y
				}
				g.InputBuffer = g.InputBuffer[1:]
			}

			// Save last tail position for growth
			lastPos := g.Player.Pieces[len(g.Player.Pieces)-1]

			// Move body from tail to head
			if len(g.Player.Pieces) > 1 {
				for i := len(g.Player.Pieces) - 1; i > 0; i-- {
					g.Player.Pieces[i] = g.Player.Pieces[i-1]
				}
			}

			// Move the head
			g.Player.Pieces[0].X += g.Player.Dx
			g.Player.Pieces[0].Y += g.Player.Dy

			checkCollision(g)

			// Check if apple is eaten
			if g.Player.Pieces[0].X == g.Apple.Pos.X && g.Player.Pieces[0].Y == g.Apple.Pos.Y {
				// Grow by adding new tail segment
				g.Player.Pieces = append(g.Player.Pieces, lastPos)
				g.Apple.Pos = getRandomApplePosition(g) // Randomize apple position
				g.Score += 1
			}
		} else {
			g.FrameCounter += 1
		}

	}
	// q to quit
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		os.Exit(0)
	}

	// Toggle pause state of the game with the space bar
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.Paused = !g.Paused
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 20, 50, 255})

	// Draw snake
	eyeSize := TILE_SIZE / 4
	eyeImg := ebiten.NewImage(eyeSize, eyeSize)
	eyeImg.Fill(color.Black)

	// Draw snake pieces
	for i, p := range g.Player.Pieces {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(p.X, p.Y)
		screen.DrawImage(g.Player.Img, op)

		// Draw eyes only on head piece
		if i == 0 {
			// Calculate eye positions based on direction
			var leftEyeX, leftEyeY, rightEyeX, rightEyeY float64

			switch {
			case g.Player.Dx < 0: // Moving left
				leftEyeX = p.X + TILE_SIZE*0.2
				leftEyeY = p.Y + TILE_SIZE*0.2
				rightEyeX = p.X + TILE_SIZE*0.2
				rightEyeY = p.Y + TILE_SIZE*0.6
			case g.Player.Dy < 0: // Moving down
				leftEyeX = p.X + TILE_SIZE*0.2
				leftEyeY = p.Y + TILE_SIZE*0.2
				rightEyeX = p.X + TILE_SIZE*0.6
				rightEyeY = p.Y + TILE_SIZE*0.2
			case g.Player.Dy > 0: // Moving up
				leftEyeX = p.X + TILE_SIZE*0.2
				leftEyeY = p.Y + TILE_SIZE*0.7
				rightEyeX = p.X + TILE_SIZE*0.6
				rightEyeY = p.Y + TILE_SIZE*0.7
			default: // Moving right or start
				leftEyeX = p.X + TILE_SIZE*0.7
				leftEyeY = p.Y + TILE_SIZE*0.2
				rightEyeX = p.X + TILE_SIZE*0.7
				rightEyeY = p.Y + TILE_SIZE*0.6
			}

			// Draw left eye
			eyeOp := &ebiten.DrawImageOptions{}
			eyeOp.GeoM.Translate(leftEyeX, leftEyeY)
			screen.DrawImage(eyeImg, eyeOp)

			// Draw right eye
			eyeOp = &ebiten.DrawImageOptions{}
			eyeOp.GeoM.Translate(rightEyeX, rightEyeY)
			screen.DrawImage(eyeImg, eyeOp)
		}
	}

	// Draw apple
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.Apple.Pos.X, g.Apple.Pos.Y)
	screen.DrawImage(g.Apple.Img, op)

	// Draw "Game Paused" message
	if g.Paused {
		op := ebiten.DrawImageOptions{}
		msg := "GAME PAUSED - press space to continue"
		op.GeoM.Translate(SCREEN_WIDTH/2, SCREEN_HEIGHT/2)
		text.Draw(screen, msg, text.NewGoXFace(basicfont.Face7x13), &text.DrawOptions{
			DrawImageOptions: op,
			LayoutOptions:    text.LayoutOptions{PrimaryAlign: text.AlignCenter},
		})
	}

	// Draw score
	text.Draw(screen, fmt.Sprintf("Score: %v", g.Score), text.NewGoXFace(basicfont.Face7x13), nil)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return SCREEN_WIDTH, SCREEN_HEIGHT
}

// Instantiate and run game
func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Snake")

	// snake(player)
	img := ebiten.NewImage(TILE_SIZE, TILE_SIZE)
	img.Fill(color.RGBA{
		R: 0,
		G: 255,
		B: 075,
		A: 255,
	})
	player := Player{
		Img:    img,
		Pieces: []Pos{Pos{X: TILE_SIZE, Y: TILE_SIZE}},
	}

	// apple
	appleImg := ebiten.NewImage(TILE_SIZE, TILE_SIZE)
	appleImg.Fill(color.RGBA{
		R: 255,
		G: 0,
		B: 0,
		A: 255,
	})

	// initial random position
	randX := float64(rand.Intn(SCREEN_WIDTH/TILE_SIZE) * TILE_SIZE)
	randY := float64(rand.Intn(SCREEN_HEIGHT/TILE_SIZE) * TILE_SIZE)
	// Make sure it is not the same position as the player
	if randX == TILE_SIZE {
		randX = TILE_SIZE * 2
	}
	apple := Apple{
		Img: appleImg,
		Pos: Pos{
			X: randX,
			Y: randY,
		},
	}

	// Run the game
	if err := ebiten.RunGame(&Game{Player: player, Apple: apple}); err != nil {
		log.Fatal(err)
	}
}
