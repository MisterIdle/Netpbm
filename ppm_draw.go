// Code created by: Alexy HOUBLOUP
// Help: "La Table", Chat GPT

package Netpbm

import (
	"math"
	"sort"
)

// Pixel represents the RGB values of a pixel.
type Pixel struct {
	R, G, B uint8
}

// Point represents a 2D point with X and Y coordinates.
type Point struct {
	X, Y int
}

// NewPPM creates a new PPM image with the specified width, height, magic number, and maximum color value.
func NewPPM(width, height int, magicNumber string, max uint8) *PPM {
	ppm := &PPM{
		width:       width,
		height:      height,
		magicNumber: magicNumber,
		max:         max,
	}

	ppm.data = make([][]Pixel, height)
	for i := range ppm.data {
		ppm.data[i] = make([]Pixel, width)
	}

	return ppm
}

///////////////////
// DRAWING ZONE ///
///////////////////

// DrawLine draws a line between two points on the PPM image.
func (ppm *PPM) DrawLine(p1, p2 Point, color Pixel) {
	dx := abs(p2.X - p1.X)
	dy := abs(p2.Y - p1.Y)
	sx, sy := 1, 1

	if p1.X > p2.X {
		sx = -1
	}
	if p1.Y > p2.Y {
		sy = -1
	}

	err := dx - dy

	for {
		if p1.X >= 0 && p1.X < ppm.width && p1.Y >= 0 && p1.Y < ppm.height {
			ppm.Set(p1.X, p1.Y, color)
		}

		if p1.X == p2.X && p1.Y == p2.Y {
			break
		}

		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			p1.X += sx
		}
		if e2 < dx {
			err += dx
			p1.Y += sy
		}

		if p1.X < 0 || p1.X >= ppm.width || p1.Y < 0 || p1.Y >= ppm.height {
			break
		}
	}
}

/////////////////////
// RECTANGLE ZONE ///
/////////////////////

// DrawRectangle draws a rectangle between two points on the PPM image.
func (ppm *PPM) DrawRectangle(p1 Point, width, height int, color Pixel) {
	// Adjust points to ensure they are within the image boundaries
	if p1.X < 0 {
		p1.X = 0
	}
	if p1.Y < 0 {
		p1.Y = 0
	}

	if p1.X+width > ppm.width {
		width = ppm.width - p1.X
	}
	if p1.Y+height > ppm.height {
		height = ppm.height - p1.Y
	}

	// Draw the rectangle using four DrawLine calls
	ppm.DrawLine(p1, Point{X: p1.X + width, Y: p1.Y}, color)
	ppm.DrawLine(Point{X: p1.X + width, Y: p1.Y}, Point{X: p1.X + width, Y: p1.Y + height}, color)
	ppm.DrawLine(Point{X: p1.X + width, Y: p1.Y + height}, Point{X: p1.X, Y: p1.Y + height}, color)
	ppm.DrawLine(Point{X: p1.X, Y: p1.Y + height}, p1, color)
}

// DrawFilledRectangle draws a filled rectangle between two points on the PPM image.
func (ppm *PPM) DrawFilledRectangle(p1 Point, width, height int, color Pixel) {
	// Adjust points to ensure they are within the image boundaries
	if p1.X < 0 {
		p1.X = 0
	}
	if p1.Y < 0 {
		p1.Y = 0
	}

	if p1.X+width > ppm.width {
		width = ppm.width - p1.X
	}
	if p1.Y+height > ppm.height {
		height = ppm.height - p1.Y
	}

	// Draw the filled rectangle using DrawLine and Set calls
	ppm.DrawLine(p1, Point{X: p1.X + width, Y: p1.Y}, color)
	ppm.DrawLine(Point{X: p1.X + width, Y: p1.Y}, Point{X: p1.X + width, Y: p1.Y + height}, color)
	ppm.DrawLine(Point{X: p1.X + width, Y: p1.Y + height}, Point{X: p1.X, Y: p1.Y + height}, color)
	ppm.DrawLine(Point{X: p1.X, Y: p1.Y + height}, p1, color)

	// Fill the rectangle by setting pixels inside the boundary
	for y := p1.Y + 1; y < p1.Y+height; y++ {
		for x := p1.X + 1; x < p1.X+width; x++ {
			ppm.Set(x, y, color)
		}
	}
}

/////////////////////
//// CIRCLE ZONE ////
/////////////////////

// DrawCircle draws a circle with the specified center and radius on the PPM image.
func (ppm *PPM) DrawCircle(center Point, radius int, color Pixel) {
	// Draw the circle using trigonometric functions
	for i := 0; i < 360; i++ {
		x := int(float64(radius)*math.Cos(float64(i)*math.Pi/180-1)) + center.X
		y := int(float64(radius)*math.Sin(float64(i)*math.Pi/180-1)) + center.Y

		ppm.Set(x, y, color)
	}
}

