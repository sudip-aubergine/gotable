package gotable

import (
	"io/ioutil"
	"path"
	"strconv"
	"strings"

	"github.com/kardianos/osext"
)

// ==========
// TEXT UTILS
// ==========

// getMultiLineText used to get multi line texts,
// from one long string which length exceeds by given column width
// it tries to split the string and store that splitted line slice such a way that
// string content fits nearly in the cell of given column width
func getMultiLineText(v string, colWidth int) ([]string, int) {
	var a []string

	// fit the content in one line whatever it is irrespective of column width
	if colWidth < 1 {
		a = append(a, v)
		return a, -1
	}

	// get multi line chunk in form of array
	sa := strings.Split(v, " ") // break up the string at the spaces
	j := 0
	maxColWidth := 0
	for i := 0; i < len(sa); i++ { // spin through all substrings
		if len(sa[i]) <= colWidth && i+1 < len(sa) { // if the width of this substring is less than the requested width, and we're not at the end of the list
			s := sa[i]                         // we know we're adding this one
			for k := i + 1; k < len(sa); k++ { // take as many as possible
				if len(s)+len(sa[k])+1 <= colWidth { // if it fits...
					s += " " + sa[k] // ...add it to the list...
					i = k            // ...and keep loop in sync
				} else {
					break // otherwise, add what we have and then go back to the outer loop
				}
			}
			a = append(a, s)
		} else {
			a = append(a, sa[i])
		}
		if len(a[j]) > maxColWidth { // if there's not enough room for the current string
			maxColWidth = len(a[j]) // then adjust the max column width we need
		}
		j++
	}
	return a, maxColWidth
}

// mkstr returns a string of n of the supplied character that is the specified length
func mkstr(n int, c byte) string {
	p := make([]byte, n)
	for i := 0; i < n; i++ {
		p[i] = c
	}
	return string(p)
}

// stringln
// For text output we want at least one "\n" at the end of a section or title.
// If the supplied string does not end in "\n", then one will be appended to it
// on the return value.
func stringln(s string) string {
	if len(s) == 0 {
		return ""
	}
	if s[len(s)-1] != '\n' {
		return s + "\n"
	}
	return s
}

// ==========
// HTML UTILS
// ==========

// getCSSMapKeyForCell format and returns key for cell for css properties usage
func getCSSMapKeyForCell(rowIndex, colIndex int) string {
	return `row:` + strconv.Itoa(rowIndex) + `-col:` + strconv.Itoa(colIndex)
}

// getCSSForClassSelector returns css string for a class
func getCSSForClassSelector(className string, cssList []*CSSProperty) string {
	var classCSS string

	// append notation for selector
	classCSS += `.` + className + `{`

	for _, cssProp := range cssList {
		// append css property name
		classCSS += cssProp.Name + `:` + cssProp.Value + `;`
	}

	// finally block ending sign
	classCSS += `}`

	// return class css string
	return classCSS
}

// getCSSForHTMLTag return css string for html tag element
func getCSSForHTMLTag(tagEl string, cssList []*CSSProperty) string {
	var classCSS string

	// append notation for selector
	classCSS += tagEl + `{`

	for _, cssProp := range cssList {
		// append css property name
		classCSS += cssProp.Name + `:` + cssProp.Value + `;`
	}

	// finally block ending sign
	classCSS += `}`

	// return class css string
	return classCSS
}

// getReportDefaultCSS reads default css from report.css
func getReportDefaultCSS() (string, error) {
	folderPath, err := osext.ExecutableFolder()
	if err != nil {
		return "", err
	}

	reportCSS := path.Join(folderPath, "report.css")

	cssString, err := ioutil.ReadFile(reportCSS)
	if err != nil {
		return "", err
	}
	return string(cssString), nil
}

// getTableTemplatePath returns the path of table template file
func getTableTemplatePath() (string, error) {
	folderPath, err := osext.ExecutableFolder()
	if err != nil {
		return "", err
	}

	tmpl := path.Join(folderPath, "table.tmpl")
	return tmpl, nil
}
