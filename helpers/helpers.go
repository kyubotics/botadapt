package helpers

import (
	"fmt"
	"reflect"
	"strings"
)

func EnsureString(v interface{}) string {
	t := reflect.TypeOf(v).String()
	if strings.HasPrefix(t, "int") || strings.HasPrefix(t, "uint") {
		return fmt.Sprintf("%d", v)
	} else if strings.HasPrefix(t, "float") {
		s := fmt.Sprintf("%f", v)
		return strings.TrimRight(s, ".0")
	} else if t == "string" {
		return v.(string)
	}
	return fmt.Sprint(v)
}
