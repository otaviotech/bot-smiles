package internal

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type GoQueryHTMLInspector struct{}

func (g *GoQueryHTMLInspector) ExtractPromotions(r io.ReadCloser) ([]Promotion, error) {
	doc, err := goquery.NewDocumentFromReader(r)

	if err != nil {
		return nil, err
	}

	promotions := []Promotion{}

	doc.Find(".promo-card").Each(func(i int, s *goquery.Selection) {
		expRaw := s.Find(".card-data-validade").AttrOr("data-msdate", "0")
		exp, err := strconv.ParseInt(expRaw, 10, 64)
		if err != nil {
			exp = 0
		}

		promotion := Promotion{
			Title:     strings.Replace(s.Find(".titulo-promo").Text(), "[sorriso]", " ", -1),
			URL:       fmt.Sprintf("%s%s", "https://smiles.com.br", s.Find(".promo-link-footer").AttrOr("href", "")),
			ExpiresIn: time.Unix(exp, 0),
		}

		promotions = append(promotions, promotion)
	})

	return promotions, nil
}
