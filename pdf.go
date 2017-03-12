package gotable

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

// WKHTMLTOPDFCMD command : html > pdf
const (
	WKHTMLTOPDFCMD = "wkhtmltopdf"
	TEMPSTORE      = "/tmp"
	DATETIMEFMT    = "_2 Jan 2006 3:04 PM IST"
)

// PDFTable struct used to prepare table in pdf version
type PDFTable struct {
	*Table
	outbuf bytes.Buffer
}

// Buffer returns the embedded output buffer used if OutputFile is empty
func (pt *PDFTable) Buffer() *bytes.Buffer {
	return &pt.outbuf
}

// Bytes returns the output byte slice from the output buffer used if OutputFile is empty
func (pt *PDFTable) Bytes() []byte {
	return pt.outbuf.Bytes()
}

// WriteFile writes the contents of the output buffer to a file
func (pt *PDFTable) WriteFile(filename string) error {
	return ioutil.WriteFile(filename, pt.Bytes(), 0666)
}

func (pt *PDFTable) getTableOutput() (string, error) {

	// get html output first
	htmlString, err := pt.Table.SprintTable(TABLEOUTHTML)
	if err != nil {
		return "", err
	}

	timeCharReplacer := strings.NewReplacer(":", "-", ".", "", "T", "")
	currentTime := timeCharReplacer.Replace(time.Now().Format(time.RFC3339Nano))

	// create temp file
	filePath := path.Join(TEMPSTORE, "tablePDF_"+currentTime)

	// only works with html file extension
	// be careful, must append it
	tempHTMLFile, err := os.Create(filePath + ".html")
	if err != nil {
		return "", err
	}
	// write html string to file
	_, err = tempHTMLFile.WriteString(htmlString)
	if err != nil {
		return "", err
	}
	tempHTMLFile.Close()

	// remove this temp file after operation
	defer os.Remove(tempHTMLFile.Name())

	// return output file path
	return pt.getPDF(filePath)
}

func (pt *PDFTable) getPDF(inputFile string) (string, error) {

	var err error

	pdfExportTime := time.Now().Format(DATETIMEFMT)
	htmlExportFile := inputFile + ".html"
	pdfExportFile := inputFile + ".pdf"

	cmdArgs := []string{
		// top margin
		"-T", "15",
		// header center content
		"--header-center", pt.Table.GetTitle(),
		// header font size
		"--header-font-size", "8",
		// header font
		"--header-font-name", "opensans",
		// header spacing
		"--header-spacing", "3",
		// bottom margin
		"-B", "15",
		// footer spacing
		"--footer-spacing", "5",
		// footer font
		"--footer-font-name", "opensans",
		// footer font size
		"--footer-font-size", "8",
		// footer left content
		"--footer-left", pdfExportTime,
		// footer right content
		"--footer-right", "Page [page] of [toPage]",
		// page size
		"--page-size", "Letter",
		// orientation
		"--orientation", "Portrait",
		// input, output
		htmlExportFile, "-",
	}

	// prepare command
	wkhtmltopdf := exec.Command(WKHTMLTOPDFCMD, cmdArgs...)

	// REF: https://github.com/aodin/go-pdf-server/blob/master/pdf_server.go

	// get output pipeline
	output, err := wkhtmltopdf.StdoutPipe()
	if err != nil {
		return "", err
	}

	// Begin the command
	if err = wkhtmltopdf.Start(); err != nil {
		return "", err
	}

	// Read the generated PDF from std out
	b, err := ioutil.ReadAll(output)
	if err != nil {
		return "", err
	}

	// End the command
	if err = wkhtmltopdf.Wait(); err != nil {
		return "", err
	}

	pt.outbuf.Write(b)

	err = pt.WriteFile(pdfExportFile)
	if err != nil {
		fmt.Println("Error WriteFile:", err)
	}

	return pdfExportFile, nil
}
