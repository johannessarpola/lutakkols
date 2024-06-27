package fetch

import (
	"github.com/gocolly/colly/v2"
	"github.com/johannessarpola/lutakkols/pkg/logger"
)

const UserAgent = "Lutakko CLI (beta)"

// collectorOptions returns the options used with customized collector
func collectorOptions() []colly.CollectorOption {
	return []colly.CollectorOption{
		colly.UserAgent(UserAgent),
	}
}

// newCollector setups the customized collector used within the application
func newCollector() *colly.Collector {
	c := colly.NewCollector(collectorOptions()...)
	c.OnRequest(func(request *colly.Request) {
		logger.Log.Debugf("visiting %s", request.URL.String())
	})
	c.OnResponse(func(r *colly.Response) {
		if r.StatusCode == 200 {
			logger.Log.Debugf("fetched %s successfully", r.Request.URL)
		} else {
			logger.Log.Errorf("fetching %s failed with %d", r.Request.URL, r.StatusCode)
		}
	})
	return c
}
