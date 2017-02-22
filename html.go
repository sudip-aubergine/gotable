package gotable

import (
	"fmt"
	"strconv"

	"github.com/dustin/go-humanize"
)

// SprintTableHTML prints the whole table in HTML form
func (t *Table) SprintTableHTML(f int) (string, error) {
	s, err := t.SprintColHdrsHTML() // get headers first
	if err != nil {
		return "", err
	}
	rs, err := t.SprintRows(f) // then append table body
	if err != nil {
		return "", err
	}
	s += rs
	return "<table>" + s + "</table>", nil // finally return HTML table layout
}

// SprintColHdrsHTML formats the requested row in HTML and returns the HTML as a string
func (t *Table) SprintColHdrsHTML() (string, error) {
	tHeader := ""
	for i := 0; i < len(t.ColDefs); i++ {
		cd := t.ColDefs[i]
		headerCell := t.ColDefs[i].ColTitle
		if cd.HTMLWidth != -1 {
			headerCell = "<th width=\"" + strconv.Itoa(cd.HTMLWidth) + "\">" + headerCell + "</th>"
		} else {
			headerCell = "<th>" + headerCell + "</th>"
		}
		tHeader += headerCell
	}
	return "<thead><tr>" + tHeader + "</tr></thead>", nil
}

// SprintRowsHTML returns all rows text string
func (t *Table) SprintRowsHTML(f int) (string, error) {
	rowsStr := ""
	for i := 0; i < t.RowCount(); i++ {
		s, err := t.SprintRow(i, f)
		if err != nil {
			return "", err
		}
		rowsStr += s
	}
	return "<tbody>" + rowsStr + "</tbody>", nil
}

// SprintRowHTML formats the requested row in HTML and returns the HTML as a string
// REF: http://stackoverflow.com/questions/21033440/disable-automatic-change-of-width-in-table-tag
func (t *Table) SprintRowHTML(row int) (string, error) {

	tRow := ""

	// fill the content in rowTextList for the first line
	for i := 0; i < len(t.Row[row].Col); i++ {

		var rowCell string
		// append content in TD
		switch t.Row[row].Col[i].Type {
		case CELLFLOAT:
			rowCell = fmt.Sprintf(t.ColDefs[i].Pfmt, humanize.FormatFloat("#,###.##", t.Row[row].Col[i].Fval))
		case CELLINT:
			rowCell = fmt.Sprintf(t.ColDefs[i].Pfmt, t.Row[row].Col[i].Ival)
		case CELLSTRING:
			// ******************************************************
			// FOR HTML, APPEND FULL STRING, THERE ARE NO
			// MULTILINE TEXT IN THIS
			// ******************************************************
			rowCell = fmt.Sprintf("%s", t.Row[row].Col[i].Sval)
		case CELLDATE:
			rowCell = fmt.Sprintf("%*.*s", t.ColDefs[i].Width, t.ColDefs[i].Width, t.Row[row].Col[i].Dval.Format(t.DateFmt))
		case CELLDATETIME:
			rowCell = fmt.Sprintf("%*.*s", t.ColDefs[i].Width, t.ColDefs[i].Width, t.Row[row].Col[i].Dval.Format(t.DateTimeFmt))
		default:
			rowCell = mkstr(t.ColDefs[i].Width, ' ')
		}

		// format td cell
		rowCell = "<td>" + rowCell + "</td>"
		tRow += rowCell
	}

	return "<tr>" + tRow + "</tr>", nil
}
