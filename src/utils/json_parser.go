package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// JSONParse parsing json from provided path
func JSONParse(path string, model interface{}) {
	pwd, _ := os.Getwd()
	jsonFile, err := os.Open(pwd + path)

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	_ = json.Unmarshal(byteValue, &model)
}
