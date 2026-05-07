package cube

import "fmt"

type Cube333 struct {
	Faces [6]*CubeFace
}

func CreateCube333() Cube333 {
	cube := Cube333{}

	cube.Faces[FaceL] = newCubeFace(Orange, 3)
	cube.Faces[FaceU] = newCubeFace(White, 3)
	cube.Faces[FaceF] = newCubeFace(Green, 3)
	cube.Faces[FaceD] = newCubeFace(Yellow, 3)
	cube.Faces[FaceR] = newCubeFace(Red, 3)
	cube.Faces[FaceB] = newCubeFace(Blue, 3)

	// cube.Faces[FaceL].Top = cube.Faces[FaceU]
	// cube.Faces[FaceL].Right = cube.Faces[FaceF]
	// cube.Faces[FaceL].Bottom = cube.Faces[FaceD]
	// cube.Faces[FaceL].Left = cube.Faces[FaceB]
	// cube.Faces[FaceL].Back = cube.Faces[FaceR]

	// cube.Faces[FaceU].Top = cube.Faces[FaceB]
	// cube.Faces[FaceU].Right = cube.Faces[FaceR]
	// cube.Faces[FaceU].Bottom = cube.Faces[FaceF]
	// cube.Faces[FaceU].Left = cube.Faces[FaceL]
	// cube.Faces[FaceU].Back = cube.Faces[FaceD]

	// cube.Faces[FaceF].Top = cube.Faces[FaceU]
	// cube.Faces[FaceF].Right = cube.Faces[FaceR]
	// cube.Faces[FaceF].Bottom = cube.Faces[FaceD]
	// cube.Faces[FaceF].Left = cube.Faces[FaceL]
	// cube.Faces[FaceF].Back = cube.Faces[FaceB]

	// cube.Faces[FaceD].Top = cube.Faces[FaceF]
	// cube.Faces[FaceD].Right = cube.Faces[FaceR]
	// cube.Faces[FaceD].Bottom = cube.Faces[FaceB]
	// cube.Faces[FaceD].Left = cube.Faces[FaceL]
	// cube.Faces[FaceD].Back = cube.Faces[FaceU]

	// cube.Faces[FaceR].Top = cube.Faces[FaceU]
	// cube.Faces[FaceR].Right = cube.Faces[FaceB]
	// cube.Faces[FaceR].Bottom = cube.Faces[FaceD]
	// cube.Faces[FaceR].Left = cube.Faces[FaceF]
	// cube.Faces[FaceR].Back = cube.Faces[FaceL]

	// cube.Faces[FaceB].Top = cube.Faces[FaceU]
	// cube.Faces[FaceB].Right = cube.Faces[FaceL]
	// cube.Faces[FaceB].Bottom = cube.Faces[FaceD]
	// cube.Faces[FaceB].Left = cube.Faces[FaceR]
	// cube.Faces[FaceB].Back = cube.Faces[FaceF]

	return cube
}

// Layout das peças em cada face:
//	0 1 2
//	3 4 5
//	6 7 8

// Índices para rotação no sentido horário
// após uma rotação de 90° no sentido horário.
var clockwiseRotation = [9]int{6, 3, 0, 7, 4, 1, 8, 5, 2}

func (c *Cube333) RotateFaceCounterClockwise(faceId CubeFaceIndex) {
	c.RotateFaceClockwise(faceId)
	c.RotateFaceClockwise(faceId)
	c.RotateFaceClockwise(faceId)
}

func (c *Cube333) RotateFace180(faceId CubeFaceIndex) {
	c.RotateFaceClockwise(faceId)
	c.RotateFaceClockwise(faceId)
}

