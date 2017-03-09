package gotable

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"
	"sort"
	"strconv"
	// "strings"
	"text/template"

	"github.com/dustin/go-humanize"
	"github.com/yosssi/gohtml"
)

// TABLECONTAINERCLASS et. al. are the constants used in the html version of table object
const (
	TABLECONTAINERCLASS = `rpt-table-container`
	TITLECLASS          = `title`
	HEADERSCLASS        = `headers`
	// DATACLASS           = `data`
	SECTION1CLASS = `section1`
	SECTION2CLASS = `section2`
)

// HTMLTable struct used to prepare table in html version
type HTMLTable struct {
	*Table
	StyleString string
}

// HTMLTemplateContext holds the context for table html template
type HTMLTemplateContext struct {
	FontSize                                    int
	HeadTitle, DefaultCSS, CustomCSS, TableHTML string
}

func (ht *HTMLTable) getTableOutput() (string, error) {
	var tContainer string

	// append title
	tContainer += ht.getTitle()

	// append section 1
	tContainer += ht.getSection1()

	// append section 2
	tContainer += ht.getSection2()

	// contains only table tag output
	var tableOut string

	// append headers
	headerStr, err := ht.getHeaders()
	if err != nil {
		return "", err
	}
	tableOut += headerStr

	// append rows
	rowsStr, err := ht.getRows()
	if err != nil {
		return "", err
	}
	tableOut += rowsStr

	// wrap headers and rows in a table
	tableOut = `<table>` + tableOut + `</table>`

	// now append to container of table output
	tContainer += tableOut

	// wrap it up in a div with a class
	tContainer = `<div class="` + TABLECONTAINERCLASS + `">` + tContainer + `</table>`

	// make context for template
	htmlContext := HTMLTemplateContext{FontSize: CSSFONTSIZE}
	htmlContext.HeadTitle = ht.Table.Title
	htmlContext.DefaultCSS, err = ht.getReportDefaultCSS()
	if err != nil {
		return "", err
	}
	htmlContext.DefaultCSS = `<style>` + htmlContext.DefaultCSS + `</style>`
	htmlContext.CustomCSS = `<style>` + ht.StyleString + `</style>`
	htmlContext.TableHTML = tContainer

	// get template string
	tableTmplPath, err := ht.getTableTemplatePath()
	if err != nil {
		return "", err
	}

	// Create a new template and parse the context in it
	tmpl := template.New("table.tmpl")
	tmpl, err = tmpl.ParseFiles(tableTmplPath)
	if err != nil {
		return "", err
	}

	// write html output in buffer
	var htmlBuffer bytes.Buffer
	err = tmpl.Execute(&htmlBuffer, htmlContext)
	if err != nil {
		return "", err
	}

	// return output
	tableHTML := gohtml.Format(htmlBuffer.String())
	// tableHTML = strings.Replace(tableHTML, "\\\n", "", -1)
	return tableHTML, nil
}

func (ht *HTMLTable) getTitle() string {
	title := ht.Table.GetTitle()

	if title != "" {
		if cssMap, ok := ht.Table.CSS[TITLECLASS]; ok {

			// list of css properties for this td cell
			var cellCSSProps []*CSSProperty
			for _, cssProp := range cssMap {
				cellCSSProps = append(cellCSSProps, cssProp)
			}

			// get css string for title
			ht.StyleString += `div.` + TABLECONTAINERCLASS + ` p`
			ht.StyleString += ht.getCSSForClassSelector(TITLECLASS, cellCSSProps)

		}

		return `<p class="` + TITLECLASS + `">` + title + `</p>`
	}

	// blank return
	return title
}

func (ht *HTMLTable) getSection1() string {
	section1 := ht.Table.GetSection1()

	if section1 != "" {
		if cssMap, ok := ht.Table.CSS[SECTION1CLASS]; ok {

			// list of css properties for this td cell
			var cellCSSProps []*CSSProperty
			for _, cssProp := range cssMap {
				cellCSSProps = append(cellCSSProps, cssProp)
			}

			// get css string for section1
			ht.StyleString += `div.` + TABLECONTAINERCLASS + ` p`
			ht.StyleString += ht.getCSSForClassSelector(SECTION1CLASS, cellCSSProps)
		}

		return `<p class="` + SECTION1CLASS + `">` + section1 + `</p>`
	}

	// blank return
	return section1
}

