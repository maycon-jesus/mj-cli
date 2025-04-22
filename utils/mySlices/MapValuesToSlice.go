package mySlices

func MapValuesToSlice[K comparable, V any](inputMap map[K]V) []V {
	var values []V
	for _, value := range inputMap {
		values = append(values, value)
	}
	return values
}
