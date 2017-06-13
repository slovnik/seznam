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
	fastMeaningsId    = "fastMeanings"
)

// parsePage parses html structure from provided reader and scrap the data from it.
// There may be 2 types of pages:
// - Page with full description for specific word, like: https://slovnik.seznam.cz/cz-ru/?q=hlavn%C3%AD
// - Mistype page that shows you that you might have spelled word wrong and showing which word it can
//   actually be, like: https://slovnik.seznam.cz/cz-ru/?q=dobr
func parsePage(pageBody io.Reader) ([]*slovnik.Word, error) {
	doc, err := html.Parse(pageBody)

	if err != nil {
		return nil, err
	}

	results := getResultsNode(doc)
	attrs := attributes(results.Attr)

	buf := new(bytes.Buffer)
	err = html.Render(buf, results)

	if err != nil {
		return nil, err
	}

	tokenizer := html.NewTokenizer(buf)

	if attrs.class() == "transl" {
		return processSingleWord(tokenizer), nil
	}

	return processMistype(tokenizer), nil
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

// processMistype is used for processing mistyped results
func processMistype(z *html.Tokenizer) (result []*slovnik.Word) {
	var w *slovnik.Word
	var prevToken html.Token

	for tokenType := z.Next(); tokenType != html.ErrorToken; {
		token := z.Token()
		if tokenType == html.TextToken {
			if prevToken.DataAtom == atom.Span {
				addTranslation(w, token.Data)
			} else if prevToken.DataAtom == atom.A && prevToken.Type == html.StartTagToken {
				w = &slovnik.Word{Word: token.Data}
				result = append(result, w)
			}
		}
		prevToken = token
		tokenType = z.Next()
	}
	return
}

// processSingleWord is used for parsing full translation of the word with samples of synonyms, antonyms, etc.
func processSingleWord(z *html.Tokenizer) (result []*slovnik.Word) {
	blockName := ""
	funcs := map[string]func(*slovnik.Word, string){
		synonymsHeader:     addSynonym,
		antonymsHeader:     addAntonym,
		derivedWordsHeader: addDerivedWord,
	}

	var w *slovnik.Word
	var prevToken html.Token

	for tokenType := z.Next(); tokenType != html.ErrorToken; {
		token := z.Token()
		attrs := attributes(token.Attr)

		switch {
		case tokenType == html.StartTagToken:
			if attrs.class() == otherMeaningClass {
				if f, ok := funcs[blockName]; ok {
					processBlock(z, w, f)
					blockName = ""
				}
			}

			if token.DataAtom == atom.Ul && attrs.id() == "fulltext" {
				sample := processSample(z)
				w.Samples = append(w.Samples, sample)
			}

			if attrs.id() == fastMeaningsId {
				processTranslations(z, w)
			}

		case tokenType == html.TextToken:
			if prevToken.DataAtom == atom.H3 && attributes(prevToken.Attr).lang() != "" {
				w = &slovnik.Word{Word: token.Data}
				result = append(result, w)
			}

			if prevToken.DataAtom == atom.Span && attributes(prevToken.Attr).class() == "morf" {
				w.WordType = token.Data
			}

			if prevToken.DataAtom == atom.P && attributes(prevToken.Attr).class() == "morf" {
				blockName = token.Data
			}
		}
		prevToken = token
		tokenType = z.Next()
	}
	return
}

func processTranslations(z *html.Tokenizer, w *slovnik.Word) {

	var prevToken html.Token
	var prevClosingToken html.Token

	for tokenType := z.Next(); tokenType != html.ErrorToken; {
		token := z.Token()

		if token.DataAtom == atom.Div {
			return
		}

		switch tokenType {
		case html.EndTagToken:
			fallthrough
		case html.SelfClosingTagToken:
			prevClosingToken = token
		case html.TextToken:
			if prevToken.DataAtom == atom.Span && prevToken.Type == html.StartTagToken {
				if attributes(prevToken.Attr).class() != "comma" {
					updateLastTranslation(w, token.Data)
				}
			}

			if prevToken.DataAtom == atom.A {
				trimmed := strings.TrimSpace(token.Data)
				if len(trimmed) > 0 {
					if prevClosingToken.DataAtom == atom.A {
						updateLastTranslation(w, trimmed)
					} else {
						addTranslation(w, trimmed)
					}
				}

			}
		}

		prevToken = token
		tokenType = z.Next()
	}
}

func processBlock(z *html.Tokenizer, w *slovnik.Word, functor func(*slovnik.Word, string)) {
	var prevToken html.Token
	for tt := z.Next(); !(tt == html.EndTagToken && z.Token().DataAtom == atom.Div); {
		token := z.Token()

		if prevToken.DataAtom == atom.A {
			functor(w, token.Data)
		}

		prevToken = token
		tt = z.Next()
	}
	return
}

func processSample(z *html.Tokenizer) (result slovnik.SampleUse) {
	spanCount := 0

	var prevToken html.Token

	for tokenType := z.Next(); tokenType != html.ErrorToken; {
		token := z.Token()

		switch tokenType {
		case html.EndTagToken:
			if token.DataAtom == atom.Ul {
				return
			}
			if token.DataAtom == atom.Span {
				spanCount = spanCount + 1
			}
		case html.TextToken:
			if prevToken.DataAtom == atom.A && prevToken.Type == html.StartTagToken {
				result.Keyword = strings.TrimSpace(token.Data)
			}

			if prevToken.DataAtom == atom.Span && spanCount == 0 {
				result.Phrase = strings.TrimSpace(token.Data)
			}

			if spanCount == 2 && len(result.Translation) == 0 {
				result.Translation = strings.TrimSpace(token.Data)
			}
		}
		prevToken = token
		tokenType = z.Next()
	}

	return
}

func addTranslation(w *slovnik.Word, data string) {
	w.Translations = append(w.Translations, data)
}

func updateLastTranslation(w *slovnik.Word, data string) {
	if len(w.Translations) > 0 {
		lastTranslation := w.Translations[len(w.Translations)-1]
		lastTranslation = lastTranslation + " " + data
		w.Translations[len(w.Translations)-1] = lastTranslation
	} else {
		addTranslation(w, data)
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
