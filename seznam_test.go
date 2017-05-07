package seznam

import (
	"os"
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
		{"привет", slovnik.Ru, "https://slovnik.seznam.cz/ru/", "q=привет&shortView=0"},
		{"sиniy", slovnik.Ru, "https://slovnik.seznam.cz/ru/", "q=sиniy&shortView=0"},
	}

	for _, c := range cases {
		resultURL := prepareQuery(c.word, c.lang)

		values, _ := url.ParseQuery(c.query)
		expectedURL, _ := url.Parse(c.url)
		expectedURL.RawQuery = values.Encode()
		if resultURL.String() != expectedURL.String() {
			t.Errorf("prepareQuery url == %q, want %q", resultURL, expectedURL)
		}
	}
}

func TestParsePage(t *testing.T) {
	f, _ := os.Open("./test/sample.html")
	w := parsePage(f)

	const expectedWord = "hlavní"

	if w.Word != expectedWord {
		t.Errorf("ParsePage word == %q, want %q", w.Word, expectedWord)
	}

	expectedTranslations := []string{
		"гла́вный",
		"основно́й",
		"центра́льный",
	}

	if len(w.Translations) != len(expectedTranslations) {
		t.Errorf("ParsePage len(translation) == %d, want %d", len(w.Translations), len(expectedTranslations))
	}

	for i, trans := range w.Translations {
		if trans != expectedTranslations[i] {
			t.Errorf("ParsePage translation == %q, want %q", trans, expectedTranslations[i])
		}
	}

	const expectedWordType = "přídavné jméno"
	if w.WordType != expectedWordType {
		t.Errorf("ParsePage wordType == %q, want %q", w.WordType, expectedWordType)
	}

	expectedSynonyms := []string{
		"ústřední",
		"podstatný",
		"základní",
		"zásadní",
	}

	if len(w.Synonyms) != len(expectedSynonyms) {
		t.Errorf("ParsePage len(synonyms) == %d, want %d", len(w.Synonyms), len(expectedSynonyms))
	}

	for i, synonym := range w.Synonyms {
		if synonym != expectedSynonyms[i] {
			t.Errorf("ParsePage synonym == %q, want %q", synonym, expectedSynonyms[i])
		}
	}

	expectedAntonyms := []string{
		"vedlejší",
		"podřadný",
		"podružný",
	}

	if len(w.Antonyms) != len(expectedAntonyms) {
		t.Errorf("ParsePage len(antonyms) == %d, want %d", len(w.Antonyms), len(expectedAntonyms))
	}

	for i, antonym := range w.Antonyms {
		if antonym != expectedAntonyms[i] {
			t.Errorf("ParsePage antonym == %q, want %q", antonym, expectedAntonyms[i])
		}
	}

	expectedDerivedWords := []string{
		"hlavně",
	}

	if len(w.DerivedWords) != len(expectedDerivedWords) {
		t.Errorf("ParsePage len(derivedWords) == %d, want %d", len(w.DerivedWords), len(expectedDerivedWords))
	}

	for i, derived := range w.DerivedWords {
		if derived != expectedDerivedWords[i] {
			t.Errorf("ParsePage derivedWord == %q, want %q", derived, expectedDerivedWords[i])
		}
	}
}

func TestParseAltPage(t *testing.T) {
	f, _ := os.Open("./test/sample_issue8.html")
	w := parsePage(f)

	const expectedWord = "soutěživý"

	if w.Word != expectedWord {
		t.Errorf("ParsePage word == %q, want %q", w.Word, expectedWord)
	}

	expectedTranslations := []string{
		"состяза́тельный",
	}

	if len(w.Translations) != len(expectedTranslations) {
		t.Errorf("ParsePage len(translation) == %d, want %d", len(w.Translations), len(expectedTranslations))
	}

	for i, trans := range w.Translations {
		if trans != expectedTranslations[i] {
			t.Errorf("ParsePage translation == %q, want %q", trans, expectedTranslations[i])
		}
	}

	expectedDerivedWords := []string{
		"soutěživost",
	}

	if len(w.DerivedWords) != len(expectedDerivedWords) {
		t.Errorf("ParsePage len(derivedWords) == %d, want %d", len(w.DerivedWords), len(expectedDerivedWords))
		return
	}

	for i, derived := range w.DerivedWords {
		if derived != expectedDerivedWords[i] {
			t.Errorf("ParsePage derivedWord == %q, want %q", derived, expectedDerivedWords[i])
		}
	}
}

func TestParseIssue7(t *testing.T) {
	f, _ := os.Open("./test/sample_issue7.html")
	w := parsePage(f)

	const expectedWord = "protože"

	if w.Word != expectedWord {
		t.Errorf("ParsePage word == %q, want %q", w.Word, expectedWord)
	}

	expectedTranslations := []string{
		"так как",
		"из-за того́",
		"потому́ что",
	}

	if len(w.Translations) != len(expectedTranslations) {
		t.Errorf("ParsePage len(translation) == %d, want %d", len(w.Translations), len(expectedTranslations))
		return
	}

	for i, trans := range w.Translations {
		if trans != expectedTranslations[i] {
			t.Errorf("ParsePage translation == %q, want %q", trans, expectedTranslations[i])
		}
	}
}
