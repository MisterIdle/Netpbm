package Netpbm

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

const (
	MagicNumberP2 = "P2"
	MagicNumberP5 = "P5"
)

type PGM struct {
	data          [][]uint8
	width, height int
	magicNumber   string
	max           uint8
}

func ReadPGM(filename string) (*PGM, error) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening the file")
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Scan()
	magicNumber := scanner.Text()

	for scanner.Scan() {
		if len(scanner.Text()) > 0 && scanner.Text()[0] == '#' {
			continue
		}
		break
	}

	scope := strings.Split(scanner.Text(), " ")
	width, _ := strconv.Atoi(scope[0])
	height, _ := strconv.Atoi(scope[1])

	scanner.Scan()

	maxValue, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Println("Error reading max value")
	}

	max := uint8(maxValue)

	data := make([][]uint8, height)
	for i := range data {
		data[i] = make([]uint8, width)
	}

	switch magicNumber {
	case MagicNumberP2:
		readP2Data(scanner, data, height, width)
	case MagicNumberP5:
		readP5Data(scanner, data, height, width)
	}

	return &PGM{data: data, width: width, height: height, magicNumber: magicNumber, max: max}, nil
}

func readP2Data(scanner *bufio.Scanner, data [][]uint8, height, width int) {
	for i := 0; i < height; i++ {
		scanner.Scan()
		line := scanner.Text()
		byteCase := strings.Fields(line)

		if len(byteCase) < width {
			break
		}

		for j := 0; j < width; j++ {
			value, _ := strconv.Atoi(byteCase[j])
			data[i][j] = uint8(value)
		}
	}
}

func readP5Data(scanner *bufio.Scanner, data [][]uint8, height, width int) {
	//Pas eu le temps de le faire
}

func (pgm *PGM) Size() (int, int) {
	return pgm.height, pgm.width
}

func (pgm *PGM) At(x, y int) uint8 {
	return pgm.data[y][x]
}

func (pgm *PGM) Set(x, y int, value uint8) {
	pgm.data[y][x] = value
}

func (pgm *PGM) Save(filename string) error {
	fileSave, error := os.Create(filename)
	if error != nil {
		return error
	}

	fmt.Fprintf(fileSave, "%s\n%d %d\n %d\n", pgm.magicNumber, pgm.width, pgm.height, pgm.max)

	for i := range pgm.data {
		for j := range pgm.data[i] {
			fmt.Fprintf(fileSave, "%d ", pgm.data[i][j])
		}
		fmt.Fprintln(fileSave)
	}
	return nil
}

func (pgm *PGM) Invert() {
	for i := range pgm.data {
		for j := range pgm.data[i] {
			pgm.data[i][j] = uint8(pgm.max) - pgm.data[i][j]
		}
	}
}

func (pgm *PGM) Flop() {
	for i := 0; i < pgm.height/2; i++ {
		pgm.data[i], pgm.data[pgm.height-i-1] = pgm.data[pgm.height-i-1], pgm.data[i]
	}
}

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

func (pgm *PGM) SetMagicNumber(magicNumber string) {
	pgm.magicNumber = magicNumber
}

func (pgm *PGM) SetMaxValue(maxValue uint8) {
	if maxValue >= 1 && maxValue <= 255 {
		pgm.max = maxValue

		for i := 0; i < pgm.height; i++ {
			for j := 0; j < pgm.width; j++ {
				pgm.data[i][j] = uint8(math.Round(float64(pgm.data[i][j]) / float64(pgm.max) * 255))
			}
		}
	} else {
		fmt.Println("Veuillez mettre une valeur comprise entre 1 et 255")
	}
}

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
