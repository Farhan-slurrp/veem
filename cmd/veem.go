package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"unicode"

	"github.com/Farhan-slurrp/veem/internal/globals"
	"github.com/Farhan-slurrp/veem/internal/screen"
	"github.com/gdamore/tcell/v2"
)

type Mode string

const (
	NORMAL Mode = "NORMAL"
	INSERT Mode = "INSERT"
)

type Cursor struct {
	x int
	y int
}

type Veem struct {
	file   *os.File
	mode   Mode
	cursor Cursor
	screen screen.Screen
}

func NewVeem(filename string) *Veem {
	var file *os.File
	s := screen.NewScreen()
	var initialContent []string
	if filename != "" {
		file, err := os.Open(filename)
		fmt.Println(file)
		if err != nil {
			panic(fmt.Sprintf("failed to read file %v", filename))
		}

		rd := bufio.NewReader(file)
		for {
			line, err := rd.ReadString('\n')
			if err != nil {
				if err == io.EOF && line == "" {
					break
				} else if err != io.EOF {
					log.Fatalf("read file line error: %v", err)
				}
			}
			initialContent = append(initialContent, line)
		}
	}
	s.InitScreen(initialContent)

	v := Veem{
		file:   file,
		mode:   NORMAL,
		cursor: Cursor{s.StartIdx, 0},
		screen: *s,
	}
	v.displayMode()
	return &v
}

func (v *Veem) GetCursor() (int, int) {
	return v.cursor.x, v.cursor.y
}

func (v *Veem) SetCursor(x int, y int) {
	newCursorPos := Cursor{x, y}
	v.cursor = newCursorPos
}

func (v *Veem) Stream() {
	quit := func() {
		maybePanic := recover()
		v.screen.Current.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

	for {
		curX, curY := v.GetCursor()
		v.screen.Current.Show()
		v.screen.Current.ShowCursor(curX, curY)

		ev := v.screen.Current.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			v.screen.Current.Sync()
		case *tcell.EventKey:
			if v.mode == NORMAL {
				v.handleNormalMode(ev, quit)
			} else if v.mode == INSERT {
				v.handleInsertMode(ev)
			}

		}
	}
}

func (v *Veem) changeMode(mode Mode) {
	v.mode = mode
	v.displayMode()
}

func (v *Veem) displayMode() {
	_, height := v.screen.Current.Size()

	for index, char := range v.mode {
		v.screen.Current.SetContent(index, height-1, char, nil, globals.DefStyle)
	}
}

func (v *Veem) handleNormalMode(ev *tcell.EventKey, quit func()) {
	curX, curY := v.GetCursor()

	if ev.Key() == tcell.KeyCtrlC {
		v.file.Close()
		quit()
	} else if ev.Rune() == rune('i') || ev.Rune() == rune('I') {
		v.changeMode(INSERT)
	} else if ev.Rune() == rune('j') || ev.Rune() == rune('J') {
		v.SetCursor(curX, curY+1)
	} else if ev.Rune() == rune('k') || ev.Rune() == rune('K') {
		v.SetCursor(curX, curY-1)
	} else if ev.Rune() == rune('h') || ev.Rune() == rune('H') {
		v.SetCursor(curX-1, curY)
	} else if ev.Rune() == rune('l') || ev.Rune() == rune('L') {
		v.SetCursor(curX+1, curY)
	}
}

func (v *Veem) handleInsertMode(ev *tcell.EventKey) {
	curX, curY := v.GetCursor()
	width, height := v.screen.Current.Size()

	if ev.Key() == tcell.KeyEscape {
		v.changeMode(NORMAL)
	} else if ev.Key() == tcell.KeyEnter {
		if curY+1 > height {
			v.screen.Current.SetSize(width, height+1)
			v.screen.Current.Sync()
		}
		v.SetCursor(2, curY+1)
	} else if ev.Key() == tcell.KeyBackspace {
		if curX-1 >= 2 {
			v.SetCursor(curX-1, curY)
			v.screen.ShiftContentLeft(curX, curY)
		} else if curY-1 >= 0 {
			lastContentIdx := 2
			for i := 2; i < width; i++ {
				currRune, _, _, _ := v.screen.Current.GetContent(i, curY-1)
				if currRune != ' ' {
					lastContentIdx = i
				}
			}
			v.SetCursor(lastContentIdx, curY-1)
			v.screen.Current.SetContent(lastContentIdx, curY-1, ' ', nil, globals.DefStyle)
		}
	} else if unicode.IsSpace(ev.Rune()) {
		v.SetCursor(curX+1, curY)
		v.screen.ShiftContentRight(curX, curY)
	} else {
		if curX+1 > width {
			v.screen.Current.SetSize(width+1, height)
			v.screen.Current.Sync()
		}
		v.screen.Current.SetContent(curX, curY, ev.Rune(), nil, globals.DefStyle)
		v.SetCursor(curX+1, curY)
	}
}

func main() {
	filename := ""
	if len(os.Args) >= 2 {
		filename = os.Args[1]
	}
	v := NewVeem(filename)
	v.Stream()
}
