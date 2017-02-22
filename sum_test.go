package gotable

import (
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

	// n := rand.Intn(48) + 3 // at least 3, no more than 50
	n := 3
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

	x := float64(0)
	j := int64(0)
	for i := 0; i < n; i++ {
		j += tbl.Geti(i, 0)
		x += tbl.Getf(i, 1)
		t.Logf("%d   %10.8f\n", tbl.Geti(i, 0), tbl.Getf(i, 1))
	}
	if j != tbl.Geti(n, 0) {
		t.Logf("sum_test: int sum failed:  %d != %d\n", j, tbl.Geti(n+1, 0))
		t.Fail()
	}
	if x != tbl.Getf(n, 1) {
		t.Logf("sum_test: float sum failed:  %10.8f != %10.8f\n", x, tbl.Getf(n+1, 1))
		t.Fail()
	}

	// t.Logf("\n\nTable values:\n")
	// for i := 0; i < tbl.RowCount(); i++ {
	// 	j := tbl.Geti(i, 0)
	// 	x := tbl.Getf(i, 1)
	// 	t.Logf("%d.  %d   %10.8f\n", i, j, x)
	// }

	// s := fmt.Sprintf("%s\n", tbl.String())
	// saveTableToFile(t, "sum_test.txt", s)
}
