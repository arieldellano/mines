package field

import "fmt"

type cellContent byte

const (
	empty cellContent = iota
	bomb
)

type cell struct {
	content cellContent
	col     int
	row     int
	visible bool
	field   *Field
	flagged bool
}

func (c *cell) String() string {

	if c.flagged {
		return "\033[1;33;41mF\033[0m"
	}

	if !c.visible {
		return "\033[1;37m░\033[0m"
	}

	if c.content == bomb {
		return fmt.Sprint("\033[1;31m@\033[0m")
	}

	if c.content == empty {
		if bombs := c.field.countNearbyBombs(c.row, c.col); bombs > 0 {

			color := 37
			switch bombs {
			case 8, 7, 6:
				color = 35
			case 5:
				color = 33
			case 4:
				color = 36
			case 3:
				color = 34
			case 2:
				color = 32
			case 1:
				color = 30
			}

			return fmt.Sprintf("\033[1;%dm%d\033[0m", color, bombs)
		}
	}

	return fmt.Sprint(" ")
}
