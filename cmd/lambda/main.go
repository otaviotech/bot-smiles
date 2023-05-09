package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/otaviotech/bot-smiles/internal"
)

func run(ctx context.Context, event interface{}) error {
	logger := log.Default()

	logger.Println("Started crawling.")

	c := buildCrawler()

	err := c.Crawl(ctx, internal.CrawlInput{
		PromotionsURL: os.Getenv("PROMOTIONS_URL"),
		UserAgent:     os.Getenv("USER_AGENT"),
		EmailToNotify: os.Getenv("EMAIL_TO_NOTIFY"),
		Regexes:       []string{os.Getenv("REGEXES")},
	})

	if err != nil {
		logger.Printf("Error crawling: %s", err.Error())
	}

	logger.Println("Finished crawling.")

	return err
}

func main() {
	lambda.Start(run)
}

func buildCrawler() internal.Crawler {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	sesClient := ses.New(sess)

	sesEmailNotifier := internal.NewSNSEmailNotifier(sesClient)

	return internal.NewCrawler(
		&http.Client{},
		&internal.GoQueryHTMLInspector{},
		&sesEmailNotifier,
	)
}
