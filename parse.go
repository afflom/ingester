package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	xml "github.com/clbanning/mxj/v2"
)

// parse.go lists supported formats and provides the parsing logic for each format

type parser interface {
	parse(string) (map[string]interface{}, error)
}

type jsonParser string
type xmlParser string
type unsupported string

func (dsc jsonParser) parse(filename string) (map[string]interface{}, error) {

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	b, _ := ioutil.ReadAll(file)
	result := make(map[string]interface{})
	json.Unmarshal(b, &result)

	return result, err

}

func (dsc xmlParser) parse(filename string) (map[string]interface{}, error) {

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	b, _ := ioutil.ReadAll(file)

	result, err := xml.NewMapXml(b)
	return result, err
}

func (dsc unsupported) parse(filename string) (result map[string]interface{}, err error) {

	return result, err

}
