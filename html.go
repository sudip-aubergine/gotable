package gotable

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/dustin/go-humanize"
)

const (
	CSSCLASSSELECTOR = `.`
	CSSBLOCKSTARTS   = `{`
	CSSBLOCKENDS     = `}`
	CSSPROPSEP       = `:`
	CSSPROPENDS      = `;`
	TABLECLASS       = `rpt-table`
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
	colSpan := strconv.Itoa(ht.Table.ColCount())
	if title != "" {
		title = `<tr class="title"><td colspan="` + colSpan + `">` + title + `</td></tr>`
	}
	return title
}

func (ht *HTMLTable) getSection1() string {
	section1 := ht.Table.GetSection1()
	colSpan := strconv.Itoa(ht.Table.ColCount())
	if section1 != "" {
		return `<tr class="section1"><td colspan="` + colSpan + `">` + section1 + `</td></tr>`
	}
	return section1
}

func (ht *HTMLTable) getSection2() string {
	section2 := ht.Table.GetSection2()
	colSpan := strconv.Itoa(ht.Table.ColCount())
	if section2 != "" {
		return `<tr class="section2"><td colspan="` + colSpan + `">` + section2 + `</td></tr>`
	}
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
		// cd := ht.Table.ColDefs[headerIndex]
		headerCell := ht.Table.ColDefs[headerIndex].ColTitle

		// TODO: handle case when custom htmlwidth passed for a cell
		// if cd.HTMLWidth != -1 {
		// 	headerCell = `<td width="` + strconv.Itoa(cd.HTMLWidth) + `">` + headerCell + `</td>`
		// } else {
		// 	headerCell = `<td>` + headerCell + `</td>`
		// }

		colWidth := strconv.Itoa(ht.Table.ColDefs[headerIndex].Width * CSSFONTSIZE)
		thClass := `header-` + strconv.Itoa(headerIndex)
		ht.StyleString += CSSCLASSSELECTOR + thClass + CSSBLOCKSTARTS +
			`width` + CSSPROPSEP + colWidth + CSSFONTUNIT + CSSPROPENDS +
			CSSBLOCKENDS

		headerCell = `<td class="` + thClass + `">` + headerCell + `</td>`
		tHeaders += headerCell
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
		g := CSSIndex{Row: rowIndex, Col: colIndex}
		if cssMap, ok := ht.Table.CSS[g]; ok {
			var cssString string
			tdClass := `cell-row-` + strconv.Itoa(rowIndex) + `-col-` + strconv.Itoa(colIndex)
			for _, cssProp := range cssMap {
				cssString += cssProp.Name + CSSPROPSEP + cssProp.Value + CSSPROPENDS
			}
			ht.StyleString += CSSCLASSSELECTOR + tdClass + CSSBLOCKSTARTS + cssString + CSSBLOCKENDS
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

// getCSSMapKeyForCell returns key for cell which has css properties
// which resides at rowIndex, colIndex
func (ht *HTMLTable) getCSSMapKeyForCell(rowIndex, colIndex int) string {
	return `row:` + strconv.Itoa(rowIndex) + `-col:` + strconv.Itoa(colIndex)
}
