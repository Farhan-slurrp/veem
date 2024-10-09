package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mattn/go-tty"
	"golang.org/x/term"
)

func main() {
	fmt.Print("\033[H\033[2J")
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

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
		charCode := fmt.Sprintf("%X", r)
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
			switch charCode {
			case "D":
				switch commands {
				case ":q":
					return
				default:
					fmt.Printf("\n")
				}
				continue
			default:
				if char == "" {
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
