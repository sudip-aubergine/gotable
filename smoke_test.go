package gotable

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"
)

func TestSmoke(t *testing.T) {
	var tbl Table
	title := "GOTABLE"
	section1 := "A Smoke Test"
	section2 := "February 21, 2017"
	tbl.Init() //sets column spacing and date format to default

	// force some edge condition errors...
	errExp := "no columns"
	// headers check
	err := tbl.HasHeaders()
	if !strings.Contains(err.Error(), errExp) {
		t.Errorf("smoke_test: Expected %q, but found: %s\n", errExp, err.Error())
	}
	// print table check
	_, err = tbl.SprintTable(TABLEOUTTEXT)
	if !strings.Contains(err.Error(), errExp) {
		t.Errorf("smoke_test: Expected %q, but found: %s\n", errExp, err.Error())
	}

	tbl.SetTitle(title)
	tbl.SetSection1(section1)
	tbl.SetSection2(section2)
	tbl.AddColumn("Name", 35, CELLSTRING, COLJUSTIFYLEFT)               // 0 Name
	tbl.AddColumn("Person Age", 3, CELLINT, COLJUSTIFYRIGHT)            // 1 Age
	tbl.AddColumn("Height (cm)", 8, CELLINT, COLJUSTIFYRIGHT)           // 2 Height in centimeters
	tbl.AddColumn("Date of Birth", 10, CELLDATE, COLJUSTIFYLEFT)        // 3 DOB
	tbl.AddColumn("Country of Birth", 14, CELLSTRING, COLJUSTIFYLEFT)   // 4 COB
	tbl.AddColumn("Winnings", 12, CELLFLOAT, COLJUSTIFYRIGHT)           // 5 total winnings
	tbl.AddColumn("Notes", 20, CELLSTRING, COLJUSTIFYLEFT)              // 6 Notes
	tbl.AddColumn("Random Date/Time", 25, CELLDATETIME, COLJUSTIFYLEFT) // 7 totally random datetime

	errExp = "no rows"
	// headers check
	err = tbl.HasData()
	if !strings.Contains(err.Error(), errExp) {
		t.Errorf("smoke_test: Expected %q, but found: %s\n", errExp, err.Error())
	}
	// print table check
	_, err = tbl.SprintTable(TABLEOUTTEXT)
	if !strings.Contains(err.Error(), errExp) {
		t.Errorf("smoke_test: Expected %q, but found: %s\n", errExp, err.Error())
	}

	const (
		Name = iota
		Age
		Height
		DOB
		COB
		Winnings
		Notest
	)

	bd1 := time.Date(1969, time.March, 2, 0, 0, 0, 0, time.UTC)
	bd2 := time.Date(1960, time.October, 4, 0, 0, 0, 0, time.UTC)
	bd3 := time.Date(1974, time.April, 10, 0, 0, 0, 0, time.UTC)
	bd4 := time.Date(1950, time.April, 21, 0, 0, 0, 0, time.UTC)
	bd5 := time.Date(1977, time.August, 6, 0, 0, 0, 0, time.UTC)

	type tdata struct {
		Name     string
		Age      int64
		Height   int64
		DOB      time.Time
		COB      string
		Winnings float64
		Notes    string
		UnxNano  int64
	}
	var d = []tdata{
		{Name: "Mary M. Oneil", Age: 47, Height: 165, DOB: bd1, COB: "United States", Winnings: float64(17633.21), Notes: "A few notes here", UnxNano: 7564006824999651664},
		{Name: "Lynette C. Allen", Age: 56, Height: 156, DOB: bd2, COB: "United States", Winnings: float64(45373.00), Notes: "A lot more notes. A whole, big, line with lots and lots and lots and lots of notes. And some more notes.", UnxNano: 7733402883116878723},
		{Name: "Stanislaus Aliyeva", Age: 42, Height: 172, DOB: bd3, COB: "Slovinia", Winnings: 106632.36, Notes: "A few notes here", UnxNano: 1584693382958379231},
		{Name: "Casandra Åberg", Age: 66, Height: 158, DOB: bd4, COB: "Sweden", Winnings: 93883.25, Notes: "2000 Seat Toledo", UnxNano: 7796987096200859545},
		{Name: "Amanda Melo Ferreira", Age: 55, Height: 174, DOB: bd5, COB: "Brazil", Winnings: 46673.42, Notes: "2006 Ford Falcon", UnxNano: 3267110399458248377},
	}

	totalsRSet := tbl.CreateRowset()
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < len(d); i++ {
		tbl.AddRow()
		tbl.AppendToRowset(totalsRSet, tbl.RowCount()-1)
		tbl.Puts(-1, Name, d[i].Name)
		tbl.Puti(-1, Age, d[i].Age)
		tbl.Puti(-1, Height, d[i].Height)
		tbl.Putd(-1, DOB, d[i].DOB)
		tbl.Puts(-1, COB, d[i].COB)
		tbl.Putf(-1, Winnings, d[i].Winnings)
		tbl.Puts(-1, Notest, d[i].Notes)
		tbl.Putdt(-1, 7, time.Unix(0, d[i].UnxNano).UTC()) // random date in the future
	}
	// Start with simple checks...
	if tbl.ColCount() != 8 {
		t.Errorf("smoke_test: Expected %d,  found %d\n", 8, tbl.ColCount())
	}
	if len(tbl.GetRowset(99)) != 0 {
		t.Errorf("smoke_test: Expected emty rowset,  found %#v\n", tbl.GetRowset(99))
	}
	if 0 != tbl.Type(999, 999) {
		t.Errorf("smoke_test: Expected %d,  found %d\n", 0, tbl.Type(999, 999))
	}
	if tbl.GetTitle() != title {
		t.Errorf("smoke_test: Expected %s,  found %s\n", tbl.GetTitle(), title)
	}
	if false != tbl.Puti(999, 999, 1) {
		t.Errorf("smoke_test: Expected return value of false, but got true\n")
	}
	if false != tbl.Putf(999, 999, 1) {
		t.Errorf("smoke_test: Expected return value of false, but got true\n")
	}
	if false != tbl.Puts(999, 999, "ignore") {
		t.Errorf("smoke_test: Expected return value of false, but got true\n")
	}
	if false != tbl.Putd(999, 999, time.Now()) {
		t.Errorf("smoke_test: Expected return value of false, but got true\n")
	}
	if tbl.GetSection1() != section1 {
		t.Errorf("smoke_test: Expected %s,  found %s\n", tbl.GetSection1(), section1)
	}
	if tbl.GetSection2() != section2 {
		t.Errorf("smoke_test: Expected %s,  found %s\n", tbl.GetSection2(), section2)
	}
	cell := tbl.Get(0, 0)
	if cell.Sval != d[0].Name {
		t.Errorf("smoke_test: Expected %s,  found %s\n", cell.Sval, d[0].Name)
	}
	if tbl.Geti(1, Age) != d[1].Age {
		t.Errorf("smoke_test: Expected %d,  found %d\n", tbl.Geti(1, Age), d[1].Age)

	}
	if tbl.Getf(1, Winnings) != d[1].Winnings {
		t.Errorf("smoke_test: Expected %f,  found %f\n", tbl.Getf(1, Winnings), d[1].Winnings)
	}
	if tbl.Gets(1, Name) != d[1].Name {
		t.Errorf("smoke_test: Expected %s,  found %s\n", tbl.Gets(1, Name), d[1].Name)
	}
	if tbl.Getd(1, DOB) != d[1].DOB {
		t.Logf("smoke_test: Expected %s,  found %s\n", tbl.Getd(1, DOB).Format("1/2/2006"), d[1].DOB.Format("1/Errorf06"))
	}
	if tbl.Type(1, Name) != CELLSTRING {
		t.Errorf("smoke_test: Expected %d,  found %d\n", tbl.Type(1, Name), CELLSTRING)
	}

	errExp = "Unrecognized"
	_, err = tbl.SprintTable(999)
	if !strings.Contains(err.Error(), errExp) {
		t.Errorf("smoke_test: Expected %q in error, but not found.  Error = %s\n", errExp, err.Error())
	}

	// Bang it a bit...
	tbl.Sort(0, tbl.RowCount()-1, DOB)
	tbl.AddLineAfter(tbl.RowCount() - 1) // a line after the last row in the table
	tbl.InsertSumRowsetCols(totalsRSet, tbl.RowCount(), []int{Winnings})

	// valid row checks
	errExp = "Row number > no of rows in table"
	err = tbl.HasValidRow(999)
	if !strings.Contains(err.Error(), errExp) {
		t.Errorf("smoke_test: Expected %q, but found: %s\n", errExp, err.Error())
	}

	errExp = "Row number is less than zero"
	err = tbl.HasValidRow(-999)
	if !strings.Contains(err.Error(), errExp) {
		t.Errorf("smoke_test: Expected %q, but found: %s\n", errExp, err.Error())
	}

	// valid column checks
	errExp = "Column number > no of columns in table"
	err = tbl.HasValidColumn(999)
	if !strings.Contains(err.Error(), errExp) {
		t.Errorf("smoke_test: Expected %q, but found: %s\n", errExp, err.Error())
	}

	errExp = "Column number is less than zero"
	err = tbl.HasValidColumn(-999)
	if !strings.Contains(err.Error(), errExp) {
		t.Errorf("smoke_test: Expected %q, but found: %s\n", errExp, err.Error())
	}

	// Now hit it hard...
	DoTextOutput(t, &tbl)
	DoCSVOutput(t, &tbl)
	// DoHTMLOutput(t, &tbl)
	DoPDFOutput(t, &tbl)
}

