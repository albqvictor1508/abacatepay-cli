package ui

import (
	"encoding/json"
	"fmt"
)

func PrintJSON(value any) error {
	bytes, err := json.MarshalIndent(value, "", " ")

	if err != nil {
		return err
	}

	fmt.Println(string(bytes))

	return nil
}
