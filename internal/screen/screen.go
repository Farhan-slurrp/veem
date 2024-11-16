package screen

import (
	"fmt"
	"log"
	"strconv"

	"github.com/Farhan-slurrp/veem/internal/constants"
	"github.com/Farhan-slurrp/veem/internal/globals"
	"github.com/Farhan-slurrp/veem/internal/utils"
	"github.com/gdamore/tcell/v2"
)

type Screen struct {
	Current  tcell.Screen
	StartIdx int
	path     string
}

func NewScreen(path string) *Screen {
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
		path:     path,
	}
}

func (s *Screen) InitScreen(mode constants.Mode) {
	_, height := s.Current.Size()
	s.StartIdx = utils.GetNumDigits(height) + 1
	s.Current.Clear()

	// paint line numbers
	for i := range height - 1 {
		num := strconv.Itoa(i + 1)
		for y, char := range num {
			s.Current.SetContent(y, i, char, nil, globals.CommentStyle)
		}
	}

	// paint content
	if s.path != "" {
		lines, err := utils.ReadLines(s.path)
		if err != nil {
			panic(fmt.Sprintf("failed to read file %v", s.path))
		}

		for yIdx, line := range lines {
			for xIdx, char := range line {
				s.Current.SetContent(xIdx+s.StartIdx, yIdx, char, nil, globals.DefStyle)
			}
		}
	}

	// paint mode
	s.DisplayMode(mode)
}

func (s *Screen) DisplayMode(mode constants.Mode) {
	_, height := s.Current.Size()
	modeDisplay := fmt.Sprintf("--%s--", mode)

	for index, char := range modeDisplay {
		s.Current.SetContent(index, height-1, char, nil, globals.DefStyle)
	}
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
