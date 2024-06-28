package writer

import (
	"context"
	"encoding/json"
	"github.com/johannessarpola/lutakkols/pkg/api/options"
	"github.com/johannessarpola/lutakkols/pkg/logger"
	"github.com/johannessarpola/lutakkols/pkg/pipes"
	"os"
	"time"
)

// WriteOption to handle output controls
type WriteOption int

const (
	_ WriteOption = iota
	PrettyPrint
)

func WriteChannel[T any](chn <-chan T, filename string, timeout time.Duration) <-chan pipes.Result[bool] {
	resultChan := make(chan pipes.Result[bool], 1)
	go func() {
		consumeCtx, cancel := context.WithTimeout(context.Background(), timeout)
		defer func() {
			close(resultChan)
			cancel()
		}()

		err := pipes.Pour(chn, func(elements []T) error {
			return WriteJson(elements, filename, PrettyPrint)
		}, consumeCtx)

		if err != nil {
			logger.Log.Error("write error", err)
		}
		resultChan <- pipes.Result[bool]{
			Val: err == nil,
			Err: err,
		}
	}()

	return resultChan
}

// WriteJson general purpose func to write generic object to a file as json
func WriteJson(data interface{}, outFile string, opts ...WriteOption) error {
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

	if options.Has(PrettyPrint, opts) {
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
