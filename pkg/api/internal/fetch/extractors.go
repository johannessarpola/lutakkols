package fetch

import (
	"github.com/gocolly/colly/v2"
	"github.com/johannessarpola/lutakkols/pkg/api/internal/fetch/selectors"
	"github.com/johannessarpola/lutakkols/pkg/api/models"
	"github.com/johannessarpola/lutakkols/pkg/logger"
	"regexp"
	"strings"
)

func extractInStock(e *colly.HTMLElement) bool {
	s := strings.TrimSpace(e.ChildText(selectors.OutOfStock))
	return len(s) == 0
}

func extractImageLink(e *colly.HTMLElement) string {
	imageLink := e.ChildAttr(selectors.ImageLink, "src")
	return imageLink
}

func stripHTML(input string) string {
	replacer := strings.NewReplacer(
		"<p>", "", "</p>", "",
		"<b>", "", "</b>", "",
		"<i>", "", "</i>", "",
	)
	return replacer.Replace(input)
}

func extractSummary(e *colly.HTMLElement) []string {
	//	summary := strings.Join(e.ChildTexts(paragraphNoClassSelector), "\n")
	var items []string
	e.ForEach(selectors.ParagraphNoClass, func(i int, element *colly.HTMLElement) {
		htmlContent, _ := element.DOM.Html()

		withLineBreaks := strings.ReplaceAll(htmlContent, "<br>", "\n")
		withLineBreaks = strings.ReplaceAll(withLineBreaks, "<br/>", "\n")

		withoutHtml := stripHTML(withLineBreaks)
		items = append(items, withoutHtml)
	})

	return items
}

func extractProdductInfo(e *colly.HTMLElement) map[string]string {
	// parse types from product info for event info
	cts := e.ChildTexts(selectors.EventProductInfoParts)
	cnt := len(cts)
	productInfo := make(map[string]string, cnt)

	for i := 0; i < len(cts)-1; i += 2 {
		k := strings.Replace(cts[i], ":", "", 1)
		v := ""
		if i+1 < len(cts) {
			v = cts[i+1]
		}
		productInfo[k] = v
	}

	return productInfo
}

func extractPlayTimes(e *colly.HTMLElement) []string {
	var playtimes []string
	e.ForEach(selectors.EventPlayTimes, func(i int, element *colly.HTMLElement) {
		playtimes = append(playtimes, element.Text)
	})
	return playtimes
}

func extractTicketPrices(e *colly.HTMLElement) models.EventTickets {
	var tickets []models.Ticket
	e.ForEach(selectors.TicketPrices, func(i int, element *colly.HTMLElement) {
		description := element.ChildText("h3")
		price := element.ChildText("bdi")
		tickets = append(tickets, models.Ticket{
			Description: description,
			Price:       price,
		})
	})

	return models.EventTickets{
		Tickets: tickets,
	}
}

func extractDoorPrice(e *colly.HTMLElement) models.DoorPrice {
	magicString := "hinta ovelta"
	cs := e.ChildTexts(selectors.ParagraphNoClass)
	for _, c := range cs {
		l := strings.ToLower(c)
		r := strings.Contains(l, magicString)
		if r {
			dp := strings.TrimSpace(strings.Replace(l, magicString, "", 1))
			return dp
		}
	}

	return ""
}

func extractEvent(e *colly.HTMLElement) models.Event {

	eventLink := e.ChildAttr(selectors.EventLink, "href")
	storeLink := e.ChildAttr(selectors.EventStoreLink, "href")
	imageLink := e.ChildAttr(selectors.EventSmallImage, "src")
	weekDay := strings.TrimSpace(e.ChildText(selectors.EventWeekDay))
	date := strings.TrimSpace(e.ChildText(selectors.EventDate))
	headline := extractHeadline(e)
	inStock := extractInStock(e)

	evt := models.Event{}
	evt.EventLink = eventLink
	evt.StoreLink = storeLink
	evt.SmallImageLink = imageLink
	evt.Weekday = weekDay
	evt.Date = date
	evt.Headline = cleanupHeadline(headline)
	evt.InStock = inStock

	id, err := createEventID(evt)
	if err != nil {
		// todo fall back
		logger.Log.Warnf("could not resolve key for event %s", err.Error())
	}
	evt.Id = id
	return evt
}

func extractHeadline(e *colly.HTMLElement) string {
	var headliners []string
	e.ForEach(selectors.EventHeadliners, func(i int, element *colly.HTMLElement) {
		headliners = append(headliners, strings.TrimSpace(element.Text))
	})

	if len(headliners) == 0 {
		return e.ChildText(selectors.EventBackupHeadliners)
	} else {
		return strings.Join(headliners, " | ")
	}
}

func cleanupHeadline(input string) string {
	re := regexp.MustCompile(`\s+`)
	output := re.ReplaceAllString(input, " ")
	return output
}
