package gotable

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/dustin/go-humanize"
)

const (
	TABLECLASS    = `rpt-table`
	TITLECLASS    = `title`
	HEADERSCLASS  = `headers`
	SECTION1CLASS = `section1`
	SECTION2CLASS = `section2`
)

// HTMLTable struct used to prepare table in html version
type HTMLTable struct {
	*Table
	StyleString string
}

func (ht *HTMLTable) getTableOutput() (string, error) {

	var tout string

	// append title
	tout += ht.getTitle()

	// append section 1
	tout += ht.getSection1()

	// append section 2
	tout += ht.getSection2()

	// append headers
	headerStr, err := ht.getHeaders()
	if err != nil {
		return "", err
	}
	tout += headerStr

	// append rows
	rowsStr, err := ht.getRows()
	if err != nil {
		return "", err
	}
	tout += rowsStr

	// return output
	return `<style>` + ht.StyleString + `</style>` + `<table class="` + TABLECLASS + `">` + tout + `</table>`, nil
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

			// get css string for this td cell
			ht.StyleString += `table.` + TABLECLASS + ` tr.` + TITLECLASS + ` `
			ht.StyleString += getCSSForHTMLTag(`td`, cellCSSProps)

		}
		colSpan := strconv.Itoa(ht.Table.ColCount())
		return `<tr class="` + TITLECLASS + `"><td colspan="` + colSpan + `">` + title + `</td></tr>`
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

			// get css string for this td cell
			ht.StyleString += `table.` + TABLECLASS + ` tr.` + SECTION1CLASS + ` `
			ht.StyleString += getCSSForHTMLTag(`td`, cellCSSProps)

		}
		colSpan := strconv.Itoa(ht.Table.ColCount())
		return `<tr class="` + SECTION1CLASS + `"><td colspan="` + colSpan + `">` + section1 + `</td></tr>`
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

			// get css string for this td cell
			ht.StyleString += `table.` + TABLECLASS + ` tr.` + SECTION2CLASS + ` `
			ht.StyleString += getCSSForHTMLTag(`td`, cellCSSProps)

		}
		colSpan := strconv.Itoa(ht.Table.ColCount())
		return `<tr class="` + SECTION2CLASS + `"><td colspan="` + colSpan + `">` + section2 + `</td></tr>`
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

		// col width
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

		// get css string for this cell
		ht.StyleString += `table.` + TABLECLASS + ` tr.` + HEADERSCLASS + ` td`
		ht.StyleString += getCSSForClassSelector(thClass, cellCSSProps)

		tHeaders += `<td class="` + thClass + `">` + headerCell.ColTitle + `</td>`
	}

	return `<tr class="headers">` + tHeaders + `</tr>`, nil
}

func (ht *HTMLTable) getRows() (string, error) {

	// check for empty data table
	blankDataErr := ht.Table.HasData()
	if blankDataErr != nil {
		return "", blankDataErr
	}

	var rowsStr string
	for i := 0; i < ht.Table.RowCount(); i++ {
		s, err := ht.getRow(i)
		if err != nil {
			return "", err
		}
		rowsStr += s
	}

	return rowsStr, nil
}

func (ht *HTMLTable) getRow(rowIndex int) (string, error) {

	// check that this passed rowIndex is valid or not
	inValidRowErr := ht.Table.HasValidRow(rowIndex)
	if inValidRowErr != nil {
		return "", inValidRowErr
	}

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
		g := getCSSMapKeyForCell(rowIndex, colIndex)
		if cssMap, ok := ht.Table.CSS[g]; ok {

			tdClass := `cell-row-` + strconv.Itoa(rowIndex) + `-col-` + strconv.Itoa(colIndex)

			// list of css properties for this td cell
			var cellCSSProps []*CSSProperty
			for _, cssProp := range cssMap {
				cellCSSProps = append(cellCSSProps, cssProp)
			}

			// get css string for this td cell
			ht.StyleString += getCSSForClassSelector(tdClass, cellCSSProps)

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
