// Command widget allows manually testing a single widget.
package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/elves/elvish/pkg/cli"
	"github.com/elves/elvish/pkg/cli/term"
	"github.com/elves/elvish/pkg/ui"
)

var (
	maxHeight  = flag.Int("max-height", 10, "maximum height")
	horizontal = flag.Bool("horizontal", false, "use horizontal listbox layout")
)

func makeWidget() cli.Widget {
	items := cli.TestItems{Prefix: "list item "}
	w := cli.NewComboBox(cli.ComboBoxSpec{
		CodeArea: cli.CodeAreaSpec{
			Prompt: func() ui.Text {
				return ui.T(" NUMBER ", ui.Bold, ui.BgMagenta).ConcatText(ui.T(" "))
			},
		},
		ListBox: cli.ListBoxSpec{
			State:       cli.ListBoxState{Items: &items},
			Placeholder: ui.T("(no items)"),
			Horizontal:  *horizontal,
		},
		OnFilter: func(w cli.ComboBox, filter string) {
			if n, err := strconv.Atoi(filter); err == nil {
				items.NItems = n
			}
		},
	})
	return w
}

func main() {
	flag.Parse()
	widget := makeWidget()

	tty := cli.NewTTY(os.Stdin, os.Stderr)
	restore, err := tty.Setup()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer restore()
	defer tty.StopInput()
	for {
		h, w := tty.Size()
		if h > *maxHeight {
			h = *maxHeight
		}
		tty.UpdateBuffer(nil, widget.Render(w, h), false)
		event, err := tty.ReadEvent()
		if err != nil {
			errBuf := term.NewBufferBuilder(w).Write(err.Error(), ui.FgRed).Buffer()
			tty.UpdateBuffer(nil, errBuf, true)
			break
		}
		handled := widget.Handle(event)
		if !handled && event == term.K('D', ui.Ctrl) {
			tty.UpdateBuffer(nil, term.NewBufferBuilder(w).Buffer(), true)
			break
		}
	}
}
