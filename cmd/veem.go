package main

import (
	"log"

	"github.com/gdamore/tcell/v2"
)

type Mode string

var (
	NORMAL Mode = "NORMAL"
	INSERT Mode = "INSERT"
)

type Cursor struct {
	x int
	y int
}

type Veem struct {
	mode   Mode
	cursor Cursor
}

func NewVeem() *Veem {
	return &Veem{
		mode:   NORMAL,
		cursor: Cursor{0, 0},
	}
}

func (v *Veem) GetCursor() (int, int) {
	return v.cursor.x, v.cursor.y
}

func (v *Veem) SetCursor(x int, y int) {
	newCursorPos := Cursor{x, y}
	v.cursor = newCursorPos
}

func main() {
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)

	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	s.EnableMouse()
	s.EnablePaste()
	s.Clear()

	ve := NewVeem()

	quit := func() {
		maybePanic := recover()
		s.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

	for {
		curX, curY := ve.GetCursor()
		s.Show()
		s.ShowCursor(curX, curY)

		ev := s.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				quit()
			} else if ev.Key() == tcell.KeyEnter {
				ve.SetCursor(0, curY+1)
			} else if ev.Key() == tcell.KeyBackspace {
				ve.SetCursor(curX-1, curY)
				s.SetContent(curX-1, curY, ' ', nil, defStyle)
			} else {
				s.SetContent(curX, curY, ev.Rune(), nil, defStyle)
				ve.SetCursor(curX+1, curY)
			}
		}
	}
}
