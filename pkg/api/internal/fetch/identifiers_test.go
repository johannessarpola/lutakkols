package fetch

import (
	"github.com/johannessarpola/lutakkols/pkg/api/models"
	"testing"
	"time"
)

func TestCreateEventID(t *testing.T) {
	a := models.Event{
		Id:             "abc",
		Order:          0,
		Headline:       "headline",
		EventLink:      "eventLink",
		SmallImageLink: "imageLink",
		Weekday:        "weekDay",
		Date:           "date",
		StoreLink:      "storeLink",
		InStock:        true,
		UpdatedAt:      time.Now(),
	}

	b := a
	b.Id = "another"

	aId, err := createEventID(a)
	if err != nil {
		t.Errorf("Error creating event ID: %s", err)
	}
	bId, err := createEventID(b)
	if err != nil {
		t.Errorf("Error creating event ID: %s", err)
	}

	if aId != bId {
		t.Errorf("Event IDs are not the same with same content, other than ID")
	}

	c := a
	c.Id = "new_id"
	c.Headline = "new headline"
	cId, err := createEventID(c)
	if err != nil {
		t.Errorf("Error creating event ID: %s", err)
	}

	if aId == cId {
		t.Errorf("Event IDs are the same with different content")
	}

}
