package gotable

import (
	"sort"
	"strings"
	"testing"
	"time"
)

func TestRowSets(t *testing.T) {
	var tbl Table
	tbl.Init() //sets column spacing and date format to default
	tbl.AddColumn("Col Title Too Big", 5, CELLFLOAT, COLJUSTIFYRIGHT)
	tbl.AddColumn("Sample", 40, CELLSTRING, COLJUSTIFYLEFT)

	// Validate access to OOB data...
	cbad := tbl.Get(50, 50)
	if cbad.Fval != 0 || cbad.Ival != 0 {
		t.Logf("rowset_test: OOB access returned unexpected: %#v\n", cbad)
		t.Fail()
	}

	if ibad := tbl.Geti(50, 50); ibad != 0 {
		t.Logf("rowset_test: OOB access returned unexpected: %#v\n", ibad)
		t.Fail()
	}

	if fbad := tbl.Getf(50, 50); fbad != 0 {
		t.Logf("rowset_test: OOB access returned unexpected: %#v\n", fbad)
		t.Fail()
	}

	if sbad := tbl.Gets(50, 50); sbad != "" {
		t.Logf("rowset_test: OOB access returned unexpected: %#v\n", sbad)
		t.Fail()
	}

	dtexp := time.Date(0, time.January, 0, 0, 0, 0, 0, time.UTC)
	if dtbad := tbl.Getd(50, 50); dtbad != dtexp {
		t.Logf("rowset_test: OOB access returned unexpected: %#v\n", dtbad)
		t.Fail()
	}

	var sa = []string{
		"random string",
		"another random string",
		"ROWSET:  first line",
		"ROWSET:  another line",
		"ROWSET:  yet another line",
		"ROWSET:  and still another line",
		"and yet another random string",
		"even another random string",
	}
	rsid := tbl.CreateRowset()
	for i := 0; i < len(sa); i++ {
		tbl.AddRow()
		tbl.Putf(-1, 0, float64(i))
		tbl.Puts(-1, 1, sa[i])
		if !strings.HasPrefix(sa[i], "ROWSET:") {
			continue
		}
		tbl.AppendToRowset(rsid, tbl.RowCount()-1)
	}
	tbl.AddLineAfter(tbl.RowCount() - 1)

	//---------------------------------------------------------------
	// Make sure the rowset is what we expect...
	//---------------------------------------------------------------
	rs := tbl.GetRowset(rsid)
	rsGood := []int{2, 3, 4, 5}
	cBefore := tbl.Sum(0) // sum of numbers in column 0
	if !compareIntSlices(rs, rsGood) {
		t.Logf("rowset_test: Expected %#v,  found %#v\n", rsGood, rs)
		t.Fail()
	}
	sort.Ints(rs)
	tbl.AddLineBefore(rs[0])
	tbl.AddLineAfter(rs[len(rs)-1])
	tbl.SetTitle("T1")
	s1 := tbl.String()

	//---------------------------------------------------------------
	// Now insert a row in the middle of the rowset, and let's see
	// what happens to the rowset row numbers...
	//---------------------------------------------------------------
	t.Logf("Before insert, rowset = %#v\n", rs)
	tbl.InsertRow(4)
	tbl.Putf(4, 0, float64(3.1))
	tbl.Puts(4, 1, "Inserted this line")
	rs = tbl.GetRowset(rsid)
	t.Logf("After insert, rowset = %#v\n", rs)
	rsExpect := []int{2, 3, 5, 6, 4}
	if !compareIntSlices(rs, rsExpect) {
		t.Logf("rowset_test: Expected %#v,  found %#v\n", rsGood, rs)
		t.Fail()
	}
	tbl.SetTitle("T2")
	s2 := tbl.String()

	//---------------------------------------------------------------
	// Now delete the line and make sure the rowset is properly adjusted...
	//---------------------------------------------------------------
	tbl.DeleteRow(4)
	rs = tbl.GetRowset(rsid)
	t.Logf("After insert, rowset = %#v\n", rs)
	if !compareIntSlices(rs, rsGood) {
		t.Logf("rowset_test: Expected %#v,  found %#v\n", rsGood, rs)
		t.Fail()
	}
	cAfter := tbl.Sum(0) // sum of numbers in column 0
	t.Logf("Sum of column1 before: %10.8f,  after: %10.8f\n", cBefore.Fval, cAfter.Fval)
	if float64(28.0) != cAfter.Fval {
		t.Logf("rowset_test: sum failed Expected %10.8f,  found %10.8f\n", float64(28.0), cAfter.Fval)
		t.Fail()
	}

	tbl.SetTitle("T3")
	s3 := tbl.String()

	//---------------------------------------------------------------
	// Finally, do a string sort, and see the result...
	//---------------------------------------------------------------
	tbl.Sort(0, tbl.RowCount()-1, 1)
	tbl.SetTitle("T SORT")
	s4 := tbl.String()

	bShowFile := true
	if bShowFile {
		saveTableToFile(t, "rowset_test.txt", s1+"\n\n"+s2+"\n\n"+s3+"\n\n"+s4)
	}
}

// returns true if they are equal in values, false otherwise
func compareIntSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