func DoTextOutput(t *testing.T, tbl *Table) {
	(*tbl).TightenColumns()
	s := fmt.Sprintf("%s\n", (*tbl))
	saveTableToFile(t, "smoke_test.txt", s)

	// now compare what we have to the known-good output
	b, _ := ioutil.ReadFile("./testdata/smoke_test.txt")
	sb := []byte(s)
	if len(b) != len(sb) {
		// fmt.Printf("smoke_test: Expected len = %d,  found len = %d\n", len(b), len(sb))
		t.Errorf("smoke_test: Expected len = %d,  found len = %d\n", len(b), len(sb))
	}
	for i := 0; i < len(b); i++ {
		if sb[i] != b[i] {
			t.Errorf("smoke_test: micompare at character %d, expected %x (%c), found %x (%c)\n", i, b[i], b[i], sb[i], sb[i])
			// fmt.Printf("smoke_test: micompare at character %d, expected %x (%c), found %x (%c)\n", i, b[i], b[i], sb[i], sb[i])
			break
		}
	}
}

func DoCSVOutput(t *testing.T, tbl *Table) {
	s, err := (*tbl).SprintTable(TABLEOUTCSV)
	if nil != err {
		t.Errorf("smoke_test: Error creating CSV output: %s\n", err.Error())
		// fmt.Printf("smoke_test: Error creating CSV output: %s\n", err.Error())
	}
	saveTableToFile(t, "smoke_test.csv", s)

	// now compare what we have to the known-good output
	b, _ := ioutil.ReadFile("./testdata/smoke_test.csv")
	sb := []byte(s)
	if len(b) != len(sb) {
		// fmt.Printf("smoke_test: Expected len = %d,  found len = %d\n", len(b), len(sb))
		t.Errorf("smoke_test: Expected len = %d,  found len = %d\n", len(b), len(sb))
	}
	for i := 0; i < len(b); i++ {
		if sb[i] != b[i] {
			t.Logf("smoke_test: micompare at character %d, expected %x (%c), found %x (%c)\n", i, b[i], b[i], sb[i], sb[i])
			// fmt.Printf("smoke_test: micompare at character %d, expected %x (%c), found %x (%c)\n", i, b[i], b[i], sb[i], sb[i])
			break
		}
	}

}

