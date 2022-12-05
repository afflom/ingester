package ingest

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"os"

	xml "github.com/clbanning/mxj/v2"
)

// parse.go lists supported formats and provides the parsing logic for each format

type parser interface {
	parse(fs.FileInfo) (map[string]interface{}, error)
}

type jsonParser string
type xmlParser string

func (dsc jsonParser) parse(filename fs.FileInfo) (map[string]interface{}, error) {

	file, err := os.Open(filename.Name())
	if err != nil {
		return nil, err
	}
	defer file.Close()

	b, _ := ioutil.ReadAll(file)
	result := make(map[string]interface{})
	json.Unmarshal(b, &result)

	return result, err

}

func (dsc xmlParser) parse(filename fs.FileInfo) (map[string]interface{}, error) {

	file, err := os.Open(filename.Name())
	if err != nil {
		return nil, err
	}
	defer file.Close()
	b, _ := ioutil.ReadAll(file)

	result, err := xml.NewMapXml(b)
	return result, err
}
