package main

import (
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Galaxy struct {
	Systems []StarSystem
	shape   int
}

// Draw the galaxy onto the world.
func (g Galaxy) Draw(i *ebiten.Image) {

}

type StarSystem struct {
	Stars   []Star
	Planets []Planet

	connections []StarSystem
	systemView  *ebiten.Image
}

func NewStarSystem(stars []Star, planets []Planet) *StarSystem {

	i := ebiten.NewImage(SCREEN_WIDTH, SCREEN_HEIGHT)

	s := &StarSystem{
		Stars:   stars,
		Planets: planets,
	}

	s.drawSystem(i)

	s.systemView = i

	return s
}

func (s StarSystem) DrawScene(i *ebiten.Image) {
	//i.DrawImage(s.systemView, nil)
	s.drawSystem(i)
}

// Draw the scene for this star system.
func (s StarSystem) drawSystem(i *ebiten.Image) {
	var (
		opts    = ebiten.DrawImageOptions{}
		star    = s.Stars[0]
		centerX = float64(i.Bounds().Dx() / 2)
		centerY = float64(i.Bounds().Dy() / 2)
	)

	opts.GeoM.Translate(centerX-star.Radius, centerY-star.Radius)
	i.DrawImage(star.Image, &opts)

	for _, p := range s.Planets {
		var (
			opts  = ebiten.DrawImageOptions{}
			hypot = float64(star.Radius + float64(p.SemimajorAxis) + float64(p.Radius()))
			r     = float64(p.OrbitPosition)
			x     = math.Cos(r) * hypot
			y     = math.Sin(r) * hypot
		)
		opts.GeoM.Translate(-p.Radius(), -p.Radius())
		opts.GeoM.Rotate(p.Rotation)
		opts.GeoM.Translate(p.Radius(), p.Radius())
		opts.GeoM.Translate(float64(centerX+x-p.Radius()), float64(centerY+y-p.Radius()))
		opts.GeoM.Apply(float64(time.Since(lastFrame)), float64(time.Since(lastFrame)))

		i.DrawImage(p.Image, &opts)
	}
}

type Star struct {
	Size int
	// Radius to the outside of the star.
	Radius float64
	Image  *ebiten.Image
}

func NewStar(size int) Star {
	var (
		// star hotness = radius ?
		stroke = float32(BASE_SYSTEM_STROKE + 10)
		// star size = radius
		radius = float32(BASE_STAR_RADIUS * size)
		l      = int((radius + stroke) * 2)
		i      = ebiten.NewImage(l, l)
	)

	i.Set(0, 0, color.RGBA{B: 255})
	vector.StrokeCircle(i, radius+stroke, radius+stroke, radius, stroke, color.White, true)

	return Star{
		Size:   size,
		Radius: float64(radius + stroke),
		Image:  i,
	}
}

type Planet struct {
	// Distance from it's star
	SemimajorAxis int
	// The position in orbit 0-360
	OrbitPosition float64
	Rotation      float64
	Size          int
	Settlement    *Settlement
	Moons         []Planet
	Spaceports    []SpacePort
	Image         *ebiten.Image
}

func (p Planet) AddSpaceElevator(level int) {

	var (
		r      = float32(p.Radius())
		length = float32(BASE_SPACE_ELEVATOR_LENGTH * level)
	)

	vector.StrokeLine(
		p.Image,
		r*2-BASE_SYSTEM_STROKE,
		r,
		r*2+length,
		r,
		BASE_SYSTEM_STROKE,
		color.White,
		true,
	)
}

// Distance to outside circle.
func (p Planet) Radius() float64 {
	return float64(BASE_SYSTEM_RADIUS*p.Size) + BASE_SYSTEM_STROKE
}

func NewPlanet(axis, size int) Planet {
	var (
		// planet pop = stroke
		stroke = float32(BASE_SYSTEM_STROKE)
		// planet size = radius
		// this is the inner radius
		radius = float32(BASE_SYSTEM_RADIUS * size)
		Radius = radius + stroke
		l      = int((radius + stroke) * 2)
		i      = ebiten.NewImage(l+BASE_PLANET_RADIUS*5, l)
	)

	i.Set(0, 0, color.RGBA{B: 255})
	vector.StrokeCircle(i, Radius, Radius, radius, stroke, color.White, true)

	return Planet{
		SemimajorAxis: axis,
		Size:          size,
		Image:         i,
		OrbitPosition: float64(rand.Intn(360)),
	}
}

type SpacePort struct{}

type Settlement struct {
	SpaceElevatorLevel int
}

// Draw the starsytem on the galaxy map
func (s StarSystem) Draw(i *ebiten.Image) {
	// number of pops = radius.
	// number of planets = stroke
	var (
		stroke = float32(BASE_SYSTEM_STROKE * len(s.Planets))
		radius = float32(BASE_SYSTEM_RADIUS)
	)

	vector.StrokeCircle(i, radius+stroke, radius+stroke, radius, stroke, color.White, true)
}
