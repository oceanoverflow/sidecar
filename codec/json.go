package codec

import (
	"fmt"
)

func WriteObject(i interface{}) []byte {
	var s string
	switch o := i.(type) {
	case string:
		s = fmt.Sprintf("\"%s\"", o)
	case []byte:
		s = fmt.Sprintf("\"%s\"", string(o))
	case map[string]string:
		s += "{"
		i := 0
		for k, v := range o {
			s += fmt.Sprintf("\"%s\":\"%s\"", k, v)
			if i != len(o)-1 {
				s += ","
			}
			i++
		}
		s += "}"
	default:
		s = fmt.Sprintf("%s", "null")
	}
	s = fmt.Sprintln(s)
	return []byte(s)
}