func DoHTMLOutput(t *testing.T, tbl *Table) {
	s, err := (*tbl).SprintTable(TABLEOUTHTML)
	if nil != err {
		t.Errorf("smoke_test: Error creating HTML output: %s\n", err.Error())
		// fmt.Printf("smoke_test: Error creating HTML output: %s\n", err.Error())
	}
	saveTableToFile(t, "smoke_test.html", s)

	// now compare what we have to the known-good output
	b, _ := ioutil.ReadFile("./testdata/smoke_test.html")
	sb := []byte(s)
	if len(b) != len(sb) {
		// fmt.Printf("smoke_test: Expected len = %d,  found len = %d\n", len(b), len(sb))
		t.Errorf("smoke_test: Expected len = %d,  found len = %d\n", len(b), len(sb))
	}
	for i := 0; i < len(b); i++ {
		if sb[i] != b[i] {
			t.Logf("smoke_test: micompare at character %d, expected %x (%c), found %x (%c)\n", i, b[i], b[i], sb[i], sb[i])
			// fmt.Printf("smoke_test: micompare at character %d, expected %x (%c), found %x (%c)\n", i, b[i], b[i], sb[i], sb[i])
			break
		}
	}
}

func DoPDFOutput(t *testing.T, tbl *Table) {
	s, err := (*tbl).SprintTable(TABLEOUTPDF)
	if nil != err {
		t.Logf("smoke_test: Error creating PDF output: %s\n", err.Error())
		// fmt.Printf("smoke_test: Error creating PDF output: %s\n", err.Error())
	}
	if len(s) > 0 {
		t.Errorf("smoke_test: Expected: `PDF output for table is not supported yet`,  found: `%s`\n", s)
	}
}

func saveTableToFile(t *testing.T, fname string, s string) error {
	// save for later inspection if anything goes wrong
	f, err := os.Create(fname)
	if nil != err {
		t.Errorf("smoke_test: Error creating file %s: %s\n", fname, err.Error())
		// fmt.Printf("smoke_test: Error creating file: %s\n", err.Error())
	}
	defer f.Close()
	fmt.Fprintf(f, "%s", s)
	return err
}
