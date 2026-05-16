package utils

func CreateDeepMap(keys []string, value interface{}) map[string]interface{} {
	if len(keys) == 0 {
		return nil
	}

	if len(keys) == 1 {
		return map[string]interface{}{
			keys[0]: value,
		}
	}

	return map[string]interface{}{
		keys[0]: CreateDeepMap(keys[1:], value),
	}
}

func MergeDeepMaps(dest, src map[string]interface{}) map[string]interface{} {
	for key, srcVal := range src {
		if destVal, exists := dest[key]; exists {
			destMap, destIsMap := destVal.(map[string]interface{})
			srcMap, srcIsMap := srcVal.(map[string]interface{})
			if destIsMap && srcIsMap {
				dest[key] = MergeDeepMaps(destMap, srcMap)
				continue
			}
		}
		dest[key] = srcVal
	}
	return dest
}
