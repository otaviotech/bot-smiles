package internal

import (
	"io"
	"net/http"
	"time"
)

type ContextKey string

type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

type Promotion struct {
	Title     string
	URL       string
	ExpiresIn time.Time
}

type HTMLInspector interface {
	ExtractPromotions(r io.ReadCloser) ([]Promotion, error)
}

type EmailNotifier interface {
	Notify(email string, promotions []Promotion) error
}
