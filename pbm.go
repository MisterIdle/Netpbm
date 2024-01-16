package Netpbm

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	MagicNumberP1 = "P1"
	MagicNumberP4 = "P4"
)

type PBM struct {
	data        [][]bool
	width       int
	height      int
	magicNumber string
}

func ReadPBM(filename string) (*PBM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Error opening the file: %v", err)
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

	scope := strings.Fields(scanner.Text())
	width, _ := strconv.Atoi(scope[0])
	height, _ := strconv.Atoi(scope[1])

	data := make([][]bool, height)
	for i := range data {
		data[i] = make([]bool, width)
	}

	switch magicNumber {
	case MagicNumberP1:
		readP1Data(scanner, data, height, width)
	case MagicNumberP4:
		readP4Data(scanner, data, height, width)
	}

	return &PBM{data: data, width: width, height: height, magicNumber: magicNumber}, nil
}

func readP1Data(scanner *bufio.Scanner, data [][]bool, height, width int) {
	for i := 0; i < height; i++ {
		scanner.Scan()
		line := scanner.Text()
		byteCase := strings.Fields(line)
		for j := 0; j < width; j++ {
			value, _ := strconv.Atoi(byteCase[j])
			data[i][j] = value == 1
		}
	}
}

func readP4Data(scanner *bufio.Scanner, data [][]bool, height, width int) {
	for i := 0; i < height; i++ {
		scanner.Scan()
		line := scanner.Text()
		data[i] = make([]bool, width)

		for j := 0; j < width; j++ {
			byteIndex := j / 8
			bitIndex := 7 - (j % 8)

			if byteIndex < len(line) {
				bit := (line[byteIndex] >> uint(bitIndex)) & 1
				data[i][j] = bit == 1
			} else {
				data[i][j] = false
			}
		}
	}
}

func (pbm *PBM) Size() (int, int) {
	return pbm.height, pbm.width
}

func (pbm *PBM) At(x, y int) bool {
	return pbm.data[y][x]
}

func (pbm *PBM) Set(x, y int, value bool) {
	pbm.data[y][x] = value
}

func (pbm *PBM) Save(filename string) error {
	fileSave, error := os.Create(filename)
	if error != nil {
		return error
	}

	fmt.Fprintf(fileSave, "%s\n%d %d\n", pbm.magicNumber, pbm.width, pbm.height)

	for _, i := range pbm.data {
		for _, j := range i {
			if j {
				fmt.Fprint(fileSave, "1 ")
			} else {
				fmt.Fprint(fileSave, "0 ")
			}
		}
		fmt.Fprintln(fileSave)
	}
	return nil
}

func (pbm *PBM) Invert() {
	for i, row := range pbm.data {
		for j := range row {
			pbm.data[i][j] = !pbm.data[i][j]
		}
	}
}

func (pbm *PBM) Flop() {
	for i := 0; i < pbm.height/2; i++ {
		pbm.data[i], pbm.data[pbm.height-i-1] = pbm.data[pbm.height-i-1], pbm.data[i]
	}
}

func (pbm *PBM) Flip() {
	for _, row := range pbm.data {
		for i, j := 0, len(row)-1; i < j; i, j = i+1, j-1 {
			row[i], row[j] = row[j], row[i]
		}
	}
}

func (pbm *PBM) SetMagicNumber(magicNumber string) {
	pbm.magicNumber = magicNumber
}
