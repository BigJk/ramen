package console

import (
	"regexp"
	"strings"

	"github.com/BigJk/ramen/concolor"

	"github.com/BigJk/ramen/t"
)

var colorSectionRegex = regexp.MustCompile("\\[\\[(([bf]):(#[0-9a-zA-Z]+))(\\|(([bf]):(#[0-9a-zA-Z]+)))?\\]\\]")

// ColorSection represents a colorized section in a text.
type ColorSection struct {
	Index       int
	Transformer []t.Transformer
}

// ColorSections represents a slice of color sections.
type ColorSections []*ColorSection

// GetCurrent gets the current color section for the given index in a string.
func (cs ColorSections) GetCurrent(index int) *ColorSection {
	for i := len(cs) - 1; i >= 0; i-- {
		if index >= cs[i].Index {
			return cs[i]
		}
	}
	return nil
}

// GetCurrentTransformer gets the transformers for the current color section
// for the given index in a string.
func (cs ColorSections) GetCurrentTransformer(index int) []t.Transformer {
	c := cs.GetCurrent(index)
	if c != nil {
		return c.Transformer
	}
	return []t.Transformer{}
}

// ParseColoredText parses the coloring annotations in a string and returns the cleaned string
// and the parsed color sections.
func ParseColoredText(text string) (string, ColorSections) {
	var results []*ColorSection
	for {
		match := colorSectionRegex.FindStringSubmatch(text)
		if len(match) <= 1 {
			break
		}

		var cs ColorSection
		cs.Index = strings.Index(text, match[0])
		text = strings.Replace(text, match[0], "", 1)

		if len(match[4]) == 0 {
			switch match[2] {
			case "f":
				cs.Transformer = append(cs.Transformer, t.Foreground(concolor.MustHex(match[3])))
			case "b":
				cs.Transformer = append(cs.Transformer, t.Background(concolor.MustHex(match[3])))
			}
		} else if len(match) == 8 {
			switch match[2] {
			case "f":
				cs.Transformer = append(cs.Transformer, t.Foreground(concolor.MustHex(match[3])))
			case "b":
				cs.Transformer = append(cs.Transformer, t.Background(concolor.MustHex(match[3])))
			}

			switch match[6] {
			case "f":
				cs.Transformer = append(cs.Transformer, t.Foreground(concolor.MustHex(match[7])))
			case "b":
				cs.Transformer = append(cs.Transformer, t.Background(concolor.MustHex(match[7])))
			}
		}

		results = append(results, &cs)
	}

	return text, results
}
