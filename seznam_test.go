package seznam

import (
	"testing"

	"net/url"

	"github.com/slovnik/slovnik"
)

func TestCorrectQuery(t *testing.T) {
	cases := []struct {
		word  string
		lang  slovnik.Language
		url   string
		query string
	}{
		{"hlavní", slovnik.Cz, "https://slovnik.seznam.cz/cz-ru/", "q=hlavní&shortView=0"},
		{"привет", slovnik.Ru, "https://slovnik.seznam.cz/ru-cz/", "q=привет&shortView=0"},
		{"sиniy", slovnik.Ru, "https://slovnik.seznam.cz/ru-cz/", "q=sиniy&shortView=0"},
	}

	for _, c := range cases {
		t.Run(c.word, func(t *testing.T) {
			resultURL, _ := prepareQuery(c.word, c.lang)

			values, _ := url.ParseQuery(c.query)
			expectedURL, _ := url.Parse(c.url)
			expectedURL.RawQuery = values.Encode()
			if resultURL.String() != expectedURL.String() {
				t.Errorf("prepareQuery url == %q, want %q", resultURL, expectedURL)
			}
		})

	}
}
