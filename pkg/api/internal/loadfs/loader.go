// Package loadfs contains the methods used in offline mode where in the data is loaded from JSON files from disk
package loadfs

import (
	"encoding/json"
	"github.com/johannessarpola/lutakkols/pkg/api/models"
	"os"
)

// Events loads events from a json file
func Events(fp string) (*models.Events, error) {
	// Open the JSON file
	file, err := os.Open(fp)
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	if err != nil {
		return nil, err
	}

	stat, err := os.Stat(fp)
	if err != nil {
		return nil, err
	}

	var events []models.Event
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&events)

	if err != nil {
		return nil, err
	}

	return &models.Events{
		Events:    events,
		UpdatedAt: stat.ModTime(),
	}, nil
}

// loadAllDetails loads all eent details from a json file
func loadAllDetails(fp string) ([]models.EventDetails, error) {
	// Open the JSON file
	file, err := os.Open(fp)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	var eventDetails []models.EventDetails
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&eventDetails)

	if err != nil {
		return nil, err
	}

	return eventDetails, nil
}

// EventDetails loads a single event detail from the details json file
func EventDetails(eventID string, fp string) (models.EventDetails, error) {
	eventDetails, err := loadAllDetails(fp)
	var ed models.EventDetails
	if err != nil {
		return ed, err
	}
	for _, ed = range eventDetails {
		if eventID == ed.ID() {
			return ed, nil
		}
	}
	return ed, notFoundException{
		ID:   eventID,
		Path: fp,
	}
}
