# Flow and Filtering Engine 

## Introduction
This project is used to run custom flows on triggers. This is a very simple implementation of a flow engine. It is not meant to be used in production. It is meant to be used as a starting point for a more complex flow engine.

## How to use the filtering engine
Just use the `flow.Filter(filter, data)` function in the package

#### More Nodes in the roadmap
- [x] Filter Step
- [ ] Transform Step
- [ ] Omit Keys Step
- [ ] HTTP Request Step
- [ ] HTTP Response Step

#### Example
```go
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
	fmt.Println(flow.Filter(filterMap, d))
}
