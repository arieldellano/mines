package field

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"virgee.com/mines/terminal"
)

type Field struct {
	cells          []cell
	rows           int
	cols           int
	crow           int
	ccol           int
	percentage     int
	bombsGenerated bool
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func NewField(rows, cols, percentage int) *Field {
	if rows < 5 {
		log.Fatal("Minimum number of rows is 5")
	}

	if cols < 5 {
		log.Fatal("Minimum number of columns is 5")
	}

	if percentage <= 0 {
		log.Fatal("Percentage must be between 1 and 100")
	}

	// seed randomizer
	rand.Seed(time.Now().UnixNano())

	// make cells
	cells := make([]cell, rows*cols)

	// initialize Field
	f := &Field{
		cells,
		rows,
		cols,
		rows / 2,
		cols / 2,
		min(percentage, 100),
		false,
	}

	f.initializeCells()
	return f
}

var won bool

func (f *Field) initializeCells() {
	for row := 0; row < f.rows; row++ {
		for col := 0; col < f.cols; col++ {
			c := f.getCell(row, col)
			c.field = f
			c.flagged = false
			c.visible = false
			c.row = row
			c.col = col
		}
	}
}

// Gets the number of cells
func (f *Field) Cells() int {
	return f.cols * f.rows
}

// Gets the number of bombs
func (f *Field) Bombs() int {
	count := int(float32(f.rows*f.cols) * (float32(f.percentage) / 100))

	if count > (len(f.cells) - 9) {
		count = (len(f.cells) - 9)
	}

	return count
}

func (f *Field) bombsLeft() int {
	return f.Bombs() - f.flagged()
}

func (f *Field) flagged() int {
	count := 0
	for _, c := range f.cells {
		if c.flagged {
			count++
		}
	}

	return count
}

func (f *Field) visible() int {
	// count visible cells
	count := 0
	for _, c := range f.cells {
		if c.visible {
			count++
		}
	}

	return count
}

func (f *Field) getCell(row, col int) *cell {
	if row < 0 || col < 0 || row >= f.rows || col >= f.cols {
		log.Fatalf(
			"Row %d/Col %d outside field bounds (0-%d/0-%d)",
			row, col, f.rows, f.cols,
		)
	}

	return &f.cells[row*f.cols+col]
}

func (f *Field) visitNeighbors(row, col int, visitor func(row, col int)) {
	for drow := -1; drow <= 1; drow++ {
		for dcol := -1; dcol <= 1; dcol++ {

			if drow == 0 && dcol == 0 {
				continue
			}

			erow := row + drow
			ecol := col + dcol

			if erow < 0 || ecol < 0 || erow >= f.rows || ecol >= f.cols {
				continue
			}

			visitor(erow, ecol)
		}
	}
}

func (f *Field) isDemilitarizedZone(row, col int) bool {
	for drow := -1; drow <= 1; drow++ {
		for dcol := -1; dcol <= 1; dcol++ {
			if row == drow+f.crow && col == dcol+f.ccol {
				return true
			}
		}
	}

	return false
}

func (f *Field) generateBombs() {
	bombCount := f.Bombs()
	for bombCount > 0 {
		row := rand.Intn(f.rows)
		col := rand.Intn(f.cols)

		cell := f.getCell(row, col)

		// if there is a bomb in the cell, continue
		if cell.content == bomb {
			continue
		}

		// make sure no bomb is placed around the cursor
		if f.isDemilitarizedZone(row, col) {
			continue
		}

		// place the bomb
		cell.content = bomb
		bombCount -= 1
	}

	f.bombsGenerated = true
}

func (f *Field) countNearbyIf(row, col int, cond func(c *cell) bool) int {
	count := 0
	f.visitNeighbors(row, col, func(erow, ecol int) {
		c := f.getCell(erow, ecol)
		if cond(c) {
			count += 1
		}
	})

	return count
}

func (f *Field) countNearbyBombs(row, col int) int {
	return f.countNearbyIf(row, col, func(c *cell) bool {
		return c.content == bomb
	})
}

func (f *Field) countNearbyFlagged(row, col int) int {
	return f.countNearbyIf(row, col, func(c *cell) bool {
		return c.flagged
	})
}

// Prints the state of the field
func (f *Field) Print(skipClear bool) {

	// status line

	str := fmt.Sprintf("%%%dd bombs left\n", (f.cols*4)-11)
	fmt.Printf(str, f.bombsLeft())

	// print field
	f.printFieldTop()
	for row := 0; row < f.rows; row++ {
		f.printRow(row)
	}
	f.printFieldBottom()

	// clear unless skipped
	if !skipClear {
		fmt.Printf("\033[%dA", 1+(f.rows*2)+1)
		fmt.Printf("\033[%dD", (f.cols*4)+1)
	}
}

func (f *Field) printRow(row int) {
	if row > 0 {
		f.printCellTop(row)
	}

	for col := 0; col < f.cols; col++ {
		f.printCell(row, col)
	}

	fmt.Println()
}

func (f *Field) printCellTop(row int) {
	for col := 0; col <= f.cols; col++ {
		if (col == 0 && f.ccol == col && row == f.crow) || (f.ccol == col && row == f.crow) {
			fmt.Print("\033[1;37m╔═══\033[0m")
		} else if (col == 0 && f.ccol == col && row == f.crow+1) || (f.ccol == col && row == f.crow+1) {
			fmt.Print("\033[1;37m╚═══\033[0m")
		} else if col == f.cols && row == f.crow && f.ccol+1 == col {
			fmt.Println("\033[1;37m╗\033[0m")
		} else if col == f.cols && row == f.crow+1 && f.ccol+1 == col {
			fmt.Println("\033[1;37m╝\033[0m")
		} else if col == f.cols {
			fmt.Println("\033[2;37m┤\033[0m")
		} else if col == f.ccol+1 && row == f.crow {
			fmt.Print("\033[1;37m╗\033[2;37m───\033[0m")
		} else if row == f.crow+1 && col == f.ccol+1 {
			fmt.Print("\033[1;37m╝\033[2;37m───\033[0m")
		} else {
			fmt.Print("\033[2;37m┼───\033[0m")
		}
	}
}

func (f *Field) printFieldTop() {
	for col := 0; col <= f.cols; col++ {
		if col == 0 && f.crow == 0 && col == f.ccol {
			fmt.Print("\033[1;37m╔═══\033[0m")
		} else if col == 0 {
			fmt.Print("\033[2;37m┌───\033[0m")
		} else if col == f.cols && f.crow == 0 && col == f.ccol+1 {
			fmt.Println("\033[1;37m╗\033[0m")
		} else if col == f.cols {
			fmt.Println("\033[2;37m┐\033[0m")
		} else if f.crow == 0 && col == f.ccol+1 {
			fmt.Print("\033[1;37m╗\033[2;37m───\033[0m")
		} else if f.crow == 0 && col == f.ccol {
			fmt.Print("\033[1;37m╔═══\033[0m")
		} else {
			fmt.Print("\033[2;37m┬───\033[0m")
		}
	}
}

func (f *Field) printCell(row int, col int) {
	c := f.getCell(row, col)
	color := 37
	opt, opt2 := 2, 2
	ch, ch2 := '│', '│'

	if row == f.crow && f.ccol == f.cols-1 && col == f.cols-1 {
		ch, ch2 = '║', '║'
		opt, opt2 = 1, 1
		color = 37
	}

	if row == f.crow && (col == f.ccol || col == f.ccol+1) {
		ch = '║'
		opt = 1
		color = 37
	}

	if col == f.cols-1 {
		fmt.Printf("\033[%d;%dm%c\033[0m %v \033[%d;%dm%c\033[0m", opt, color, ch, c, opt2, color, ch2)
	} else {
		fmt.Printf("\033[%d;%dm%c\033[0m %v ", opt, color, ch, c)
	}
}

func (f *Field) printFieldBottom() {
	for col := 0; col <= f.cols; col++ {
		switch col {
		case 0:
			if f.crow == f.rows-1 && col == f.ccol {
				fmt.Print("\033[1;37m╚═══\033[0m")
			} else {
				fmt.Print("\033[2;37m└───\033[0m")
			}
		case f.cols:
			if f.crow == f.rows-1 && f.ccol == f.cols-1 {
				fmt.Println("\033[1;37m╝\033[0m")
			} else {
				fmt.Println("\033[2;37m┘\033[0m")
			}
		default:
			if f.crow == f.rows-1 && col == f.ccol {
				fmt.Print("\033[1;37m╚═══\033[0m")
			} else if f.crow == f.rows-1 && col == f.ccol+1 {
				fmt.Print("\033[1;37m╝\033[2;37m───\033[0m")
			} else {
				fmt.Print("\033[2;37m┴───\033[0m")
			}
		}
	}
}

func (f *Field) up() {
	f.crow -= 1
	if f.crow < 0 {
		f.crow = f.rows - 1
	}
}

func (f *Field) down() {
	f.crow += 1
	if f.crow >= f.rows {
		f.crow = 0
	}
}

func (f *Field) left() {
	f.ccol -= 1
	if f.ccol < 0 {
		f.ccol = f.cols - 1
	}
}

func (f *Field) right() {
	f.ccol += 1
	if f.ccol >= f.cols {
		f.ccol = 0
	}
}

func (f *Field) panic(v ...any) {
	f.showBombs()
	f.Print(true)
	fmt.Fprintln(os.Stderr, v...)
	terminal.ResetTattr()
	os.Exit(0)
}

func (f *Field) tapCell(row, col int) {
	c := f.getCell(row, col)
	flagged := f.countNearbyFlagged(c.row, c.col)
	nearbyBombCount := f.countNearbyBombs(c.row, c.col)

	// if the cell is visible...
	if c.visible {

		// ...and have at least 1 bomb nearby
		// then we can give the player a little help
		if nearbyBombCount > 0 && flagged == nearbyBombCount {
			f.visitNeighbors(c.row, c.col, func(erow, ecol int) {
				c := f.getCell(erow, ecol)
				if !c.visible && !c.flagged {
					f.tapCell(erow, ecol)
				}
			})
		}

		// otherwise there is nothing to do
		return
	}

	// otherwise, make it visible
	c.visible = true
	c.flagged = false

	// then, we check for its content
	// if its a bomb, we panic
	if c.content == bomb {
		f.panic("Boooooom!!! You stepped on a bomb and you died.")
	}

	// ok, we are safe
	// but how safe? how many bombs around us?
	if f.countNearbyBombs(row, col) > 0 {
		// if at least one bomb,
		// there is nothing we can do to help you
		return
	}

	// but if there are  no bombs around you
	// let's give you a little help and tap them for you
	f.visitNeighbors(row, col, func(erow, ecol int) {
		f.tapCell(erow, ecol)
	})
}

func (f *Field) showBombs() error {
	for row := 0; row < f.rows; row++ {
		for col := 0; col < f.cols; col++ {
			c := f.getCell(row, col)

			if c.content == bomb {
				c.visible = true
			}
		}
	}

	return nil
}

func (f *Field) DidWin() bool {
	// count visible cells
	untappedCellCount := len(f.cells) - f.visible() - f.flagged()
	if untappedCellCount >= 1 {
		return false
	}

	// everything seems ok so far, so...
	// make sure all flagged cells have a  bomb?
	for _, c := range f.cells {
		if c.flagged && c.content != bomb {
			return false
		}
	}

	return true
}

func (f *Field) toggleFlag() {
	c := f.getCell(f.crow, f.ccol)

	// can't flag if there are no bombs left
	// or if the cell is already visible
	if f.bombsLeft() <= 0 || c.visible {
		return
	}

	c.flagged = !c.flagged
}

func (f *Field) Play() bool {
	reader := bufio.NewReader(os.Stdin)
	ch, _, err := reader.ReadRune()
	if err != nil {
		f.panic(err)
	}

	// handle keyed char
	switch ch {
	case 'w':
		f.up()
	case 'a':
		f.left()
	case 's':
		f.down()
	case 'd':
		f.right()
	case 'f':
		f.toggleFlag()
	case ' ':
		if !f.bombsGenerated {
			f.generateBombs()
		}

		f.tapCell(f.crow, f.ccol)
	case 'q':
		return true
	case 27:
		_, _, err := reader.ReadRune()
		if err != nil {
			log.Fatal(err)
		}

		char, _, err := reader.ReadRune()
		if err != nil {
			log.Fatal(err)
		}

		switch char {
		case 65:
			f.up()
		case 66:
			f.down()
		case 67:
			f.right()
		case 68:
			f.left()
		}
	}

	return false
}
