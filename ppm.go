package Netpbm

// Ce code est très différent par rapport à PGM & PBM car il n'a pas était testé par code de mentor

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type PPM struct {
	data          [][]Pixel
	width, height int
	magicNumber   string
	max           uint8
}

type Pixel struct {
	R, G, B uint8
}

type Point struct {
	X, Y int
}

/////////////////////
// DEBUG FUNCTIONS //
/////////////////////

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

//////////////////////
// NORMAL FUNCTIONS //
//////////////////////

func ReadPPM(filename string) (*PPM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	skipLine := func(line string) bool {
		return strings.HasPrefix(line, "#")
	}

	for scanner.Scan() {
		if !skipLine(scanner.Text()) {
			break
		}
	}

	magicNumber := scanner.Text()

	for scanner.Scan() {
		if !skipLine(scanner.Text()) {
			break
		}
	}
	sizeLine := scanner.Text()

	size := strings.Fields(sizeLine)
	width, err := strconv.Atoi(size[0])
	if err != nil {
		return nil, err
	}
	height, err := strconv.Atoi(size[1])
	if err != nil {
		return nil, err
	}

	for scanner.Scan() {
		if !skipLine(scanner.Text()) {
			break
		}
	}
	max, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return nil, err
	}

	ppm := &PPM{
		width:       width,
		height:      height,
		magicNumber: magicNumber,
		max:         uint8(max),
	}

	ppm.data = make([][]Pixel, height)
	for i := range ppm.data {
		ppm.data[i] = make([]Pixel, width)
	}

	if magicNumber == "P3" {
		for i := 0; i < height; i++ {
			for scanner.Scan() {
				if !skipLine(scanner.Text()) {
					break
				}
			}
			line := scanner.Text()

			tokens := strings.Fields(line)

			for j := 0; j < len(tokens); j += 3 {
				r, _ := strconv.Atoi(tokens[j])
				g, _ := strconv.Atoi(tokens[j+1])
				b, _ := strconv.Atoi(tokens[j+2])

				ppm.data[i][j/3] = Pixel{R: uint8(r), G: uint8(g), B: uint8(b)}
			}
		}
	} else if magicNumber == "P6" {
		dataSize := width * height * 3
		data := make([]byte, dataSize)
		_, err := file.Read(data)
		if err != nil {
			return nil, err
		}

		index := 0
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				ppm.data[y][x] = Pixel{R: data[index], G: data[index+1], B: data[index+2]}
				index += 3
			}
		}
	}

	return ppm, nil
}

func (ppm *PPM) Size() (int, int) {
	return ppm.width, ppm.height
}

func (ppm *PPM) At(x, y int) Pixel {
	return ppm.data[y][x]
}

func (ppm *PPM) Set(x, y int, value Pixel) {
	ppm.data[y][x] = value
}

func (ppm *PPM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	_, err = writer.WriteString(ppm.magicNumber + "\n")
	if err != nil {
		return err
	}

	_, err = writer.WriteString(strconv.Itoa(ppm.width) + " " + strconv.Itoa(ppm.height) + "\n")
	if err != nil {
		return err
	}

	_, err = writer.WriteString(strconv.Itoa(int(ppm.max)) + "\n")
	if err != nil {
		return err
	}

	if ppm.magicNumber == "P6" {
		for y := 0; y < ppm.height; y++ {
			line := ""
			for x := 0; x < ppm.width; x++ {
				line += fmt.Sprintf("0x%02X 0x%02X 0x%02X ", ppm.data[y][x].R, ppm.data[y][x].G, ppm.data[y][x].B)
			}
			_, err = writer.WriteString(line + "\n")
			if err != nil {
				return err
			}
		}
	} else {
		for y := 0; y < ppm.height; y++ {
			line := ""
			for x := 0; x < ppm.width; x++ {
				line += strconv.Itoa(int(ppm.data[y][x].R)) + " " +
					strconv.Itoa(int(ppm.data[y][x].G)) + " " +
					strconv.Itoa(int(ppm.data[y][x].B)) + " "
			}
			_, err = writer.WriteString(line + "\n")
			if err != nil {
				return err
			}
		}
	}

	err = writer.Flush()
	if err != nil {
		return err
	}

	return nil
}