// DrawFilledCircle draws a filled circle with the specified center and radius on the PPM image.
func (ppm *PPM) DrawFilledCircle(center Point, radius int, color Pixel) {
	// Fill the circle by setting pixels within the circle's boundary
	for y := center.Y - radius; y < center.Y+radius; y++ {
		for x := center.X - radius; x < center.X+radius; x++ {
			if (x-center.X)*(x-center.X)+(y-center.Y)*(y-center.Y) < radius*radius {
				ppm.Set(x, y, color)
			}
		}
	}
}

/////////////////////
// TRIANGLE ZONE ////
/////////////////////

// DrawTriangle draws a triangle between three points on the PPM image.
func (ppm *PPM) DrawTriangle(p1, p2, p3 Point, color Pixel) {
	// Draw the triangle by connecting three lines
	ppm.DrawLine(p1, p2, color)
	ppm.DrawLine(p2, p3, color)
	ppm.DrawLine(p3, p1, color)
}

// DrawFilledTriangle draws a filled triangle between three points on the PPM image.
func (ppm *PPM) DrawFilledTriangle(p1, p2, p3 Point, color Pixel) {
	// Sort vertices based on Y coordinates
	vertices := []Point{p1, p2, p3}
	sort.Slice(vertices, func(i, j int) bool {
		return vertices[i].Y < vertices[j].Y
	})

	// Fill the triangle by interpolating between the top and bottom vertices
	for y := vertices[0].Y; y <= vertices[2].Y; y++ {
		x1 := interpolate(vertices[0], vertices[2], y)
		x2 := interpolate(vertices[1], vertices[2], y)

		ppm.DrawLine(Point{X: int(x1), Y: y}, Point{X: int(x2), Y: y}, color)
	}
}

/////////////////////
// POLYGON ZONE /////
/////////////////////

// DrawPolygon draws a polygon between the specified points on the PPM image.
func (ppm *PPM) DrawPolygon(points []Point, color Pixel) {
	// Draw the polygon outline by connecting consecutive points
	for i := 0; i < len(points)-1; i++ {
		ppm.DrawLine(points[i], points[i+1], color)
	}

	ppm.DrawLine(points[len(points)-1], points[0], color)
}

// DrawFilledPolygon draws a filled polygon between the specified points on the PPM image.
func (ppm *PPM) DrawFilledPolygon(points []Point, color Pixel) {
	// Draw the polygon outline
	ppm.DrawPolygon(points, color)

	// Find the bounding box of the polygon
	minY := points[0].Y
	maxY := points[0].Y
	for _, point := range points {
		if point.Y < minY {
			minY = point.Y
		}
		if point.Y > maxY {
			maxY = point.Y
		}
	}

	// Fill the polygon by scanning lines within the bounding box
	for y := minY + 1; y < maxY; y++ {
		intersectionPoints := []int{}

		// Find intersection points with polygon edges
		for i := 0; i < len(points); i++ {
			p1 := points[i]
			p2 := points[(i+1)%len(points)]

			if (p1.Y <= y && p2.Y > y) || (p2.Y <= y && p1.Y > y) {
				x := int(interpolate(p1, p2, y))
				intersectionPoints = append(intersectionPoints, x)
			}
		}

		// Sort intersection points and fill between pairs
		sort.Ints(intersectionPoints)

		for i := 0; i < len(intersectionPoints)-1; i += 2 {
			start, end := intersectionPoints[i], intersectionPoints[i+1]

			for x := start + 1; x < end; x++ {
				ppm.Set(x, y, color)
			}
		}
	}
}

/////////////////////
// KOCH ZONE ///////
/////////////////////

// DrawKochSnowflake draws a Koch snowflake with the specified center, recursion depth, and size.
// It uses the KochSnowflake function for recursive drawing.

func (ppm *PPM) DrawKochSnowflake(n int, start Point, size int, color Pixel) {
	height := int(math.Sqrt(3) * float64(size) / 2)
	p1 := start
	p2 := Point{X: start.X + size, Y: start.Y}
	p3 := Point{X: start.X + size/2, Y: start.Y + height}

	ppm.KochSnowflake(n, p1, p2, color)
	ppm.KochSnowflake(n, p2, p3, color)
	ppm.KochSnowflake(n, p3, p1, color)
}

