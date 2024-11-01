package screen

import (
	"log"
	"strconv"

	"github.com/Farhan-slurrp/veem/internal/globals"
	"github.com/gdamore/tcell/v2"
)

type Screen struct {
	Current tcell.Screen
}

func NewScreen() *Screen {
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	initScreen(s)

	return &Screen{
		Current: s,
	}
}

func initScreen(s tcell.Screen) {
	s.EnableMouse()
	s.EnablePaste()
	s.Clear()

	_, height := s.Size()

	// paint page number
	for i := range height - 1 {
		num := strconv.Itoa(i + 1)
		s.SetContent(0, i, []rune(num)[0], nil, globals.CommentStyle)
	}
}
