package internal

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
)

const (
	Sender   = "aotavio097@gmail.com"
	Subject  = "Smiles BOT - Your promotions are here!"
	TextBody = "This email was sent with Amazon SES using the AWS SDK for Go."
	CharSet  = "UTF-8"
)

type SNSEmailNotifier struct {
	ses *ses.SES
}

func NewSNSEmailNotifier(sesClient *ses.SES) SNSEmailNotifier {
	return SNSEmailNotifier{ses: sesClient}
}

func (s *SNSEmailNotifier) Notify(email string, promotions []Promotion) error {
	htmlMessage, err := s.getHTMLMessage(email, promotions)

	if err != nil {
		return err
	}

	input := s.getSESMessageInput(email, htmlMessage)

	_, err = s.ses.SendEmail(&input)

	if err != nil {
		return err
	}

	return nil
}

func (s *SNSEmailNotifier) getSESMessageInput(recepient, htmlMessage string) ses.SendEmailInput {
	input := ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(recepient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(htmlMessage),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(Subject),
			},
		},
		Source: aws.String(Sender),
	}

	return input
}

func (s *SNSEmailNotifier) getHTMLMessage(email string, promotions []Promotion) (string, error) {
	body := strings.Builder{}

	_, err := body.WriteString("<p>Olá, encontramos as seguintes promoções:</p>")
	if err != nil {
		return "", err
	}

	for i, p := range promotions {
		entry := fmt.Sprintf(`<p>%d. <a href="%s">%s</a></p>`, i+1, p.URL, p.Title)

		_, err = body.WriteString(entry)
		if err != nil {
			return "", err
		}
	}

	return body.String(), nil
}
