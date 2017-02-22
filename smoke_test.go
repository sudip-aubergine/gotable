package gotable

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestSmoke(t *testing.T) {
	var tbl Table
	title := "GOTABLE\n"
	section1 := "A Smoke Test\n"
	section2 := "February 21, 2017\n"
	tbl.Init() //sets column spacing and date format to default
	tbl.SetTitle(title)
	tbl.SetSection1(section1)
	tbl.SetSection2(section2)
	tbl.AddColumn("Name", 35, CELLSTRING, COLJUSTIFYLEFT)             // 0 Name
	tbl.AddColumn("Age", 3, CELLINT, COLJUSTIFYRIGHT)                 // 1 Age
	tbl.AddColumn("Height (cm)", 8, CELLINT, COLJUSTIFYRIGHT)         // 2 Height in centimeters
	tbl.AddColumn("Date of Birth", 10, CELLDATE, COLJUSTIFYLEFT)      // 3 DOB
	tbl.AddColumn("Country of Birth", 14, CELLSTRING, COLJUSTIFYLEFT) // 4 COB
	tbl.AddColumn("Winnings", 12, CELLFLOAT, COLJUSTIFYRIGHT)         // 5 total winnings
	tbl.AddColumn("Notes", 40, CELLSTRING, COLJUSTIFYLEFT)            // 6 Notes

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
	}
	var d = []tdata{
		{Name: "Mary M. Oneil", Age: 47, Height: 165, DOB: bd1, COB: "United States", Winnings: float64(17633.21), Notes: "A few notes here"},
		{Name: "Lynette C. Allen", Age: 56, Height: 156, DOB: bd2, COB: "United States", Winnings: float64(45373.00), Notes: "A lot more notes. A whole, big, line with lots and lots and lots and lots of notes. And some more notes."},
		{Name: "Stanislaus Aliyeva", Age: 42, Height: 172, DOB: bd3, COB: "Slovinia", Winnings: 106632.36, Notes: "A few notes here"},
		{Name: "Casandra Åberg", Age: 66, Height: 158, DOB: bd4, COB: "Sweden", Winnings: 93883.25, Notes: "2000 Seat Toledo"},
		{Name: "Amanda Melo Ferreira", Age: 55, Height: 174, DOB: bd5, COB: "Brazil", Winnings: 46673.42, Notes: "2006 Ford Falcon"},
	}

	totalsRSet := tbl.CreateRowset()

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
	}
	// Start with simple checks...
	if tbl.GetTitle() != title {
		t.Logf("smoke_test: Expected %s,  found %s\n", tbl.GetTitle(), title)
		t.Fail()
	}
	if tbl.GetSection1() != section1 {
		t.Logf("smoke_test: Expected %s,  found %s\n", tbl.GetSection1(), section1)
		t.Fail()
	}
	if tbl.GetSection2() != section2 {
		t.Logf("smoke_test: Expected %s,  found %s\n", tbl.GetSection2(), section2)
		t.Fail()
	}

	cell := tbl.Get(0, 0)
	if cell.Sval != d[0].Name {
		t.Logf("smoke_test: Expected %s,  found %s\n", cell.Sval, d[0].Name)
		t.Fail()
	}
	if tbl.Geti(1, Age) != d[1].Age {
		t.Logf("smoke_test: Expected %d,  found %d\n", tbl.Geti(1, Age), d[1].Age)
		t.Fail()
	}
	if tbl.Getf(1, Winnings) != d[1].Winnings {
		t.Logf("smoke_test: Expected %d,  found %d\n", tbl.Getf(1, Winnings), d[1].Winnings)
		t.Fail()
	}
	if tbl.Gets(1, Name) != d[1].Name {
		t.Logf("smoke_test: Expected %d,  found %d\n", tbl.Gets(1, Name), d[1].Name)
		t.Fail()
	}
	if tbl.Getd(1, DOB) != d[1].DOB {
		t.Logf("smoke_test: Expected %d,  found %d\n", tbl.Getd(1, DOB), d[1].DOB)
		t.Fail()
	}
	if tbl.Type(1, Name) != CELLSTRING {
		t.Logf("smoke_test: Expected %d,  found %d\n", tbl.Type(1, Name), CELLSTRING)
		t.Fail()
	}

	// Bang it a bit...
	tbl.Sort(0, tbl.RowCount()-1, DOB)
	tbl.AddLineAfter(tbl.RowCount() - 1) // a line after the last row in the table
	tbl.InsertSumRowsetCols(totalsRSet, tbl.RowCount(), []int{Winnings})

	// Now hit it hard...
	DoTextOutput(t, &tbl)
	DoCSVOutput(t, &tbl)
	DoHTMLOutput(t, &tbl)
}

