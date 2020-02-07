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
		if floatVal, err := strconv.ParseFloat(v, 64); err == nil {
			return floatVal, true
		}

		return 0, false

	default:
		return 0, false
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
		} else if strings.EqualFold(v, "false") || v == "0" {
			return false, true
		}

		return false, false
	case bool:
		return v, true

	default:
		return false, false
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
			hasVal = valuesMatch(subjectVal, termVal)

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
	return value == nil || value == "" || value == 0 || value == float64(0) || value == false
}

//empty values will count as equivalent if the field value is nil
//nil or empty string will count as equivalent is field is an empty slice
//a term slice is considered a match for a subject slice if the subject contains all items from the term (order doesn't matter)
func valuesMatch(subject interface{}, term interface{}) bool {
	switch subjVal := subject.(type) {
	case int:
		return valuesMatch(float64(subjVal), term)
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
		if termSlice, isSlice := term.([]interface{}); isSlice {
			return sliceContains(subjVal, termSlice)
		}

		if len(subjVal) == 0 {
			return term == nil || term ==""
		}

		return false
	case string:
		return stringVal(term) == subjVal
	default:
		if subjVal == nil {
			return valueIsEmpty(term)
		}

		return subject == term
	}
}
