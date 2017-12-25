//logdisplay package displays data processed by logreader.LogReader to the configured output (io.Writer)
package logdisplay

import (
	"github.com/oskanaan/golog/logreader"
	"text/tabwriter"
	"fmt"
	"github.com/jroimartin/gocui"
	"log"
	"bytes"
)

type LogDisplay struct {
	logReader logreader.LogReader
	currentPage *[][]string
}

//Returns a new instance of a LogDisplay
func NewLogDisplay(logReader logreader.LogReader) LogDisplay{
	var l LogDisplay
	l.logReader = logReader
	return l
}

func (l LogDisplay) DisplayUI() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	g.Cursor = true
	g.Mouse = true

	g.SetManagerFunc(l.layout)

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func (l LogDisplay) DisplayStdout() {
	fmt.Print(l.getFormattedLog())
}

func (l LogDisplay) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("main", 0, -1, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintf(v, "%s", l.getFormattedLog())
		v.Editable = false
		v.Wrap = false
		if _, err := g.SetCurrentView("main"); err != nil {
			return err
		}
	}
	return nil
}

//Writes the log data to the terminal
//The data will be read from the logReader and displayed in a column format using tabwriter
func (l LogDisplay) getFormattedLog() string {
	tabWriter := new(tabwriter.Writer)
	var output bytes.Buffer
	tabWriter.Init(&output, 0, 8, 0, '\t', tabwriter.TabIndent)

	l.writeHeader(tabWriter)

	l.writeBody(tabWriter)

	tabWriter.Flush()

	return output.String()
}

//Writes the header of the log file.
func (l LogDisplay) writeHeader(writer *tabwriter.Writer) {
	var header string
	for index, columnHeader := range l.logReader.GetHeaders() {
		header = header + l.formatColumnText(columnHeader, index)
		if index < len(l.logReader.GetHeaders())-1 {
			header = header + "\t"
		}
	}
	fmt.Fprintln(writer, colorizeHeader(header))
}

func (l LogDisplay) writeBody(tabWriter *tabwriter.Writer) {
	severityIndex := l.getSeverityColumnIndex()
	for _, row := range *l.Tail() {
		var rowText string
		for index, columnText := range row {
			formattedColText := l.formatColumnText(columnText, index)
			rowText = rowText + formattedColText
			if index < len(columnText)-1 {
				rowText = rowText + "\t"
			}
		}

		if rowText == "" {
			continue
		}
		fmt.Fprintln(tabWriter, colorizeLogEntry(rowText, row[severityIndex]))
	}
}

func (l LogDisplay) getSeverityColumnIndex() int {
	for index, header := range l.logReader.GetHeaders() {
		if header == l.logReader.GetSeverityColumnName() {
			return index
		}
	}

	return -1
}

//Returns the tail data based on the "capacity" configuration passed to the program
func (l LogDisplay) Tail() *[][]string{
	tail := l.logReader.Tail()
	return &tail
}

//Applies column formatting based on the program parameters.
//Currently there is only the column size, which shows N characters from the end of the text where N is the configured value.
//If the column size is negative, it will return the text as is.
//If the text length is less than the column size, it adds some right padding to match the column size
func (l LogDisplay) formatColumnText(text string, columnIndex int) string {
	configuredSize := l.logReader.GetColumnSizes()[columnIndex]
	startIndex := len(text) - l.logReader.GetColumnSizes()[columnIndex]
	formattedText := text
	if configuredSize >= 0 && configuredSize <= len(text) {
		formattedText = text[startIndex:]
	}

	if configuredSize > 0 && len(formattedText) < configuredSize {
		formattedText = func () string {
			padded := formattedText
			paddedLength := len(padded)
			for i := 0 ; i < configuredSize - paddedLength ; i++ {
				padded = padded + " "
			}
			return padded
		} ()
	}

	return formattedText
}

//Binds keyboard keys and mouse buttons to actions
//CTRL-C : quit
//Mouse Left: show log entry details
//Arrow Down: scroll down
//Arrow Up: scroll up
func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("main", gocui.MouseLeft, gocui.ModNone, showLogEntryDetails); err != nil {
		return err
	}

	if err := g.SetKeybinding("msg", gocui.MouseLeft, gocui.ModNone, hideLogEntryDetails); err != nil {
		return err
	}

	if err := g.SetKeybinding("main", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}

	if err := g.SetKeybinding("main", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}

	return nil
}

//Quits the application
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

//Shows the full details of the log entry in a popup window
func showLogEntryDetails(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	if _, err := g.SetCurrentView(v.Name()); err != nil {
		return err
	}

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	maxX, maxY := g.Size()
	if v, err := g.SetView("msg", maxX/2-100, maxY/2-10, maxX/2+100, maxY/2+10); err != nil {
		v.Wrap = true
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, l)
	}
	return nil
}

//Hides the log entry details popup
func hideLogEntryDetails(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("msg"); err != nil {
		return err
	}
	return nil
}

//Scroll down the log file lines
func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("msg"); err != nil {
		return err
	}
	return nil
}

//Scroll up the log file lines
func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("msg"); err != nil {
		return err
	}
	return nil
}