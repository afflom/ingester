package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/imdario/mergo"
	"github.com/nqd/flat"
	"github.com/tidwall/gjson"
	clientapi "github.com/uor-framework/uor-client-go/api/client/v1alpha1"
	"github.com/uor-framework/uor-client-go/content"
	"github.com/uor-framework/uor-client-go/content/layout"
)

// Ingest
func (ingester Ingester) Ingest() error {

	foundAttributes := make(map[string]interface{})

	// Create a dataset-config
	dsc := clientapi.DataSetConfiguration{
		TypeMeta: clientapi.TypeMeta{
			Kind:       clientapi.DataSetConfigurationKind,
			APIVersion: clientapi.GroupVersion,
		},
	}

	store, err := layout.New(ingester.CacheDir)
	if err != nil {
		fmt.Println(err)
	}
	schemaMap := make(map[string]interface{})
	schema, err := fetchJSONSchema(context.Background(), ingester.Schema, store)
	if err != nil {
		fmt.Println(err)
	}
	json.Unmarshal(schema, &schemaMap)
	json.Marshal(schemaMap)

	flatSchema, err := flat.Flatten(schemaMap, nil)
	if err != nil {
		fmt.Println(err)
	}

	// schemaStructs, err := rts.DoRaw("", string(schema))
	// if err != nil {
	//   fmt.Println(err)
	// }
	// fmt.Printf("schemastructs: %v", string(schemaStructs))
	// Add the files, links, and annotations to the dataset-config
	err = filepath.Walk(ingester.Workspace, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
		}
		file := clientapi.File{
			File: path,
		}
		dsc.Collection.Files = append(dsc.Collection.Files, file)
		// Get the mediatype of the file
		mtype, err := mimetype.DetectFile(path)
		if err != nil {
			fmt.Println(err)
		}

		var mt parser
		switch mtype.String() {
		case "application/json":
			fmt.Printf("File: %s, is %s\n", path, mtype.String())
			mt = jsonParser("")
		case "text/xml":
			fmt.Printf("File: %s, is %s\n", path, mtype.String())
			mt = xmlParser("")
		default:
			fmt.Printf("File: %s, is %s\n", path, mtype.String())
			mt = unsupported("unsupported")
		}
		if mt == unsupported("unsupported") {
			return nil
		}
		fmt.Println("Starting parser")
		p, err := mt.parse(path)
		if err != nil {
			fmt.Printf("Parsing Error: %v\n", err)
		}
		parsed, err := json.Marshal(p)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("parsed: %v\n", string(parsed))

		for jsonPath, _ := range flatSchema {
			fmt.Printf("searching content for: %s\n", jsonPath)
			jsonPath := strings.TrimSuffix(jsonPath, ".type")
			value := gjson.Get(string(parsed), jsonPath)
			fmt.Printf("Match: %s\n", value.String())
			if value.String() != "" {
				foundPair := map[string]interface{}{jsonPath: value}
				out, err := flat.Unflatten(foundPair, nil)
				if err != nil {
					fmt.Println(err)
				}
				if err := mergo.Merge(&foundAttributes, out); err != nil {
					fmt.Println(err)
				}
				fmt.Printf("foundAttributes: %s\n", foundAttributes)

			}
		}

		// Search foundAttributes for filenames within the workspace
		// If files are found in foundAttributes, get the json object that filename occurs in and
		// add that object to the corresponding file's attributes in the dataset config
		// Any attributes remaining in the foundAttributes map are applied to the file that they
		// were found in.

		return nil
	})
	if err != nil {
		fmt.Println(err)
	}

	return nil

}

func fetchJSONSchema(ctx context.Context, schemaAddress string, store content.AttributeStore) ([]byte, error) {
	desc, err := store.AttributeSchema(ctx, schemaAddress)
	if err != nil {
		return nil, err
	}

	schemaReader, err := store.Fetch(ctx, desc)
	if err != nil {
		return nil, fmt.Errorf("error fetching schema from store: %w", err)
	}
	schemaBytes, err := ioutil.ReadAll(schemaReader)
	if err != nil {
		return nil, err
	}

	return schemaBytes, err
}
