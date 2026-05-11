package cube

import (
	"errors"
	"fmt"
	"strings"

	"github.com/maycon-jesus/mj-cli/pkg/tint"
)

type Cube333 struct {
	Faces [6]*CubeFace
}

var Cube333MovesCollection = MovesCollection{
	"U": MoveSet{
		PieceMove{FaceU, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, FaceU, []int{3, 6, 9, 2, 5, 8, 1, 4, 7}},
		{FaceL, []int{1, 2, 3}, FaceB, []int{1, 2, 3}},
		{FaceB, []int{1, 2, 3}, FaceR, []int{1, 2, 3}},
		{FaceR, []int{1, 2, 3}, FaceF, []int{1, 2, 3}},
		{FaceF, []int{1, 2, 3}, FaceL, []int{1, 2, 3}},
	},
	"D": MoveSet{
		PieceMove{FaceD, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, FaceD, []int{3, 6, 9, 2, 5, 8, 1, 4, 7}},
		{FaceL, []int{7, 8, 9}, FaceF, []int{7, 8, 9}},
		{FaceF, []int{7, 8, 9}, FaceR, []int{7, 8, 9}},
		{FaceR, []int{7, 8, 9}, FaceB, []int{7, 8, 9}},
		{FaceB, []int{7, 8, 9}, FaceL, []int{7, 8, 9}},
	},
	"F": MoveSet{
		PieceMove{FaceF, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, FaceF, []int{3, 6, 9, 2, 5, 8, 1, 4, 7}},
		{FaceU, []int{7, 8, 9}, FaceR, []int{1, 4, 7}},
		{FaceR, []int{1, 4, 7}, FaceD, []int{3, 2, 1}},
		{FaceD, []int{1, 2, 3}, FaceL, []int{3, 6, 9}},
		{FaceL, []int{3, 6, 9}, FaceU, []int{9, 8, 7}},
	},
	"B": MoveSet{
		PieceMove{FaceB, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, FaceB, []int{3, 6, 9, 2, 5, 8, 1, 4, 7}},
		{FaceR, []int{3, 6, 9}, FaceU, []int{1, 2, 3}},
		{FaceU, []int{1, 2, 3}, FaceL, []int{7, 4, 1}},
		{FaceL, []int{1, 4, 7}, FaceD, []int{7, 8, 9}},
		{FaceD, []int{7, 8, 9}, FaceR, []int{9, 6, 3}},
	},
	"R": MoveSet{
		PieceMove{FaceR, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, FaceR, []int{3, 6, 9, 2, 5, 8, 1, 4, 7}},
		{FaceU, []int{3, 6, 9}, FaceB, []int{7, 4, 1}},
		{FaceB, []int{1, 4, 7}, FaceD, []int{9, 6, 3}},
		{FaceD, []int{3, 6, 9}, FaceF, []int{3, 6, 9}},
		{FaceF, []int{3, 6, 9}, FaceU, []int{3, 6, 9}},
	},
	"L": MoveSet{
		PieceMove{FaceL, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, FaceL, []int{3, 6, 9, 2, 5, 8, 1, 4, 7}},
		{FaceU, []int{1, 4, 7}, FaceF, []int{1, 4, 7}},
		{FaceF, []int{1, 4, 7}, FaceD, []int{1, 4, 7}},
		{FaceD, []int{1, 4, 7}, FaceB, []int{9, 6, 3}},
		{FaceB, []int{3, 6, 9}, FaceU, []int{7, 4, 1}},
	},
}

func CreateCube333() Cube333 {
	cube := Cube333{}

	cube.Faces[FaceL] = newCubeFace(Orange, 3)
	cube.Faces[FaceU] = newCubeFace(White, 3)
	cube.Faces[FaceF] = newCubeFace(Green, 3)
	cube.Faces[FaceD] = newCubeFace(Yellow, 3)
	cube.Faces[FaceR] = newCubeFace(Red, 3)
	cube.Faces[FaceB] = newCubeFace(Blue, 3)

	return cube
}

