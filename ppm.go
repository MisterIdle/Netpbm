// Code created by: Alexy HOUBLOUP
// Help: "La Table", Chat GPT

package Netpbm

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// Constants representing PPM magic numbers
const (
	MagicNumberP3 = "P3"
	MagicNumberP6 = "P6"
)

// Pixel represents a color pixel in a PPM image.
type PPM struct {
	data          [][]Pixel // 2D slice to store pixel data
	width, height int       // Width and height of the image
	magicNumber   string    // Magic number identifying the PPM format
	max           uint8     // Maximum color value for the image
}

// Variables to store global width, height, and max values
var (
	max           uint8
	width, height int
)

// ReadPPM reads a PPM image from a file.
func ReadPPM(filename string) (*PPM, error) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a buffered reader for efficient reading
	read := bufio.NewReader(file)

	// Read and parse the magic number
	magicNumber, err := read.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading magic number: %v", err)
	}
	magicNumber = strings.TrimSpace(magicNumber)
	if magicNumber != MagicNumberP3 && magicNumber != MagicNumberP6 {
		return nil, fmt.Errorf("invalid magic number: %s", magicNumber)
	}

	// Read and parse dimensions (width and height)
	dim, err := read.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading dimensions: %v", err)
	}
	_, err = fmt.Sscanf(strings.TrimSpace(dim), "%d %d", &width, &height)
	if err != nil {
		return nil, fmt.Errorf("invalid dimensions: %v", err)
	}
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("invalid dimensions: width and height must be positive")
	}

	// Read and parse the maximum color value
	maxValue, err := read.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading max value: %v", err)
	}
	maxValue = strings.TrimSpace(maxValue)
	_, err = fmt.Sscanf(maxValue, "%d", &max)
	if err != nil {
		return nil, fmt.Errorf("invalid max value: %v", err)
	}

	// Initialize a 2D slice to store pixel data
	data := make([][]Pixel, height)
	expectedBytesPerPixel := 3

	// Read pixel data based on the magic number
	if magicNumber == MagicNumberP3 {
		// Read ASCII data for P3 format
		for y := 0; y < height; y++ {
			line, err := read.ReadString('\n')
			if err != nil {
				return nil, fmt.Errorf("error reading data at row %d: %v", y, err)
			}
			fields := strings.Fields(line)
			rowData := make([]Pixel, width)
			for x := 0; x < width; x++ {
				if x*3+2 >= len(fields) {
					return nil, fmt.Errorf("index out of range at row %d, column %d", y, x)
				}
				var pixel Pixel
				_, err := fmt.Sscanf(fields[x*3], "%d", &pixel.R)
				if err != nil {
					return nil, fmt.Errorf("error parsing Red value at row %d, column %d: %v", y, x, err)
				}
				_, err = fmt.Sscanf(fields[x*3+1], "%d", &pixel.G)
				if err != nil {
					return nil, fmt.Errorf("error parsing Green value at row %d, column %d: %v", y, x, err)
				}
				_, err = fmt.Sscanf(fields[x*3+2], "%d", &pixel.B)
				if err != nil {
					return nil, fmt.Errorf("error parsing Blue value at row %d, column %d: %v", y, x, err)
				}
				rowData[x] = pixel
			}
			data[y] = rowData
		}
	} else if magicNumber == MagicNumberP6 {
		// Read binary data for P6 format
		for y := 0; y < height; y++ {
			row := make([]byte, width*expectedBytesPerPixel)
			n, err := read.Read(row)
			if err != nil {
				if err == io.EOF {
					return nil, fmt.Errorf("unexpected end of file at row %d", y)
				}
				return nil, fmt.Errorf("error reading pixel data at row %d: %v", y, err)
			}
			if n < width*expectedBytesPerPixel {
				return nil, fmt.Errorf("unexpected end of file at row %d, expected %d bytes, got %d", y, width*expectedBytesPerPixel, n)
			}

			rowData := make([]Pixel, width)
			for x := 0; x < width; x++ {
				pixel := Pixel{R: row[x*expectedBytesPerPixel], G: row[x*expectedBytesPerPixel+1], B: row[x*expectedBytesPerPixel+2]}
				rowData[x] = pixel
			}
			data[y] = rowData
		}
	}

	// Create and return the PPM object
	return &PPM{data, width, height, magicNumber, uint8(max)}, nil
}

// Size returns the width and height of a PPM image.
func (ppm *PPM) Size() (int, int) {
	return ppm.width, ppm.height
}

// At returns the color of the pixel at coordinates (x, y).
func (ppm *PPM) At(x, y int) Pixel {
	return ppm.data[y][x]
}

// Set sets the color of the pixel at coordinates (x, y).
func (ppm *PPM) Set(x, y int, value Pixel) {
	ppm.data[y][x] = value
}

