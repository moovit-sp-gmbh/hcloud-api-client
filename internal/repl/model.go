package repl

import (
	"bytes"
	"fmt"
	"os"
	"slices"
	"strings"
	"unicode"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command"
)

// prompt precedes the editable line. The cursor is placed relative to it, so
// its width has to stay in sync with what [model.View] renders.
const prompt = "> "

type result struct {
	command string
	output  string
}

type model struct {
	root command.Executor
	init command.Executor

	commands []result
	command  []rune

	// pos is where the next rune is inserted, as a rune offset into command.
	// It is always in [0, len(command)]; len(command) means end of line.
	pos int

	idx int

	helpMessage string
}

// setCommand replaces the editable line and parks the cursor at its end.
func (m *model) setCommand(cmd []rune) {
	m.command = cmd
	m.pos = len(cmd)
}

// wordLeft returns the offset of the start of the word preceding the cursor,
// skipping any run of whitespace in between. It returns 0 at the start of the
// line.
func (m *model) wordLeft() int {
	i := m.pos
	for i > 0 && unicode.IsSpace(m.command[i-1]) {
		i--
	}
	for i > 0 && !unicode.IsSpace(m.command[i-1]) {
		i--
	}
	return i
}

// wordRight returns the offset just past the end of the word following the
// cursor, skipping any run of whitespace in between. It returns len(command)
// at the end of the line.
func (m *model) wordRight() int {
	i := m.pos
	for i < len(m.command) && unicode.IsSpace(m.command[i]) {
		i++
	}
	for i < len(m.command) && !unicode.IsSpace(m.command[i]) {
		i++
	}
	return i
}

func New(root, init command.Executor) tea.Model {
	return &model{
		root:     root,
		init:     init,
		commands: []result{},
	}
}

type errMsg struct{ err error }
type outputMsg struct{ output fmt.Stringer }

// Init implements [tea.Model].
func (m *model) Init() tea.Cmd {
	if m.init == nil {
		return nil
	}

	return func() tea.Msg {
		output, err := m.init.Execute(os.Args)
		if err != nil {
			return errMsg{err}
		}
		return outputMsg{output}
	}
}

// Update implements [tea.Model].
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "up":
			m.idx = min(len(m.commands), m.idx+1)
			if m.idx > 0 {
				m.setCommand([]rune(m.commands[len(m.commands)-m.idx].command))
			}
		case "down":
			m.idx = max(0, m.idx-1)
			if m.idx <= 0 {
				m.setCommand(nil)
			} else {
				m.setCommand([]rune(m.commands[len(m.commands)-m.idx].command))
			}
		case "left":
			m.pos = max(0, m.pos-1)
		case "right":
			m.pos = min(len(m.command), m.pos+1)
		case "ctrl+left":
			m.pos = m.wordLeft()
		case "ctrl+right":
			m.pos = m.wordRight()
		case "ctrl+c", "ctrl+d":
			return m, tea.Quit
		case "tab":
			m.helpMessage = strings.Join(m.help(), " ")
		case "enter":
			m.idx = 0
			m.helpMessage = ""
			return m, m.runCommand()
		case "backspace":
			m.helpMessage = ""
			if m.pos > 0 {
				m.command = slices.Delete(m.command, m.pos-1, m.pos)
				m.pos--
			}
		default:
			m.helpMessage = ""
			code := msg.Code
			if msg.Key().Mod == tea.ModShift {
				code = msg.ShiftedCode
			}
			m.command = slices.Insert(m.command, m.pos, code)
			m.pos++
		}
	case errMsg:
		m.commands = append(m.commands, result{
			command: string(m.command),
			output:  msg.err.Error(),
		})
		m.setCommand(nil)
	case outputMsg:
		out := ""
		if msg.output != nil {
			out = msg.output.String()
		}
		m.commands = append(m.commands, result{
			command: string(m.command),
			output:  out,
		})
		m.setCommand(nil)
	}
	return m, nil
}

// View implements [tea.Model].
func (m *model) View() tea.View {
	var buf bytes.Buffer

	for _, command := range m.commands {
		fmt.Fprintf(&buf, "%s%s\n%s\n", prompt, command.command, command.output)
	}

	// Every entry above ends in a newline, so the prompt lands on the line
	// after them.
	row := bytes.Count(buf.Bytes(), []byte("\n"))
	col := ansi.StringWidth(prompt + string(m.command[:m.pos]))

	fmt.Fprintf(&buf, "%s%s", prompt, string(m.command))

	if m.helpMessage != "" {
		fmt.Fprintf(&buf, "\n%s", m.helpMessage)
	}

	view := tea.NewView(buf.String())
	view.Cursor = tea.NewCursor(col, row)
	view.Cursor.Shape = tea.CursorBar
	view.Cursor.Blink = false

	return view
}

func (m *model) runCommand() tea.Cmd {
	return func() tea.Msg {
		output, err := m.root.Execute(strings.Split(string(m.command), " "))
		if err != nil {
			return errMsg{err}
		}
		return outputMsg{output}
	}
}

func (m *model) help() []string {
	cmd, rest, err := m.root.Traverse(strings.Split(string(m.command), " "))
	if err != nil {
		return nil
	}

	return cmd.Help(rest)
}
