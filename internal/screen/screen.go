package screen

import (
	"bufio"
	"log"
	"os"
	"strconv"

	"github.com/Farhan-slurrp/veem/internal/globals"
	"github.com/gdamore/tcell/v2"
)

type Screen struct {
	Current tcell.Screen
}

func NewScreen(file *os.File) *Screen {
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	initScreen(s, file)

	return &Screen{
		Current: s,
	}
}

func initScreen(s tcell.Screen, file *os.File) {
	s.EnableMouse()
	s.EnablePaste()
	s.Clear()

	_, height := s.Size()

	// paint line numbers
	for i := range height - 1 {
		num := strconv.Itoa(i + 1)
		s.SetContent(0, i, []rune(num)[0], nil, globals.CommentStyle)
	}

	// paint content
	scanner := bufio.NewScanner(file)
	line := 0
	for scanner.Scan() {
		for xIdx, _ := range scanner.Text() {
			s.SetContent(xIdx+2, line, 'a', nil, globals.DefStyle)
		}
		line += 1
	}
}
