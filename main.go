package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/math/f64"
)

const (
	ZOOM_RATIO                 = 10
	CAMERA_MOVE_RATIO          = 10
	WORLD_WIDTH                = 10000
	WORLD_HEIGHT               = 10000
	SCREEN_WIDTH               = 1920
	SCREEN_HEIGHT              = 1080
	CAMERA_MOVE_SPEED          = 10
	BASE_RADIUS                = 8.0
	BASE_SYSTEM_RADIUS         = BASE_RADIUS
	BASE_SYSTEM_STROKE         = BASE_RADIUS
	BASE_PLANET_RADIUS         = BASE_RADIUS
	BASE_STAR_RADIUS           = 20.0
	BASE_SPACE_ELEVATOR_LENGTH = 1.0
	MAX_SPACE_ELEVATOR_LENGTH  = BASE_SPACE_ELEVATOR_LENGTH * 5
)

var (
	DEBUG_MODE  = *(flag.Bool("debug", true, "enable debug mode"))
	DEBUG_COLOR = color.RGBA{B: 255}
	lastFrame   = time.Now()
	G           = 6.674 * math.Pow10(-11)
)

type Game struct {
	camera Camera
	world  *ebiten.Image

	Galaxy *Galaxy
}

func NewGame(worldWidth, worldHeight int) *Game {
	var (
		camera = Camera{
			ViewPort: f64.Vec2{
				SCREEN_WIDTH,
				SCREEN_HEIGHT,
			},
			Position: f64.Vec2{
				WORLD_WIDTH/2 - SCREEN_WIDTH/2,
				WORLD_HEIGHT/2 - SCREEN_HEIGHT/2,
			},
			ZoomFactor: -100,
		}
		world = ebiten.NewImage(worldWidth, worldHeight)
	)
	g := GenGalaxy(100)

	return &Game{
		camera: camera,
		world:  world,
		Galaxy: g,
	}
}

func (g *Game) Update() error {
	_, dy := ebiten.Wheel()
	g.camera.ZoomFactor = g.camera.ZoomFactor + int(dy*ZOOM_RATIO)

	switch true {
	case ebiten.IsKeyPressed(ebiten.KeyQ):
		return fmt.Errorf("quit")
	case ebiten.IsKeyPressed(ebiten.KeyArrowLeft):
		g.camera.Position[0] -= CAMERA_MOVE_SPEED
	case ebiten.IsKeyPressed(ebiten.KeyArrowRight):
		g.camera.Position[0] += CAMERA_MOVE_SPEED
	case ebiten.IsKeyPressed(ebiten.KeyArrowUp):
		g.camera.Position[1] -= CAMERA_MOVE_SPEED
	case ebiten.IsKeyPressed(ebiten.KeyArrowDown):
		g.camera.Position[1] += CAMERA_MOVE_SPEED
	case ebiten.IsKeyPressed(ebiten.KeySpace):
		g.camera.Reset()
	}

	// mouse button todo
	//if ebiten.IsMouseButtonPressed(ebiten.MouseButton3) {

	g.Galaxy.Tick()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	g.world.Clear()
	if DEBUG_MODE {
		screen.Fill(DEBUG_COLOR)
	}

	g.Galaxy.Systems[0].DrawScene(g.world)

	g.camera.Render(g.world, screen)

	_, _ = g.camera.ScreenToWorld(ebiten.CursorPosition())

	cx, cy := ebiten.CursorPosition()

	ebitenutil.DebugPrint(
		screen,
		fmt.Sprintf(
			"FPS: %f\nKey Press: %v\nCursor Positino: %v, %v\n",
			ebiten.ActualFPS(),
			ebiten.InputChars(),
			cx, cy,
		),
	)
	lastFrame = time.Now()
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return SCREEN_WIDTH, SCREEN_HEIGHT
}

func GenGalaxy(numSystems int) *Galaxy {
	var (
		systems = make([]StarSystem, 0, numSystems)
	)

	for i := 0; i <= numSystems; i++ {
		var (
			star       = NewStar(rand.Intn(10))
			numPlanets = rand.Intn(10)
			planets    = make([]Planet, 0, numPlanets)
		)

		for j := 0; j <= numPlanets; j++ {
			planets = append(planets, NewPlanet(rand.Intn(1000), rand.Intn(10)))
		}

		systems = append(systems, *(NewStarSystem([]Star{star}, planets)))
	}

	return &Galaxy{
		Systems: systems,
	}
}

func main() {
	ebiten.SetWindowSize(SCREEN_WIDTH, SCREEN_HEIGHT)
	ebiten.SetWindowTitle("game")

	g := NewGame(WORLD_WIDTH, WORLD_HEIGHT)

	err := ebiten.RunGame(g)
	if err != nil {
		log.Fatal(err)
	}
}
