package cube

type CubeColor int
type CubeFaceIndex int

const (
	Orange CubeColor = iota + 1
	White
	Green
	Yellow
	Red
	Blue
)

const (
	FaceL CubeFaceIndex = iota
	FaceU
	FaceF
	FaceD
	FaceR
	FaceB
)

type CubeFace struct {
	Color  CubeColor
	Pieces []CubeColor
}

func (f *CubeFace) Clone() *CubeFace {
	clone := &CubeFace{
		Color:  f.Color,
		Pieces: make([]CubeColor, len(f.Pieces)),
	}
	copy(clone.Pieces, f.Pieces)
	return clone
}

func newCubeFace(color CubeColor, size int) *CubeFace {
	face := &CubeFace{
		Color:  color,
		Pieces: make([]CubeColor, size*size),
	}
	for i := range face.Pieces {
		face.Pieces[i] = color
	}
	return face
}

// Layout das peças em cada face:
//
//	1 2 3
//	4 5 6
//	7 8 9
type PieceMove struct {
	SrcFace CubeFaceIndex
	SrcPos  []int
	DstFace CubeFaceIndex
	DstPos  []int
}

type MoveSet = []PieceMove

type MovesCollection = map[string]MoveSet

func ResolveMoveSet(faces [6]*CubeFace, moveSet MoveSet) {
	facesSrc := make([]*CubeFace, len(faces))
	for i, face := range faces {
		facesSrc[i] = face.Clone()
	}

	for _, move := range moveSet {
		dstFace := faces[move.DstFace]
		srcFace := facesSrc[move.SrcFace]
		for i := range move.SrcPos {
			// Subtrai 1 para ajustar a posição de 1-9 para 0-8
			dstFace.Pieces[move.DstPos[i]-1] = srcFace.Pieces[move.SrcPos[i]-1]
		}
	}
}
