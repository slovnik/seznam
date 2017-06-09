package seznam

import (
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Attributes []html.Attribute

func (attrs Attributes) get(name atom.Atom) string {
	for _, a := range attrs {
		if a.Key == name.String() {
			return a.Val
		}
	}

	return ""
}

func (attrs Attributes) id() string {
	return attrs.get(atom.Id)
}

func (attrs Attributes) class() string {
	return attrs.get(atom.Class)
}

func (attrs Attributes) lang() string {
	return attrs.get(atom.Lang)
}
