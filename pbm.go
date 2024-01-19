// Code created by: Alexy HOUBLOUP
// Help: "La Table", Chat GPT

package Netpbm

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	MagicNumberP1 = "P1"
	MagicNumberP4 = "P4"
)

// PBM represents a Netpbm PBM (Portable BitMap) image.
type PBM struct {
	data          [][]bool
	width, height int
	magicNumber   string
}

// ReadPBM reads a PBM image from a file.
func ReadPBM(filename string) (*PBM, error) {

	var width, height int

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	read := bufio.NewReader(file)

	// Read magic number
	magicNumber, err := read.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading magic number: %v", err)
	}
	magicNumber = strings.TrimSpace(magicNumber)
	if magicNumber != MagicNumberP1 && magicNumber != MagicNumberP4 {
		return nil, fmt.Errorf("invalid magic number: %s", magicNumber)
	}

	// Read dim
	dim, err := read.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading dimensions: %v", err)
	}
	_, err = fmt.Sscanf(strings.TrimSpace(dim), "%d %d", &width, &height)
	if err != nil {
		return nil, fmt.Errorf("invalid dimensions: %v", err)
	}

	data := make([][]bool, height)

	for i := range data {
		data[i] = make([]bool, width)
	}

	if magicNumber == MagicNumberP1 {
		// Read binary pixel data
		for y := 0; y < height; y++ {
			line, err := read.ReadString('\n')
			if err != nil {
				return nil, fmt.Errorf("error reading data at row %d: %v", y, err)
			}
			fields := strings.Fields(line)
			for x, field := range fields {
				if x >= width {
					return nil, fmt.Errorf("index out of range at row %d", y)
				}
				data[y][x] = field == "1"
			}
		}

	} else if magicNumber == MagicNumberP4 {
		// Read binary pixel data (compressed)
		expectedBytesPerRow := (width + 7) / 8
		for y := 0; y < height; y++ {
			row := make([]byte, expectedBytesPerRow)
			n, err := read.Read(row)
			if err != nil {
				if err == io.EOF {
					return nil, fmt.Errorf("unexpected end of file at row %d", y)
				}
				return nil, fmt.Errorf("error reading pixel data at row %d: %v", y, err)
			}
			if n < expectedBytesPerRow {
				return nil, fmt.Errorf("unexpected end of file at row %d, expected %d bytes, got %d", y, expectedBytesPerRow, n)
			}

			// Extract individual bits and store in the data matrix
			for x := 0; x < width; x++ {
				byteIndex := x / 8
				bitIndex := 7 - (x % 8)

				decimalValue := int(row[byteIndex])
				bitValue := (decimalValue >> bitIndex) & 1

				data[y][x] = bitValue != 0
			}
		}
	}

	return &PBM{data, width, height, magicNumber}, nil
}

// Size returns the dimensions (width and height) of the PBM image.
func (pbm *PBM) Size() (int, int) {
	return pbm.height, pbm.width
}

// At returns the value of the pixel at the specified coordinates.
func (pbm *PBM) At(x, y int) bool {
	return pbm.data[y][x]
}

// Set sets the value of the pixel at the specified coordinates.
func (pbm *PBM) Set(x, y int, value bool) {
	pbm.data[y][x] = value
}

// Save saves the PBM image to a file.
func (pbm *PBM) Save(filename string) error {
	if pbm == nil {
		return errors.New("cannot save a nil PBM")
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write magic number, width, and height
	fmt.Fprintf(file, "%s\n%d %d\n", pbm.magicNumber, pbm.width, pbm.height)

	// Choose the appropriate method based on the magic number
	switch pbm.magicNumber {
	case MagicNumberP1:
		// Write binary pixel data
		for i := 0; i < pbm.height; i++ {
			for j := 0; j < pbm.width; j++ {
				if pbm.data[i][j] {
					fmt.Fprint(file, "1")
				} else {
					fmt.Fprint(file, "0")
				}

				if j < pbm.width-1 {
					fmt.Fprint(file, " ")
				}
			}
			fmt.Fprintln(file)
		}
	case MagicNumberP4:
		// Write binary pixel data (compressed)
		expectedBytesPerRow := (pbm.width + 7) / 8
		for y := 0; y < pbm.height; y++ {
			row := make([]byte, expectedBytesPerRow)
			for x := 0; x < pbm.width; x++ {
				byteIndex := x / 8
				bitIndex := 7 - (x % 8)
				if pbm.data[y][x] {
					row[byteIndex] |= 1 << bitIndex
				}
			}
			_, err := file.Write(row)
			if err != nil {
				return fmt.Errorf("error writing pixel data at row %d: %v", y, err)
			}
		}
	default:
		return fmt.Errorf("unsupported magic number: %s", pbm.magicNumber)
	}

	return nil
}

// Invert inverts the colors of the PBM image.
func (pbm *PBM) Invert() {
	for i, row := range pbm.data {
		for j := range row {
			pbm.data[i][j] = !pbm.data[i][j]
		}
	}
}

// Flop flips the PBM image horizontally.
func (pbm *PBM) Flop() {
	for i := 0; i < pbm.height/2; i++ {
		pbm.data[i], pbm.data[pbm.height-i-1] = pbm.data[pbm.height-i-1], pbm.data[i]
	}
}

// Flip flips the PBM image vertically.
func (pbm *PBM) Flip() {
	for _, row := range pbm.data {
		for i, j := 0, len(row)-1; i < j; i, j = i+1, j-1 {
			row[i], row[j] = row[j], row[i]
		}
	}
}

// SetMagicNumber sets the magic number of the PBM image.
func (pbm *PBM) SetMagicNumber(magicNumber string) {
	pbm.magicNumber = magicNumber
}
