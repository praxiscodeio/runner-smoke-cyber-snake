package main

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

type Point struct {
	x, y int
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

	// Голова змейки
	head := g.snake[0]
	termbox.SetCell(head.x, head.y, '■', termbox.ColorLightMagenta, termbox.ColorDefault)

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
	game.draw()

	for {
		ev := termbox.PollEvent()
		if ev.Type == termbox.EventKey && ev.Key == termbox.KeyEsc {
			return
		}
	}
}
