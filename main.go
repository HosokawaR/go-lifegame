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
	ch := make(chan bool)

	for x := 0; x < world.Width; x++ {
		for y := 0; y < world.Height; y++ {
			index := y*world.Width + x
			// Using Goroutines here is not a good idea
			// but they are being used for practice purposes.
			go cellAlive(world, index, ch)
		}
	}

	for x := 0; x < world.Width; x++ {
		for y := 0; y < world.Height; y++ {
			index := y*world.Width + x
			nextCells[index].Alive = <-ch
		}
	}

	nextWorld := World{
		Cells:  nextCells,
		Width:  world.Width,
		Height: world.Height,
	}

	return nextWorld
}

func cellAlive(world World, targetCellIndex int, ch chan bool) {
	var alive bool

	count := countAroundAliveCells(world, targetCellIndex)
	if count == 3 {
		alive = true
	} else if count == 2 {
		alive = world.Cells[targetCellIndex].Alive
	} else {
		alive = false
	}

	ch <- alive
}

func countAroundAliveCells(world World, centerCellIndex int) int {
	dummyHeavyTask()

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

func dummyHeavyTask() {
	time.Sleep(1 * time.Millisecond)
}
