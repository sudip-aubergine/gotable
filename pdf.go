package gotable

import (
	"fmt"
)

// PDFTable struct used to prepare table in pdf version
type PDFTable struct {
	*Table
}

// getTableOutput return table output in pdf form
func (ht *PDFTable) getTableOutput() (string, error) {
	return "", fmt.Errorf("%s", "PDF output for table is not supported yet")
}

func (ht *PDFTable) getHeaders() (string, error) {
	panic("Implementation Error")
}

func (ht *PDFTable) getRows() (string, error) {
	panic("Implementation Error")
}

func (ht *PDFTable) getRow(row int) (string, error) {
	panic("Implementation Error")
}
