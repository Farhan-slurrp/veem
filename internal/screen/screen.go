package screen

import (
	"log"
	"strconv"

	"github.com/Farhan-slurrp/veem/internal/globals"
	"github.com/Farhan-slurrp/veem/internal/utils"
	"github.com/gdamore/tcell/v2"
)

type Screen struct {
	Current  tcell.Screen
	StartIdx int
}

func NewScreen() *Screen {
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

	return &Screen{
		Current:  s,
		StartIdx: 2,
	}
}

func (s *Screen) InitScreen(initialContent []string) {
	width, height := s.Current.Size()
	s.StartIdx = utils.GetNumDigits(height) + 1

	if len(initialContent) > height {
		s.Current.SetSize(width, len(initialContent))
	}

	// paint line numbers
	go func() {
		for i := range height - 1 {
			num := strconv.Itoa(i + 1)
			for y, char := range num {
				s.Current.SetContent(y, i, char, nil, globals.CommentStyle)
			}
		}
	}()

	// paint content
	go func() {
		for yIdx, line := range initialContent {
			for xIdx, char := range line {
				s.Current.SetContent(xIdx+s.StartIdx, yIdx, char, nil, globals.DefStyle)
			}
		}
	}()
}

func (s *Screen) ShiftContentRight(curX int, curY int, content rune) {
	width, _ := s.Current.Size()
	currRune, _, _, _ := s.Current.GetContent(curX, curY)
	s.Current.SetContent(curX, curY, content, nil, globals.DefStyle)
	for i := curX; i < width; i++ {
		nextRune, _, _, _ := s.Current.GetContent(i+1, curY)
		s.Current.SetContent(i+1, curY, currRune, nil, globals.DefStyle)
		currRune = nextRune
	}
}

func (s *Screen) ShiftContentLeft(curX int, curY int, content rune) {
	width, _ := s.Current.Size()
	currRune, _, _, _ := s.Current.GetContent(curX, curY)
	s.Current.SetContent(curX, curY, content, nil, globals.DefStyle)
	for i := curX; i < width; i++ {
		nextRune, _, _, _ := s.Current.GetContent(i+1, curY)
		s.Current.SetContent(i-1, curY, currRune, nil, globals.DefStyle)
		currRune = nextRune
	}
}
