package game

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/THPTUHA/repeatword/audio"
	"github.com/THPTUHA/repeatword/config"
	"github.com/THPTUHA/repeatword/db"
	"github.com/THPTUHA/repeatword/logger"
	"github.com/THPTUHA/repeatword/vocab"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/timer"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sirupsen/logrus"
)

var (
	correctStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#008000"))
	wrongStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))

	highLighColor = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00"))
	normalColor   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
)

const (
	START_STATUS       = 0
	PLAYING_STATUS     = 1
	FINISH_STATUS      = 2
	DICT_STATUS        = 3
	VIEW_ANSWER_STATUS = 4
)

const (
	PLAY = iota
	PLAY_AGAIN
	PLAY_AGAIN_WORD_WRONG
)

const (
	VOB_NO_ANSWER = iota
	VOB_ANSWER_CORRECT
	VOB_ANSWER_WRONG
)

const (
	NORMAL_MODE = iota
	INFINITY_MODE
)

type Game struct {
	Vobs     []*vocab.Vocabulary `json:"vobs"`
	Mode     uint                `json:"mode"`
	BeginAt  int                 `json:"begin_at"`
	FinishAt int                 `json:"finish_at"`

	status     int
	currentIdx int
	vobDict    *vocab.Vocabulary

	readyView  bool
	viewport   viewport.Model
	input      textinput.Model
	timer      timer.Model
	resultView table.Model
	keymap     keymap
	help       help.Model

	debug *os.File

	stopCh  chan int
	lock    sync.RWMutex
	err     error
	queries *db.Queries
	logger  *logrus.Entry
}

type keymap struct {
	status        key.Binding
	playAgain     key.Binding
	playWordWrong key.Binding

	showDict key.Binding
}

type Config struct {
	Root *config.Configs

	CollectionID uint64
	Limit        uint64
	PlayMode     uint
	Mode         uint

	RecentDayNum int
	Logger       *logrus.Entry
	game         *Game
}

var currentConfig *Config = nil

func (g *Game) Debug(str string) {
	if g.debug != nil {
		g.debug.Write([]byte(str))
	}
}

