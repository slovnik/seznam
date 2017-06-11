package seznam

import (
	"bytes"
	"io"

	"github.com/slovnik/slovnik"

	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const (
	synonymsHeader     = "Synonyma"
	antonymsHeader     = "Antonyma"
	derivedWordsHeader = "Odvozená slova"

	otherMeaningClass = "other-meaning"
	fastMeaningsClass = "fastMeanings"
)

var f func(w *slovnik.Word, data string)

func parsePage(pageBody io.Reader) []*slovnik.Word {
	doc, _ := html.Parse(pageBody)

	results := getResultsNode(doc)
	attrs := attributes(results.Attr)

	buf := new(bytes.Buffer)
	html.Render(buf, results)
	tokenizer := html.NewTokenizer(buf)

	if attrs.class() == "transl" {
		return processSingleWord(tokenizer)
	}

	return processMistype(tokenizer)
}

// getResultsNode parses page HTML to find node containing results of translation
func getResultsNode(document *html.Node) (results *html.Node) {
	var traverse func(*html.Node)

	traverse = func(n *html.Node) {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.DataAtom == atom.Div && attributes(c.Attr).id() == "results" {
				results = c
				return
			}
			traverse(c)
		}
	}
	traverse(document)
	return
}

func processMistype(z *html.Tokenizer) []*slovnik.Word {
	result := []*slovnik.Word{}
	var w *slovnik.Word

	for tokenType := z.Next(); tokenType != html.ErrorToken; {
		token := z.Token()

		switch {
		case tokenType == html.StartTagToken:
			if token.DataAtom == atom.Span {
				f = addTranslation
			}

			if token.DataAtom == atom.A {
				w = new(slovnik.Word)
				result = append(result, w)
				f = addWord
			}

		case tokenType == html.TextToken:
			if f != nil {
				f(w, token.Data)
				f = nil
			}
		}

		tokenType = z.Next()
	}
	return result
}

func processSingleWord(z *html.Tokenizer) []*slovnik.Word {
	inTranslations := false
	foundSynonymsHeader := false
	inSynonymsBlock := false
	foundAntonymsHeader := false
	inAntonymsBlock := false
	foundDerivedWordsHeader := false
	inDerivedWordsBlock := false

	prevTag := atom.Body

	result := []*slovnik.Word{}
	var w *slovnik.Word

	for tokenType := z.Next(); tokenType != html.ErrorToken; {
		token := z.Token()
		attrs := attributes(token.Attr)

		switch {
		case tokenType == html.StartTagToken:
			if token.DataAtom == atom.H3 {

				lang := attrs.lang()
				if lang == "cs" || lang == "ru" {
					w = &slovnik.Word{}
					result = append(result, w)
					f = addWord
				}
			}

			if token.DataAtom == atom.Div {
				inTranslations = attrs.id() == fastMeaningsClass
				inSynonymsBlock = attrs.class() == otherMeaningClass && foundSynonymsHeader
				inAntonymsBlock = attrs.class() == otherMeaningClass && foundAntonymsHeader
				inDerivedWordsBlock = attrs.class() == otherMeaningClass && foundDerivedWordsHeader
			}

			if token.DataAtom == atom.Span && inTranslations {
				if attrs.class() != "comma" {
					f = addTranslation
				}
			}

			if token.DataAtom == atom.Ul && attrs.id() == "fulltext" {
				sample := processSample(z)
				w.Samples = append(w.Samples, sample)
			}

			if token.DataAtom == atom.A {
				if inTranslations {
					if prevTag == atom.A {
						f = updateLastTranslation
					} else {
						f = addTranslation
					}
				} else if inSynonymsBlock {
					f = addSynonym
				} else if inAntonymsBlock {
					f = addAntonym
				} else if inDerivedWordsBlock {
					f = addDerivedWord
				}
			}

			if token.DataAtom == atom.Span && attrs.class() == "morf" {
				f = addWordType
			}

			prevTag = token.DataAtom

		case tokenType == html.SelfClosingTagToken:
			prevTag = token.DataAtom

		case tokenType == html.EndTagToken:
			if token.DataAtom == atom.Div {
				if inTranslations {
					inTranslations = false
				} else if inSynonymsBlock {
					inSynonymsBlock = false
					foundSynonymsHeader = false
				} else if inAntonymsBlock {
					inAntonymsBlock = false
					foundAntonymsHeader = false
				} else if inDerivedWordsBlock {
					inDerivedWordsBlock = false
					foundDerivedWordsHeader = false
				}
			}

		case tokenType == html.TextToken:
			if f != nil {
				f(w, token.Data)
				f = nil
			}

			switch token.Data {
			case synonymsHeader:
				foundSynonymsHeader = true
			case antonymsHeader:
				foundAntonymsHeader = true
			case derivedWordsHeader:
				foundDerivedWordsHeader = true
			}
		}

		tokenType = z.Next()
	}
	return result
}

func processSample(z *html.Tokenizer) slovnik.SampleUse {
	inWord := false
	inSpan := false
	spanCount := 0

	result := slovnik.SampleUse{}

loop:
	for tokenType := z.Next(); tokenType != html.ErrorToken; {
		token := z.Token()

		switch tokenType {
		case html.StartTagToken:
			if token.DataAtom == atom.A {
				inWord = true
			}

			if token.DataAtom == atom.Span {
				inSpan = true
			}

		case html.EndTagToken:
			if token.DataAtom == atom.A && inWord {
				inWord = false
			}

			if token.DataAtom == atom.Span {
				inSpan = false
				spanCount = spanCount + 1
			}

			if token.DataAtom == atom.Ul {
				break loop
			}
		case html.TextToken:

			if inWord {
				result.Keyword = strings.TrimSpace(token.Data)
			}

			if inSpan && spanCount == 0 {
				result.Phrase = strings.TrimSpace(token.Data)
			}

			if spanCount == 2 && len(result.Translation) == 0 {
				result.Translation = strings.TrimSpace(token.Data)
			}
		}

		tokenType = z.Next()
	}

	return result
}

func addWord(w *slovnik.Word, data string) {
	w.Word = data
}

func addWordType(w *slovnik.Word, data string) {
	w.WordType = data
}

func addTranslation(w *slovnik.Word, data string) {
	w.Translations = append(w.Translations, data)
}

func updateLastTranslation(w *slovnik.Word, data string) {
	if len(w.Translations) > 0 {
		lastTranslation := w.Translations[len(w.Translations)-1]
		lastTranslation = lastTranslation + " " + data
		w.Translations[len(w.Translations)-1] = lastTranslation
	}
}

func addSynonym(w *slovnik.Word, data string) {
	w.Synonyms = append(w.Synonyms, data)
}

func addAntonym(w *slovnik.Word, data string) {
	w.Antonyms = append(w.Antonyms, data)
}

func addDerivedWord(w *slovnik.Word, data string) {
	w.DerivedWords = append(w.DerivedWords, data)
}