func (c *Cube333) ApplyMove(move string) error {
	if len(move) == 0 {
		return errors.New("invalid move: empty string")
	}
	if len(move) > 2 {
		return errors.New("invalid move: too long")
	}

	double := false
	reverse := false
	face := string(move[0])

	switch move[len(move)-1] {
	case '2':
		double = true
	case '\'':
		reverse = true
	}

	moveSet, exist := Cube333MovesCollection[face]
	if !exist {
		return errors.New("invalid move: " + move)
	}

	ResolveMoveSet(c.Faces, moveSet)
	if double {
		ResolveMoveSet(c.Faces, moveSet)
	}
	if reverse {
		ResolveMoveSet(c.Faces, moveSet)
		ResolveMoveSet(c.Faces, moveSet)
	}

	return nil
}

func (c *Cube333) Print() {
	fmt.Print(c.Render())
}

func (c *Cube333) Render() string {
	const reset = "\033[0m"
	const pad = "       "

	var sb strings.Builder

	for i := 0; i < 3; i++ {
		sb.WriteString(pad)
		writeRow(&sb, c.Faces[FaceU].Pieces[i*3:i*3+3])
		sb.WriteString(reset)
		sb.WriteByte('\n')
	}

	for i := 0; i < 3; i++ {
		writeRow(&sb, c.Faces[FaceL].Pieces[i*3:i*3+3])
		sb.WriteByte(' ')
		writeRow(&sb, c.Faces[FaceF].Pieces[i*3:i*3+3])
		sb.WriteByte(' ')
		writeRow(&sb, c.Faces[FaceR].Pieces[i*3:i*3+3])
		sb.WriteByte(' ')
		writeRow(&sb, c.Faces[FaceB].Pieces[i*3:i*3+3])
		sb.WriteString(reset)
		sb.WriteByte('\n')
	}

	for i := 0; i < 3; i++ {
		sb.WriteString(pad)
		writeRow(&sb, c.Faces[FaceD].Pieces[i*3:i*3+3])
		sb.WriteString(reset)
		sb.WriteByte('\n')
	}

	return sb.String()
}

func writeRow(sb *strings.Builder, row []CubeColor) {
	for j := 0; j < 3; j++ {
		sb.WriteString(colorCode(row[j]).Render("  "))
	}
}

func colorCode(c CubeColor) *tint.Style {
	switch c {
	case White:
		return tint.NewStyle().Background(tint.CompleteColor{
			TrueColor: tint.TrueColor{R: 255, G: 255, B: 255},
			Ansi256:   tint.ColorAnsi256{Code: 15},
			Ansi:      tint.ColorAnsi{Code: tint.ANSI_COLOR_WHITE},
		})
	case Yellow:
		return tint.NewStyle().Background(tint.CompleteColor{
			TrueColor: tint.TrueColor{R: 255, G: 255, B: 0},
			Ansi256:   tint.ColorAnsi256{Code: 11},
			Ansi:      tint.ColorAnsi{Code: tint.ANSI_COLOR_YELLOW},
		})
	case Orange:
		return tint.NewStyle().Background(tint.CompleteColor{
			TrueColor: tint.TrueColor{R: 255, G: 165, B: 0},
			Ansi256:   tint.ColorAnsi256{Code: 208},
		})
	case Red:
		return tint.NewStyle().Background(tint.CompleteColor{
			TrueColor: tint.TrueColor{R: 255, G: 0, B: 0},
			Ansi256:   tint.ColorAnsi256{Code: 9},
			Ansi:      tint.ColorAnsi{Code: tint.ANSI_COLOR_RED},
		})
	case Green:
		return tint.NewStyle().Background(tint.CompleteColor{
			TrueColor: tint.TrueColor{R: 0, G: 255, B: 0},
			Ansi256:   tint.ColorAnsi256{Code: 10},
			Ansi:      tint.ColorAnsi{Code: tint.ANSI_COLOR_GREEN},
		})
	case Blue:
		return tint.NewStyle().Background(tint.CompleteColor{
			TrueColor: tint.TrueColor{R: 0, G: 0, B: 255},
			Ansi256:   tint.ColorAnsi256{Code: 12},
			Ansi:      tint.ColorAnsi{Code: tint.ANSI_COLOR_BLUE},
		})
	}
	return tint.NewStyle()
}
