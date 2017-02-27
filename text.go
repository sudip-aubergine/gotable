package gotable

import (
	"fmt"
	"sort"

	"github.com/dustin/go-humanize"
)

// TextTable struct used to prepare table in text version
type TextTable struct {
	*Table
	TextColSpace int
}

// prepareForTextPrint
// For text output we want at least one "\n" at the end of a section or title.
// If the supplied string does not end in "\n", then one will be appended to it
// on the return value.
func prepareForTextPrint(s string) string {
	if len(s) == 0 {
		return ""
	}
	if s[len(s)-1] != '\n' {
		return s + "\n"
	}
	return s
}

// getTableOutput prints the whole table in text form
func (tt *TextTable) getTableOutput() (string, error) {
	// get headers first
	et, err := tt.getHeaders()
	if err != nil {
		return "", err
	}

	// then append table body
	rs, err := tt.getRows()
	if err != nil {
		return "", err
	}
	et += rs

	return et, nil
}

// SprintColHdrsText formats the column headers as text and returns the string
func (tt *TextTable) getHeaders() (string, error) {

	// check for blank headers
	blankHdrsErr := tt.Table.HasHeaders()
	if blankHdrsErr != nil {
		return "", blankHdrsErr
	}

	tt.Table.AdjustAllColumnHeaders()

	s := ""
	for j := 0; j < len(tt.Table.ColDefs[0].Hdr); j++ {
		for i := 0; i < len(tt.Table.ColDefs); i++ {
			sf := ""
			lft := ""
			if tt.Table.ColDefs[i].Justify == COLJUSTIFYLEFT {
				lft += "-"
			}
			sf += fmt.Sprintf("%%%s%ds", lft, tt.Table.ColDefs[i].Width)
			s += fmt.Sprintf(sf, tt.Table.ColDefs[i].Hdr[j])
			if i < len(tt.Table.ColDefs)-1 {
				s += mkstr(tt.TextColSpace, ' ')
			}
		}
		s += "\n"
	}

	for i := 0; i < len(tt.Table.ColDefs); i++ {
		s += fmt.Sprintf("%s", mkstr(tt.Table.ColDefs[i].Width, '-'))
		if i < len(tt.Table.ColDefs)-1 {
			s += mkstr(tt.TextColSpace, ' ')
		}
	}
	s += "\n"
	return s, nil
}

func (tt *TextTable) getRows() (string, error) {
	// check for empty data table
	blankDataErr := tt.Table.HasData()
	if blankDataErr != nil {
		return "", blankDataErr
	}

	var rowsStr string
	for i := 0; i < tt.Table.RowCount(); i++ {
		s, err := tt.getRow(i)
		if err != nil {
			return "", err
		}
		rowsStr += s
	}

	return rowsStr, nil
}