func Init(config *Config) *Game {
	if config == nil {
		config = &Config{
			CollectionID: 1,
			Limit:        10,
			Logger:       logger.InitLogger(logrus.DebugLevel.String()),
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
	g.stopCh = make(chan int)
	d, err := db.ConnectMysql()
	if err != nil {
		g.logger.Fatal(err)
	}
	queries := db.New(d)

	g.queries = queries
	g.logger = config.Logger
	g.Mode = config.Mode

	// file, err := os.Create("debug.txt")
	// g.debug = file

	if err != nil {
		g.logger.Fatal(err)
	}

	ctx := context.Background()

	Vobs := make([]*vocab.Vocabulary, 0)
	switch config.PlayMode {
	case PLAY_AGAIN_WORD_WRONG:
		for _, v := range config.game.Vobs {
			if v.Status == VOB_ANSWER_WRONG || v.Status == VOB_NO_ANSWER {
				Vobs = append(Vobs, v)
			}
		}
	default:
		result, err := queries.GetVobsRandom(ctx, db.GetVobsRandomParams{
			Getvobsrandom:   config.CollectionID,
			Getvobsrandom_2: config.Limit,
			Getvobsrandom_3: config.RecentDayNum,
		})
		if err != nil {
			g.logger.Fatal(err)
		}
		if result == nil {
			g.logger.Fatal("no Vobs")
		}
		err = json.Unmarshal(result.([]byte), &Vobs)
		if err != nil {
			g.logger.Fatal(err)
		}
	}

	g.Vobs = Vobs

	return g
}

const timeInterval = 5 * time.Second
const timeDelayViewAnswer = 1 * time.Second
const answerPlentyExtra = 3

func (g *Game) Init() tea.Cmd {
	return nil
}

var windowSize tea.WindowSizeMsg

func (g *Game) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
		cmds  []tea.Cmd
	)

	g.input, tiCmd = g.input.Update(msg)
	g.viewport, vpCmd = g.viewport.Update(msg)
	// cmds = append(cmds, tiCmd, vpCmd)
	switch msg := msg.(type) {
	case timer.TickMsg:
		var cmd tea.Cmd
		g.timer, cmd = g.timer.Update(msg)
		if g.timer.Timedout() {
			g.Debug(g.Vobs[g.currentIdx].Word.String + "timeout\n")
			return g.nextQuiz(cmd, true)
		}
		return g, cmd

	case timer.StartStopMsg:
		var cmd tea.Cmd
		g.timer, cmd = g.timer.Update(msg)
		return g, cmd
	case tea.WindowSizeMsg:
		windowSize = msg
		if !g.readyView {
			if g.status == DICT_STATUS {
				g.viewport = viewport.New(msg.Width, msg.Height-50)
				g.viewport.YPosition = 100
				g.readyView = true
			}

		} else {
			// g.viewport.Width = msg.Width
			g.viewport.Height = msg.Height - 10
			g.viewport.Width = msg.Width
			cmds = append(cmds, viewport.Sync(g.viewport))
		}

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlS:
			if g.status == START_STATUS {
				g.status = PLAYING_STATUS
				g.viewport.SetContent("Type a answer and press Enter to check.")
				go g.playAudio()
				return g, g.timer.Init()
			} else if g.status == FINISH_STATUS {
				currentConfig.PlayMode = PLAY_AGAIN
				g = Init(currentConfig)
				g.status = PLAYING_STATUS
				go g.playAudio()
				return g, g.timer.Init()
			}
		case tea.KeyCtrlD:
			if g.status != PLAYING_STATUS {
				g.status = DICT_STATUS
				g.input.Placeholder = "Search word..."
				g.viewport.SetContent("")
			}
		case tea.KeyCtrlW:
			if g.status == FINISH_STATUS {
				for _, v := range g.Vobs {
					if v.Status == VOB_NO_ANSWER || v.Status == VOB_ANSWER_WRONG {
						currentConfig.PlayMode = PLAY_AGAIN_WORD_WRONG
						currentConfig.game = g
						g = Init(currentConfig)
						g.status = PLAYING_STATUS
						go g.playAudio()
						return g, g.timer.Init()
					}
				}
			}
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(g.input.Value())
			return g, tea.Quit
		case tea.KeyEnter:
			if g.status == PLAYING_STATUS {
				vob := g.Vobs[g.currentIdx]
				var checkAnswer string
				if vob.Word.String == g.input.Value() {
					vob.Status = VOB_ANSWER_CORRECT
					checkAnswer = correctStyle.Render(g.input.Value())
				} else {
					vob.Status = VOB_ANSWER_WRONG
					checkAnswer = wrongStyle.Render(g.input.Value())
				}
				g.viewport.SetContent(checkAnswer)

				if vob.Status == VOB_ANSWER_CORRECT {
					vob.AnswerNum++
					vob.Remand--
					g.Debug(g.Vobs[g.currentIdx].Word.String + "next quiz\n")
					return g.nextQuiz(tea.Batch(tiCmd, vpCmd), false)
				}
			} else if g.status == DICT_STATUS && g.input.Value() != "" {

				result, err := g.queries.GetWordDict(context.Background(), g.input.Value())
				if err != nil {
					log.Fatalln(err)
				}
				if result == nil {
					g.viewport.SetContent(fmt.Sprintf("not found '%s'", g.input.Value()))
				} else {
					var vob vocab.Vocabulary
					err = json.Unmarshal(result.([]byte), &vob)
					if err != nil {
						g.logger.Fatal(err)
					}
					g.vobDict = &vob
					str := vob.Word.String
					for _, p := range vob.Parts {
						str += "\n" + p.Title.String + "\n" + p.Type.String + "\n"
						for _, pn := range p.Pronounces {
							str += highLighColor.Render(pn.Region.String) + " " + normalColor.Render(pn.Pro.String) + "\t"
						}
						for _, g := range p.Means {
							str += "\n----------------\n"
							str += highLighColor.Render(g.Meaning.String) + "\n"
							for idx, e := range g.Examples {
								str += fmt.Sprint(idx+1) + "." + e.Example.String + "\n"
							}
						}
					}
					g.viewport.Width = windowSize.Width
					g.viewport.Height = windowSize.Height - 10
					g.viewport.MouseWheelEnabled = true
					g.viewport.SetContent(str)
				}

			}
			g.input.Reset()
			g.viewport.GotoBottom()
		}

	case error:
		g.err = msg
		return g, nil
	}
	cmds = append(cmds, tiCmd, vpCmd)
	return g, tea.Batch(cmds...)
}