func (ppm *PPM) KochSnowflake(n int, p1, p2 Point, color Pixel) {
	if n == 0 {
		ppm.DrawLine(p1, p2, color)
	} else {
		p1Third := Point{
			X: p1.X + (p2.X-p1.X)/3,
			Y: p1.Y + (p2.Y-p1.Y)/3,
		}
		p2Third := Point{
			X: p1.X + 2*(p2.X-p1.X)/3,
			Y: p1.Y + 2*(p2.Y-p1.Y)/3,
		}

		angle := math.Pi / 3
		cosTheta := math.Cos(angle)
		sinTheta := math.Sin(angle)

		p3 := Point{
			X: int(float64(p1Third.X-p2Third.X)*cosTheta-float64(p1Third.Y-p2Third.Y)*sinTheta) + p2Third.X,
			Y: int(float64(p1Third.X-p2Third.X)*sinTheta+float64(p1Third.Y-p2Third.Y)*cosTheta) + p2Third.Y,
		}

		ppm.KochSnowflake(n-1, p1, p1Third, color)
		ppm.KochSnowflake(n-1, p1Third, p3, color)
		ppm.KochSnowflake(n-1, p3, p2Third, color)
		ppm.KochSnowflake(n-1, p2Third, p2, color)
	}
}

/////////////////////
// SIERPINSKI ZONE //
/////////////////////

// DrawSierpinskiTriangle draws a Sierpinski triangle with the specified center, recursion depth, and size.
// It uses the sierpinskiTriangle function for recursive drawing.

func (ppm *PPM) DrawSierpinskiTriangle(n int, start Point, width int, color Pixel) {
	height := int(math.Sqrt(3) * float64(width) / 2)
	p1 := start
	p2 := Point{X: start.X + width, Y: start.Y}
	p3 := Point{X: start.X + width/2, Y: start.Y + height}

	ppm.sierpinskiTriangle(n, p1, p2, p3, color)
}

func (ppm *PPM) sierpinskiTriangle(n int, p1, p2, p3 Point, color Pixel) {
	if n == 0 {
		ppm.DrawFilledTriangle(p1, p2, p3, color)
	} else {
		mid1 := Point{X: (p1.X + p2.X) / 2, Y: (p1.Y + p2.Y) / 2}
		mid2 := Point{X: (p2.X + p3.X) / 2, Y: (p2.Y + p3.Y) / 2}
		mid3 := Point{X: (p3.X + p1.X) / 2, Y: (p3.Y + p1.Y) / 2}

		ppm.sierpinskiTriangle(n-1, p3, mid2, mid3, color)
		ppm.sierpinskiTriangle(n-1, mid2, mid1, p2, color)
		ppm.sierpinskiTriangle(n-1, mid1, p1, mid3, color)
	}
}

/////////////////////
// PERLIN NOISE FUNC //
/////////////////////

// DrawPerlinNoise draws a Perlin noise image using the Perlin noise function and color interpolation.

func (ppm *PPM) DrawPerlinNoise(color1 Pixel, color2 Pixel) {
	frequency := 0.02
	amplitude := 50.0

	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			noiseValue := perlinNoise(float64(x)*frequency, float64(y)*frequency) * amplitude
			normalizedValue := (noiseValue + amplitude) / (2 * amplitude)
			interpolatedColor := interpolateColors(color1, color2, normalizedValue)
			ppm.Set(x, y, interpolatedColor)
		}
	}
}

// perlinNoise generates Perlin noise for a given (x, y) coordinate.
func perlinNoise(x, y float64) float64 {
	n := int(x) + int(y)*57
	n = (n << 13) ^ n
	return (1.0 - ((float64((n*(n*n*15731+789221)+1376312589)&0x7fffffff)/1073741824.0)+1.0)/2.0)
}

/////////////////////
/// UTIL FUNCTIONS //
/////////////////////

// InterpolateColors interpolates between two colors based on a parameter t.
// It returns a new color resulting from the interpolation.

func interpolateColors(color1 Pixel, color2 Pixel, t float64) Pixel {
	r := uint8(float64(color1.R)*(1-t) + float64(color2.R)*t)
	g := uint8(float64(color1.G)*(1-t) + float64(color2.G)*t)
	b := uint8(float64(color1.B)*(1-t) + float64(color2.B)*t)

	return Pixel{R: r, G: g, B: b}
}

// interpolate calculates the interpolation value for a point between two given points.
// It returns the interpolated value for a given y coordinate.

func interpolate(p1, p2 Point, y int) float64 {
	return float64(p1.X) + float64(y-p1.Y)*(float64(p2.X-p1.X)/float64(p2.Y-p1.Y))
}

// abs returns the absolute value of an integer.

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Spectump is a very good game !
// Try it ! (misteridle.itch.io/spectump)
