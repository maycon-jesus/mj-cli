package cube

import (
	"fmt"
)

func main() {
	println("CubeX - Genetic Algorithm for Rubik's Cube")
	println("Estado atual:")
	c := CreateCube333()
	c.Print()

	// scramble := "U F L R U F L R U F L R U F L R U F L R"
	// scramble := "R U L D F B R U L F B D R U F L D B R U "
	scramble := "R' U2 F"
	currentMove := ""

	for _, move := range scramble {
		switch move {
		case ' ':
			moveRaw(&c, currentMove)
			currentMove = ""
		default:
			currentMove += string(move)
		}
	}
	if currentMove != "" {
		moveRaw(&c, currentMove)
	}
}

func moveRaw(c *Cube333, currentMove string) {
	fmt.Printf("Rotação %s\n", currentMove)
	applyMove(c, currentMove)
	c.Print()
	fmt.Println("---------")
}

func applyMove(c *Cube333, notation string) {
	if notation == "" {
		return
	}

	face, ok := faceFromNotation(notation[0])
	if !ok {
		return
	}

	switch notation[1:] {
	case "":
		c.RotateFaceClockwise(face)
	case "'":
		c.RotateFaceCounterClockwise(face)
	case "2":
		c.RotateFace180(face)
	}
}

func faceFromNotation(c byte) (CubeFaceIndex, bool) {
	switch c {
	case 'R':
		return FaceR, true
	case 'L':
		return FaceL, true
	case 'U':
		return FaceU, true
	case 'D':
		return FaceD, true
	case 'F':
		return FaceF, true
	case 'B':
		return FaceB, true
	}
	return 0, false
}
