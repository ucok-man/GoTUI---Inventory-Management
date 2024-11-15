package tui

import (
	"fmt"
	"log"
	"strconv"

	"github.com/rivo/tview"
	"github.com/ucok-man/go-tui-inventory-management/src/internal/data"
)

type TUI struct {
	app  *tview.Application
	page struct {
		list       *tview.TextView
		addform    *tview.Form
		deleteform *tview.Form
	}
	components struct {
		navigation *tview.Form
	}
	rightlayout *tview.Flex // this is save so that the layout can be switched

	model *data.InventoryModel
}

func NewTUI(model *data.InventoryModel) *TUI {
	// tui :=
	tui := TUI{
		app:   tview.NewApplication(),
		model: model,
	}

	// navigation components
	tui.components.navigation = tview.NewForm().
		AddButton("Prev", func() {
			tui.rightlayout.Clear()
			tui.rightlayout.
				AddItem(tui.page.addform, 0, 3, true).
				AddItem(tui.components.navigation, 0, 1, false)
		}).
		AddButton("Next", func() {
			tui.rightlayout.Clear()
			tui.rightlayout.
				AddItem(tui.page.deleteform, 0, 3, true).
				AddItem(tui.components.navigation, 0, 1, false)
		}).
		AddButton("Quit", func() {
			tui.app.Stop()
		}).
		SetButtonsAlign(tview.AlignCenter)
	tui.components.navigation.
		SetTitle("Navigation").
		SetBorder(true).
		SetBorderPadding(3, 3, 1, 1)

	// list page config
	tui.page.list = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetScrollable(true)
	tui.page.list.SetBorder(true).
		SetTitle("List Inventory").
		SetBorderPadding(1, 2, 2, 2)

	// add form page config
	namefield := tview.NewInputField().SetLabel("Item Name: ")
	stockfield := tview.NewInputField().SetLabel("Stock: ")
	tui.page.addform = tview.NewForm().
		AddFormItem(namefield).
		AddFormItem(stockfield).
		AddButton("Create", func() {
			if namefield.GetText() == "" || stockfield.GetText() == "" {
				return
			}

			nameinput := namefield.GetText()
			stockinput, err := strconv.Atoi(stockfield.GetText())
			if err != nil {
				fmt.Fprintln(tui.page.list, "Error: You must assign positive number value on stock")
				return
			}

			if stockinput < 0 {
				fmt.Fprintln(tui.page.list, "Error: You must assign positive number value on stock")
				return
			}

			err = tui.model.Add(data.Item{
				Name:  nameinput,
				Stock: stockinput,
			})
			if err != nil {
				tui.app.Stop()
				log.Fatalf("Sorry we have problem in our server: %v\n", err)
			}

			tui.refresh()

			namefield.SetText("")
			stockfield.SetText("")
		})
	tui.page.addform.
		SetBorder(true).
		SetTitle("Create Inventory")

	// delete form page config
	deletefield := tview.NewInputField().SetLabel("Delete Item (id): ")
	tui.page.deleteform = tview.NewForm().
		AddFormItem(deletefield).
		AddButton("Delete", func() {
			inputid := deletefield.GetText()
			if inputid == "" {
				return
			}

			inputidnum, err := strconv.Atoi(inputid)
			if err != nil {
				fmt.Fprintln(tui.page.list, "Error: You must assign positive number value id")
				return
			}

			if inputidnum < 1 || inputidnum > len(tui.model.Get()) {
				fmt.Fprintln(tui.page.list, "Error: ID is not defined in the list")
				return
			}

			if err := tui.model.Delete(inputidnum); err != nil {
				tui.app.Stop()
				log.Fatalf("Sorry we have problem in our server: %v\n", err)
			}
			tui.refresh()

			deletefield.SetText("")
		})
	tui.page.deleteform.
		SetBorder(true).
		SetTitle("Delete Inventory")

	return &tui
}

func (t *TUI) layoutWrapper() *tview.Flex {
	t.rightlayout = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(t.page.addform, 0, 3, true).
		AddItem(t.components.navigation, 0, 1, false)

	flexmain := tview.NewFlex().
		AddItem(t.page.list, 0, 1, false).
		AddItem(t.rightlayout, 0, 1, true)

	return flexmain
}

func (t *TUI) refresh() {
	t.page.list.Clear()
	if len(t.model.Get()) == 0 {
		fmt.Fprintf(t.page.list, "No inventory is added :)")
		return
	}

	for i, item := range t.model.Get() {
		fmt.Fprintf(t.page.list, "[%d] %s (Stock: %d)\n", i+1, item.Name, item.Stock)
	}
}

func (t *TUI) Run() error {
	t.refresh()
	return t.app.SetRoot(t.layoutWrapper(), true).
		EnableMouse(true).
		Run()
}
