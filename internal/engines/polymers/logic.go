package polymers

func removeByValue(slice []string, value string) []string {
	for i := 0; i < len(slice); i++ {
		if slice[i] == value {
			slice = append(slice[:i], slice[i+1:]...)
			break
		}
	}
	return slice
}

// Keys returns the keys of the map m.
// The keys will be an indeterminate order.
func Keys[M ~map[K]V, K comparable, V any](m M) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
}
