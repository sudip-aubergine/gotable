package gotable

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/dustin/go-humanize"
)

// CSSProperty holds css property to be used as inline css
type CSSProperty struct {
	Name, Value string
}

// HTMLTable struct used to prepare table in html version
type HTMLTable struct {
	*Table
}

func (ht *HTMLTable) getTableOutput() (string, error) {
	var tout string

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
	return `<table class="rpt-table">` + tout + `</table>`, nil
}

func (ht *HTMLTable) getTitle() string {
	title := ht.Table.GetTitle()
	colSpan := strconv.Itoa(ht.Table.ColCount())
	if title != "" {
		title = `<tr class="title"><th colspan="` + colSpan + `">` + title + `</th></tr>`
	}
	return title
}

func (ht *HTMLTable) getSection1() string {
	section1 := ht.Table.GetSection1()
	colSpan := strconv.Itoa(ht.Table.ColCount())
	if section1 != "" {
		return `<tr class="section1"><th colspan="` + colSpan + `">` + section1 + `</th></tr>`
	}
	return section1
}

func (ht *HTMLTable) getSection2() string {
	section2 := ht.Table.GetSection2()
	colSpan := strconv.Itoa(ht.Table.ColCount())
	if section2 != "" {
		return `<tr class="section2"><th colspan="` + colSpan + `">` + section2 + `</th></tr>`
	}
	return section2
}

func (ht *HTMLTable) getHeaders() (string, error) {

	// check for blank headers
	blankHdrsErr := ht.Table.HasHeaders()
	if blankHdrsErr != nil {
		return "", blankHdrsErr
	}

	// htmlHeader includes title, section1, section2, header of table struct
	// this all going to be part of thead tag
	var htmlHeader string

	// append title
	htmlHeader += ht.getTitle()

	// append section 1
	htmlHeader += ht.getSection1()

	// append section 2
	htmlHeader += ht.getSection2()

	// format headers
	var tHeaders string

	for i := 0; i < len(ht.Table.ColDefs); i++ {
		// cd := ht.Table.ColDefs[i]
		headerCell := ht.Table.ColDefs[i].ColTitle

		// TODO: handle case when custom htmlwidth passed for a cell
		// if cd.HTMLWidth != -1 {
		// 	headerCell = `<th width="` + strconv.Itoa(cd.HTMLWidth) + `">` + headerCell + `</th>`
		// } else {
		// 	headerCell = `<th>` + headerCell + `</th>`
		// }

		headerCell = `<th>` + headerCell + `</th>`
		tHeaders += headerCell
	}

	htmlHeader += `<tr class="headers">` + tHeaders + `</tr>`

	// finally return html header with thead
	return `<thead>` + htmlHeader + `</thead>`, nil
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

	return `<tbody>` + rowsStr + `</tbody>`, nil
}

func (ht *HTMLTable) getRow(row int) (string, error) {

	// check that this passed row is valid or not
	inValidRowErr := ht.Table.HasValidRow(row)
	if inValidRowErr != nil {
		return "", inValidRowErr
	}

	// format table row
	var tRow string
	var trClass string

	if len(ht.Table.LineBefore) > 0 {
		j := sort.SearchInts(ht.Table.LineBefore, row)
		if j < len(ht.Table.LineBefore) && row == ht.Table.LineBefore[j] {
			trClass += `top-line`
		}
	}

	// fill the content in rowTextList for the first line
	for i := 0; i < len(ht.Table.Row[row].Col); i++ {

		var rowCell string
		// append content in TD
		switch ht.Table.Row[row].Col[i].Type {
		case CELLFLOAT:
			rowCell = fmt.Sprintf(ht.Table.ColDefs[i].Pfmt, humanize.FormatFloat("#,###.##", ht.Table.Row[row].Col[i].Fval))
		case CELLINT:
			rowCell = fmt.Sprintf(ht.Table.ColDefs[i].Pfmt, ht.Table.Row[row].Col[i].Ival)
		case CELLSTRING:
			// ******************************************************
			// FOR HTML, APPEND FULL STRING, THERE ARE NO
			// MULTILINE TEXT IN THIS
			// ******************************************************
			rowCell = fmt.Sprintf("%s", ht.Table.Row[row].Col[i].Sval)
		case CELLDATE:
			rowCell = fmt.Sprintf("%*.*s", ht.Table.ColDefs[i].Width, ht.Table.ColDefs[i].Width, ht.Table.Row[row].Col[i].Dval.Format(ht.Table.DateFmt))
		case CELLDATETIME:
			rowCell = fmt.Sprintf("%*.*s", ht.Table.ColDefs[i].Width, ht.Table.ColDefs[i].Width, ht.Table.Row[row].Col[i].Dval.Format(ht.Table.DateTimeFmt))
		default:
			rowCell = mkstr(ht.Table.ColDefs[i].Width, ' ')
		}

		// format td cell
		rowCell = `<td>` + rowCell + `</td>`
		tRow += rowCell
	}

	if len(ht.Table.LineAfter) > 0 {
		j := sort.SearchInts(ht.Table.LineAfter, row)
		if j < len(ht.Table.LineAfter) && row == ht.Table.LineAfter[j] {
			trClass += `bottom-line`
		}
	}

	if trClass != "" {
		return `<tr class="` + trClass + `">` + tRow + `</tr>`, nil
	}
	return `<tr>` + tRow + `</tr>`, nil

}
