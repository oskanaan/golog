//logdisplay package displays data processed by logreader.LogReader to a gocui view
package logdisplay

import (
	"bytes"
	"fmt"
	"github.com/jroimartin/gocui"
	"github.com/oskanaan/golog/logreader"
	"log"
	"sync"
	"text/tabwriter"
	"time"
	"github.com/atotto/clipboard"
	"strconv"
)

var wg sync.WaitGroup
const detailsView = "details"
const searchField = "searchField"
const progress = "progress"
const tailIndicator = "tailIndicator"
const searchButton = "searchButton"
const mainView = "mainView"

type LogDisplayConfig struct {
	Severities []Severity
	Files      []LogFile `yaml:files`
	Search     Search    `yaml:search`
}

type LogFile struct {
	LogFile string `yaml:"file"`
	Name string `yaml:"name"`
}

type Severity struct {
	Severity string `yaml:"severity"`
	Colors []interface{} `yaml:"colors"`
}

type Search struct {
	HighlightColor []interface{} `yaml:"highlightColor"`
}

type LogDisplay struct {
	logReader   *logreader.LogReader
	currentPage *[][]string
	searchResultLocation int
	tailOn      []*bool
	searchOn    *bool
	logFileIndex int
	logdisplayConfig *LogDisplayConfig
}

//Returns a new instance of a LogDisplay
func NewLogDisplay(logReader *logreader.LogReader, logdisplayConfig *LogDisplayConfig) LogDisplay {
	var l LogDisplay
	l.logReader = logReader
	l.logdisplayConfig = logdisplayConfig
	l.tailOn = make([]*bool, len(l.logdisplayConfig.Files))
	for index := 0 ; index < len(l.logdisplayConfig.Files) ; index ++ {
		l.tailOn[index] = &[]bool{true}[0]
	}
	l.searchResultLocation = -1
	return l
}

//Displays the log using gocui for the UI
func (l *LogDisplay) DisplayUI() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	g.Cursor = true
	g.Mouse = true
	//l.tail()

	g.SetManagerFunc(l.layout)

	if err := l.keybindings(g); err != nil {
		log.Panicln(err)
	}

	wg.Add(1)
	go func(g *gocui.Gui, l *LogDisplay) {
		for {
			if *l.tailOn[l.logFileIndex] {

				g.Update(func(g *gocui.Gui) error {
					_, err := g.View(mainView)
					if err != nil {
						return err
					}
					l.currentPage = l.logReader.Tail()
					l.rerender(g)
					return nil
				})
			}

			time.Sleep(1 * time.Second)
		}
	}(g, l)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

	wg.Wait()
}

//Prints the log to the stdout, used for debugging purposes only
func (l LogDisplay) DisplayStdout() {
	l.logReader.SetCapacity(50)
	l.tail()
	fmt.Print(l.formattedLog())
}

func (l LogDisplay) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	//Set the capacity based on the terminal size
	l.logReader.SetCapacity(maxY - 4)
	if v, err := g.SetView(mainView, 0, -1, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = false
		v.Wrap = false
		if _, err := g.SetCurrentView(mainView); err != nil {
			return err
		}
	}

	if v, err := g.SetView(progress, maxX-6, maxY-3, maxX-2, maxY-1); err != nil {
		v.Wrap = false
		v.Editable = true

		v.Title = "%"
		if err != gocui.ErrUnknownView {
			return err
		}
	}

	if v, err := g.SetView(tailIndicator, maxX-15, maxY-3, maxX-8, maxY-1); err != nil {
		v.Wrap = false
		v.Editable = true

		v.Title = "Tail"
		if err != gocui.ErrUnknownView {
			return err
		}
	}

	for index, file := range l.logdisplayConfig.Files {
		x0 := 4 + (20 * index)
		if v, err := g.SetView(file.Name, x0, maxY - 3, x0 + 15 , maxY - 1); err != nil {
			v.Wrap = false
			v.Editable = true
			v.Title = `F`+strconv.Itoa(index + 1)
			var fileName string

			if l.logFileIndex == index {
				fileName = colorizeActive(file.Name, l.logdisplayConfig)
			} else {
				fileName = file.Name
			}

			fmt.Fprint(v, fileName)
			if err != gocui.ErrUnknownView {
				return err
			}
		}
	}

	return nil
}

//Writes the log data to a string
//The data will be read from the logReader and displayed in a column format using tabwriter
//Returns a string containing the current log page formatted using tabwriter
func (l LogDisplay) formattedLog() string {
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
	fmt.Fprintln(writer, colorizeHeader(header, l.logdisplayConfig))
}

