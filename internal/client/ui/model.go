package ui

import (
	"encoding/binary"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type State uint8

const (
	AddrChoice State = iota
	Play
	Spectator
)

type Model struct {
	st        State
	addrInput textinput.Model
	addr      string
	cm        *ConnectionManager
	Plot      []byte
	err       error
}

func NewModel() *Model {
	ti := textinput.New()
	ti.Placeholder = "host:port"
	ti.Focus()

	return &Model{
		addrInput: ti,
		st:        AddrChoice,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

// isAlive проверяет, есть ли на поле хоть одна клетка с нашим playerID.
func isAlive(plot []byte, playerID int) bool {
	const cellSize = 4
	for i := 0; i+cellSize <= len(plot); i += cellSize {
		if int(binary.LittleEndian.Uint16(plot[i:i+2])) == playerID {
			return true
		}
	}
	return false
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case FrameMsg:
		m.Plot = []byte(msg)
		if m.st == Play && !isAlive(m.Plot, m.cm.playerID) {
			m.st = Spectator
		}
		return m, m.cm.Listen()

	case WsErrMsg:
		m.err = msg.Err
		m.st = AddrChoice
		m.cm = nil
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			if m.cm != nil {
				m.cm.Close()
			}
			return m, tea.Quit

		case "enter":
			if m.st == AddrChoice {
				m.addr = m.addrInput.Value()
				cm, err := NewConnectionManager(m.addr)
				if err != nil {
					m.err = err
					return m, nil
				}
				m.cm = cm
				m.err = nil
				m.st = Play
				return m, m.cm.Listen()
			}

		case "r":
			// Возрождение: закрываем старое соединение и переходим к вводу адреса.
			if m.st == Spectator {
				m.cm.Close()
				m.cm = nil
				m.Plot = nil
				m.err = nil
				m.st = AddrChoice
				m.addrInput.SetValue(m.addr) // оставляем последний адрес для удобства
				m.addrInput.Focus()
				return m, nil
			}

		default:
			if m.st == Play {
				dir := MapMove(msg.String())
				if dir != 0xFF {
					return m, m.cm.SendMove(dir)
				}
			}
		}
	}

	if m.st == AddrChoice {
		m.addrInput, cmd = m.addrInput.Update(msg)
	}

	return m, cmd
}

func (m *Model) View() tea.View {
	var content string

	switch m.st {
	case AddrChoice:
		content = lipgloss.JoinVertical(
			lipgloss.Top,
			m.headerView(),
			m.addrInput.View(),
			m.footerView(),
		)
	case Play:
		content = RenderPlot(m.Plot, m.cm.playerID, m.cm.mapSize)
	case Spectator:
		content = lipgloss.JoinVertical(
			lipgloss.Top,
			RenderPlot(m.Plot, m.cm.playerID, m.cm.mapSize),
			"\nYou died. Press [r] to reconnect or [ctrl+c] to quit.",
		)
	}

	if m.err != nil {
		content += "\n\nError: " + m.err.Error()
	}

	v := tea.NewView(content)

	if m.st == AddrChoice && !m.addrInput.VirtualCursor() {
		c := m.addrInput.Cursor()
		c.Y += lipgloss.Height(m.headerView())
		v.Cursor = c
	}

	return v
}

func (m *Model) headerView() string { return "Enter server address:\n" }
func (m *Model) footerView() string { return "\n(ctrl+c to quit)" }
