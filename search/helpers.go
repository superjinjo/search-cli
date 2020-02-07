package search

import (
	"fmt"
	"strconv"
	"strings"
)

func floatVal(value interface{}) (float64, bool) {
	switch v := value.(type) {
	case int:
		return float64(v), true

	case float64:
		return v, true

	case string:
		if v == "" {
			return 0, true
		}

		if floatVal, err := strconv.ParseFloat(v, 64); err == nil {
			return floatVal, true
		}

		return 0, false

	default:
		return 0, v == nil
	}
}

func boolVal(value interface{}) (bool, bool) {
	//FYI: in Go there is no fallthrough for type assertions
	switch v := value.(type) {
	case int:
		return v > 0, true
	case float64:
		return v > 0, true
	case string:
		if strings.EqualFold(v, "true") || v == "1" {
			return true, true
		} else if strings.EqualFold(v, "false") || v == "0" || v == "" {
			return false, true
		}

		return false, false
	case bool:
		return v, true

	default:
		return false, v == nil
	}
}

func sliceVal(value interface{}) []interface{} {
	switch v := value.(type) {
	case []interface{}:
		return v
	case string:
		stringSlice := []interface{}{}
		if v == "" {
			return stringSlice
		}

		splitVal := strings.Split(v, ",")
		//after splitting we have to loop through and put it in the interface slice
		//because a string slice is not convertable to any other type of slice
		for _, nextVal := range splitVal {
			stringSlice = append(stringSlice, nextVal)
		}

		return stringSlice
	default:
		if value == nil {
			return []interface{}{}
		}

		return []interface{}{value}
	}
}

func stringVal(value interface{}) string {
	if value == nil {
		return ""
	}

	return fmt.Sprintf("%v", value)
}

//subjectSlice must have ALL values in term slice unless the term  slice is empty, in which case they must be equal
func sliceContains(subjectSlice []interface{}, termSlice []interface{}) bool {
	if len(termSlice) == 0 {
		return len(subjectSlice) == 0
	}

	if len(termSlice) > len(subjectSlice) {
		return false
	}

	for _, termVal := range termSlice {
		hasVal := false
		for _, subjectVal := range subjectSlice {
			hasVal = SearchValueMatches(subjectVal, termVal)

			if hasVal {
				break
			}
		}

		if !hasVal {
			return false
		}
	}

	return true
}

//only checking values on supported data types
func valueIsEmpty(value interface{}) bool {
	if valSlice, isSlice := value.([]interface{}); isSlice {
		return len(valSlice) == 0
	}

	return value == nil || value == "" || value == 0 || value == float64(0) || value == false
}

//ValueMatcher functions dictate the rules for deciding if two values match
type ValueMatcher func(subject interface{}, term interface{}) bool

//empty values will count as equivalent if the field value is nil
//nil or empty string will count as equivalent is field is an empty slice
//a term slice is considered a match for a subject slice if the subject contains all items from the term (order doesn't matter)
func SearchValueMatches(subject interface{}, term interface{}) bool {
	switch subjVal := subject.(type) {
	case int:
		return SearchValueMatches(float64(subjVal), term)
	case float64:
		termVal, isFloat := floatVal(term)
		if isFloat {
			return subjVal == termVal
		}

		return false
	case bool:
		termVal, isBool := boolVal(term)
		if isBool {
			return subjVal == termVal
		}

		return false
	case []interface{}:
		return sliceContains(subjVal, sliceVal(term))
	case string:
		return stringVal(term) == subjVal
	default:
		if subjVal == nil {
			return valueIsEmpty(term)
		}

		return subject == term
	}
}
