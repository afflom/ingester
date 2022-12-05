package ingester

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gabriel-vasile/mimetype"
	"github.com/uor-framework/uor-client-go/api/client/v1alpha1"

	"github.com/uor-framework/uor-client-go/content"
	"github.com/uor-framework/uor-client-go/content/layout"
	"github.com/uor-framework/uor-client-go/schema"
)

// Ingest
func Ingest(schemaAddress string, workspace string) {
	// create empty object of type schemaid

	attributes := make(map[string]interface{})

	// Create a dataset-config
	dsc := v1alpha1.DataSetConfiguration{}
	// Add the files, links, and annotations to the dataset-config
	err := filepath.Walk(workspace, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		file := v1alpha1.File{
			File: path,
		}
		dsc.Collection.Files = append(dsc.Collection.Files, file)
		// Get the mediatype of the file
		mtype, err := mimetype.DetectFile(path)
		if err != nil {
			fmt.Println(err)
			return err
		}

		var mt parser
		switch mtype.String() {
		case "application/json":
			mt = jsonParser("")
		case "text/xml":
			mt = xmlParser("")
		}
		parsed, err := mt.parse(info)
		if err != nil {
			fmt.Println(err)
		}

		// for each attribute in the schema, search the parsed file for that attribute key.
		store, err := layout.New("./content-store")
		if err != nil {
			fmt.Println(err)
		}

		schema, _, err := fetchJSONSchema(context.Background(), schemaAddress, store)
		for k, v := range schema {

		}
		// if the attribute exists in the file, add it to the object of type schemaid
		// When done, search object of type schemaid for each filename in the workspace
		// if a filename from the workspace is found,
		// write the object that it occurs within to the attributes of the file in the dataset-config

		return nil
	})
	if err != nil {
		fmt.Println(err)
	}

	// Add attributes to each file in the dataset config

}

func fetchJSONSchema(ctx context.Context, schemaAddress string, store content.AttributeStore) (schema.Schema, string, error) {
	desc, err := store.AttributeSchema(ctx, schemaAddress)
	if err != nil {
		return schema.Schema{}, "", err
	}

	var schemaID string
	node, err := v2.NewNode(desc.Digest.String(), desc)
	if err != nil {
		return schema.Schema{}, "", err
	}
	props := node.Properties
	if props.IsASchema() {
		schemaID = props.Schema.ID
	}

	schemaReader, err := store.Fetch(ctx, desc)
	if err != nil {
		return schema.Schema{}, "", fmt.Errorf("error fetching schema from store: %w", err)
	}
	schemaBytes, err := ioutil.ReadAll(schemaReader)
	if err != nil {
		return schema.Schema{}, "", err
	}
	loader, err := schema.FromBytes(schemaBytes)
	if err != nil {
		return schema.Schema{}, "", err
	}

	sc, err := schema.New(loader)
	return sc, schemaID, err
}
