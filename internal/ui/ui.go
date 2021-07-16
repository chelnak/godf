package ui

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/spf13/viper"

	. "github.com/ahmetb/go-linq/v3"
	df "github.com/chelnak/godf/internal/factory"
)

var (
	grid  *tview.Grid
	table *tview.Table
	app   *tview.Application
	stats *tview.TextView
)

func write() {

	runs := df.GetPipelineRuns()

	for i, run := range runs {

		i = i + 1
		colour := tcell.ColorWhite

		if run.Status == "Succeeded" {
			colour = tcell.ColorGreen
		}

		if run.Status == "InProgress" {
			colour = tcell.ColorYellow
		}

		if run.Status == "Failed" {
			colour = tcell.ColorRed
		}

		// r := run.Name
		// if len(run.Name) > 30 {
		// 	r = fmt.Sprintf("%s..", run.Name[0:30])
		// }

		//run.Name

		table.SetCell(i, 0, tview.NewTableCell(fmt.Sprintf("test-pipeline-%d", i)).SetMaxWidth(40))
		table.SetCell(i, 1, tview.NewTableCell(run.Start.Format(time.RFC3339)))

		runEndVal := ""
		if run.End != nil {
			runEndVal = run.End.Format(time.RFC3339)
		}
		table.SetCell(i, 2, tview.NewTableCell(runEndVal))
		table.SetCell(i, 3, tview.NewTableCell(run.Status).SetTextColor(colour))

	}

	inProgress := From(runs).Where(func(c interface{}) bool {
		return c.(df.PipelineRun).Status == "InProgress"
	}).Count()

	succeeded := From(runs).Where(func(c interface{}) bool {
		return c.(df.PipelineRun).Status == "Succeeded"
	}).Count()

	failed := From(runs).Where(func(c interface{}) bool {
		return c.(df.PipelineRun).Status == "Failed"
	}).Count()

	statText := `[yellow]InProgress[white] %d | [green]Succeeded[white] %d | [red]Failed[white] %d | Last Updated: %s`

	stats.SetText(fmt.Sprintf(statText, inProgress, succeeded, failed, time.Now().Format("Mon Jan 2 15:04:05")))
}

func refresh() {

	refreshIntervalSeconds := viper.GetInt64("RefreshIntervalSeconds")
	refreshInterval := time.Duration(refreshIntervalSeconds) * time.Second
	tick := time.NewTicker(refreshInterval)

	for {
		select {
		case <-tick.C:
			app.QueueUpdateDraw(func() {
				write()
			})
		}
	}
}

func newTextView(text string, title string, borders bool, align int) *tview.TextView {

	textView := tview.NewTextView()
	textView.SetTextAlign(align)
	textView.SetText(text)
	textView.SetBorder(borders)
	textView.SetTitle(title)
	textView.SetDynamicColors(true)
	textView.SetRegions(true)
	return textView
}

func Draw() {
	app = tview.NewApplication()

	// Pipelines table
	table = tview.NewTable()
	table.SetTitle("Pipelines")
	table.SetBorders(false)
	table.SetBorder(true)
	table.SetBorderPadding(1, 1, 5, 5)
	table.SetSeparator('\t')
	table.SetFixed(1, 0)
	table.SetCell(0, 0, tview.NewTableCell("Name").SetAttributes(tcell.AttrBold))
	table.SetCell(0, 1, tview.NewTableCell("Start").SetAttributes(tcell.AttrBold))
	table.SetCell(0, 2, tview.NewTableCell("End").SetAttributes(tcell.AttrBold))
	table.SetCell(0, 3, tview.NewTableCell("Status").SetAttributes(tcell.AttrBold))

	// Text views
	title := newTextView("Data Factory Monitor", "", false, tview.AlignCenter)
	stats = newTextView("", "Stats", true, tview.AlignCenter)
	stats.SetBorderPadding(1, 1, 0, 0)

	// Layout
	grid = tview.NewGrid()

	grid.SetRows(3, 0, 5).
		SetColumns(0, 150, 0).
		SetBorders(false).
		AddItem(title, 0, 0, 1, 3, 0, 0, false)

	// Layout for screens narrower than 100 cells (menu and side bar are hidden).
	// AddItem(p tview.Primitive, row int, column int, rowSpan int, colSpan int, minGridHeight int, minGridWidth int, focus bool)
	// grid.AddItem(tview.NewBox(), 0, 0, 0, 0, 0, 0, false).
	// 	AddItem(table, 1, 0, 1, 3, 0, 0, false).
	// 	AddItem(tview.NewBox(), 0, 0, 0, 0, 0, 0, false)

	// Layout for screens wider than 100 cells.
	grid.AddItem(tview.NewBox(), 1, 0, 2, 1, 0, 0, false).
		AddItem(table, 1, 1, 1, 1, 0, 0, false).
		AddItem(stats, 2, 1, 1, 1, 0, 0, false).
		AddItem(tview.NewBox(), 1, 2, 2, 1, 0, 0, false)

	// Write the initial screen
	write()

	// Update screen on ticker
	go refresh()

	if err := app.SetRoot(grid, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

}
