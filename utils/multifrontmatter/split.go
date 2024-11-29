package multifrontmatter

import "strings"

type (
	Data struct {
		Fm, Md string
	}
	Split[T any] struct {
		Pre   string
		Pages []T
	}
)

type (
	splitter struct{ Lang, Tag string }
	Splitter interface{ Split(mds string) Split[Data] }
)

var _ Splitter = (*splitter)(nil)

type splitStep int

const (
	pre splitStep = iota
	fm
	md
)

func NewSplitter(lang, tag string) Splitter {
	return splitter{Lang: lang, Tag: tag}
}

func (s splitter) Split(mds string) Split[Data] {
	opening := "```" + s.Lang + " :" + s.Tag + ":"
	closing := "```"

	openingLen, closingLen := len(opening), len(closing)
	mdsLen := len(mds)

	step := pre
	split := Split[Data]{}

	lastBreak := 0

	for i, ch := range mds {
		switch step {
		case pre:
			if (i == 0 && mdsLen > i+openingLen && mds[i:i+openingLen] == opening) ||
				(ch == '\n' && mdsLen > i+1+openingLen && mds[i+1:i+1+openingLen] == opening) {
				step = fm
				split.Pre = strings.Trim(mds[:i], "\n ")
				lastBreak = i
			}
		case fm:
			if ch == '\n' && mdsLen >= i+1+closingLen && mds[i+1:i+1+closingLen] == closing {
				step = md
				split.Pages = append(split.Pages, Data{Fm: strings.Trim(mds[lastBreak+openingLen+1:i], "\n "), Md: ""})
				lastBreak = i
			}
		case md:
			if ch == '\n' && mdsLen > i+1+openingLen && mds[i+1:i+1+openingLen] == opening {
				step = fm
				split.Pages[len(split.Pages)-1].Md = strings.Trim(mds[lastBreak+closingLen+1:i], "\n ")
				lastBreak = i
			}
		}
	}

	if len(split.Pages) > 0 {
		split.Pages[len(split.Pages)-1].Md = strings.Trim(mds[lastBreak+closingLen+1:], "\n ")
	} else {
		split.Pre = mds
	}

	return split
}
