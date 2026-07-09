package repl

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command"
)

type result struct {
	command string
	output  string
}

type model struct {
	root command.Executor
	init command.Executor

	commands []result
	command  []rune

	idx int

	helpMessage string
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
				m.command = []rune(m.commands[len(m.commands)-m.idx].command)
			}
		case "down":
			m.idx = max(0, m.idx-1)
			if m.idx <= 0 {
				m.command = nil
			} else {
				m.command = []rune(m.commands[len(m.commands)-m.idx].command)
			}
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
			if len(m.command) > 0 {
				m.command = m.command[:len(m.command)-1]
			}
		default:
			m.helpMessage = ""
			if msg.Key().Mod == tea.ModShift {
				m.command = append(m.command, msg.ShiftedCode)
			} else {
				m.command = append(m.command, msg.Code)
			}
		}
	case errMsg:
		m.commands = append(m.commands, result{
			command: string(m.command),
			output:  msg.err.Error(),
		})
		m.command = nil
	case outputMsg:
		out := ""
		if msg.output != nil {
			out = msg.output.String()
		}
		m.commands = append(m.commands, result{
			command: string(m.command),
			output:  out,
		})
		m.command = nil
	}
	return m, nil
}

// View implements [tea.Model].
func (m *model) View() tea.View {
	var buf bytes.Buffer

	for _, command := range m.commands {
		fmt.Fprintf(&buf, "> %s\n%s\n", command.command, command.output)
	}

	fmt.Fprintf(&buf, "> %s", string(m.command))

	if m.helpMessage != "" {
		fmt.Fprintf(&buf, "\n%s", m.helpMessage)
	}

	return tea.NewView(buf.String())
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