//Writes the formatted current page of the log to a tabwriter
func (l LogDisplay) writeBody(tabWriter *tabwriter.Writer) {
	for index, row := range *l.currentPage {
		var rowText string
		//If this is a stack trace or some debugging information then no parsing is needed, display as is
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

		if index == l.searchResultLocation {
			fmt.Fprintln(tabWriter, colorizeLogEntry(rowText, l.logdisplayConfig, true))
		} else {
			fmt.Fprintln(tabWriter, colorizeLogEntry(rowText, l.logdisplayConfig, false))
		}
	}
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
	if columnIndex >= len(l.logReader.GetColumnSizes()) {
		return text
	}
	configuredSize := l.logReader.GetColumnSizes()[columnIndex]
	startIndex := len(text) - l.logReader.GetColumnSizes()[columnIndex]
	formattedText := text
	if configuredSize >= 0 && configuredSize <= len(text) {
		formattedText = text[startIndex:]
	}

	if configuredSize > 0 && len(formattedText) < configuredSize {
		formattedText = func() string {
			padded := formattedText
			paddedLength := len(padded)
			for i := 0; i < configuredSize-paddedLength; i++ {
				padded = padded + " "
			}
			return padded
		}()
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
//Key home: navigates to the beginning of the log
//End: tails and follows the log
func (l *LogDisplay) keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.SetKeybinding(mainView, gocui.MouseLeft, gocui.ModNone, l.showLogEntryDetails); err != nil {
		return err
	}

	if err := g.SetKeybinding(detailsView, gocui.MouseRight, gocui.ModNone, hideLogEntryDetails); err != nil {
		return err
	}

	if err := g.SetKeybinding(mainView, gocui.KeyPgdn, gocui.ModNone, l.pageDown); err != nil {
		return err
	}

	if err := g.SetKeybinding(mainView, gocui.KeyPgup, gocui.ModNone, l.pageUp); err != nil {
		return err
	}

	if err := g.SetKeybinding(mainView, gocui.KeyArrowUp, gocui.ModNone, l.arrowUp); err != nil {
		return err
	}

	if err := g.SetKeybinding(mainView, gocui.KeyArrowDown, gocui.ModNone, l.arrowDown); err != nil {
		return err
	}

	if err := g.SetKeybinding(mainView, gocui.KeyHome, gocui.ModNone, l.home); err != nil {
		return err
	}

	if err := g.SetKeybinding(mainView, gocui.KeyEnd, gocui.ModNone, l.end); err != nil {
		return err
	}

	if err := g.SetKeybinding(mainView, gocui.KeyF1, gocui.ModNone, l.switchToFile1); err != nil {
		return err
	}

	if err := g.SetKeybinding(mainView, gocui.KeyF2, gocui.ModNone, l.switchToFile2); err != nil {
		return err
	}

	if err := g.SetKeybinding(mainView, gocui.KeyF3, gocui.ModNone, l.switchToFile3); err != nil {
		return err
	}

	if err := g.SetKeybinding(mainView, gocui.KeyF4, gocui.ModNone, l.switchToFile4); err != nil {
		return err
	}

	if err := g.SetKeybinding(mainView, gocui.KeyF5, gocui.ModNone, l.switchToFile5); err != nil {
		return err
	}

	if err := g.SetKeybinding(mainView, gocui.KeyF6, gocui.ModNone, l.switchToFile6); err != nil {
		return err
	}

	if err := g.SetKeybinding(mainView, gocui.KeyF7, gocui.ModNone, l.switchToFile7); err != nil {
		return err
	}

	if err := g.SetKeybinding(mainView, gocui.KeyF8, gocui.ModNone, l.switchToFile8); err != nil {
		return err
	}

	if err := g.SetKeybinding(mainView, gocui.KeyF9, gocui.ModNone, l.switchToFile9); err != nil {
		return err
	}

	if err := g.SetKeybinding(mainView, gocui.KeyF10, gocui.ModNone, l.switchToFile10); err != nil {
		return err
	}

	if err := g.SetKeybinding(mainView, gocui.KeySpace, gocui.ModNone, l.search); err != nil {
		return err
	}

	if err := g.SetKeybinding(searchField, gocui.KeyEnd, gocui.ModNone, l.exitSearch); err != nil {
		return err
	}

	if err := g.SetKeybinding(searchField, gocui.KeyEnter, gocui.ModNone, l.performSearch); err != nil {
		return err
	}

	return nil
}

//Quits the application
func quit(g *gocui.Gui, v *gocui.View) error {
	wg.Done()
	return gocui.ErrQuit
}

//Shows the full details of the log entry in a popup window
//Returns an error if the detailsView view cannot be found
func (l LogDisplay) showLogEntryDetails(g *gocui.Gui, v *gocui.View) error {
	if _, err := g.SetCurrentView(v.Name()); err != nil {
		return err
	}

	_, lineNum := v.Cursor()
	message := l.logReader.Message(lineNum)
	//Write to clipboard
	clipboard.WriteAll(message)

	maxX, maxY := g.Size()
	if v, err := g.SetView(detailsView, 5, 5, maxX-5, maxY-5); err != nil {
		v.Wrap = true
		v.Editable = true

		v.Title = "Line details"
		if err != gocui.ErrUnknownView {
			return err
		}

		fmt.Fprintln(v, message)
	}
	return nil
}

//Hides the log entry details popup
//Returns an error if the detailsView view cannot be found
func hideLogEntryDetails(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView(detailsView); err != nil {
		return err
	}
	return nil
}

//Navigates to the beginning of the file
//Returns an error if the main view cannot be found
func (l *LogDisplay) home(g *gocui.Gui, v *gocui.View) error {
	l.tailOn[l.logFileIndex] = &[]bool{false}[0]
	l.exitSearchMode()
	g.Update(func(g *gocui.Gui) error {
		_, err := g.View(mainView)
		if err != nil {
			return err
		}
		l.currentPage = l.logReader.Head()
		l.rerender(g)
		return nil
	})

	return nil

}

//Navigates to the end of the file
//Returns an error if the main view cannot be found
func (l *LogDisplay) end(g *gocui.Gui, v *gocui.View) error {
	l.exitSearchMode()
	g.Update(func(g *gocui.Gui) error {
		_, err := g.View(mainView)
		if err != nil {
			return err
		}
		l.currentPage = l.logReader.Tail()
		l.rerender(g)
		return nil
	})
	l.tailOn[l.logFileIndex] = &[]bool{true}[0]

	return nil

}

//Sets first log file to the active one
func (l *LogDisplay) switchToFile1(g *gocui.Gui, v *gocui.View) error {
	l.switchToFile(g, v, 0)
	return nil
}

//Sets second log file to the active one
func (l *LogDisplay) switchToFile2(g *gocui.Gui, v *gocui.View) error {
	l.switchToFile(g, v, 1)
	return nil
}

//Sets third log file to the active one
func (l *LogDisplay) switchToFile3(g *gocui.Gui, v *gocui.View) error {
	l.switchToFile(g, v,2)
	return nil
}

//Sets fourth log file to the active one
func (l *LogDisplay) switchToFile4(g *gocui.Gui, v *gocui.View) error {
	l.switchToFile(g, v,3)
	return nil
}

//Sets fifth log file to the active one
func (l *LogDisplay) switchToFile5(g *gocui.Gui, v *gocui.View) error {
	l.switchToFile(g, v,4)
	return nil
}

//Sets sixth log file to the active one
func (l *LogDisplay) switchToFile6(g *gocui.Gui, v *gocui.View) error {
	l.switchToFile(g, v,5)
	return nil
}

//Sets seventh log file to the active one
func (l *LogDisplay) switchToFile7(g *gocui.Gui, v *gocui.View) error {
	l.switchToFile(g, v,6)
	return nil
}

//Sets eighth log file to the active one
func (l *LogDisplay) switchToFile8(g *gocui.Gui, v *gocui.View) error {
	l.switchToFile(g, v,7)
	return nil
}

//Sets ninth log file to the active one
func (l *LogDisplay) switchToFile9(g *gocui.Gui, v *gocui.View) error {
	l.switchToFile(g, v,8)
	return nil
}

//Sets tenth log file to the active one
func (l *LogDisplay) switchToFile10(g *gocui.Gui, v *gocui.View) error {
	l.switchToFile(g, v,9)
	return nil
}

//Switches to the log file at the specified index
func (l *LogDisplay) switchToFile(g *gocui.Gui, v *gocui.View, index int) error {
	l.exitSearchMode()
	previousFileIndex := l.logFileIndex
	l.logFileIndex = index
	l.logReader.FileIndex = index

	g.Update(func(g *gocui.Gui) error {
		prev, _ := g.View(l.logdisplayConfig.Files[previousFileIndex].Name)
		active, _ := g.View(l.logdisplayConfig.Files[index].Name)
		prev.Clear()
		fmt.Fprint(prev, l.logdisplayConfig.Files[previousFileIndex].Name)
		active.Clear()
		fmt.Fprint(active, colorizeActive(l.logdisplayConfig.Files[index].Name, l.logdisplayConfig))

		_, err := g.View(mainView)
		if err != nil {
			return err
		}
		l.currentPage = l.logReader.Refresh()
		l.rerender(g)
		return nil
	})
	return nil
}

//Displays a box for entering a search term and performing search
//Navigates to the location of the search result
func (l *LogDisplay) search(g *gocui.Gui, v *gocui.View) error {
	l.tailOn[l.logFileIndex] = &[]bool{false}[0]
	maxX, maxY := g.Size()
	if v, err := g.SetView(searchField, 5, maxY-3, maxX-40, maxY-1); err != nil {
		v.Wrap = false
		v.Editable = true

		v.Title = "Search"
		if err != gocui.ErrUnknownView {
			return err
		}

	}
	g.SetCurrentView(searchField)
	return nil
}

//Exits the search mode
//Returns an error if the searchView view cannot be found
func (l LogDisplay) exitSearch(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView(v.Name()); err != nil {
		return err
	}
	g.SetCurrentView(mainView)

	return nil
}

//Sets the search mode flag to true and calls the search func on logreader
func (l LogDisplay) performSearch(g *gocui.Gui, v *gocui.View) error {
	l.searchOn = &[]bool{true}[0]
	l.tailOn[l.logFileIndex] = &[]bool{false}[0]
	searchTerm, _ := v.Line(0)
	l.currentPage, l.searchResultLocation = l.logReader.Search(searchTerm, l.searchResultLocation)
	details, _ := g.View(mainView)
	details.Clear()
	l.rerender(g)

	return nil
}

//Navigates one page down where the page size equals the "capacity" configuration
//Returns an error if the main view cannot be found
func (l *LogDisplay) pageDown(g *gocui.Gui, v *gocui.View) error {
	l.tailOn[l.logFileIndex] = &[]bool{false}[0]
	l.exitSearchMode()
	g.Update(func(g *gocui.Gui) error {
		_, err := g.View(mainView)
		if err != nil {
			return err
		}
		l.currentPage = l.logReader.PageDown()
		l.rerender(g)
		return nil
	})

	return nil

}

//Navigates one page up where the page size equals the "capacity" configuration
//Returns an error if the main view cannot be found
func (l *LogDisplay) pageUp(g *gocui.Gui, v *gocui.View) error {
	l.tailOn[l.logFileIndex] = &[]bool{false}[0]
	l.exitSearchMode()
	g.Update(func(g *gocui.Gui) error {
		_, err := g.View(mainView)
		if err != nil {
			return err
		}
		l.currentPage = l.logReader.PageUp()
		l.rerender(g)
		return nil
	})

	return nil
}

//Navigates one line down
//Returns an error if the main view cannot be found
func (l *LogDisplay) arrowDown(g *gocui.Gui, v *gocui.View) error {
	l.exitSearchMode()
	g.Update(func(g *gocui.Gui) error {
		_, err := g.View(mainView)
		if err != nil {
			return err
		}
		l.currentPage = l.logReader.Down()
		l.rerender(g)
		return nil
	})

	return nil

}

//Navigates one line up
//Returns an error if the main view cannot be found
func (l *LogDisplay) arrowUp(g *gocui.Gui, v *gocui.View) error {
	l.tailOn[l.logFileIndex] = &[]bool{false}[0]
	l.exitSearchMode()
	g.Update(func(g *gocui.Gui) error {
		_, err := g.View(mainView)
		if err != nil {
			return err
		}
		l.currentPage = l.logReader.Up()
		l.rerender(g)
		return nil
	})

	return nil

}

func (l *LogDisplay) rerender(g *gocui.Gui){
	mainViewWidget, _ := g.View(mainView)
	mainViewWidget.Clear()
	fmt.Fprintf(mainViewWidget, "%s", l.formattedLog())

	progressWidget, _ := g.View(progress)
	progressWidget.Clear()
	fmt.Fprintf(progressWidget, "%d%%", l.logReader.Progress())

	tailWidget, _ := g.View(tailIndicator)
	tailWidget.Clear()
	if *l.tailOn[l.logFileIndex] {
		fmt.Fprintf(tailWidget, " \033[3%d;%d;1m%s\033[0m", 2, 4, "ON")
	} else {
		fmt.Fprintf(tailWidget, " \033[3%d;%d;1m%s\033[0m", 1, 4, "OFF")
	}

}

func (l *LogDisplay) exitSearchMode() {
	l.searchOn = &[]bool{false}[0]
	l.searchResultLocation = -1
}