package seznam

import (
	"net/http"
	"net/url"

	"github.com/slovnik/slovnik"
)

// Map of urls used for translation
var urls = map[slovnik.Language]string{
	slovnik.Cz: "https://slovnik.seznam.cz/cz-ru/",
	slovnik.Ru: "https://slovnik.seznam.cz/ru/",
}

// Translate word using slovnik.seznam.cz
func Translate(word string, language slovnik.Language) (slovnik.Word, error) {
	query := prepareQuery(word, language)
	return getTranslations(query)
}

func prepareQuery(word string, language slovnik.Language) *url.URL {
	query, _ := url.Parse(urls[language])

	p := url.Values{}
	p.Add("q", word)
	p.Add("shortView", "0")

	query.RawQuery = p.Encode()
	return query
}

func getTranslations(url *url.URL) (slovnik.Word, error) {
	resp, err := http.Get(url.String())
	if err != nil {
		return slovnik.Word{}, err
	}
	return parsePage(resp.Body), nil
}
