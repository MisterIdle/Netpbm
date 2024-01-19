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

const (
	MagicNumberP2 = "P2"
	MagicNumberP5 = "P5"
)

// PGM represents a grayscale image in PGM format.
type PGM struct {
	data          [][]uint8
	width, height int
	magicNumber   string
	max           uint8
}

// ReadPGM reads a PGM file and returns a PGM struct.
func ReadPGM(filename string) (*PGM, error) {

	var width, height int

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	read := bufio.NewReader(file)

	magicNumber, err := read.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading magic number: %v", err)
	}
	magicNumber = strings.TrimSpace(magicNumber)
	if magicNumber != MagicNumberP2 && magicNumber != MagicNumberP5 {
		return nil, fmt.Errorf("invalid magic number: %s", magicNumber)
	}

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

	maxValue, err := read.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading max value: %v", err)
	}
	maxValue = strings.TrimSpace(maxValue)
	var max int
	_, err = fmt.Sscanf(maxValue, "%d", &max)
	if err != nil {
		return nil, fmt.Errorf("invalid max value: %v", err)
	}

	data := make([][]uint8, height)
	expectedBytesPerPixel := 1

	if magicNumber == MagicNumberP2 {
		for y := 0; y < height; y++ {
			line, err := read.ReadString('\n')
			if err != nil {
				return nil, fmt.Errorf("error reading data at row %d: %v", y, err)
			}
			fields := strings.Fields(line)
			rowData := make([]uint8, width)
			for x, field := range fields {
				if x >= width {
					return nil, fmt.Errorf("index out of range at row %d", y)
				}
				var pixelValue uint8
				_, err := fmt.Sscanf(field, "%d", &pixelValue)
				if err != nil {
					return nil, fmt.Errorf("error parsing pixel value at row %d, column %d: %v", y, x, err)
				}
				rowData[x] = pixelValue
			}
			data[y] = rowData
		}
	} else if magicNumber == MagicNumberP5 {
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

			rowData := make([]uint8, width)
			for x := 0; x < width; x++ {
				pixelValue := uint8(row[x*expectedBytesPerPixel])
				rowData[x] = pixelValue
			}
			data[y] = rowData
		}
	}

	return &PGM{data, width, height, magicNumber, uint8(max)}, nil
}

// Size returns the dimensions (width and height) of the PGM image.
func (pgm *PGM) Size() (int, int) {
	return pgm.height, pgm.width
}

// At returns the pixel value at the specified coordinates (x, y) in the PGM image.
func (pgm *PGM) At(x, y int) uint8 {
	return pgm.data[y][x]
}

// Set sets the pixel value at the specified coordinates (x, y) in the PGM image.
func (pgm *PGM) Set(x, y int, value uint8) {
	pgm.data[y][x] = value
}

// Save saves the PGM image to the specified file.
func (pgm *PGM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = fmt.Fprintln(writer, pgm.magicNumber)
	if err != nil {
		return fmt.Errorf("error writing magic number: %v", err)
	}

	_, err = fmt.Fprintf(writer, "%d %d\n", pgm.width, pgm.height)
	if err != nil {
		return fmt.Errorf("error writing dimensions: %v", err)
	}

	_, err = fmt.Fprintln(writer, pgm.max)
	if err != nil {
		return fmt.Errorf("error writing max value: %v", err)
	}

	if pgm.magicNumber == MagicNumberP2 {
		for y := 0; y < pgm.height; y++ {
			for x := 0; x < pgm.width; x++ {
				_, err := fmt.Fprint(writer, pgm.data[y][x])
				if err != nil {
					return fmt.Errorf("error writing pixel data at row %d, column %d: %v", y, x, err)
				}

				if x < pgm.width-1 {
					_, err = fmt.Fprint(writer, " ")
					if err != nil {
						return fmt.Errorf("error writing space after pixel at row %d, column %d: %v", y, x, err)
					}
				}
			}
			_, err := fmt.Fprintln(writer)
			if err != nil {
				return fmt.Errorf("error writing newline after row %d: %v", y, err)
			}
		}
	} else if pgm.magicNumber == MagicNumberP5 {
		for y := 0; y < pgm.height; y++ {
			row := make([]byte, pgm.width)
			for x := 0; x < pgm.width; x++ {
				row[x] = byte(pgm.data[y][x])
			}
			_, err := writer.Write(row)
			if err != nil {
				return fmt.Errorf("error writing pixel data at row %d: %v", y, err)
			}
		}
	}

	return writer.Flush()
}

// Invert inverts the colors of the PGM image.
func (pgm *PGM) Invert() {
	for i := range pgm.data {
		for j := range pgm.data[i] {
			pgm.data[i][j] = uint8(pgm.max) - pgm.data[i][j]
		}
	}
}

// Flop flips the PGM image horizontally.
func (pgm *PGM) Flop() {
	for i := 0; i < pgm.height/2; i++ {
		pgm.data[i], pgm.data[pgm.height-i-1] = pgm.data[pgm.height-i-1], pgm.data[i]
	}
}

// Flip flips the PGM image vertically.
func (pgm *PGM) Flip() {
	for i := 0; i < pgm.height; i++ {
		count := pgm.width - 1
		for j := 0; j < pgm.width/2; j++ {
			valTemp := pgm.data[i][j]
			pgm.data[i][j] = pgm.data[i][count]
			pgm.data[i][count] = valTemp
			count--
		}
	}
}

// Sets the magic number of the PGM image.
func (pgm *PGM) SetMagicNumber(magicNumber string) {
	pgm.magicNumber = magicNumber
}

// SetMaxValue sets the maximum value of the PGM image.
func (pgm *PGM) SetMaxValue(maxValue uint8) {
	for y := 0; y < pgm.height; y++ {
		for x := 0; x < pgm.width; x++ {
			scaledValue := float64(pgm.data[y][x]) * float64(maxValue) / float64(pgm.max)
			newValue := uint8(scaledValue)
			pgm.data[y][x] = newValue
		}
	}

	// Update the max value
	pgm.max = maxValue
}

// Rotate90CW rotates the PGM image 90 degrees clockwise.
func (pgm *PGM) Rotate90CW() {

	rotateData := make([][]uint8, pgm.width)
	for i := range rotateData {
		rotateData[i] = make([]uint8, pgm.height)
	}

	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			d := pgm.height - j - 1
			rotateData[i][d] = pgm.data[j][i]
		}
	}

	pgm.width, pgm.height = pgm.height, pgm.width
	pgm.data = rotateData
}

func (pgm *PGM) ToPBM() *PBM {
	pbm := &PBM{
		width:       pgm.width,
		height:      pgm.height,
		magicNumber: "P1",
	}

	pbm.data = make([][]bool, pgm.height)
	for i := range pbm.data {
		pbm.data[i] = make([]bool, pgm.width)
	}

	threshold := uint8(pgm.max / 2)
	for y := 0; y < pgm.height; y++ {
		for x := 0; x < pgm.width; x++ {
			pbm.data[y][x] = pgm.data[y][x] > threshold
		}
	}

	return pbm
}
