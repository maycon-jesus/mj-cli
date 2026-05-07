package cube

func CreateCubeResolved(faces [6][9]CubeColor) [6][9]CubeColor {
	resolved := faces
	for faceIndex := 0; faceIndex < 6; faceIndex++ {
		fillFace(&resolved[faceIndex], CubeColor(faceIndex+1))
	}
	return resolved
}

func fillFace(face *[9]CubeColor, color CubeColor) {
	for i := 0; i < 9; i++ {
		face[i] = color
	}
}