func (tt *TextTable) getRow(row int) (string, error) {

	// check that this passed row is valid or not
	inValidRowErr := tt.Table.HasValidRow(row)
	if inValidRowErr != nil {
		return "", inValidRowErr
	}

	// format table row
	var s string

	if len(tt.Table.LineBefore) > 0 {
		j := sort.SearchInts(tt.Table.LineBefore, row)
		if j < len(tt.Table.LineBefore) && row == tt.Table.LineBefore[j] {
			s += tt.sprintLineText()
		}
	}

	rowColumns := len(tt.Table.Row[row].Col)

	// columns string chunk map, each column holds list of string
	// that fits in one line at best
	colStringChunkMap := map[int][]string{}

	// get Height of row that require to fit the content of max cell string content
	// by default table has no all the data in string format, so that we need to add
	// logic here only, to support multi line functionality
	for i := 0; i < rowColumns; i++ {
		if tt.Table.Row[row].Col[i].Type == CELLSTRING {
			cd := tt.Table.ColDefs[i]
			a, _ := getMultiLineText(tt.Table.Row[row].Col[i].Sval, cd.Width)

			colStringChunkMap[i] = a

			if len(a) > tt.Table.Row[row].Height {
				tt.Table.Row[row].Height = len(a)
			}
		}
	}

	// rowTextList holds the 2D array, containing data for each block
	// to achieve multiline row
	rowTextList := [][]string{}

	// init rowTextList with empty values
	for k := 0; k < tt.Table.Row[row].Height; k++ {
		temp := make([]string, rowColumns)
		for i := 0; i < rowColumns; i++ {
			// assign default string with length of column width
			temp = append(temp, mkstr(tt.Table.ColDefs[i].Width, ' '))
		}
		rowTextList = append(rowTextList, temp)
	}

	// fill the content in rowTextList for the first line
	for i := 0; i < rowColumns; i++ {
		switch tt.Table.Row[row].Col[i].Type {
		case CELLFLOAT:
			rowTextList[0][i] = fmt.Sprintf(tt.Table.ColDefs[i].Pfmt, humanize.FormatFloat("#,###.##", tt.Table.Row[row].Col[i].Fval))
		case CELLINT:
			rowTextList[0][i] = fmt.Sprintf(tt.Table.ColDefs[i].Pfmt, tt.Table.Row[row].Col[i].Ival)
		case CELLSTRING:
			rowTextList[0][i] = fmt.Sprintf(tt.Table.ColDefs[i].Pfmt, colStringChunkMap[i][0])
		case CELLDATE:
			rowTextList[0][i] = fmt.Sprintf("%*.*s", tt.Table.ColDefs[i].Width, tt.Table.ColDefs[i].Width, tt.Table.Row[row].Col[i].Dval.Format(tt.Table.DateFmt))
		case CELLDATETIME:
			rowTextList[0][i] = fmt.Sprintf("%*.*s", tt.Table.ColDefs[i].Width, tt.Table.ColDefs[i].Width, tt.Table.Row[row].Col[i].Dval.Format(tt.Table.DateTimeFmt))
		default:
			rowTextList[0][i] = mkstr(tt.Table.ColDefs[i].Width, ' ')
		}
	}

	// rowTextList to string
	for k := 0; k < tt.Table.Row[row].Height; k++ {
		for i := 0; i < rowColumns; i++ {

			// if not first row then process multi line format
			if k > 0 {
				if tt.Table.Row[row].Col[i].Type == CELLSTRING {
					if k >= len(colStringChunkMap[i]) {
						rowTextList[k][i] = fmt.Sprintf(tt.Table.ColDefs[i].Pfmt, "")
					} else {
						rowTextList[k][i] = fmt.Sprintf(tt.Table.ColDefs[i].Pfmt, colStringChunkMap[i][k])
					}
				}
			}

			// if blank then append string of column width with blank
			if rowTextList[k][i] == "" {
				rowTextList[k][i] = mkstr(tt.Table.ColDefs[i].Width, ' ')
			}
			s += rowTextList[k][i]

			// if it is not last block then
			if i < (rowColumns - 1) {
				s += mkstr(tt.TextColSpace, ' ')
			}
		}
		s += "\n"
	}

	if len(tt.Table.LineAfter) > 0 {
		j := sort.SearchInts(tt.Table.LineAfter, row)
		if j < len(tt.Table.LineAfter) && row == tt.Table.LineAfter[j] {
			s += tt.sprintLineText()
		}
	}
	return s, nil
}

// SprintLineText returns a line across all rows in the table
func (tt *TextTable) sprintLineText() string {
	var s string
	for i := 0; i < len(tt.Table.ColDefs); i++ {
		// draw line with hyphen `-` char
		s += mkstr(tt.Table.ColDefs[i].Width, '-')

		// separate text columns
		s += mkstr(tt.TextColSpace, ' ')
	}
	// remove last textColSpace from s
	s = s[0 : len(s)-tt.TextColSpace]
	return s + "\n"
}