func (ppm *PPM) Invert() {
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			ppm.data[y][x].R = uint8(ppm.max) - ppm.data[y][x].R
			ppm.data[y][x].G = uint8(ppm.max) - ppm.data[y][x].G
			ppm.data[y][x].B = uint8(ppm.max) - ppm.data[y][x].B
		}
	}
}

func (ppm *PPM) Flip() {
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width/2; x++ {
			ppm.data[y][x], ppm.data[y][ppm.width-x-1] = ppm.data[y][ppm.width-x-1], ppm.data[y][x]
		}
	}
}

func (ppm *PPM) Flop() {
	for y := 0; y < ppm.height/2; y++ {
		for x := 0; x < ppm.width; x++ {
			ppm.data[y][x], ppm.data[ppm.height-y-1][x] = ppm.data[ppm.height-y-1][x], ppm.data[y][x]
		}
	}
}

func (ppm *PPM) SetMagicNumber(magicNumber string) {
	ppm.magicNumber = magicNumber
}

func (ppm *PPM) SetMaxValue(maxValue uint8) {
	ppm.max = maxValue
}

func (ppm *PPM) Rotate90CW() {
	rotated := make([][]Pixel, ppm.width)
	for i := range rotated {
		rotated[i] = make([]Pixel, ppm.height)
	}

	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			rotated[x][ppm.height-y-1] = ppm.data[y][x]
		}
	}

	ppm.width, ppm.height = ppm.height, ppm.width
	ppm.data = rotated
}

func (ppm *PPM) ToPGM() *PGM {
	pgm := &PGM{
		width:       ppm.width,
		height:      ppm.height,
		magicNumber: "P2",
		max:         ppm.max,
	}

	pgm.data = make([][]uint8, ppm.height)
	for i := range pgm.data {
		pgm.data[i] = make([]uint8, ppm.width)
	}

	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			gray := uint8(0.299*float64(ppm.data[y][x].R) + 0.587*float64(ppm.data[y][x].G) + 0.114*float64(ppm.data[y][x].B))
			pgm.data[y][x] = gray
		}
	}

	return pgm
}

func (ppm *PPM) ToPBM() *PBM {
	pbm := &PBM{
		width:       ppm.width,
		height:      ppm.height,
		magicNumber: "P1",
	}

	pbm.data = make([][]bool, ppm.height)
	for i := range pbm.data {
		pbm.data[i] = make([]bool, ppm.width)
	}

	threshold := uint8(ppm.max / 2)
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			average := (uint16(ppm.data[y][x].R) + uint16(ppm.data[y][x].G) + uint16(ppm.data[y][x].B)) / 3
			pbm.data[y][x] = average > uint16(threshold)
		}
	}

	return pbm
}

///////////////////
// DRAWING ZONE ///
///////////////////

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
		// Check if the current point is within the image boundaries
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

		// Additional check to break the loop if the point goes out of bounds
		if p1.X < 0 || p1.X >= ppm.width || p1.Y < 0 || p1.Y >= ppm.height {
			break
		}
	}
}

/////////////////////
// RECTANGLE ZONE ///
/////////////////////

func (ppm *PPM) DrawRectangle(p1 Point, width, height int, color Pixel) {
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

	ppm.DrawLine(p1, Point{X: p1.X + width, Y: p1.Y}, color)
	ppm.DrawLine(Point{X: p1.X + width, Y: p1.Y}, Point{X: p1.X + width, Y: p1.Y + height}, color)
	ppm.DrawLine(Point{X: p1.X + width, Y: p1.Y + height}, Point{X: p1.X, Y: p1.Y + height}, color)
	ppm.DrawLine(Point{X: p1.X, Y: p1.Y + height}, p1, color)
}

