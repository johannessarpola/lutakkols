package downloader

import (
	"encoding/json"
	"github.com/johannessarpola/go-lutakko-gigs/pkg/api/options"
	"os"
)

// writeJson general purpose func to write generic object to a file as json
func writeJson(data interface{}, outFile string, opts ...options.WriteOpton) error {
	// Open a file for writing
	file, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	// Marshal the array of structs into JSON
	var (
		jsonData []byte
	)

	if options.Has(options.PrettyPrint, opts) {
		jsonData, err = json.MarshalIndent(data, "", "  ")
	} else {
		jsonData, err = json.Marshal(data)
	}
	if err != nil {
		return err
	}

	// Write the JSON data to the file
	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}
	return nil
}