func DoTextOutput(t *testing.T, tbl *Table) {
	s := fmt.Sprintf("%s\n", (*tbl))
	(*tbl).TightenColumns()
	s += fmt.Sprintf("%s\n", (*tbl))

	// save output for later inspection if anything goes wrong
	f, err := os.Create("smoke_test.out")
	if nil != err {
		t.Logf("smoke_test: Error creating file: %s\n", err.Error())
		t.Fail()
		// fmt.Printf("smoke_test: Error creating file: %s\n", err.Error())
	}
	defer f.Close()
	fmt.Fprintf(f, "%s", s)
	f.Sync()

	// now compare what we have to the known-good output
	b, _ := ioutil.ReadFile("./testdata/smoke_test.txt")
	sb := []byte(s)
	if len(b) != len(sb) {
		// fmt.Printf("smoke_test: Expected len = %d,  found len = %d\n", len(b), len(sb))
		t.Logf("smoke_test: Expected len = %d,  found len = %d\n", len(b), len(sb))
		t.Fail()
	}
	for i := 0; i < len(b); i++ {
		if sb[i] != b[i] {
			t.Logf("smoke_test: micompare at character %d, expected %x (%c), found %x (%c)\n", i, b[i], b[i], sb[i], sb[i])
			t.Fail()
			// fmt.Printf("smoke_test: micompare at character %d, expected %x (%c), found %x (%c)\n", i, b[i], b[i], sb[i], sb[i])
			break
		}
	}
}

func DoCSVOutput(t *testing.T, tbl *Table) {
	s, err := (*tbl).SprintTable(TABLEOUTCSV)
	if nil != err {
		t.Logf("smoke_test: Error creating CSV output: %s\n", err.Error())
		t.Fail()
		// fmt.Printf("smoke_test: Error creating CSV output: %s\n", err.Error())
	}

	// save for later inspection if anything goes wrong
	f, err := os.Create("smoke_test.csv")
	if nil != err {
		t.Logf("smoke_test: Error creating file: %s\n", err.Error())
		t.Fail()
		// fmt.Printf("smoke_test: Error creating file: %s\n", err.Error())
	}
	defer f.Close()
	fmt.Fprintf(f, "%s", s)
	f.Sync()

	// now compare what we have to the known-good output
	b, _ := ioutil.ReadFile("./testdata/smoke_test.csv")
	sb := []byte(s)
	if len(b) != len(sb) {
		// fmt.Printf("smoke_test: Expected len = %d,  found len = %d\n", len(b), len(sb))
		t.Logf("smoke_test: Expected len = %d,  found len = %d\n", len(b), len(sb))
		t.Fail()
	}
	for i := 0; i < len(b); i++ {
		if sb[i] != b[i] {
			t.Logf("smoke_test: micompare at character %d, expected %x (%c), found %x (%c)\n", i, b[i], b[i], sb[i], sb[i])
			t.Fail()
			// fmt.Printf("smoke_test: micompare at character %d, expected %x (%c), found %x (%c)\n", i, b[i], b[i], sb[i], sb[i])
			break
		}
	}

}

func DoHTMLOutput(t *testing.T, tbl *Table) {
	s, err := (*tbl).SprintTable(TABLEOUTHTML)
	if nil != err {
		t.Logf("smoke_test: Error creating HTML output: %s\n", err.Error())
		t.Fail()
		// fmt.Printf("smoke_test: Error creating HTML output: %s\n", err.Error())
	}

	// save for later inspection if anything goes wrong
	f, err := os.Create("smoke_test.html")
	if nil != err {
		t.Logf("smoke_test: Error creating file: %s\n", err.Error())
		t.Fail()
		// fmt.Printf("smoke_test: Error creating file: %s\n", err.Error())
	}
	defer f.Close()
	fmt.Fprintf(f, "%s", s)
	f.Sync()

	// now compare what we have to the known-good output
	b, _ := ioutil.ReadFile("./testdata/smoke_test.html")
	sb := []byte(s)
	if len(b) != len(sb) {
		// fmt.Printf("smoke_test: Expected len = %d,  found len = %d\n", len(b), len(sb))
		t.Logf("smoke_test: Expected len = %d,  found len = %d\n", len(b), len(sb))
		t.Fail()
	}
	for i := 0; i < len(b); i++ {
		if sb[i] != b[i] {
			t.Logf("smoke_test: micompare at character %d, expected %x (%c), found %x (%c)\n", i, b[i], b[i], sb[i], sb[i])
			t.Fail()
			// fmt.Printf("smoke_test: micompare at character %d, expected %x (%c), found %x (%c)\n", i, b[i], b[i], sb[i], sb[i])
			break
		}
	}
}
