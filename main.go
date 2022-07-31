package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"virgee.com/mines/field"
)

func usage() {
	fmt.Println("mines - a minesweeper clone for terminal nerds")
	fmt.Println("USAGE:")
	fmt.Println("\tmines <rows> <cols> [percentage]")
	fmt.Println("\n\trows\t\t- required, number of rows in the field")
	fmt.Println("\tcols\t\t- required, number of columns in the field")
	fmt.Println("\tpercentage\t- optional, percentage of the field to populate with mines (20% default)")
	fmt.Println("\nEXAMPLES:")
	fmt.Println("\tmines 10 25\t - a 10 rows by 25 columns field, 20% mine population")
	fmt.Println("\tmines 10 10 25\t - a 10 rows by 10 columns field, 25% mine population")

	os.Exit(0)
}

func main() {

	argc := len(os.Args[1:])
	percentage := 20
	if argc < 2 || argc > 3 {
		usage()
	}

	rows, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	cols, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	if argc == 3 {
		p, err := strconv.Atoi(os.Args[3])
		if err != nil {
			log.Fatal(err)
		}
		percentage = p
	}

	quit := false

	f := field.NewField(rows, cols, percentage)

	f.Print(false)
	for !quit {

		// make your move
		quit = f.Play()

		// update terminal
		didWin := f.DidWin()
		f.Print(didWin)

		// did you win?
		if didWin {

			fmt.Println("OMG! You did it!! You won the game!!!!!")
			switch {
			case percentage < 10:
				fmt.Printf("But that was super easy, you only had to flag %d bombs out of %d cells. Come on! you can do better.\n", f.Bombs(), f.Cells())
			case percentage >= 10 && percentage < 20:
				fmt.Printf("Most people do ok at %d%%, try increasing your percentage. You had to flag just %d bombs out of %d cells.\n", percentage, f.Bombs(), f.Cells())
			case percentage >= 20 && percentage < 30:
				fmt.Printf("This is getting serious... at %d%% you are above average. You flaggged %d bombs out of %d cells.\n", percentage, f.Bombs(), f.Cells())
			case percentage >= 30 && percentage < 40:
				fmt.Printf("Are you for real? at %d%% you are way above average. You flaggged %d bombs out of %d cells.\n", percentage, f.Bombs(), f.Cells())
			case percentage >= 40 && percentage < 50:
				fmt.Printf("Seriously, are you human? at %d%% you are at elite status! You flaggged %d bombs out of %d cells.\n", percentage, f.Bombs(), f.Cells())
			case percentage >= 50:
				fmt.Printf("You must be an AI, right? at %d%% you are at god level! You flaggged %d bombs out of %d cells.\n", percentage, f.Bombs(), f.Cells())
			}

			quit = true
		}
	}
}