func (c *Cube333) RotateFaceClockwise(faceId CubeFaceIndex) {
	face := c.Faces[faceId]
	src := append([]CubeColor{}, face.Pieces...)
	for i := 0; i < 9; i++ {
		face.Pieces[i] = src[clockwiseRotation[i]]
	}

	switch faceId {
	case FaceU:
		moves := MoveSet{
			{FaceL, []int{1, 2, 3}, FaceB, []int{1, 2, 3}},
			{FaceB, []int{1, 2, 3}, FaceR, []int{1, 2, 3}},
			{FaceR, []int{1, 2, 3}, FaceF, []int{1, 2, 3}},
			{FaceF, []int{1, 2, 3}, FaceL, []int{1, 2, 3}},
		}
		ResolveMoveSet(c.Faces, moves)
	case FaceD:
		moves := MoveSet{
			{FaceL, []int{7, 8, 9}, FaceF, []int{7, 8, 9}},
			{FaceF, []int{7, 8, 9}, FaceR, []int{7, 8, 9}},
			{FaceR, []int{7, 8, 9}, FaceB, []int{7, 8, 9}},
			{FaceB, []int{7, 8, 9}, FaceL, []int{7, 8, 9}},
		}
		ResolveMoveSet(c.Faces, moves)
	case FaceF:
		moves := MoveSet{
			{FaceU, []int{7, 8, 9}, FaceR, []int{1, 4, 7}},
			{FaceR, []int{1, 4, 7}, FaceD, []int{3, 2, 1}},
			{FaceD, []int{1, 2, 3}, FaceL, []int{3, 6, 9}},
			{FaceL, []int{3, 6, 9}, FaceU, []int{9, 8, 7}},
		}
		ResolveMoveSet(c.Faces, moves)
	case FaceB:
		moves := MoveSet{
			{FaceR, []int{3, 6, 9}, FaceU, []int{1, 2, 3}},
			{FaceU, []int{1, 2, 3}, FaceL, []int{7, 4, 1}},
			{FaceL, []int{1, 4, 7}, FaceD, []int{7, 8, 9}},
			{FaceD, []int{7, 8, 9}, FaceR, []int{9, 6, 3}},
		}
		ResolveMoveSet(c.Faces, moves)
	case FaceL:
		moves := MoveSet{
			{FaceU, []int{1, 4, 7}, FaceF, []int{1, 4, 7}},
			{FaceF, []int{1, 4, 7}, FaceD, []int{1, 4, 7}},
			{FaceD, []int{1, 4, 7}, FaceB, []int{9, 6, 3}},
			{FaceB, []int{3, 6, 9}, FaceU, []int{7, 4, 1}},
		}
		ResolveMoveSet(c.Faces, moves)
	case FaceR:
		moves := MoveSet{
			{FaceU, []int{3, 6, 9}, FaceB, []int{7, 4, 1}},
			{FaceB, []int{1, 4, 7}, FaceD, []int{9, 6, 3}},
			{FaceD, []int{3, 6, 9}, FaceF, []int{3, 6, 9}},
			{FaceF, []int{3, 6, 9}, FaceU, []int{3, 6, 9}},
		}
		ResolveMoveSet(c.Faces, moves)
	}
}

func (c *Cube333) Print() {
	const reset = "\033[0m"
	pad := "       "

	for i := 0; i < 3; i++ {
		fmt.Print(pad)
		printRow(c.Faces[FaceU].Pieces[i*3 : i*3+3])
		fmt.Println(reset)
	}

	for i := 0; i < 3; i++ {
		printRow(c.Faces[FaceL].Pieces[i*3 : i*3+3])
		fmt.Print(" ")
		printRow(c.Faces[FaceF].Pieces[i*3 : i*3+3])
		fmt.Print(" ")
		printRow(c.Faces[FaceR].Pieces[i*3 : i*3+3])
		fmt.Print(" ")
		printRow(c.Faces[FaceB].Pieces[i*3 : i*3+3])
		fmt.Println(reset)
	}

	for i := 0; i < 3; i++ {
		fmt.Print(pad)
		printRow(c.Faces[FaceD].Pieces[i*3 : i*3+3])
		fmt.Println(reset)
	}
}

func printFace(face *CubeFace) {
	for i := 0; i < 9; i++ {
		fmt.Print(colorCode(face.Pieces[i]) + "  " + "\033[0m")
		if (i+1)%3 == 0 {
			fmt.Println()
		}
	}
}

func printRow(row []CubeColor) {
	for j := 0; j < 3; j++ {
		fmt.Print(colorCode(row[j]) + "  " + "\033[0m")
	}
}

func colorCode(c CubeColor) string {
	switch c {
	case White:
		return "\033[48;5;15m"
	case Yellow:
		return "\033[48;5;11m"
	case Orange:
		return "\033[48;5;208m"
	case Red:
		return "\033[48;5;9m"
	case Green:
		return "\033[48;5;10m"
	case Blue:
		return "\033[48;5;12m"
	}
	return "\033[48;5;0m"
}
