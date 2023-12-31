package main

import (
	"fmt"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/math/f64"
)

const (
	ZOOM_RATIO                 = 10
	CAMERA_MOVE_RATIO          = 10
	SCREEN_WIDTH               = 1920
	SCREEN_HEIGHT              = 1080
	CAMERA_MOVE_SPEED          = 10
	BASE_RADIUS                = 10.0
	BASE_SYSTEM_RADIUS         = BASE_RADIUS
	BASE_SYSTEM_STROKE         = BASE_RADIUS
	BASE_PLANET_RADIUS         = BASE_RADIUS
	BASE_STAR_RADIUS           = 20.0
	BASE_SPACE_ELEVATOR_LENGTH = 5.0
)

type Game struct {
	camera Camera
	world  *ebiten.Image

	Galaxy Galaxy
}

// todo remove
var g Galaxy

func NewGame(worldWidth, worldHeight int) *Game {
	var (
		world  = ebiten.NewImage(worldWidth, worldHeight)
		camera = Camera{ViewPort: f64.Vec2{SCREEN_WIDTH, SCREEN_HEIGHT}}
	)

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
		g.camera.Position[0] += CAMERA_MOVE_SPEED
	case ebiten.IsKeyPressed(ebiten.KeyArrowRight):
		g.camera.Position[0] -= CAMERA_MOVE_SPEED
	case ebiten.IsKeyPressed(ebiten.KeyArrowUp):
		g.camera.Position[1] += CAMERA_MOVE_SPEED
	case ebiten.IsKeyPressed(ebiten.KeyArrowDown):
		g.camera.Position[1] -= CAMERA_MOVE_SPEED
	case ebiten.IsKeyPressed(ebiten.KeySpace):
		g.camera.Reset()
	}

	// mouse button todo
	//if ebiten.IsMouseButtonPressed(ebiten.MouseButton3) {

	g.Galaxy.Systems[0].Planets[0].OrbitPosition += 0.01
	g.Galaxy.Systems[0].Planets[1].OrbitPosition += 0.01
	g.Galaxy.Systems[0].Planets[1].Rotation += 0.1

	return nil
}

var lastFrame time.Time = time.Now()

func (g *Game) Draw(screen *ebiten.Image) {

	g.world.Clear()
	g.Galaxy.Systems[0].DrawScene(g.world)

	g.camera.Render(g.world, screen)

	_, _ = g.camera.ScreenToWorld(ebiten.CursorPosition())

	ebitenutil.DebugPrint(
		screen,
		fmt.Sprintf(
			"FPS: %f\nKey Press: %v",
			ebiten.ActualFPS(),
			ebiten.InputChars(),
		),
	)
	lastFrame = time.Now()
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return SCREEN_WIDTH, SCREEN_HEIGHT
}
func init() {

	star := NewStar(5)
	p1 := NewPlanet(50, 2)
	p2 := NewPlanet(100, 5)
	p2.AddSpaceElevator(5)
	s := NewStarSystem([]Star{star}, []Planet{p1, p2})

	g = Galaxy{
		Systems: []StarSystem{
			*s,
		},
	}

}

func main() {
	ebiten.SetWindowSize(1920, 1080)
	ebiten.SetWindowTitle("game")

	g := NewGame(1920, 1080)

	err := ebiten.RunGame(g)
	if err != nil {
		log.Fatal(err)
	}
}
