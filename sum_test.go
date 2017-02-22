package gotable

import (
	"math"
	"math/rand"
	"testing"
	"time"
)

func TestSum(t *testing.T) {
	var tbl Table
	tbl.Init() //sets column spacing and date format to default
	tbl.AddColumn("int", 20, CELLINT, COLJUSTIFYRIGHT)
	tbl.AddColumn("float", 20, CELLFLOAT, COLJUSTIFYRIGHT)
	now := time.Now()
	rand.Seed(now.UnixNano())
	var ibuf []int64
	var fbuf []float64

	n := rand.Intn(48) + 3 // at least 3, no more than 50
	for i := 0; i < n; i++ {
		ibuf = append(ibuf, int64(rand.Intn(1000000)))
		fbuf = append(fbuf, rand.Float64())
	}
	t.Logf("BEFORE ADDING ROWS: tbl.RowCount = %d\n", tbl.RowCount())
	for i := 0; i < n; i++ {
		if i%2 == 0 {
			tbl.AddRow()
			tbl.Puti(-1, 0, ibuf[i])
			tbl.Putf(-1, 1, fbuf[i])
		} else {
			x := tbl.RowCount() - 1
			t.Logf("Insert Row at x = %d\n", x)
			tbl.InsertRow(x)
			tbl.Puti(x, 0, ibuf[i])
			tbl.Putf(x, 1, fbuf[i])
		}
		t.Logf("i = %d. tbl.RowCount = %d\n", i, tbl.RowCount())
	}

	// before adding summary rows, etc. we can hit some parameter correction code
	// now. The only rows in the table at the moment are the data rows. So the sum
	// will be valid even after the summary rows are added.
	csave := tbl.SumRows(1, -1, n+1000)

	tbl.Sort(0, tbl.RowCount()-1, 1) // sort by the floats first
	for i := 0; i < n-1; i++ {
		if tbl.Getf(i, 1) > tbl.Getf(i+1, 1) {
			t.Logf("sum_test: float sort failed:  %f > %f\n", tbl.Getf(i, 1), tbl.Getf(i+1, 1))
			t.Fail()
		}
	}
	tbl.Sort(0, tbl.RowCount()-1, 0) // sort by the ints
	for i := 0; i < n-1; i++ {
		if tbl.Getf(i, 0) > tbl.Getf(i+1, 0) {
			t.Logf("sum_test: int sort failed:  %d > %d\n", tbl.Geti(i, 0), tbl.Geti(i+1, 0))
			t.Fail()
		}
	}
	tbl.AddLineAfter(tbl.RowCount() - 1)
	tbl.InsertSumRow(-1, 0, tbl.RowCount()-1, []int{0, 1})

	tbl.SetTitle("T1")
	s1 := tbl.String()

	// add a row, then delete it.  This should have cause the totals not to change
	// but it validates that
	tbl.InsertRow(1) // insert at row 1
	tbl.SetTitle("T1.1")
	s11 := tbl.String()
	tbl.DeleteRow(1) // delete it

	tbl.SetTitle("T2")
	s2 := tbl.String()

	x := float64(0)
	j := int64(0)
	tri := tbl.CreateRowset() // a rowset to be used to validate SumRowset()
	for i := 0; i < n; i++ {
		j += tbl.Geti(i, 0)
		x += tbl.Getf(i, 1)
		tbl.AppendToRowset(tri, i)
		t.Logf("%d   %10.8f\n", tbl.Geti(i, 0), tbl.Getf(i, 1))
	}
	if j != tbl.Geti(n, 0) {
		t.Logf("sum_test: int sum failed:  %d != %d\n", j, tbl.Geti(n+1, 0))
		t.Fail()
	}
	c := tbl.SumRowset(tri, 0)
	if j != c.Ival {
		t.Logf("sum_test: SumRowset:  %d != %d\n", j, tbl.Geti(n+1, 0))
		t.Fail()
	}

	if x != tbl.Getf(n, 1) {
		t.Logf("sum_test: float sum failed:  %10.8f != %10.8f\n", x, tbl.Getf(n+1, 1))
		t.Fail()
	}

	// now go back and validate the sum computed from above
	dx := math.Abs(csave.Fval - x)
	if dx > 0.000001 {
		t.Logf("sum_test: float sum failed:  %10.8f != %10.8f.  The difference is: %10.8f\n", csave.Fval, x, csave.Fval-x)
		t.Fail()
	}

	tbl.SetTitle("T3")
	s3 := tbl.String()

	// add a double line before the summary row...  so add line before last row
	tbl.AddLineBefore(tbl.RowCount() - 1)
	tbl.SetTitle("T4")
	s4 := tbl.String()

	// t.Logf("\n\nTable values:\n")
	// for i := 0; i < tbl.RowCount(); i++ {
	// 	j := tbl.Geti(i, 0)
	// 	x := tbl.Getf(i, 1)
	// 	t.Logf("%d.  %d   %10.8f\n", i, j, x)
	// }

	saveTableToFile(t, "sum_test.txt", s1+"\n\n"+s11+"\n\n"+s2+"\n\n"+s3+"\n\n"+s4)
}
