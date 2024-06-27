package offline

import (
	"fmt"
	"os"
	"testing"
)

func TestOfflineProvider(t *testing.T) {

	wd, _ := os.Getwd()
	fmt.Printf("running on %s", wd)
	edp := "test_data/event_details_test.json"
	ep := "test_data/events_test.json"
	placeholderGen := func(_ string, _ string) string { return "" }

	ofp := New(ep, edp, placeholderGen)

	events, err := ofp.GetEvents()
	if events == nil {
		t.Errorf("events is nil")
		return
	}

	if len(events.Events) == 0 || err != nil {
		t.Errorf("No event found for %s", ep)
		return
	}

	event := events.Events[0]
	details, err := ofp.GetDetails(event.ID(), event.EventURL())
	if details == nil || err != nil {
		t.Errorf("err getting details for %s", event.Headline)
	}

}