func (ppm *PPM) DrawFilledRectangle(p1 Point, width, height int, color Pixel) {
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

	ppm.DrawLine(p1, Point{X: p1.X + width, Y: p1.Y}, color)
	ppm.DrawLine(Point{X: p1.X + width, Y: p1.Y}, Point{X: p1.X + width, Y: p1.Y + height}, color)
	ppm.DrawLine(Point{X: p1.X + width, Y: p1.Y + height}, Point{X: p1.X, Y: p1.Y + height}, color)
	ppm.DrawLine(Point{X: p1.X, Y: p1.Y + height}, p1, color)

	for y := p1.Y + 1; y < p1.Y+height; y++ {
		for x := p1.X + 1; x < p1.X+width; x++ {
			ppm.Set(x, y, color)
		}
	}
}

/////////////////////
//// CIRCLE ZONE ////
/////////////////////

func (ppm *PPM) DrawCircle(center Point, radius int, color Pixel) {

	// Avec cos et sin seulement le contour du cercle est dessiné
	for i := 0; i < 360; i++ {
		x := int(float64(radius)*math.Cos(float64(i)*math.Pi/180-1)) + center.X
		y := int(float64(radius)*math.Sin(float64(i)*math.Pi/180-1)) + center.Y

		ppm.Set(x, y, color)
	}
}

func (ppm *PPM) DrawFilledCircle(center Point, radius int, color Pixel) {
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

func (ppm *PPM) DrawTriangle(p1, p2, p3 Point, color Pixel) {
	ppm.DrawLine(p1, p2, color)
	ppm.DrawLine(p2, p3, color)
	ppm.DrawLine(p3, p1, color)
}

func (ppm *PPM) DrawFilledTriangle(p1, p2, p3 Point, color Pixel) {
	vertices := []Point{p1, p2, p3}
	sort.Slice(vertices, func(i, j int) bool {
		return vertices[i].Y < vertices[j].Y
	})

	for y := vertices[0].Y; y <= vertices[2].Y; y++ {
		x1 := interpolate(vertices[0], vertices[2], y)
		x2 := interpolate(vertices[1], vertices[2], y)

		ppm.DrawLine(Point{X: int(x1), Y: y}, Point{X: int(x2), Y: y}, color)
	}
}

/////////////////////
// POLYGON ZONE /////
/////////////////////

func (ppm *PPM) DrawPolygon(points []Point, color Pixel) {
	for i := 0; i < len(points)-1; i++ {
		ppm.DrawLine(points[i], points[i+1], color)
	}

	ppm.DrawLine(points[len(points)-1], points[0], color)
}

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

	for y := minY + 1; y < maxY; y++ {
		intersectionPoints := []int{}

		for i := 0; i < len(points); i++ {
			p1 := points[i]
			p2 := points[(i+1)%len(points)]

			if (p1.Y <= y && p2.Y > y) || (p2.Y <= y && p1.Y > y) {
				x := int(interpolate(p1, p2, y))
				intersectionPoints = append(intersectionPoints, x)
			}
		}

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

func perlinNoise(x, y float64) float64 {
	n := int(x) + int(y)*57
	n = (n << 13) ^ n
	return (1.0 - ((float64((n*(n*n*15731+789221)+1376312589)&0x7fffffff)/1073741824.0)+1.0)/2.0)
}

/////////////////////
/// UTIL FUNCTIONS //
/////////////////////

func interpolateColors(color1 Pixel, color2 Pixel, t float64) Pixel {
	r := uint8(float64(color1.R)*(1-t) + float64(color2.R)*t)
	g := uint8(float64(color1.G)*(1-t) + float64(color2.G)*t)
	b := uint8(float64(color1.B)*(1-t) + float64(color2.B)*t)

	return Pixel{R: r, G: g, B: b}
}

func interpolate(p1, p2 Point, y int) float64 {
	return float64(p1.X) + float64(y-p1.Y)*(float64(p2.X-p1.X)/float64(p2.Y-p1.Y))
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
