package mySlices

func Filter[T any](slice []T, checker func(T) bool) []T {
	nSlice := make([]T, 0)
	for _, item := range slice {
		if checker(item) {
			nSlice = append(nSlice, item)
		}
	}
	return nSlice
}
