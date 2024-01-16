package Netpbm

import (
	"fmt"
	"log"
)

func TestAll(pbm, pgm, ppm, to, draw bool) {
	if pbm {
		TestPBM()
	}
	if pgm {
		TestPGM()
	}
	if ppm {
		TestPPM()
	}
	if to {
		TestTo()
	}
	if draw {
		Drawing()
	}
}

func TestPBM() {
	filename := "./Debug/PBM/debug.pbm"

	pbm, err := ReadPBM(filename)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(" ")
	fmt.Println("--------------------")
	fmt.Println("Test PBM")
	fmt.Println("--------------------")

	fmt.Print("Flip !" + "\n")
	pbm.Flip()
	pbm.Save("./Debug/PBM/debugFlip.pbm")

	fmt.Print("Flop !" + "\n")
	pbm.Flop()
	pbm.Save("./Debug/PBM/debugFlop.pbm")

	fmt.Print("Inversion !" + "\n")
	pbm.Invert()
	pbm.Save("./Debug/PBM/debugInvert.pbm")

	fmt.Println("--------------------")
	fmt.Println(" ")

}

func TestPGM() {
	filename := "./Debug/PGM/debug.pgm"

	pgm, err := ReadPGM(filename)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(" ")
	fmt.Println("--------------------")
	fmt.Println("Test PGM")
	fmt.Println("--------------------")

	fmt.Print("Flip !" + "\n")
	pgm.Flip()
	pgm.Save("./Debug/PGM/debugFlip.pgm")

	fmt.Print("Flop !" + "\n")
	pgm.Flop()
	pgm.Save("./Debug/PGM/debugFlop.pgm")

	fmt.Print("Rotation !" + "\n")
	pgm.Rotate90CW()
	pgm.Save("./Debug/PGM/debugRotate.pgm")

	fmt.Print("Inversion !" + "\n")
	pgm.Invert()
	pgm.Save("./Debug/PGM/debugInvert.pgm")

	fmt.Println("--------------------")
	fmt.Println(" ")
}

func TestPPM() {
	filename := "./Debug/PPM/debug.ppm"

	ppm, err := ReadPPM(filename)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(" ")
	fmt.Println("--------------------")
	fmt.Println("Test PPM")
	fmt.Println("--------------------")

	fmt.Print("Flip !" + "\n")
	ppm.Flip()
	ppm.Save("./Debug/PPM/debugFlip.ppm")

	fmt.Print("Flop !" + "\n")
	ppm.Flop()
	ppm.Save("./Debug/PPM/debugFlop.ppm")

	fmt.Print("Rotation !" + "\n")
	ppm.Rotate90CW()
	ppm.Save("./Debug/PPM/debugRotate.ppm")

	fmt.Print("Inversion !" + "\n")
	ppm.Invert()
	ppm.Save("./Debug/PPM/debugInvert.ppm")

	fmt.Println("--------------------")
	fmt.Println(" ")
}

func TestTo() {

	fmt.Println(" ")
	fmt.Println("--------------------")
	fmt.Println("Test To")
	fmt.Println("--------------------")

	filename := "./Debug/PGM/debug.pgm"

	pgm, err := ReadPGM(filename)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("To PBM !" + "\n")
	pgm.ToPBM().Save("./Debug/To/debugPGMToPBM.pbm")

	filename = "./Debug/PPM/debug.ppm"
	ppm, err := ReadPPM(filename)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("To PGM !" + "\n")
	ppm.ToPGM().Save("./Debug/To/debugPPMToPGM.pgm")

	fmt.Print("To PBM !" + "\n")
	ppm.ToPBM().Save("./Debug/To/debugPPMToPBM.pbm")

	fmt.Println("--------------------")
	fmt.Println(" ")
}

func Drawing() {
	fmt.Print("Line !" + "\n")
	PPMLine := NewPPM(100, 100, "P3", 255)
	PPMLine.DrawLine(Point{X: 0, Y: 0}, Point{X: 99, Y: 99}, Pixel{R: 255, G: 0, B: 0})
	PPMLine.Save("./Debug/Draw/PPMLine.ppm")

	fmt.Print("Rectangle !" + "\n")
	PPMRectangle := NewPPM(100, 100, "P3", 255)
	PPMRectangle.DrawRectangle(Point{X: 0, Y: 0}, 50, 10, Pixel{R: 255, G: 0, B: 0})
	PPMRectangle.Save("./Debug/Draw/PPMRectangle.ppm")

	fmt.Print("Filled Rectangle !" + "\n")
	PPMRectangleFilled := NewPPM(100, 100, "P3", 255)
	PPMRectangleFilled.DrawFilledRectangle(Point{X: 0, Y: 0}, 50, 10, Pixel{R: 255, G: 0, B: 0})
	PPMRectangleFilled.Save("./Debug/Draw/PPMRectangleFilled.ppm")

	fmt.Print("Circle !" + "\n")
	PPMCircle := NewPPM(100, 100, "P3", 255)
	PPMCircle.DrawCircle(Point{X: 50, Y: 50}, 20, Pixel{R: 255, G: 0, B: 0})
	PPMCircle.Save("./Debug/Draw/PPMCircle.ppm")

	fmt.Print("Filled Circle !" + "\n")
	PPMCircleFilled := NewPPM(100, 100, "P3", 255)
	PPMCircleFilled.DrawFilledCircle(Point{X: 50, Y: 50}, 20, Pixel{R: 255, G: 0, B: 0})
	PPMCircleFilled.Save("./Debug/Draw/PPMCircleFilled.ppm")

	fmt.Print("Triangle !" + "\n")
	PPMTriangle := NewPPM(100, 100, "P3", 255)
	PPMTriangle.DrawTriangle(Point{X: 0, Y: 0}, Point{X: 99, Y: 99}, Point{X: 0, Y: 99}, Pixel{R: 255, G: 0, B: 0})
	PPMTriangle.Save("./Debug/Draw/PPMTriangle.ppm")

	//Not work !
	//fmt.Print("Filled Triangle !" + "\n")
	//PPMTriangleFilled := NewPPM(200, 200, "P3", 255)
	//PPMTriangleFilled.DrawFilledTriangle(Point{X: 0, Y: 0}, Point{X: 99, Y: 99}, Point{X: 0, Y: 99}, Pixel{R: 255, G: 0, B: 0})
	//PPMTriangleFilled.Save("./Debug/Draw/PPMTriangleFilled.ppm")

	fmt.Print("Koch Snowflake !" + "\n")
	PPMKoch := NewPPM(800, 800, "P3", 255)
	PPMKoch.DrawKochSnowflake(3, Point{X: 100, Y: 200}, 300, Pixel{R: 255, G: 0, B: 0})
	PPMKoch.Save("./Debug/Draw/PPMKoch.ppm")

	fmt.Print("sierpinskiTriangle !" + "\n")
	PPMSierpinski := NewPPM(800, 800, "P3", 255)
	PPMSierpinski.DrawSierpinskiTriangle(3, Point{X: 100, Y: 200}, 300, Pixel{R: 255, G: 0, B: 0})
	PPMSierpinski.Flop()
	PPMSierpinski.Save("./Debug/Draw/PPMSierpinski.ppm")

	fmt.Print("Perlin Noise !" + "\n")
	PPMPerlin := NewPPM(1000, 1000, "P3", 255)
	PPMPerlin.DrawPerlinNoise(Pixel{R: 255, G: 255, B: 255}, Pixel{R: 0, G: 0, B: 0})
	PPMPerlin.Save("./Debug/Draw/PPMPerlin.ppm")
}
