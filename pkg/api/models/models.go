// Package models contains the model definitions used within the application
package models

import "time"

// Event is an event with some basic information scraped from the shorter description
type Event struct {
	Id             string    `json:"id"`
	Order          int32     `json:"order"`
	Headline       string    `json:"headline"`
	EventLink      string    `json:"event_link"`
	SmallImageLink string    `json:"image_link"`
	Weekday        string    `json:"week_day"`
	Date           string    `json:"date"`
	StoreLink      string    `json:"store_ink"`
	InStock        bool      `json:"in_stock"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
}

// EventDetails is a more rich information related one to one to event and fetched from the events' page
type EventDetails struct {
	EventID     string            `json:"event_id"`
	Description []string          `json:"summary"`
	ImageLink   string            `json:"image_link"`
	ProductInfo map[string]string `json:"product_info"`
	PlayTimes   []string          `json:"play_times"`
	Tickets     EventTickets      `json:"tickets"`
	DoorPrice   DoorPrice         `json:"door_price"`
	UpdatedAt   time.Time         `json:"updated_at,omitempty"`
}

// HasID interface to allow to use general access to ID variable (usually event.ID)
type HasID interface {
	ID() string
}

func (e Event) ID() string {
	return e.Id
}
func (e Event) EventURL() string { return e.EventLink }

func (ed EventDetails) ID() string {
	return ed.EventID
}
func (ed EventDetails) ImageURL() string {
	return ed.ImageLink
}

type EventPartial interface {
	HasID
	EventURL() string
}

type EventDetailsPartial interface {
	HasID
	ImageURL() string
}

// Ticket single ticket for event
type Ticket struct {
	Description string `json:"description"`
	Price       string `json:"price"`
}

// EventTickets tickets for the event, can be emppty as well
type EventTickets struct {
	Tickets []Ticket `json:"tickets"`
}

// DoorPrice is alias of string for now
type DoorPrice = string

// EventAscii is a container for events' image which is converted into string
type EventAscii struct {
	Ascii     string    `json:"ascii"`
	EventID   string    `json:"event_id"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// Events simple wrapper for event list with a timestamp
type Events struct {
	Events    []Event   `json:"events"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
