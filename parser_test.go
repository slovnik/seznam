package seznam

import (
	"os"
	"testing"
	"fmt"
)

func TestParsePage(t *testing.T) {
	f, _ := os.Open("./test/sample.html")
	result, _ := parsePage(f)

	const expectedWord = "hlavní"
	w := result[0]

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

	if len(w.Samples) != 31 {
		t.Errorf("ParsePage len(Samples) == %d, want 31", len(w.Samples))
	}
}

func TestParseAltPage(t *testing.T) {
	f, _ := os.Open("./test/sample_issue8.html")
	result, _ := parsePage(f)
	w := result[0]

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
	result, _ := parsePage(f)
	w := result[0]

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

func TestMultipleResults(t *testing.T) {
	f, _ := os.Open("./test/sample_multiple_results.html")
	result, _ := parsePage(f)

	const expectedCount = 9

	for _, value := range result {
		fmt.Println(value)
	}

	//fmt.Println(result)

	if len(result) != expectedCount {
		t.Errorf("ParsePage len(w) == %d, want %d", len(result), expectedCount)

		return
	}

	expectedWords := []string{
		"dobrat se",
		"doba",
		"do",
		"dobrý",
		"dobro",
		"dobré",
		"dobrat",
		"obr",
		"bobr",
	}

	expectedTranslations := []string{
		"добра́ться",
		"вре́мя",
		"в",
		"хоро́ший",
		"добро́",
		"добро́",
		"израсхо́довать",
		"гига́нт",
		"бобр",
	}

	if len(result) != len(expectedWords) {
		t.Errorf("ParsePage len(result) == %d, want %d", len(result), len(expectedWords))
	}

	for i, w := range result {
		if w.Word != expectedWords[i] {
			t.Errorf("ParsePage w.Word == %s, want %s", w.Word, expectedWords[i])

			return
		}

		if len(w.Translations) == 0 {
			t.Errorf("ParsePage len(w.Translations) == %d, want 1", len(w.Translations))

			return
		}

		if w.Translations[0] != expectedTranslations[i] {
			t.Errorf("ParsePage w.Translation == %s, want %s", w.Translations[0], expectedTranslations[i])

			return
		}
	}
}
