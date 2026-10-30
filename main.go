package main

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

type Point struct {
	x, y int
}

func (p Point) ToRune() rune {
	switch {
	case p.x == 0 && p.y == -1:
		return '▲'
	case p.x == 0 && p.y == 1:
		return '▼'
	case p.x == -1 && p.y == 0:
		return '◀'
	case p.x == 1 && p.y == 0:
		return '▶'
	default:
		return '●'
	}
}

type Game struct {
	snake    []Point
	food     Point
	malware  []Point
	dir      Point
	score    int
	level    int
	gameOver bool
	width    int
	height   int
	quit     chan struct{}
}

func NewGame(width, height int) *Game {
	g := &Game{
		snake:  []Point{{width / 2, height / 2}},
		dir:    Point{1, 0},
		level:  1,
		width:  width,
		height: height,
		quit:   make(chan struct{}),
	}
	return g
}

func (g *Game) handleInput(ev termbox.Event) {
	if ev.Type != termbox.EventKey {
		return
	}

	newDir := g.dir
	switch ev.Key {
	case termbox.KeyArrowUp:
		newDir = Point{0, -1}
	case termbox.KeyArrowDown:
		newDir = Point{0, 1}
	case termbox.KeyArrowLeft:
		newDir = Point{-1, 0}
	case termbox.KeyArrowRight:
		newDir = Point{1, 0}
	case termbox.KeyEsc:
		g.stop()
		return
	}
	switch ev.Ch {
	case 'w', 'W':
		newDir = Point{0, -1}
	case 's', 'S':
		newDir = Point{0, 1}
	case 'a', 'A':
		newDir = Point{-1, 0}
	case 'd', 'D':
		newDir = Point{1, 0}
	case 'q', 'Q':
		g.stop()
		return
	}

	// Запрет разворота на 180°
	if newDir.x != -g.dir.x || newDir.y != -g.dir.y {
		g.dir = newDir
	}
}

// stop закрывает канал quit один раз: повторное нажатие Q или Esc не паникует.
func (g *Game) stop() {
	select {
	case <-g.quit:
	default:
		close(g.quit)
	}
}

func (g *Game) draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	// Границы
	for x := 0; x < g.width; x++ {
		termbox.SetCell(x, 0, '─', termbox.ColorWhite, termbox.ColorDefault)
		termbox.SetCell(x, g.height-1, '─', termbox.ColorWhite, termbox.ColorDefault)
	}
	for y := 0; y < g.height; y++ {
		termbox.SetCell(0, y, '│', termbox.ColorWhite, termbox.ColorDefault)
		termbox.SetCell(g.width-1, y, '│', termbox.ColorWhite, termbox.ColorDefault)
	}
	termbox.SetCell(0, 0, '┌', termbox.ColorWhite, termbox.ColorDefault)
	termbox.SetCell(g.width-1, 0, '┐', termbox.ColorWhite, termbox.ColorDefault)
	termbox.SetCell(0, g.height-1, '└', termbox.ColorWhite, termbox.ColorDefault)
	termbox.SetCell(g.width-1, g.height-1, '┘', termbox.ColorWhite, termbox.ColorDefault)

	// Змейка
	for i, p := range g.snake {
		if i == 0 {
			termbox.SetCell(p.x, p.y, g.dir.ToRune(), termbox.ColorLightMagenta, termbox.ColorDefault)
		} else {
			termbox.SetCell(p.x, p.y, '○', termbox.ColorLightMagenta, termbox.ColorDefault)
		}
	}

	// Счёт и уровень
	scoreStr := fmt.Sprintf("Score: %d  Level: %d", g.score, g.level)
	for i, ch := range scoreStr {
		termbox.SetCell(i+1, 0, ch, termbox.ColorWhite, termbox.ColorDefault)
	}

	termbox.Flush()
}

func main() {
	err := termbox.Init()
	if err != nil {
		fmt.Printf("Ошибка инициализации termbox: %v\n", err)
		return
	}
	defer termbox.Close()

	game := NewGame(40, 20)

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	game.draw()
	for {
		select {
		case ev := <-eventQueue:
			game.handleInput(ev)
			game.draw()
		case <-game.quit:
			return
		}
	}
}
