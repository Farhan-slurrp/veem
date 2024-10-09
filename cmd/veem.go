package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mattn/go-tty"
	"golang.org/x/term"
)

type Mode string

var (
	NORMAL Mode = "NORMAL"
	INSERT Mode = "INSERT"
)

type Veem struct {
	mode            Mode
	currentCommands string
	recordsCommands string
}

func main() {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	fmt.Print("\033[H\033[2J")
	t, err := tty.Open()
	if err != nil {
		panic(err)
	}
	defer t.Close()

	go func() {
		for ws := range t.SIGWINCH() {
			fmt.Println("Resized", ws.W, ws.H)
		}
	}()

	clean, err := t.Raw()
	if err != nil {
		log.Fatal(err)
	}
	defer clean()

	recordCommands := false
	commands := ""
	for {
		r, err := t.ReadRune()
		if err != nil {
			log.Fatal(err)
		}
		if r == 0 {
			continue
		}

		char := fmt.Sprintf("%c", r)

		switch char {
		case ":":
			_, height, err := t.Size()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("\033[%d;%dH", height, 0)
			recordCommands = true
			commands = ":"
			fmt.Print(commands)
		default:
			switch r {
			case 13:
				switch commands {
				case ":q":
					return
				default:
					fmt.Printf("\n")
				}
			default:
				if strings.TrimSpace(char) == "" {
					continue
				}
				if recordCommands {
					commands += char
				}
				fmt.Print(char)

			}
		}
	}
}
