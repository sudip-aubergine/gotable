package gotable

import (
	"fmt"
	"strings"

	"github.com/dustin/go-humanize"
)

// CSVTable struct used to prepare table in html version
type CSVTable struct {
	*Table
	CellSep string
}

// getTableOutput return the table header in csv layout
func (ct *CSVTable) getTableOutput() (string, error) {
	// get headers first
	et, err := ct.getHeaders()
	if err != nil {
		return "", err
	}

	// then append table body
	rs, err := ct.getRows()
	if err != nil {
		return "", err
	}
	et += rs

	// fmt.Println(strings.Replace(s, "\\\"", "'", -1))
	return et, nil
}

func (ct *CSVTable) getHeaders() (string, error) {
	// check for blank headers
	blankHdrsErr := ct.Table.HasHeaders()
	if blankHdrsErr != nil {
		return "", blankHdrsErr
	}

	// format headers
	var tHeader []string

	for i := 0; i < len(ct.Table.ColDefs); i++ {
		tHeader = append(tHeader, fmt.Sprintf("%q", ct.Table.ColDefs[i].ColTitle))
	}

	// append last newLine char
	return strings.Join(tHeader, ct.CellSep) + "\n", nil
}

func (ct *CSVTable) getRows() (string, error) {
	// check for empty data table
	blankDataErr := ct.Table.HasData()
	if blankDataErr != nil {
		return "", blankDataErr
	}

	var rowsStr string
	for i := 0; i < ct.Table.Rows(); i++ {
		s, err := ct.getRow(i)
		if err != nil {
			return "", err
		}
		rowsStr += s
	}

	return rowsStr, nil
}

func (ct *CSVTable) getRow(row int) (string, error) {
	// check that this passed row is valid or not
	inValidRowErr := ct.Table.HasValidRow(row)
	if inValidRowErr != nil {
		return "", inValidRowErr
	}

	// format table row
	var tRow []string

	for i := 0; i < len(ct.Table.Row[row].Col); i++ {

		switch ct.Table.Row[row].Col[i].Type {
		case CELLFLOAT:
			tRow = append(tRow, fmt.Sprintf(ct.Table.ColDefs[i].Pfmt, humanize.FormatFloat("#,###.##", ct.Table.Row[row].Col[i].Fval)))
		case CELLINT:
			tRow = append(tRow, fmt.Sprintf(ct.Table.ColDefs[i].Pfmt, ct.Table.Row[row].Col[i].Ival))
		case CELLSTRING:
			// FOR CSV, APPEND FULL STRING, THERE ARE NO MULTILINE STRING IN THIS
			tRow = append(tRow, fmt.Sprintf("%q", ct.Table.Row[row].Col[i].Sval))
		case CELLDATE:
			tRow = append(tRow, fmt.Sprintf("%*.*s", ct.Table.ColDefs[i].Width, ct.Table.ColDefs[i].Width, ct.Table.Row[row].Col[i].Dval.Format(ct.Table.DateFmt)))
		case CELLDATETIME:
			tRow = append(tRow, fmt.Sprintf("%*.*s", ct.Table.ColDefs[i].Width, ct.Table.ColDefs[i].Width, ct.Table.Row[row].Col[i].Dval.Format(ct.Table.DateTimeFmt)))
		default:
			tRow = append(tRow, mkstr(ct.Table.ColDefs[i].Width, ' '))
		}
	}

	// append newline char at last
	return strings.Join(tRow, ct.CellSep) + "\n", nil
}
