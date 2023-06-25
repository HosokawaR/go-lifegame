package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

type World struct {
	Cells  []Cell
	Width  int
	Height int
}

type Cell struct {
	Alive bool
}

const (
	Width     = 40
	Height    = 20
	AliveChar = 'o'
	DeadChar  = '_'
)

func main() {
	generation := 0
	world := generateWorld(Width, Height)

	for {
		print(world, generation)
		world = next(world)
		generation++
		time.Sleep(1 * time.Second)
	}
}

func generateWorld(width int, height int) World {
	cells := make([]Cell, width*height)

	for i := range cells {
		if rand.Float64() < 0.2 {
			cells[i].Alive = true
		} else {
			cells[i].Alive = false
		}
	}

	world := World{
		Cells:  cells,
		Width:  width,
		Height: height,
	}

	return world
}

func next(world World) World {
	nextCells := make([]Cell, world.Width*world.Height)

	for x := 0; x < world.Width; x++ {
		for y := 0; y < world.Height; y++ {
			index := y*world.Width + x
			count := countAroundAliveCells(world, index)
			if count == 3 {
				nextCells[index].Alive = true
			} else if count == 2 {
				nextCells[index].Alive = world.Cells[index].Alive
			} else {
				nextCells[index].Alive = false
			}
		}
	}

	nextWorld := World{
		Cells:  nextCells,
		Width:  world.Width,
		Height: world.Height,
	}

	return nextWorld
}

func countAroundAliveCells(world World, centerCellIndex int) int {
	count := 0

	for xdiff := -1; xdiff <= 1; xdiff++ {
		for ydiff := -1; ydiff <= 1; ydiff++ {
			if xdiff == 0 && ydiff == 0 {
				continue
			}
			targetCellIndex := centerCellIndex + (world.Width * ydiff) + xdiff
			if targetCellIndex < 0 || len(world.Cells)-1 < targetCellIndex {
				continue
			}
			if world.Cells[targetCellIndex].Alive {
				count++
			}
		}
	}

	return count
}

func print(world World, generation int) {
	clear()

	fmt.Printf("Generation: %d\n", generation)

	for i, cell := range world.Cells {
		if cell.Alive {
			fmt.Printf("%c", AliveChar)
		} else {
			fmt.Printf("%c", DeadChar)
		}

		if i%world.Width == world.Width-1 {
			fmt.Printf("\n")
		}
	}
}

func clear() {
	// Note: The clear command is available only in Unix systems.
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
