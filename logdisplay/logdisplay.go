package logdisplay

import (
	"log"
	t "github.com/gizak/termui"
	"github.com/oskanaan/golog/logreader"
)

type LogDisplay struct {
	logReader logreader.LogReader
	currentPage *[][]string
}

func NewLogDisplay(logReader logreader.LogReader) LogDisplay{
	var l LogDisplay
	l.logReader = logReader
	return l
}

func (l LogDisplay) Display() {
	err := t.Init()
	if err != nil {
		log.Fatalln("Cannot initialize Golog", err)
	}
	defer t.Close()

	l.Tail()
	l.renderBody()

	t.Handle("/sys/wnd/resize", func(t.Event) {
		l.renderBody()
	})

	t.Handle("/sys/kbd/C-c", func(t.Event) {
		t.StopLoop()
	})

	t.Loop()

}

func (l LogDisplay) getHeader() *t.Row{
	var cols []*t.Row
	for index, val := range l.logReader.GetColumnSizes() {
		column := t.NewList()
		column.Height = t.TermHeight()
		column.BorderLabel = l.logReader.GetHeaders()[index]
		column.BorderLabelFg = t.ColorGreen
		column.BorderFg = t.ColorGreen
		column.ItemFgColor = t.ColorWhite
		l.Tail()

		column.Items = func () []string {
			var colData []string
			for _, val := range *l.Tail() {
				colData = append(colData, val[index])
			}

			return colData
		}()
		cols = append(cols, t.NewCol(val, 0, column))
	}
	return t.NewRow(cols...)
}

func (l LogDisplay) renderBody() {
	t.Clear()
	t.Body.Rows = []*t.Row{}

	t.Body.AddRows(
		t.NewRow(
			l.getHeader()),
	)

	t.Body.Align()
	t.Render(t.Body)
}

func (l LogDisplay) Tail() *[][]string{
	tail := l.logReader.Tail()
	return &tail
}