func (ht *HTMLTable) getSection2() string {
	section2 := ht.Table.GetSection2()

	if section2 != "" {
		if cssMap, ok := ht.Table.CSS[SECTION2CLASS]; ok {

			// list of css properties for this td cell
			var cellCSSProps []*CSSProperty
			for _, cssProp := range cssMap {
				cellCSSProps = append(cellCSSProps, cssProp)
			}

			// get css string for section2
			// get css string for section1
			ht.StyleString += `div.` + TABLECONTAINERCLASS + ` p`
			ht.StyleString += ht.getCSSForClassSelector(SECTION2CLASS, cellCSSProps)
		}

		return `<p class="` + SECTION2CLASS + `">` + section2 + `</p>`
	}

	// blank return
	return section2
}

func (ht *HTMLTable) getHeaders() (string, error) {

	// check for blank headers
	blankHdrsErr := ht.Table.HasHeaders()
	if blankHdrsErr != nil {
		return "", blankHdrsErr
	}

	// format headers
	var tHeaders string

	for headerIndex := 0; headerIndex < len(ht.Table.ColDefs); headerIndex++ {

		headerCell := ht.Table.ColDefs[headerIndex]

		// css class for this header cell
		thClass := `header-` + strconv.Itoa(headerIndex)

		// list of css property for th cells
		var cellCSSProps []*CSSProperty

		// decide align property
		alignProp := &CSSProperty{Name: "text-align"}
		if headerCell.Justify == COLJUSTIFYRIGHT {
			alignProp.Value = "right"
		} else if headerCell.Justify == COLJUSTIFYLEFT {
			alignProp.Value = "left"
		}
		// append align css property
		cellCSSProps = append(cellCSSProps, alignProp)
		// apply this property to all cells belong to this column
		ht.Table.SetColCSS(headerIndex, cellCSSProps)

		// NOTE: width calculatation should be done after alignment
		// width only needs to be set on header cells only not on all
		// cells belong to column
		var colWidth string

		if headerCell.HTMLWidth != -1 {
			// calculate column width based on characters with font size
			colWidth = strconv.Itoa(headerCell.HTMLWidth) + CSSFONTUNIT
		} else {
			// calculate column width based on characters with font size
			colWidth = strconv.Itoa(ht.Table.ColDefs[headerIndex].Width*CSSFONTSIZE) + CSSFONTUNIT
		}

		// append width css property
		cellCSSProps = append(cellCSSProps, &CSSProperty{Name: "width", Value: colWidth})

		if cssMap, ok := ht.Table.CSS[HEADERSCLASS]; ok {
			// list of css properties for this td cell
			for _, cssProp := range cssMap {
				cellCSSProps = append(cellCSSProps, cssProp)
			}
		}

		// get css string for headers
		ht.StyleString += `div.` + TABLECONTAINERCLASS + ` table thead tr th`
		// ht.StyleString += `div.` + TABLECONTAINERCLASS + ` table thead.` + HEADERSCLASS + ` tr th`
		ht.StyleString += ht.getCSSForClassSelector(thClass, cellCSSProps)

		tHeaders += `<th class="` + thClass + `">` + headerCell.ColTitle + `</th>`
	}

	return `<thead><tr>` + tHeaders + `</tr></thead>`, nil
	// return `<thead class="` + HEADERSCLASS + `"><tr>` + tHeaders + `</tr></thead>`, nil
}

func (ht *HTMLTable) getRows() (string, error) {

	// check for empty data table
	blankDataErr := ht.Table.HasData()
	if blankDataErr != nil {
		return "", blankDataErr
	}

	var rowsStr string
	for i := 0; i < ht.Table.RowCount(); i++ {
		// for valid row, we will never get an error
		s, _ := ht.getRow(i)
		rowsStr += s
	}

	return `<tbody>` + rowsStr + `</tbody>`, nil
}

