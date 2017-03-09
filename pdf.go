package gotable

import (
	"fmt"
)

// PDFTable struct used to prepare table in pdf version
type PDFTable struct {
	*Table
}

func (pt *PDFTable) getTableOutput() (string, error) {
	return "", fmt.Errorf("PDF output for table is not supported yet")
}

func (pt *PDFTable) getTitle() string {
	panic("Implementation Error")
}

func (pt *PDFTable) getSection1() string {
	panic("Implementation Error")
}

func (pt *PDFTable) getSection2() string {
	panic("Implementation Error")
}

func (pt *PDFTable) getHeaders() (string, error) {
	panic("Implementation Error")
}

func (pt *PDFTable) getRows() (string, error) {
	panic("Implementation Error")
}

func (pt *PDFTable) getRow(row int) (string, error) {
	panic("Implementation Error")
}
