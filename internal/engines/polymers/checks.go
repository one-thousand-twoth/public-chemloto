package polymers

import "reflect"

type Checks struct {
	Fields map[string]map[string][]map[string]int
}

func checkFields(checks Checks, field string, name string, data map[string]int) bool {
	var eq bool
	for _, entry := range checks.Fields[field][name] {
		eq = reflect.DeepEqual(removeZeroValues(entry), data)
		if eq {
			break
		}
	}
	return eq
}

// type Field struct {
// 	entries
// }

// // type Entry struct {
// // 	structure
// // }
