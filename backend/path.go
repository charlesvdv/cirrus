package cirrus

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
)

// ParsePath takes a raw string and returns either a path or an error
// if the path is not correct.
// This function does not check if the path actually exists.
func ParsePath(rawPath string) (Path, error) {
	var elements []string

	var currentElement []rune
	escaped := false
	for _, runeValue := range rawPath {
		switch runeValue {
		case '/':
			if escaped {
				escaped = false
				currentElement = append(currentElement, runeValue)
				continue
			}

			if len(currentElement) != 0 {
				elements = append(elements, string(currentElement))
				currentElement = make([]rune, 0)
			}
		case '\\':
			if escaped {
				escaped = false
				currentElement = append(currentElement, runeValue)
			} else {
				escaped = true
			}
		default:
			if escaped {
				return Path{}, fmt.Errorf("invalid escape character: '\\%v'", runeValue)
			}
			currentElement = append(currentElement, runeValue)
		}
	}

	if len(currentElement) != 0 {
		elements = append(elements, string(currentElement))
	}

	return Path{
		elements: elements,
	}, nil
}

// MustParsePath is the same function as `ParsePath` but panic if
// a path is invalid.
// This should only be used in test code.
func MustParsePath(rawPath string) Path {
	path, err := ParsePath(rawPath)
	if err != nil {
		log.Fatal().Err(err).Msg("invalid path")
	}
	return path
}

// Path describes a filesystem path.
type Path struct {
	elements []string
}

// Element gets the element at a give index.
// The index must be known to be valid (0 < i < ElementCount())
func (p Path) Element(i int) string {
	return p.elements[i]
}

// Elements gets a list of path elements.
func (p Path) Elements() []string {
	return append([]string{}, p.elements...)
}

// ElementCount gets the number of path element.
func (p Path) ElementCount() int {
	return len(p.elements)
}

// IsRoot checks if the path is a root path (ie. the top-level of the filesystem)
func (p Path) IsRoot() bool {
	return p.ElementCount() == 0
}

// AppendChild returns a path resolving to the path/elem
func (p Path) AppendChild(elem string) Path {
	return Path{
		elements: append(p.elements, elem),
	}
}

func (p Path) String() string {
	return strings.Join(p.elements, "/")
}
