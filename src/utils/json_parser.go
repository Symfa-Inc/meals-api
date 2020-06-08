package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func JsonParse(path string, model interface{}) {
	pwd, _ := os.Getwd()
	jsonFile, err := os.Open(pwd + path)

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	_ = json.Unmarshal(byteValue, &model)
}
