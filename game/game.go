package game

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/THPTUHA/repeatword/audio"
	"github.com/THPTUHA/repeatword/db"
	"github.com/THPTUHA/repeatword/vocab"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/timer"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	correctStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#008000"))
	wrongStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
)

const (
	START_STATUS   = 0
	PLAYING_STATUS = 1
	FINISH_STATUS  = 2
)

type Game struct {
	vobs       []*vocab.Vocabulary
	status     int
	currentIdx int

	viewport viewport.Model
	input    textinput.Model
	timer    timer.Model
	result   table.Model
	keymap   keymap
	help     help.Model

	stopCh chan bool
	err    error
}

type keymap struct {
	status key.Binding
}

type Config struct {
	CollectionID uint64
	Limit        uint64
}

var currentConfig *Config = nil

func Init(config *Config) *Game {
	if config == nil {
		config = &Config{
			CollectionID: 1,
			Limit:        10,
		}
	}

	if config.CollectionID == 0 {
		config.CollectionID = 1
	}

	if config.Limit == 0 {
		config.Limit = 10
	}
	currentConfig = config

	g := initialModel()
	g.stopCh = make(chan bool)
	d, err := db.ConnectMysql()
	if err != nil {
		log.Fatal(err)
	}
	queries := db.New(d)
	ctx := context.Background()

	vobs := make([]*vocab.Vocabulary, 0)
	result, err := queries.GetVobsRandom(ctx, db.GetVobsRandomParams{
		Getvobsrandom:   config.CollectionID,
		Getvobsrandom_2: config.Limit,
	})
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(result.([]byte), &vobs)
	if err != nil {
		log.Fatal(err)
	}
	g.vobs = vobs
	// for _, v := range vobs {
	// 	fmt.Println(v.String())
	// }
	return g
}

const timeInterval = 5 * time.Second

func (g *Game) Init() tea.Cmd {
	go g.playAudio()
	return textinput.Blink
}

func (g *Game) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	g.input, tiCmd = g.input.Update(msg)
	g.viewport, vpCmd = g.viewport.Update(msg)

	switch msg := msg.(type) {
	case timer.TickMsg:
		var cmd tea.Cmd
		g.timer, cmd = g.timer.Update(msg)
		if g.timer.Timedout() {
			return g.nextQuiz(cmd)
		}
		return g, cmd

	case timer.StartStopMsg:
		var cmd tea.Cmd
		g.timer, cmd = g.timer.Update(msg)
		return g, cmd

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlS:
			if g.status == START_STATUS {
				g.status = PLAYING_STATUS
				g.viewport.SetContent("Type a answer and press Enter to check.")
				return g, g.timer.Init()
			} else if g.status == FINISH_STATUS {
				g = Init(currentConfig)
				g.status = PLAYING_STATUS
				go g.playAudio()
				return g, g.timer.Init()
			}
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(g.input.Value())
			return g, tea.Quit
		case tea.KeyEnter:
			if g.status == PLAYING_STATUS {
				vob := g.vobs[g.currentIdx]
				var checkAnswer string
				if vob.Word.String == g.input.Value() {
					vob.Correct = true
					checkAnswer = correctStyle.Render(g.input.Value())
				} else {
					checkAnswer = wrongStyle.Render(g.input.Value())
				}
				g.viewport.SetContent(checkAnswer)

				if vob.Correct {
					return g.nextQuiz(tea.Batch(tiCmd, vpCmd))
				}

			}
			g.input.Reset()
			g.viewport.GotoBottom()
		}

	case error:
		g.err = msg
		return g, nil
	}

	return g, tea.Batch(tiCmd, vpCmd)
}

func (g *Game) View() string {
	if g.status == FINISH_STATUS {
		return fmt.Sprintf("%s\n%s", g.result.View(), g.help.ShortHelpView([]key.Binding{
			g.keymap.status,
		}))
	}

	if g.status == PLAYING_STATUS {
		return fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			g.timer.View(),
			g.viewport.View(),
			g.input.View(),
		) + "\n\n"
	}
	return fmt.Sprintf("%s\n%v", g.viewport.View(), g.help.ShortHelpView([]key.Binding{
		g.keymap.status,
	}))
}

func initialModel() *Game {
	ti := textinput.New()
	ti.Placeholder = "Send a answer..."
	ti.Focus()

	ti.Prompt = "┃ "
	ti.CharLimit = 280

	vp := viewport.New(100, 1)
	vp.SetContent(`Welcome to the reapeatword!`)

	return &Game{
		keymap: keymap{
			status: key.NewBinding(
				key.WithKeys("ctrl+s"),
				key.WithHelp("ctrl+s", "status"),
			),
		},
		help:     help.New(),
		timer:    timer.NewWithInterval(timeInterval, time.Millisecond),
		input:    ti,
		viewport: vp,
		err:      nil,
	}
}

func (g *Game) nextQuiz(cmd tea.Cmd) (tea.Model, tea.Cmd) {
	if g.timer.Running() {
		g.timer.Stop()
	}

	if g.currentIdx >= len(g.vobs)-1 {
		g.viewport.SetContent("Finish")
		g.status = FINISH_STATUS

		columns := []table.Column{
			{Title: "Word", Width: 30},
			{Title: "Answer", Width: 10},
		}

		rows := []table.Row{}
		for _, vob := range g.vobs {
			status := correctStyle.Render(fmt.Sprint(vob.Correct))
			wordStatus := correctStyle.Render(fmt.Sprint(vob.Word.String))
			if !vob.Correct {
				status = wrongStyle.Render(fmt.Sprint(vob.Correct))
				wordStatus = wrongStyle.Render(fmt.Sprint(vob.Word.String))
			}
			rows = append(rows, table.Row{wordStatus, status})
		}
		g.result = table.New(
			table.WithColumns(columns),
			table.WithRows(rows),
			table.WithFocused(false),
			table.WithHeight(len(rows)+1))
		return g, cmd
	}

	g.timer = timer.NewWithInterval(timeInterval, time.Millisecond)
	g.currentIdx++
	g.viewport.SetContent("Type a answer and press Enter to check.")
	g.input.SetValue("")
	return g, g.timer.Init()
}

func (g *Game) playAudio() {
	for {
		select {
		case <-g.stopCh:
			return
		default:
			if g.status == PLAYING_STATUS {
				vob := g.vobs[g.currentIdx]
				idx := randomInt(len(vob.Parts[0].Pronounces) - 1)
				audio.PlayAudio(vob.Parts[0].Pronounces[idx].LocalFile.String)
			}
		}
	}
}

func randomInt(n int) int {
	rand.Seed(time.Now().UnixNano())

	return rand.Intn(n) + 1
}

func (g *Game) Play() {
	p := tea.NewProgram(g)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

}
