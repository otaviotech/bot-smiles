package internal

import (
	"context"
	"log"
	"net/http"
	"regexp"
)

type Crawler struct {
	httpClient    Doer
	htmlInspector HTMLInspector
	emailNotifier EmailNotifier
}

func NewCrawler(httpClient Doer, htmlInspector HTMLInspector, emailNotifier EmailNotifier) Crawler {
	return Crawler{
		httpClient:    httpClient,
		htmlInspector: htmlInspector,
		emailNotifier: emailNotifier,
	}
}

type CrawlInput struct {
	PromotionsURL string
	UserAgent     string
	EmailToNotify string
	Regexes       []string
}

func (c *Crawler) Crawl(ctx context.Context, input CrawlInput) error {
	logger := log.Default()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, input.PromotionsURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", input.UserAgent)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	promotions, err := c.htmlInspector.ExtractPromotions(res.Body)
	if err != nil {
		return err
	}

	promotionsToNotify := []Promotion{}

	for _, promotion := range promotions {
		for _, regex := range input.Regexes {
			matched, _ := regexp.MatchString(regex, promotion.Title)

			if matched {
				promotionsToNotify = append(promotionsToNotify, promotion)
				break
			}
		}
	}

	logger.Printf("Found %d promotions, %d matched the regexes.", len(promotions), len(promotionsToNotify))

	if len(promotionsToNotify) == 0 {
		return nil
	}

	logger.Printf("Notifying %d promotions to %s.", len(promotionsToNotify), input.EmailToNotify)

	return c.emailNotifier.Notify(input.EmailToNotify, promotionsToNotify)
}
