package mySlices

func MapKeysToSlice[K comparable, V any](inputMap map[K]V) []K {
	var keys []K
	for key := range inputMap {
		keys = append(keys, key)
	}
	return keys
}
