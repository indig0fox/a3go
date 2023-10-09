package a3interface

import (
	"encoding/json"
	"errors"
	"strings"
)

func RemoveEscapeQuotes(input string) string {
	// Remove leading and trailing double quotes if they exist.
	input = strings.TrimPrefix(input, `"`)
	input = strings.TrimSuffix(input, `"`)

	// Replace all double-double quotes with a single double quote.
	input = strings.ReplaceAll(input, `""`, `"`)
	input = strings.ReplaceAll(input, `""`, `\"`)
	return input
}

func ParseSQF(input string) (interface{}, error) {
	preparedInput := RemoveEscapeQuotes(input)

	var result interface{}
	err := json.Unmarshal([]byte(preparedInput), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func ParseSQFHashMap(input interface{}) (map[string]interface{}, error) {
	array, ok := input.([]interface{})
	if !ok {
		return nil, errors.New("invalid format")
	}

	result := make(map[string]interface{})
	for _, item := range array {
		pair, ok := item.([]interface{})
		if !ok || len(pair) != 2 {
			return nil, errors.New("invalid key-value pair")
		}

		key, ok := pair[0].(string)
		if !ok {
			return nil, errors.New("invalid key type")
		}
		value := pair[1]

		// Check if the value is another key-value pair array (i.e., nested HashMap)
		if subArray, ok := value.([]interface{}); ok && len(subArray) > 0 {
			if subPair, ok := subArray[0].([]interface{}); ok && len(subPair) == 2 {
				subResult, err := ParseSQFHashMap(value)
				if err != nil {
					return nil, err
				}
				result[key] = subResult
				continue
			}
		}

		result[key] = value
	}

	return result, nil
}

// func main() {
// 	example1 := `"[1, 6, 3, [""message"", ""Could not finish""]]"`
// 	example2 := `"[[""test_element"", 1], ["data", [[""name"", ""Danny""], [""time_in_session"", 370.2]]]]"`
// 	example3 := `"[[""name"", ""Rick""], [""age"", 29], [""latest_visits"", [[""most_recent"", 20], [""secondmost_recent"", 35], [""thirdmost_recent"", 50]]]]"`

// 	r1, err := parseSQF(example1)
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 	} else {
// 		fmt.Println(r1)
// 	}

// 	// Parsing HashMaps now
// 	r2Parsed, err := parseSQF(example2)
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		return
// 	}
// 	r2, err := parseSQFHashMap(r2Parsed)
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 	} else {
// 		fmt.Println(r2)
// 		fmt.Println(r2["data"])
// 	}

// 	r3Parsed, err := parseSQF(example3)
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		return
// 	}
// 	r3, err := parseSQFHashMap(r3Parsed)
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 	} else {
// 		fmt.Println(r3)
// 		fmt.Println(r3["latest_visits"].(map[string]interface{})["most_recent"])
// 	}
// }
