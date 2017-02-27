package gotable

import (
	"fmt"
	// "strconv"

	"github.com/dustin/go-humanize"
)

// HTMLTable struct used to prepare table in html version
type HTMLTable struct {
	*Table
}

// getTableOutput return table outputin html form
func (ht *HTMLTable) getTableOutput() (string, error) {
	// get headers first
	et, err := ht.getHeaders()
	if err != nil {
		return "", err
	}

	// then append table body
	rs, err := ht.getRows()
	if err != nil {
		return "", err
	}
	et += rs

	// finally return HTML table layout
	return "<table class='rpt-table'>" + et + "</table>", nil
}

func (ht *HTMLTable) getHeaders() (string, error) {

	// check for blank headers
	blankHdrsErr := ht.Table.HasHeaders()
	if blankHdrsErr != nil {
		return "", blankHdrsErr
	}

	// format headers
	var tHeader string

	for i := 0; i < len(ht.Table.ColDefs); i++ {
		// cd := ht.Table.ColDefs[i]
		headerCell := ht.Table.ColDefs[i].ColTitle

		// TODO: handle case when custom htmlwidth passed for a cell
		// if cd.HTMLWidth != -1 {
		// 	headerCell = "<th width=\"" + strconv.Itoa(cd.HTMLWidth) + "\">" + headerCell + "</th>"
		// } else {
		// 	headerCell = "<th>" + headerCell + "</th>"
		// }

		headerCell = "<th>" + headerCell + "</th>"
		tHeader += headerCell
	}

	return "<thead><tr>" + tHeader + "</tr></thead>", nil
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

	return "<tbody>" + rowsStr + "</tbody>", nil
}

func (ht *HTMLTable) getRow(row int) (string, error) {

	// check that this passed row is valid or not
	inValidRowErr := ht.Table.HasValidRow(row)
	if inValidRowErr != nil {
		return "", inValidRowErr
	}

	// format table row
	var tRow string

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
		rowCell = "<td>" + rowCell + "</td>"
		tRow += rowCell
	}

	return "<tr>" + tRow + "</tr>", nil
}
