package terminal

import (
	"fmt"
	"log"
	"os"

	"github.com/mattn/go-isatty"
	"github.com/pkg/term/termios"
	"golang.org/x/sys/unix"
)

var tattr unix.Termios
var savedTattr unix.Termios

func SetupTerminal() {
	// make sure stdin is a terminal
	if !isatty.IsTerminal(os.Stdin.Fd()) {
		log.Fatal("ERROR: this is not a terminal")
	}

	// save the terminal attributes so we can restore them later
	termios.Tcgetattr(os.Stdin.Fd(), &savedTattr)

	// turn off ICANON/ECHO
	termios.Tcgetattr(os.Stdin.Fd(), &tattr)
	tattr.Lflag &^= unix.ICANON | unix.ECHO
	tattr.Cc[unix.VMIN] = 1
	tattr.Cc[unix.VTIME] = 0
	termios.Tcsetattr(os.Stdin.Fd(), termios.TCSAFLUSH, &tattr)

	// hide cursor
	fmt.Printf("\033[?25l")
}

func ResetTattr() {
	// restore terminal
	termios.Tcsetattr(os.Stdin.Fd(), termios.TCSAFLUSH, &savedTattr)

	// show cursor
	fmt.Print("\033[?25h")
}
