package writer

import (
	"context"
	"encoding/json"
	"fmt"
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

// WriteChannel is a method to write elements from a channel into a file respecting a timeout, it returns a channel
// which either signals Success or Error (buffered to 1)
func WriteChannel[T any](chn <-chan T, filename string, timeout time.Duration) <-chan pipes.Result[bool] {
	resultChan := make(chan pipes.Result[bool], 1)
	go func() {
		defer close(resultChan)

		var head T
	initialWait:
		for {
			select {
			case <-time.After(timeout):
				// we will first wait for first messages for timeout
				logger.Log.Errorf("timeout after %v", timeout)
				return
			case head = <-chn:
				// we receive first messages before timeout, continue
				fmt.Println("!!!!!!!!!!!!!! Received message")
				break initialWait
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		err := pipes.Pour(head, chn, func(elements []T) error {
			return WriteJson(elements, filename, PrettyPrint)
		}, ctx)

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
