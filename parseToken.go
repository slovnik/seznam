package seznam

import (
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type parseToken struct {
	html.Token
}

func (pt parseToken) getAttr(name atom.Atom) string {
	for _, a := range pt.Attr {
		if a.Key == name.String() {
			return a.Val
		}
	}

	return ""
}

func (pt parseToken) id() string {
	return pt.getAttr(atom.Id)
}

func (pt parseToken) class() string {
	return pt.getAttr(atom.Class)
}

func (pt parseToken) lang() string {
	return pt.getAttr(atom.Lang)
}
