package models

// HasID interface to allow to use general access to ID variable (usually event.ID)
type HasID interface {
	ID() string
}

type EventPartial interface {
	HasID
	EventURL() string
}

type EventDetailsPartial interface {
	HasID
	ImageURL() string
}

type EventAsciiPartial interface {
	HasID
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

func (ea EventAscii) ID() string {
	return ea.EventID
}
