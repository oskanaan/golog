//logdisplay package displays data processed by logreader.LogReader to the configured output (io.Writer)
package logdisplay

import (
	"github.com/oskanaan/golog/logreader"
	"text/tabwriter"
	"fmt"
	"github.com/jroimartin/gocui"
	"log"
	"bytes"
	"time"
	"sync"
)

var wg sync.WaitGroup

type LogDisplay struct {
	logReader logreader.LogReader
	currentPage *[][]string
	tailOn *bool
}

//Returns a new instance of a LogDisplay
func NewLogDisplay(logReader logreader.LogReader) LogDisplay{
	var l LogDisplay
	l.logReader = logReader
	l.tailOn = &[]bool{true}[0]
	return l
}

func (l *LogDisplay) DisplayUI() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	g.Cursor = true
	g.Mouse = true
	l.tail()

	g.SetManagerFunc(l.layout)

	if err := l.keybindings(g); err != nil {
		log.Panicln(err)
	}

	wg.Add(1)
	go func(g *gocui.Gui, l *LogDisplay){
		for {
			if *l.tailOn {

				g.Update(func(g *gocui.Gui) error {
					v, err := g.View("main")
					if err != nil {
						return err
					}
					l.currentPage = l.logReader.Tail()
					v.Clear()
					fmt.Fprintf(v, "%s", l.getFormattedLog())
					return nil
				})
			}

			time.Sleep(500 * time.Millisecond)
		}
	}(g, l)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

	wg.Wait()
}

func (l LogDisplay) DisplayStdout() {
	l.tail()
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
	for _, row := range *l.currentPage {
		var rowText string
		if len(row) == 1 {
			rowText = row[0]
		} else {
			for index, columnText := range row {
				formattedColText := l.formatColumnText(columnText, index)
				rowText = rowText + formattedColText
				if index < len(columnText)-1 {
					rowText = rowText + "\t"
				}
			}
		}

		if rowText == "" {
			continue
		}

		severity := ""
		if len(row) > severityIndex {
			severity = row[severityIndex]
		}

		fmt.Fprintln(tabWriter, colorizeLogEntry(rowText, severity))
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
func (l *LogDisplay) tail() {
	l.currentPage = l.logReader.Tail()
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
//Page Down: scroll one page down
//Page Up: scroll one page up
//Arrow Down: scroll down
//Arrow Up: scroll up
func (l *LogDisplay) keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("main", gocui.MouseLeft, gocui.ModNone, showLogEntryDetails); err != nil {
		return err
	}

	if err := g.SetKeybinding("msg", gocui.MouseLeft, gocui.ModNone, hideLogEntryDetails); err != nil {
		return err
	}

	if err := g.SetKeybinding("main", gocui.KeyPgdn, gocui.ModNone, l.pageDown); err != nil {
		return err
	}

	if err := g.SetKeybinding("main", gocui.KeyPgup, gocui.ModNone, l.pageUp); err != nil {
		return err
	}

	if err := g.SetKeybinding("main", gocui.KeyArrowUp, gocui.ModNone, l.arrowUp); err != nil {
		return err
	}

	if err := g.SetKeybinding("main", gocui.KeyArrowDown, gocui.ModNone, l.arrowDown); err != nil {
		return err
	}

	if err := g.SetKeybinding("main", gocui.KeyHome, gocui.ModNone, l.home); err != nil {
		return err
	}

	if err := g.SetKeybinding("main", gocui.KeyEnd, gocui.ModNone, l.end); err != nil {
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

func (l *LogDisplay) home(g *gocui.Gui, v *gocui.View) error {
	g.Update(func(g *gocui.Gui) error {
		v, err := g.View("main")
		if err != nil {
			return err
		}
		l.currentPage = l.logReader.Head()
		v.Clear()
		fmt.Fprintf(v, "%s", l.getFormattedLog())
		//fmt.Fprintf(v, "%s", l.getFormattedLog())
		return nil
	})

	return nil

}

func (l *LogDisplay) end(g *gocui.Gui, v *gocui.View) error {
	g.Update(func(g *gocui.Gui) error {
		v, err := g.View("main")
		if err != nil {
			return err
		}
		l.currentPage = l.logReader.Tail()
		v.Clear()
		fmt.Fprintf(v, "%s", l.getFormattedLog())
		//fmt.Fprintf(v, "%s", l.getFormattedLog())
		return nil
	})
	l.tailOn = &[]bool{true}[0]

	return nil

}


func (l *LogDisplay) pageDown(g *gocui.Gui, v *gocui.View) error {
	g.Update(func(g *gocui.Gui) error {
		v, err := g.View("main")
		if err != nil {
			return err
		}
		l.currentPage = l.logReader.PageDown()
		v.Clear()
		fmt.Fprintf(v, "%s", l.getFormattedLog())
		//fmt.Fprintf(v, "%s", l.getFormattedLog())
		return nil
	})

	return nil

}

func (l *LogDisplay) pageUp(g *gocui.Gui, v *gocui.View) error {
	l.tailOn = &[]bool{false}[0]

	g.Update(func(g *gocui.Gui) error {
		v, err := g.View("main")
		if err != nil {
			return err
		}
		l.currentPage = l.logReader.PageUp()
		v.Clear()
		fmt.Fprintf(v, "%s", l.getFormattedLog())
		return nil
	})

	return nil
}

func (l *LogDisplay) arrowDown(g *gocui.Gui, v *gocui.View) error {
	g.Update(func(g *gocui.Gui) error {
		v, err := g.View("main")
		if err != nil {
			return err
		}
		l.currentPage = l.logReader.Down()
		v.Clear()
		fmt.Fprintf(v, "%s", l.getFormattedLog())
		//fmt.Fprintf(v, "%s", l.getFormattedLog())
		return nil
	})

	return nil

}

func (l *LogDisplay) arrowUp(g *gocui.Gui, v *gocui.View) error {
	l.tailOn = &[]bool{false}[0]

	g.Update(func(g *gocui.Gui) error {
		v, err := g.View("main")
		if err != nil {
			return err
		}
		l.currentPage = l.logReader.Up()
		v.Clear()
		fmt.Fprintf(v, "%s", l.getFormattedLog())
		//fmt.Fprintf(v, "%s", l.getFormattedLog())
		return nil
	})

	return nil

}
