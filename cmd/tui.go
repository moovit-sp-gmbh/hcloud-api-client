package cmd

import (
	_ "embed"
	"hcloud-api-client/config"
	"hcloud-api-client/pkg"

	"github.com/gdamore/tcell/v2"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/high5"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/idp"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:    "tui",
	Short:  "run terimal user interface",
	Run:    initTui,
	Hidden: true,
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}

func initTui(cmd *cobra.Command, args []string) {
	format = "yaml"
	app := tview.NewApplication()
	var tree *tview.TreeView
	root := tview.NewTreeNode("helmut.cloud").SetSelectable(false).SetColor(tcell.ColorGreen)
	service := tview.NewTreeNode("service").
		SetSelectable(true).
		SetReference("service").
		SetColor(tcell.ColorWhite).
		AddChild(tview.NewTreeNode("idp").AddChild(tview.NewTreeNode("authorize").SetColor(tcell.ColorBlue).SetSelectedFunc(func() {
			ctx := config.Config.GetActiveContext()
			idp := idp.New(hcloud.New(&hcloud.ClientConfig{Api: ctx.Server, Token: ctx.Token}))
			user, err := idp.Authorize()
			if err != nil {
				pkg.PrintErr(err)
			}
			app.Stop()
			pkg.Print(user)
		}))).
		AddChild(tview.NewTreeNode("high5").AddChild(tview.NewTreeNode("app").AddChild(tview.NewTreeNode("list").SetColor(tcell.ColorBlue).SetSelectedFunc(func() {
			ctx := config.Config.GetActiveContext()
			high5 := high5.New(hcloud.New(&hcloud.ClientConfig{Api: ctx.Server, Token: ctx.Token}))

			apps, err := high5.GetApps(1000, 0)
			if err != nil {
				pkg.PrintErr(err)
			}
			app.Stop()
			pkg.Print(apps)
		})).AddChild(tview.NewTreeNode("create").SetColor(tcell.ColorBlue).SetSelectedFunc(func() {
			form := tview.NewForm()
			input := tview.NewInputField()
			form.AddFormItem(input)
			form.AddButton("Save", func() {
				ctx := config.Config.GetActiveContext()
				high5 := high5.New(hcloud.New(&hcloud.ClientConfig{Api: ctx.Server, Token: ctx.Token}))

				apps, err := high5.CreateApp(input.GetText())
				if err != nil {
					app.Stop()
					pkg.PrintErr(err)
				}
				app.Stop()
				pkg.Print(apps)
			}).
				AddButton("Cancel", func() {
					app.SetRoot(tree, true)
				})
			form.SetBorder(true).SetTitle("Enter some data").SetTitleAlign(tview.AlignLeft)
			modal := func(p tview.Primitive, width, height int) tview.Primitive {
				return tview.NewGrid().
					SetColumns(0, width, 0).
					SetRows(0, height, 0).
					AddItem(p, 1, 1, 1, 1, 0, 0, true)
			}
			app.SetRoot(modal(form, 100, 100), true)
		}))))
	version := tview.NewTreeNode("version").
		SetSelectable(true).
		SetReference("version").
		SetColor(tcell.ColorWhite).
		AddChild(tview.NewTreeNode("local").SetColor(tcell.ColorBlue).SetSelectedFunc(func() {
			app.Stop()
			localVersion(nil, nil)
		})).
		AddChild(tview.NewTreeNode("remote").SetColor(tcell.ColorBlue).SetSelectedFunc(func() {
			app.Stop()
			remoteVersion(nil, nil)
		}))

	tree = tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(service)
	tree.GetRoot().AddChild(service)
	tree.GetRoot().AddChild(version)

	if err := app.SetRoot(tree, true).Run(); err != nil {
		panic(err)
	}
}
