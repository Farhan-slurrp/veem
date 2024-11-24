package main

import (
	"bufio"
	"os"
	"unicode"

	"github.com/Farhan-slurrp/veem/internal/constants"
	"github.com/Farhan-slurrp/veem/internal/cursor"
	"github.com/Farhan-slurrp/veem/internal/globals"
	"github.com/Farhan-slurrp/veem/internal/screen"
	"github.com/Farhan-slurrp/veem/internal/utils"
	"github.com/gdamore/tcell/v2"
)

type Veem struct {
	path   string
	mode   constants.Mode
	cursor cursor.Cursor
	screen screen.Screen
}

func NewVeem(path string) *Veem {
	var lines []string
	if path != "" {
		file, err := os.Open(path)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
	}

	s := screen.NewScreen(lines)
	s.InitScreen(constants.NORMAL)

	return &Veem{
		path:   path,
		mode:   constants.NORMAL,
		cursor: *cursor.NewCursor(s.StartXIdx, 0),
		screen: *s,
	}
}

func (v *Veem) Stream() {
	quit := func() {
		maybePanic := recover()
		v.screen.Current.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	for {
		curX, curY := v.cursor.GetCursor()
		v.screen.Current.Show()
		v.screen.Current.ShowCursor(curX, curY)

		ev := v.screen.Current.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			v.screen.Current.Sync()
			v.screen.InitScreen(v.mode)
		case *tcell.EventKey:
			if v.mode == constants.NORMAL {
				v.handleNormalMode(ev, quit)
			} else if v.mode == constants.INSERT {
				v.handleInsertMode(ev)
			}

		}
	}
}

func (v *Veem) changeMode(mode constants.Mode) {
	v.mode = mode
	v.screen.DisplayMode(mode)
}

func (v *Veem) handleNormalMode(ev *tcell.EventKey, quit func()) {
	if ev.Key() == tcell.KeyCtrlC {
		quit()
	} else if ev.Key() == tcell.KeyCtrlS {
		v.saveFile()
	} else if ev.Rune() == rune('i') || ev.Rune() == rune('I') {
		v.changeMode(constants.INSERT)
	} else if ev.Rune() == rune('j') || ev.Rune() == rune('J') {
		v.moveDown()
	} else if ev.Rune() == rune('k') || ev.Rune() == rune('K') {
		v.moveUp()
	} else if ev.Rune() == rune('h') || ev.Rune() == rune('H') {
		v.moveLeft()
	} else if ev.Rune() == rune('l') || ev.Rune() == rune('L') {
		v.moveRight()
	}
}

func (v *Veem) handleInsertMode(ev *tcell.EventKey) {
	curX, curY := v.cursor.GetCursor()
	width, _ := v.screen.Current.Size()

	if ev.Key() == tcell.KeyEscape {
		v.changeMode(constants.NORMAL)
	} else if ev.Key() == tcell.KeyEnter {
		v.cursor.SetCursor(v.screen.StartXIdx, curY+1)
	} else if ev.Key() == tcell.KeyUp {
		v.moveUp()
	} else if ev.Key() == tcell.KeyDown {
		v.moveDown()
	} else if ev.Key() == tcell.KeyRight {
		v.moveRight()
	} else if ev.Key() == tcell.KeyLeft {
		v.moveLeft()
	} else if ev.Key() == tcell.KeyBackspace {
		if curX-1 >= v.screen.StartXIdx {
			v.cursor.SetCursor(curX-1, curY)
			v.screen.ShiftContentLeft(curX, curY, ' ')
		} else if curY-1 >= 0 {
			lastContentIdx := v.screen.StartXIdx
			for i := v.screen.StartXIdx; i < width; i++ {
				currRune, _, _, _ := v.screen.Current.GetContent(i, curY-1)
				if currRune != ' ' {
					lastContentIdx = i
				}
			}
			v.cursor.SetCursor(lastContentIdx, curY-1)
			v.screen.Current.SetContent(lastContentIdx, curY-1, ' ', nil, globals.DefStyle)
		}
	} else if unicode.IsSpace(ev.Rune()) {
		v.cursor.SetCursor(curX+1, curY)
		v.screen.ShiftContentRight(curX, curY, ' ')
	} else {
		v.screen.ShiftContentRight(curX, curY, ev.Rune())
		v.cursor.SetCursor(curX+1, curY)
	}
}

func (v *Veem) moveUp() {
	nextY := v.cursor.Y - 1
	if nextY < v.screen.StartYIdx {
		if v.screen.StartYIdx == 0 {
			return
		}
		v.screen.StartYIdx = nextY
		v.screen.InitScreen(v.mode)
	}
	v.cursor.SetCursor(v.cursor.X, nextY)
}

func (v *Veem) moveDown() {
	_, height := v.screen.Current.Size()
	nextY := v.cursor.Y + 1
	if nextY > height-2 {
		v.screen.StartYIdx += 1
		v.screen.InitScreen(v.mode)
		nextY = height - 2
	}
	v.cursor.SetCursor(v.cursor.X, nextY)
}

func (v *Veem) moveRight() {
	nextX := v.cursor.X + 1
	width, _ := v.screen.Current.Size()
	if nextX > width-1 {
		return
	}
	v.cursor.SetCursor(nextX, v.cursor.Y)
}

func (v *Veem) moveLeft() {
	nextX := v.cursor.X - 1
	if nextX < v.screen.StartXIdx {
		return
	}
	v.cursor.SetCursor(nextX, v.cursor.Y)
}

func (v *Veem) saveFile() {
	width, height := v.screen.Current.Size()
	lines := make([]string, height-1)
	for yIdx := range height - 1 {
		for xIdx := range width {

			currRune, _, _, _ := v.screen.Current.GetContent(xIdx+v.screen.StartXIdx, yIdx)
			lines[yIdx] += string(currRune)
		}
	}

	err := utils.WriteLines(lines, v.path)
	if err != nil {
		panic(err)
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
