package stiebeleltron

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type (
	ISGClient struct {
		Options ClientOptions
		client  http.Client
	}
	ClientOptions struct {
		BaseURL string
		Headers http.Header
	}
	Property interface {
		GetGroup() string
		GetSearchString() string
		SetValue(v float64)
	}
	properties []Property
	ParseError struct {
		Property Property
		RawText  string
		Error    error
	}
)

var (
	PropertyTableQueryExpression = "form#werte table.info tbody"
	NumberRegex                  = regexp.MustCompile("([-.,\\d]+)")
)

func (p properties) findProperty(group, searchString string) Property {
	for _, prop := range p {
		if prop.GetGroup() == group && prop.GetSearchString() == searchString {
			return prop
		}
	}
	return nil
}

// NewISGClient constructs a client for interacting with Stiebel Eltron ISG.
func NewISGClient(options ClientOptions) (*ISGClient, error) {
	return &ISGClient{
		Options: options,
		client:  http.Client{},
	}, nil
}

func (c *ISGClient) ParsePage(urlPath string, properties []Property) ([]ParseError, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", c.Options.BaseURL, urlPath), nil)
	if err != nil {
		return nil, err
	}
	req.Header = c.Options.Headers

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	return c.findValues(doc, properties), nil
}

func (c *ISGClient) findValues(doc *goquery.Document, properties properties) []ParseError {
	var p []ParseError
	doc.Find(PropertyTableQueryExpression).Each(func(i int, selection *goquery.Selection) {
		group := selection.Find("th").Text()
		selection.Find("tr.even,tr.odd").Each(func(i int, selection *goquery.Selection) {
			key := selection.Find("td.key").Text()

			property := properties.findProperty(group, key)
			if property == nil {
				p = append(p, ParseError{
					Error: fmt.Errorf("property found in document but not processed: %s/%s", group, key),
				})
				return
			}

			cellText := strings.TrimSpace(selection.Find("td.value").Text())
			parsed, err := c.findNumericValueInCell(cellText)
			if err != nil {
				p = append(p, ParseError{
					Property: property,
					RawText:  cellText,
					Error:    err,
				})
			}
			property.SetValue(parsed)
		})
	})
	return p
}

func (c *ISGClient) findNumericValueInCell(str string) (float64, error) {
	res := NumberRegex.FindAllStringSubmatch(str, -1)
	if len(res) == 0 {
		return 0, fmt.Errorf("could not find a match: " + NumberRegex.String())
	}
	match := res[0]
	replacedComma := strings.ReplaceAll(match[1], ",", ".")
	return strconv.ParseFloat(replacedComma, 64)
}
