package gotable

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/dustin/go-humanize"
)

// CSVTable struct used to prepare table in html version
type CSVTable struct {
	*Table
	CellSep string
	outbuf  bytes.Buffer
}

func (ct *CSVTable) writeTableOutput(w io.Writer) error {
	var tout string

	// append title
	tout += ct.getTitle()

	// append section 1
	tout += ct.getSection1()

	// append section 2
	tout += ct.getSection2()

	// append section 3
	tout += ct.getSection3()

	var tableOut string

	// append headers
	if headerStr, err := ct.getHeaders(); err != nil {
		tableOut += stringln(fmt.Sprintf("%q", err.Error()))
	} else {

		// append rows
		if rowsStr, err := ct.getRows(); err != nil {
			tableOut += stringln(fmt.Sprintf("%q", err.Error()))
		} else {
			tableOut += headerStr
			tableOut += rowsStr
		}
	}

	// // render errorlist
	// tout += ct.getErrorSection()

	if tableOut != "" {
		tout += tableOut
	}

	// return output
	if _, err := ct.outbuf.WriteString(tout); err != nil {
		return err
	}
	// write output to passed io.Writer interface object
	_, err := w.Write(ct.outbuf.Bytes())
	return err
}

func (ct *CSVTable) getTitle() string {
	return stringln(fmt.Sprintf("%q", ct.Table.GetTitle()))
}

func (ct *CSVTable) getSection1() string {
	return stringln(fmt.Sprintf("%q", ct.Table.GetSection1()))
}

func (ct *CSVTable) getSection2() string {
	return stringln(fmt.Sprintf("%q", ct.Table.GetSection2()))
}

func (ct *CSVTable) getSection3() string {
	return stringln(fmt.Sprintf("%q", ct.Table.GetSection3()))
}

// func (ct *CSVTable) getErrorSection() string {
// 	errSection := ""

// 	errList := ct.Table.GetErrorList()
// 	if len(errList) > 0 {
// 		for _, errStr := range errList {
// 			errSection += stringln(fmt.Sprintf("%q", errStr))
// 		}
// 		errSection = NEWLINE + errSection + NEWLINE
// 	}

// 	// blank return
// 	return errSection
// }

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
	return stringln(strings.Join(tHeader, ct.CellSep)), nil
}

func (ct *CSVTable) getRows() (string, error) {
	// check for empty data table
	blankDataErr := ct.Table.HasData()
	if blankDataErr != nil {
		return "", blankDataErr
	}

	var rowsStr string
	for i := 0; i < ct.Table.RowCount(); i++ {
		// for valid row, we will never get an error
		s, _ := ct.getRow(i)
		rowsStr += s
	}

	return rowsStr, nil
}

func (ct *CSVTable) getRow(row int) (string, error) {

	// This method is only called by internal instance of TextTable
	// in getRows method, so we should avoid following error check
	// unless we make it as export

	// // check that this passed row is valid or not
	// inValidRowErr := ct.Table.HasValidRow(row)
	// if inValidRowErr != nil {
	// 	return "", inValidRowErr
	// }

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
	return stringln(strings.Join(tRow, ct.CellSep)), nil
}
