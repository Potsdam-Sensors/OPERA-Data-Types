package operadatatypes

import "math"

func RemoveNaN(m map[string]any) {
	keysToRemove := []string{}
	for k, v := range m {
		switch rv := v.(type) {
		case float32:
			if math.IsNaN(float64(rv)) {
				keysToRemove = append(keysToRemove, k)
			}
		case float64:
			if math.IsNaN(rv) {
				keysToRemove = append(keysToRemove, k)
			}
		default:
		}
	}
	for _, k := range keysToRemove {
		delete(m, k)
	}
}
