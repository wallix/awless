package utils

import "reflect"

func Intersect(a, b interface{}) []string {
	return intersect(extractField(a, "Id"), extractField(b, "Id"))
}

func Substraction(a, b interface{}) []string {
	return substraction(extractField(a, "Id"), extractField(b, "Id"))
}

func intersect(a, b []string) []string {
	var inter []string

	for i := 0; i < len(a); i++ {
		for j := 0; j < len(b); j++ {
			if a[i] == b[j] {
				inter = append(inter, a[i])
			}
		}
	}

	return inter
}

func substraction(a, b []string) []string {
	var sub []string

	for i := 0; i < len(a); i++ {
		var found bool
		for j := 0; j < len(b); j++ {
			if a[i] == b[j] {
				found = true
			}
		}
		if !found {
			sub = append(sub, a[i])
		}
	}

	return sub
}

func extractField(i interface{}, field string) []string {
	var fields []string

	value := reflect.ValueOf(i)

	if value.Kind() == reflect.Slice {
		for i := 0; i < value.Len(); i++ {
			s1 := value.Index(i).Elem()
			fields = append(fields, s1.FieldByName(field).Interface().(string))
		}

	}

	return fields
}
