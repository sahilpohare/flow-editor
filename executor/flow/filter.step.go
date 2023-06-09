package flow

import (
	"encoding/json"
	"fmt"
)

// import (
// 	"encoding/json"
// 	"fmt"
// )

// type FilterStep struct {
// 	Filter map[string]interface{}
// 	Init   func(data interface{}) *Step
// }

// func NewFilterStep(filter string) *FilterStep {
// 	filterMap := make(map[string]interface{})
// 	err := json.Unmarshal([]byte(filter), &filterMap)
// 	if err != nil {
// 		return nil
// 	}

// 	return &FilterStep{
// 		Filter: filterMap,
// 	}
// }

// func (f *FilterStep) Run(ctx *FlowContext, data map[string]interface{}) (map[string]interface{}, error) {
// 	d := map[string]interface{}{
// 		"ctx":  ctx.Data,
// 		"body": data,
// 	}

// 	// Check if the filter matches the data
// 	if !Filter(f.Filter, d) {
// 		ctx.ExitFlow()
// 	}

// 	return data, nil
// }

// func Run() {
// 	d := map[string]interface{}{
// 		"test":  "test",
// 		"test2": "test2",
// 		"test3": "test3",
// 	}

// 	filter := `
// 		{
// 			"test" : "test"
// 		}
// 	`

// 	filterMap := make(map[string]interface{})

// 	// Parse the rendered filter string
// 	err := json.Unmarshal([]byte(filter), &filterMap)
// 	if err != nil {
// 		return
// 	}

//		// Check if the filter matches the data
//		fmt.Println(Filter(filterMap, d))
//	}

func Example() {
	d := map[string]interface{}{
		"test":  "test",
		"test2": "test2",
		"test3": "test3",
	}

	filter := `
		{
			"test" : {
				"$in" : ["test1", "test2"]
			}
		}
	`

	filterMap := make(map[string]interface{})

	// Parse the rendered filter string
	err := json.Unmarshal([]byte(filter), &filterMap)
	if err != nil {
		return
	}

	// Check if the filter matches the data
	fmt.Println(Filter(filterMap, d))
}

func Filter(filter interface{}, data interface{}) bool {
	filterMap, ok := filter.(map[string]interface{})
	if !ok {
		fmt.Println("Unknown type in filter string")
		return false
	}

	for key, value := range filterMap {
		var operand interface{}
		op, ok := data.(map[string]interface{})
		if ok {
			operand = op[key]
		} else {
			operand = data
		}

		fmt.Println("Unknown type in filter string")
		switch key {
		case "$and":
			return AndOperator(value.([]interface{}), data)
		case "$or":
			return OrOperator(value.([]interface{}), data)
		case "$in":
			return InOperator(value.([]interface{}), operand)
		case "$nin":
			return !InOperator(value.([]interface{}), operand)
		case "$gt":
			return GtOperator(value, operand)
		case "$gte":
			return GteOperator(value, operand)
		case "$lt":
			return LtOperator(value, operand)
		case "$lte":
			return LteOperator(value, operand)
		default:
			fmt.Println(value, operand)
			if data.(map[string]interface{})[key] == nil {
				return false
			}
			switch value.(type) {
			case string:
				if operand != value {
					return false
				}
			case int16, int32, int64, int, float32, float64:
				op := operand.(float64)
				v := value.(float64)
				if op != v {
					return false
				}
			//If value is a nested object, then we need to recursively call Filte
			case map[string]interface{}:
				if !Filter(value, operand) {
					return false
				}

			case bool:
				if operand != value {
					return false
				}
			default:
				return false
			}
		}

	}

	return true
}

func AndOperator(values []interface{}, data interface{}) bool {
	for _, value := range values {
		if !Filter(value.(map[string]interface{}), data) {
			return false
		}
	}

	return true
}

func OrOperator(values []interface{}, data interface{}) bool {
	for _, value := range values {
		if Filter(value.(map[string]interface{}), data) {
			return true
		}
	}

	return false
}

func EqOperator(value interface{}, data interface{}) bool {
	switch value.(type) {
	case string:
		if data != value {
			return false
		}
	case int16, int32, int64, int, float32, float64:
		op, err := PerformOperatorOnNumber(value, data, func(op float64, v float64) bool {
			return op == v
		})

		if err != nil {
			fmt.Println("Error Performing Operation")
			return false
		}
		return op
	case bool:
		if data != value {
			return false
		}
	case []interface{}:
		if !InOperator(value.([]interface{}), data) {
			return false
		}
	default:
		fmt.Println("Unknown type in filter string")
		return false

	}

	return true
}

func InOperator(values []interface{}, data interface{}) bool {
	for _, value := range values {
		if value == data {
			return true
		}
	}

	return false
}

func GtOperator(value interface{}, data interface{}) bool {
	out, err := PerformOperatorOnNumber(value, data, func(op float64, v float64) bool {
		return op > v
	})
	if err != nil {
		fmt.Println("Error Performing Operation")
		return false
	}
	return out
}

func GteOperator(value interface{}, data interface{}) bool {
	out, err := PerformOperatorOnNumber(value, data, func(op float64, v float64) bool {
		return op <= v
	})
	if err != nil {
		fmt.Println("Error Performing Operation")
		return false
	}
	return out
}

func LtOperator(value interface{}, data interface{}) bool {
	out, err := PerformOperatorOnNumber(value, data, func(op float64, v float64) bool {
		return op > v
	})
	if err != nil {
		fmt.Println("Error Performing Operation")
		return false
	}
	return out
}

func LteOperator(value interface{}, data interface{}) bool {
	out, err := PerformOperatorOnNumber(value, data, func(op float64, v float64) bool {
		return op >= v
	})
	if err != nil {
		fmt.Println("Error Performing Operation")
		return false
	}
	return out
}

func IsNumber(value interface{}) bool {
	switch value.(type) {
	case int16, int32, int64, int, float32, float64:
		return true
	}

	return false
}

func convertToFloat64(value interface{}) (float64, error) {
	switch value := value.(type) {
	case int:
		return float64(value), nil
	case int8:
		return float64(value), nil
	case int16:
		return float64(value), nil
	case int32:
		return float64(value), nil
	case int64:
		return float64(value), nil
	case uint:
		return float64(value), nil
	case uint8:
		return float64(value), nil
	case uint16:
		return float64(value), nil
	case uint32:
		return float64(value), nil
	case uint64:
		return float64(value), nil
	case float32:
		return float64(value), nil
	case float64:
		return value, nil
	default:
		return 0, fmt.Errorf("conversion error: value is not a number")
	}
}

func PerformOperatorOnNumber(a interface{}, b interface{}, operator func(a, b float64) bool) (bool, error) {
	ok := IsNumber(a)
	if !ok {
		return false, fmt.Errorf("Invalid Type for A")
	}

	ok = IsNumber(b)
	if !ok {
		return false, fmt.Errorf("Invalid Type for B")
	}

	op, err := convertToFloat64(a)
	v, err := convertToFloat64(b)
	if err != nil {
		fmt.Println("Unknown type in filter string")
		return false, err
	}

	return operator(op, v), nil
}