func (g *Game) View() string {
	if g.status == VIEW_ANSWER_STATUS {
		g.viewport.SetContent(highLighColor.Render(g.Vobs[g.currentIdx].Word.String))
		return fmt.Sprintf(
			"%s\n\n%s\n",
			g.timer.View(),
			g.viewport.View(),
		) + "\n\n"
	}

	if g.status == FINISH_STATUS {
		keys := make([]key.Binding, 0)
		keys = append(keys, g.keymap.playAgain, g.keymap.showDict)
		for _, v := range g.Vobs {
			if v.Status == VOB_ANSWER_WRONG || v.Status == VOB_NO_ANSWER {
				keys = append(keys, g.keymap.playWordWrong)
				break
			}
		}
		return fmt.Sprintf("%s\n%s", g.resultView.View(), g.help.ShortHelpView(keys))
	}

	if g.status == PLAYING_STATUS {
		return fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			g.timer.View(),
			g.viewport.View(),
			g.input.View(),
		) + "\n\n"
	}

	if g.status == DICT_STATUS {
		return fmt.Sprintf("%s\n%s\n%v", g.input.View(), g.viewport.View(), g.help.ShortHelpView([]key.Binding{
			g.keymap.status,
			g.keymap.showDict,
		}))
	}

	return fmt.Sprintf("%s\n%v", g.viewport.View(), g.help.ShortHelpView([]key.Binding{
		g.keymap.status,
		g.keymap.showDict,
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
		BeginAt: int(time.Now().UnixMilli() / 1000),
		lock:    sync.RWMutex{},
		keymap: keymap{
			status: key.NewBinding(
				key.WithKeys("ctrl+s"),
				key.WithHelp("ctrl+s", "to start"),
			),
			playAgain: key.NewBinding(
				key.WithKeys("ctrl+s"),
				key.WithHelp("ctrl+s", "play again"),
			),
			playWordWrong: key.NewBinding(
				key.WithKeys("ctrl+w"),
				key.WithHelp("ctrl+w", "play again word wrong"),
			),
			showDict: key.NewBinding(
				key.WithKeys("ctrl+d"),
				key.WithHelp("ctrl+d", "show dict word"),
			),
		},
		help:     help.New(),
		timer:    timer.NewWithInterval(timeInterval, time.Millisecond),
		input:    ti,
		viewport: vp,
		err:      nil,
	}
}

func (g *Game) saveRecord() error {
	js, err := json.Marshal(g)
	if err != nil {
		return err
	}
	return g.queries.SaveRecord(context.Background(), string(js))
}

func (g *Game) nextQuiz(cmd tea.Cmd, timeout bool) (tea.Model, tea.Cmd) {
	if g.timer.Running() {
		g.timer.Stop()
	}

	if g.currentIdx >= len(g.Vobs)-1 && g.Mode == NORMAL_MODE {
		g.lock.Lock()
		if g.status != FINISH_STATUS {
			g.status = FINISH_STATUS
			g.FinishAt = int(time.Now().UnixMilli() / 1000)
			if err := g.saveRecord(); err != nil {
				log.Fatalln(err)
			}
		}
		g.lock.Unlock()

		columns := []table.Column{
			{Title: "Number", Width: 10},
			{Title: "Word", Width: 30},
			{Title: "Answer", Width: 10},
		}

		rows := []table.Row{}

		var num, status, wordStatus string
		for idx, vob := range g.Vobs {
			switch vob.Status {
			case VOB_ANSWER_WRONG:
				num = wrongStyle.Render(fmt.Sprint(idx + 1))
				status = wrongStyle.Render("false")
				wordStatus = wrongStyle.Render(fmt.Sprint(vob.Word.String))
			case VOB_NO_ANSWER:
				num = highLighColor.Render(fmt.Sprint(idx + 1))
				status = highLighColor.Render("empty")
				wordStatus = highLighColor.Render(fmt.Sprint(vob.Word.String))
			case VOB_ANSWER_CORRECT:
				num = correctStyle.Render(fmt.Sprint(idx + 1))
				status = correctStyle.Render("true")
				wordStatus = correctStyle.Render(fmt.Sprint(vob.Word.String))
			}

			rows = append(rows, table.Row{num, wordStatus, status})
		}
		g.resultView = table.New(
			table.WithColumns(columns),
			table.WithRows(rows),
			table.WithFocused(false),
			table.WithHeight(len(rows)+1))
		return g, cmd
	}

	if timeout && g.Mode == INFINITY_MODE {
		if g.status == PLAYING_STATUS {
			v := g.Vobs[g.currentIdx]
			if v.Status == VOB_ANSWER_WRONG || v.Status == VOB_NO_ANSWER {
				v.Status = VOB_NO_ANSWER
				v.Remand += answerPlentyExtra
			}
			g.Debug(g.Vobs[g.currentIdx].Word.String + " show_answer\n")
			g.status = VIEW_ANSWER_STATUS
			g.timer = timer.NewWithInterval(timeDelayViewAnswer, time.Millisecond)
			return g, g.timer.Init()
		}
	}

	if g.Mode == INFINITY_MODE {
		conti := false
		for _, v := range g.Vobs {
			g.Debug(fmt.Sprintf("%s %d\n", v.Word.String, v.Remand))
		}
		for i := 1; i <= len(g.Vobs); i++ {
			idx := (g.currentIdx + i) % len(g.Vobs)
			v := g.Vobs[idx]
			if v.Remand >= 0 {
				g.currentIdx = idx
				conti = true
				break
			}
		}
		if !conti {
			g.lock.Lock()
			if g.status != FINISH_STATUS {
				g.status = FINISH_STATUS
				g.FinishAt = int(time.Now().UnixMilli() / 1000)
				if err := g.saveRecord(); err != nil {
					log.Fatalln(err)
				}
			}
			g.lock.Unlock()

			columns := []table.Column{
				{Title: "Number", Width: 10},
				{Title: "Word", Width: 30},
				{Title: "AnswerNumber", Width: 10},
			}

			rows := []table.Row{}

			var num, status, wordStatus string
			for idx, vob := range g.Vobs {
				switch vob.Status {
				case VOB_ANSWER_WRONG:
					num = wrongStyle.Render(fmt.Sprint(idx + 1))
					status = wrongStyle.Render(fmt.Sprint(vob.AnswerNum))
					wordStatus = wrongStyle.Render(fmt.Sprint(vob.Word.String))
				case VOB_NO_ANSWER:
					num = highLighColor.Render(fmt.Sprint(idx + 1))
					status = highLighColor.Render(fmt.Sprint(vob.AnswerNum))
					wordStatus = highLighColor.Render(fmt.Sprint(vob.Word.String))
				case VOB_ANSWER_CORRECT:
					num = correctStyle.Render(fmt.Sprint(idx + 1))
					status = correctStyle.Render(fmt.Sprint(vob.AnswerNum))
					wordStatus = correctStyle.Render(fmt.Sprint(vob.Word.String))
				}

				rows = append(rows, table.Row{num, wordStatus, status})
			}

			g.resultView = table.New(
				table.WithColumns(columns),
				table.WithRows(rows),
				table.WithFocused(false),
				table.WithHeight(len(rows)+1))

			return g, cmd
		}
	} else {
		g.currentIdx++
	}

	g.status = PLAYING_STATUS
	g.timer = timer.NewWithInterval(timeInterval, time.Millisecond)
	g.viewport.SetContent(fmt.Sprintf("Type a answer and press Enter to check. remand %d", g.Vobs[g.currentIdx].Remand))
	g.input.SetValue("")
	return g, g.timer.Init()
}

func (g *Game) playAudio() {
	for {
		select {
		case status := <-g.stopCh:
			if status == FINISH_STATUS {
				return
			}
		default:
			if g.status == PLAYING_STATUS {
				vob := g.Vobs[g.currentIdx]
				idx := randomInt(len(vob.Parts[0].Pronounces) - 1)
				audio.PlayAudio(currentConfig.Root.DataDir, vob.Parts[0].Pronounces[idx].LocalFile.String)
			} else if g.status == FINISH_STATUS {
				return
			}
		}
	}
}

func randomInt(n int) int {
	rand.Seed(time.Now().UnixNano())

	return rand.Intn(n) + 1
}

func (g *Game) Play() {
	p := tea.NewProgram(g, tea.WithAltScreen(), tea.WithMouseCellMotion())

	if _, err := p.Run(); err != nil {
		g.logger.Fatal(err)
	}
}