func (ht *HTMLTable) getRow(rowIndex int) (string, error) {

	// This method is only called by internal instance of TextTable
	// in getRows method, so we should avoid following error check
	// unless we make it as export

	// // check that this passed rowIndex is valid or not
	// inValidRowErr := ht.Table.HasValidRow(rowIndex)
	// if inValidRowErr != nil {
	// 	return "", inValidRowErr
	// }

	// format table rows
	var tRow string
	var trClass string

	if len(ht.Table.LineBefore) > 0 {
		j := sort.SearchInts(ht.Table.LineBefore, rowIndex)
		if j < len(ht.Table.LineBefore) && rowIndex == ht.Table.LineBefore[j] {
			trClass += `top-line`
		}
	}

	// fill the content in rowTextList for the first line
	for colIndex := 0; colIndex < len(ht.Table.Row[rowIndex].Col); colIndex++ {

		var rowCell string
		// append content in TD
		switch ht.Table.Row[rowIndex].Col[colIndex].Type {
		case CELLFLOAT:
			rowCell = fmt.Sprintf(ht.Table.ColDefs[colIndex].Pfmt, humanize.FormatFloat("#,###.##", ht.Table.Row[rowIndex].Col[colIndex].Fval))
		case CELLINT:
			rowCell = fmt.Sprintf(ht.Table.ColDefs[colIndex].Pfmt, ht.Table.Row[rowIndex].Col[colIndex].Ival)
		case CELLSTRING:
			// ******************************************************
			// FOR HTML, APPEND FULL STRING, THERE ARE NO
			// MULTILINE TEXT IN THIS
			// ******************************************************
			rowCell = fmt.Sprintf("%s", ht.Table.Row[rowIndex].Col[colIndex].Sval)
		case CELLDATE:
			rowCell = fmt.Sprintf("%*.*s", ht.Table.ColDefs[colIndex].Width, ht.Table.ColDefs[colIndex].Width, ht.Table.Row[rowIndex].Col[colIndex].Dval.Format(ht.Table.DateFmt))
		case CELLDATETIME:
			rowCell = fmt.Sprintf("%*.*s", ht.Table.ColDefs[colIndex].Width, ht.Table.ColDefs[colIndex].Width, ht.Table.Row[rowIndex].Col[colIndex].Dval.Format(ht.Table.DateTimeFmt))
		default:
			rowCell = mkstr(ht.Table.ColDefs[colIndex].Width, ' ')
		}

		// format td cell with custom class if exists for it
		g := ht.Table.getCSSMapKeyForCell(rowIndex, colIndex)
		if cssMap, ok := ht.Table.CSS[g]; ok {

			tdClass := `cell-row-` + strconv.Itoa(rowIndex) + `-col-` + strconv.Itoa(colIndex)

			// list of css properties for this td cell
			var cellCSSProps []*CSSProperty
			for _, cssProp := range cssMap {
				cellCSSProps = append(cellCSSProps, cssProp)
			}

			// get css string for a row
			ht.StyleString += `div.` + TABLECONTAINERCLASS + ` table tbody tr td`
			ht.StyleString += ht.getCSSForClassSelector(tdClass, cellCSSProps)

			rowCell = `<td class="` + tdClass + `">` + rowCell + `</td>`
		} else {
			rowCell = `<td>` + rowCell + `</td>`
		}

		tRow += rowCell
	}

	if len(ht.Table.LineAfter) > 0 {
		j := sort.SearchInts(ht.Table.LineAfter, rowIndex)
		if j < len(ht.Table.LineAfter) && rowIndex == ht.Table.LineAfter[j] {
			trClass += `bottom-line`
		}
	}

	if trClass != "" {
		return `<tr class="` + trClass + `">` + tRow + `</tr>`, nil
	}
	return `<tr>` + tRow + `</tr>`, nil

}

// getCSSForClassSelector returns css string for a class
func (ht *HTMLTable) getCSSForClassSelector(className string, cssList []*CSSProperty) string {
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
func (ht *HTMLTable) getCSSForHTMLTag(tagEl string, cssList []*CSSProperty) string {
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
func (ht *HTMLTable) getReportDefaultCSS() (string, error) {
	reportCSS := path.Join(ht.Table.Container, "report.css")

	cssString, err := ioutil.ReadFile(reportCSS)
	if err != nil {
		return "", err
	}
	return string(cssString), nil
}

// getTableTemplatePath returns the path of table template file
func (ht *HTMLTable) getTableTemplatePath() (string, error) {
	tmpl := path.Join(ht.Table.Container, "table.tmpl")
	return tmpl, nil
}
