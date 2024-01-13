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

func (g *Galaxy) Tick() {
	for i := 0; i < len(g.Systems); i++ {
		g.Systems[i].Tick()
	}
}

type StarSystem struct {
	Stars   []Star
	Planets []Planet

	connections []StarSystem
	systemView  *ebiten.Image
}

// Creates a new star system.
// Currently only supports single star system.
func NewStarSystem(stars []Star, planets []Planet) *StarSystem {

	s := &StarSystem{
		Stars:   stars,
		Planets: planets,
	}

	return s
}

func (s *StarSystem) Tick() {
	for i := 0; i < len(s.Planets); i++ {
		s.Planets[i].Tick(s.Stars[0])
	}
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
		centerX = float64(WORLD_WIDTH / 2)
		centerY = float64(WORLD_HEIGHT / 2)
	)

	opts.GeoM.Translate(centerX-star.Radius, centerY-star.Radius)
	i.DrawImage(star.Image, &opts)

	for _, p := range s.Planets {
		var (
			opts  = ebiten.DrawImageOptions{}
			hypot = float64(star.Radius + float64(p.semimajorAxis) + float64(p.Radius()))
			r     = float64(p.orbitPosition)
			x     = math.Cos(r) * hypot
			y     = math.Sin(r) * hypot
		)

		// rotate the planet
		opts.GeoM.Translate(-p.Radius(), -p.Radius())
		opts.GeoM.Rotate(p.rotation)
		opts.GeoM.Translate(p.Radius(), p.Radius())

		// orbit around the sun
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
		stroke = float32(BASE_SYSTEM_STROKE + 30)
		// star size = radius
		radius = float32(BASE_STAR_RADIUS * size)
		l      = int((radius + stroke) * 2)
		i      = ebiten.NewImage(l, l)
	)

	if DEBUG_MODE {
		i.Set(0, 0, DEBUG_COLOR)
	}

	vector.StrokeCircle(i, radius+stroke, radius+stroke, radius, stroke, color.White, true)

	return Star{
		Size:   size,
		Radius: float64(radius + stroke),
		Image:  i,
	}
}

type Planet struct {
	// Distance from it's star
	semimajorAxis int
	size          int
	// The position in orbit 0-360
	orbitPosition float64
	rotation      float64

	Settlement *Settlement
	Moons      []Planet
	Spaceport  SpacePort
	Image      *ebiten.Image
}

func (p *Planet) Mass() float64 {
	return float64(p.size * 100)
}

func (s *Star) Mass() float64 {
	return math.Pow(float64(s.Size), math.Pow10(10))
}

// OrbitalVeloctiy will cacluate the planet's orbital velocity given the star it is in orbit of.
func (p *Planet) OrbitalVelocity(s Star) float64 {
	return math.Sqrt(G * math.Pow10(10) / (float64(p.semimajorAxis) + p.Radius() + s.Radius))
}

func (p *Planet) Tick(s Star) {
	// TODO calcualte period velocity from mass.
	//p.orbitPosition = math.Mod((p.orbitPosition + (2 * math.Pi * float64(p.semimajorAxis) / 100000)), 360.0)
	p.orbitPosition = p.orbitPosition + p.OrbitalVelocity(s)
	// TODO calculate roation speed from mass.
	p.rotation = p.rotation + 0.1
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

// Distance to outside the circle.
func (p Planet) Radius() float64 {
	return float64(BASE_SYSTEM_RADIUS*p.size) + BASE_SYSTEM_STROKE
}

// NewPlanet creates a new planet given the length of its semimajor axis, and
// the size of the planet.
func NewPlanet(axis, planetSize int) Planet {
	var (
		// planet pop = stroke
		stroke = float32(BASE_SYSTEM_STROKE)
		// planet size = radius
		radius = float32(BASE_PLANET_RADIUS * float64(planetSize))
		l      = int((radius + stroke) * 2)
		i      = ebiten.NewImage(l+MAX_SPACE_ELEVATOR_LENGTH, l)
	)

	if DEBUG_MODE {
		i.Set(0, 0, color.RGBA{B: 255})
	}

	vector.StrokeCircle(i, radius+stroke, radius+stroke, radius, stroke, color.White, true)

	return Planet{
		semimajorAxis: axis,
		size:          planetSize,
		rotation:      0.0,
		Image:         i,
		orbitPosition: float64(rand.Intn(360)),
	}
}

func addrOf[T any](v T) *T {
	return &v
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
