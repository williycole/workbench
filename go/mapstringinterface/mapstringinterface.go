package mapstringinterface

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
)

/*
Assignment

Sometimes you have to deal with JSON data of unknown or varying structure in Go. In those instances
map[string]interface{} offers a flexible way to handle it without predefined structs.
*/

/*
getResources()
getResources takes a url string and returns a slice of maps []map[string]interface{} and an error.
Decode the response body into a slice of maps []map[string]interface{} and return it.
*/
func getResources(url string) ([]map[string]any, error) {
	var resources []map[string]any

	res, err := http.Get(url)
	if err != nil {
		return resources, err
	}
	defer res.Body.Close()

	var data []byte
	if data, err = io.ReadAll(res.Body); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &resources); err != nil {
		return nil, err
	}

	return resources, nil
}

/*
logResources()
logResources takes a slice of maps []map[string]interface{} and prints its keys and values to the console.
Because maps are unsorted we will be adding formatted strings to a slice of strings []string which is then sorted.

Iterate over the slice of map[string]interface{}
For each map[string]interface{} get its keys and values using range and append it to formattedStrings as
Key: %s - Value: %v, where %s is the key and %v is the value.
*/
func logResources(resources []map[string]any) {
	var formattedStrings []string

	for _, resource := range resources {
		for k, v := range resource {
			formattedStrings = append(formattedStrings, fmt.Sprintf("Key: %s - Value: %v", k, v))
		}
	}

	sort.Strings(formattedStrings)

	for _, str := range formattedStrings {
		fmt.Println(str)
	}
}
