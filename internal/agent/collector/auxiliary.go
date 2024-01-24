package collector

import (
	"fmt"
	"reflect"
	"runtime"
)

// GetMemStatsFieldValueByName returns float64 representation of runtime.MemStats field value if such field name exists
// and field type is appropriate.
func GetMemStatsFieldValueByName(memstats *runtime.MemStats, fieldName string) (float64, error) {
	reflectedMemStats := reflect.ValueOf(*memstats)
	value := reflectedMemStats.FieldByName(fieldName)

	if value.IsValid() {
		switch {
		case value.CanFloat():
			return value.Float(), nil
		case value.CanUint():
			return float64(value.Uint()), nil
		default:
			return 0, fmt.Errorf("inappropriate field type")
		}
	}

	return 0, fmt.Errorf("no such field")
}
