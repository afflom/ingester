package ingest

import (
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"github.com/uor-framework/uor-client-go/api/client/v1alpha1"
	"mime"
	"os"
	"path/filepath"
)

func Ingest(shema string, workspace string) {
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

		// Call parse by mediatype
		switch mtype {
		case "application/json":
			attributes := parseJson(path)
		case "text/xml":
			attributes := parseXML(path)
		case 3:
			fmt.Println("three")
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}

	// Add attributes to each file in the dataset config

}
