package seznam

import (
	"net/http"
	"net/url"

	"fmt"

	"github.com/slovnik/slovnik"
)

// Map of urls used for translation
var urls = map[slovnik.Language]string{
	slovnik.Cz: "https://slovnik.seznam.cz/cz-ru/",
	slovnik.Ru: "https://slovnik.seznam.cz/ru/",
}

// Translate word using slovnik.seznam.cz
func Translate(word string, language slovnik.Language) ([]*slovnik.Word, error) {
	query, err := prepareQuery(word, language)

	if err != nil {
		return nil, err
	}

	return getTranslations(query)
}

func prepareQuery(word string, language slovnik.Language) (*url.URL, error) {
	query, err := url.Parse(urls[language])

	if err != nil {
		return nil, err
	}

	p := url.Values{}
	p.Add("q", word)
	p.Add("shortView", "0")

	query.RawQuery = p.Encode()
	return query, nil
}

func getTranslations(url fmt.Stringer) ([]*slovnik.Word, error) {
	resp, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}
	return parsePage(resp.Body), nil
}
