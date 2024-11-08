// Package selectors contains the selectors for source HTML
package selectors

const (
	EventSummary           = ".summary"
	EventProductInfoMobile = ".product-info-mobile"
	EventProductInfo       = ".product-info"
	Events                 = ".products li"
	OutOfStock             = "p.out-of-stock"
	ParagraphNoClass       = "p:not([class])"
	ImageLink              = "img"
	EventProductInfoParts  = "span"
	EventLink              = "a:first-child"
	EventStoreLink         = "a:nth-child(2)"
	Headline2              = "h2:first-child"
	BulletPoints           = "div:last-child > span"
	EventSmallImage        = "img"
	EventWeekDay           = "p.datetime > span:first-child"
	EventDate              = "p.datetime > span.date"
	EventHeadliners        = "h2 span"
	EventBackupHeadliners  = "h2"
	EventPlayTimes         = ".play-times > div > p"
	EventTickets           = ".variations_form"
	DoorPrice              = ".add-to-cart-wrapper"
	TicketPrices           = ".single-variation"
)