// Save saves a PPM image to a file.
func (ppm *PPM) Save(filename string) error {
	// Create or overwrite the file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write header information to the file
	if ppm.magicNumber == MagicNumberP3 || ppm.magicNumber == MagicNumberP6 {
		fmt.Fprintf(file, "%s\n%d %d\n%d\n", ppm.magicNumber, ppm.width, ppm.height, ppm.max)
	} else {
		err = fmt.Errorf("magic number error")
		return err
	}

	// Write pixel data to the file
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			pixel := ppm.data[y][x]
			if ppm.magicNumber == MagicNumberP6 {
				file.Write([]byte{pixel.R, pixel.G, pixel.B})
			} else if ppm.magicNumber == MagicNumberP3 {
				fmt.Fprintf(file, "%d %d %d ", pixel.R, pixel.G, pixel.B)
			}
		}
		// Add newline for P3 format
		if ppm.magicNumber == MagicNumberP3 {
			fmt.Fprint(file, "\n")
		}
	}

	return nil
}

// Invert inverts the colors of a PPM image.
func (ppm *PPM) Invert() {
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			// Invert each color channel
			ppm.data[y][x].R = uint8(ppm.max) - ppm.data[y][x].R
			ppm.data[y][x].G = uint8(ppm.max) - ppm.data[y][x].G
			ppm.data[y][x].B = uint8(ppm.max) - ppm.data[y][x].B
		}
	}
}

// Flips a PPM image horizontally.
func (ppm *PPM) Flip() {
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width/2; x++ {
			// Swap pixels horizontally
			ppm.data[y][x], ppm.data[y][ppm.width-x-1] = ppm.data[y][ppm.width-x-1], ppm.data[y][x]
		}
	}
}

// Flops a PPM image vertically.
func (ppm *PPM) Flop() {
	for y := 0; y < ppm.height/2; y++ {
		for x := 0; x < ppm.width; x++ {
			// Swap pixels vertically
			ppm.data[y][x], ppm.data[ppm.height-y-1][x] = ppm.data[ppm.height-y-1][x], ppm.data[y][x]
		}
	}
}

// Sets the magic number of a PPM image.
func (ppm *PPM) SetMagicNumber(magicNumber string) {
	ppm.magicNumber = magicNumber
}

// Sets the maximum value of a PPM image.
func (ppm *PPM) SetMaxValue(maxValue uint8) {
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			// Scale each color channel to the new max value
			ppm.data[y][x].R = uint8(float64(ppm.data[y][x].R) * float64(maxValue) / float64(ppm.max))
			ppm.data[y][x].G = uint8(float64(ppm.data[y][x].G) * float64(maxValue) / float64(ppm.max))
			ppm.data[y][x].B = uint8(float64(ppm.data[y][x].B) * float64(maxValue) / float64(ppm.max))
		}
	}
	ppm.max = maxValue
}

// Rotate90CW rotates a PPM image 90 degrees clockwise.
func (ppm *PPM) Rotate90CW() {
	// Create a new rotated 2D slice
	rotated := make([][]Pixel, ppm.width)
	for i := range rotated {
		rotated[i] = make([]Pixel, ppm.height)
	}

	// Populate the rotated slice with the transposed data
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			rotated[x][ppm.height-y-1] = ppm.data[y][x]
		}
	}

	// Update width, height, and data of the original PPM object
	ppm.width, ppm.height = ppm.height, ppm.width
	ppm.data = rotated
}

// ToPGM converts a PPM image to a PGM image.
func (ppm *PPM) ToPGM() *PGM {
	// Create a new PGM object with the same width, height, and max value
	pgm := &PGM{
		width:       ppm.width,
		height:      ppm.height,
		magicNumber: "P2",
		max:         ppm.max,
	}

	// Initialize a 2D slice for PGM data
	pgm.data = make([][]uint8, ppm.height)
	for i := range pgm.data {
		pgm.data[i] = make([]uint8, ppm.width)
	}

	// Convert each pixel to grayscale and store in the PGM data
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			gray := uint8((int(ppm.data[y][x].R) + int(ppm.data[y][x].G) + int(ppm.data[y][x].B)) / 3)
			pgm.data[y][x] = gray
		}
	}

	return pgm
}

// ToPBM converts a PPM image to a PBM image.
func (ppm *PPM) ToPBM() *PBM {
	// Create a new PBM object with the same width, height, and magic number
	pbm := &PBM{
		width:       ppm.width,
		height:      ppm.height,
		magicNumber: "P1",
	}

	// Initialize a 2D slice for PBM data
	pbm.data = make([][]bool, ppm.height)
	for i := range pbm.data {
		pbm.data[i] = make([]bool, ppm.width)
	}

	// Convert each pixel to binary based on a threshold and store in the PBM data
	threshold := uint8(ppm.max / 2)
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			average := (uint16(ppm.data[y][x].R) + uint16(ppm.data[y][x].G) + uint16(ppm.data[y][x].B)) / 3
			pbm.data[y][x] = average > uint16(threshold)
		}
	}

	return pbm
}
