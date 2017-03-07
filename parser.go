package seznam

import (
	"io"

	"github.com/slovnik/slovnik"

	"golang.org/x/net/html"
)

var f func(w *slovnik.Word, data string)

func parsePage(pageBody io.Reader) slovnik.Word {
	z := html.NewTokenizer(pageBody)

	inTranslations := false
	foundSynonymsHeader := false
	inSynonymsBlock := false
	foundAntonymsHeader := false
	inAntonymsBlock := false
	foundDerivedWordsHeader := false
	inDerivedWordsBlock := false

	prevTag := ""

	w := slovnik.Word{}
	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			return w

		case tt == html.StartTagToken:
			t := z.Token()

			if t.Data == "h3" {
				lang := getAttr(t.Attr, "lang")
				if lang == "cs" || lang == "ru" {
					f = addWord
				}
			}

			if t.Data == "div" {
				inTranslations = getAttr(t.Attr, "id") == "fastMeanings"
				inSynonymsBlock = getAttr(t.Attr, "class") == "other-meaning" && foundSynonymsHeader
				inAntonymsBlock = getAttr(t.Attr, "class") == "other-meaning" && foundAntonymsHeader
				inDerivedWordsBlock = getAttr(t.Attr, "class") == "other-meaning" && foundDerivedWordsHeader
			}

			if t.Data == "span" && inTranslations {
				if getAttr(t.Attr, "class") != "comma" {
					f = addTranslation
				}
			}

			if t.Data == "a" {
				if inTranslations {
					if prevTag == "a" {
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

			if t.Data == "span" && getAttr(t.Attr, "class") == "morf" {
				f = addWordType
			}

			prevTag = t.Data

			break

		case tt == html.SelfClosingTagToken:
			t := z.Token()
			prevTag = t.Data
			break

		case tt == html.EndTagToken:
			t := z.Token()
			if t.Data == "div" {
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
			break

		case tt == html.TextToken:
			t := z.Token()
			if f != nil {
				f(&w, t.Data)
				f = nil
			}

			switch t.Data {
			case "Synonyma":
				foundSynonymsHeader = true
				break
			case "Antonyma":
				foundAntonymsHeader = true
				break
			case "Odvozená slova":
				foundDerivedWordsHeader = true
				break
			}

			break
		}
	}
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

func getAttr(attrs []html.Attribute, name string) string {
	for _, a := range attrs {
		if a.Key == name {
			return a.Val
		}
	}

	return ""
}
