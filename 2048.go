package main

import (
	"fmt"
	"math/rand"
)

const PROBABILITY_OF_SPAWNING_FOUR = 0.25
const NUM_ROWS = 4
const NUM_COLUMNS = 4
const NUM_SPACES = NUM_ROWS * NUM_COLUMNS

func getGridIndex(is_column_vector bool, iterate_backwards bool, vector int, item int) int {
	if is_column_vector {
		if iterate_backwards {
			item = NUM_ROWS - item - 1
		}
		return vector + item*NUM_COLUMNS
	} else {
		if iterate_backwards {
			item = NUM_COLUMNS - item - 1
		}
		return vector*NUM_COLUMNS + item
	}
}

func spawnNewNumber(grid *[NUM_SPACES]int) {
	var open_spaces [NUM_SPACES]int
	num_open_spaces := 0
	for space_index := 0; space_index < NUM_SPACES; space_index++ {
		if grid[space_index] == 0 {
			open_spaces[num_open_spaces] = space_index
			num_open_spaces++
		}
	}
	spawn_index := open_spaces[rand.Intn(num_open_spaces)]
	if rand.Float32() < PROBABILITY_OF_SPAWNING_FOUR {
		grid[spawn_index] = 4
	} else {
		grid[spawn_index] = 2
	}
}

func printGrid(grid *[NUM_SPACES]int) {
	for space_index := 0; space_index < NUM_SPACES; space_index++ {
		fmt.Print(grid[space_index])
		if space_index%NUM_COLUMNS == NUM_COLUMNS-1 {
			fmt.Println()
		} else {
			fmt.Print("\t")
		}
	}
}

func processMove(grid [NUM_SPACES]int, is_column_vector bool, iterate_backwards bool) [NUM_SPACES]int {
	var num_vectors int
	var vector_size int
	if is_column_vector {
		num_vectors = NUM_COLUMNS
		vector_size = NUM_ROWS
	} else {
		num_vectors = NUM_ROWS
		vector_size = NUM_COLUMNS
	}

	for vector := 0; vector < num_vectors; vector++ {
		for item := 0; item < vector_size-1; item++ {
			grid_index := getGridIndex(is_column_vector, iterate_backwards, vector, item)
			if grid[grid_index] == 0 {
				for next_item := item + 1; next_item < vector_size; next_item++ {
					next_grid_index := getGridIndex(is_column_vector, iterate_backwards, vector, next_item)
					if grid[next_grid_index] != 0 {
						grid[grid_index] = grid[next_grid_index]
						grid[next_grid_index] = 0
						break
					}
				}
			}
			space_value := grid[grid_index]
			if space_value != 0 {
				for next_item := item + 1; next_item < vector_size; next_item++ {
					next_grid_index := getGridIndex(is_column_vector, iterate_backwards, vector, next_item)
					if grid[next_grid_index] == space_value {
						grid[grid_index] *= 2
						grid[next_grid_index] = 0
						break
					} else if grid[next_grid_index] != 0 {
						break
					}
				}
			}
		}
	}
	return grid
}

func addMoveOptionIfValid(move_options_map *map[string]([NUM_SPACES]int), grid *[NUM_SPACES]int, name string,
	is_column_vector bool, iterate_backwards bool) {
	result_grid := processMove(*grid, is_column_vector, iterate_backwards)
	if result_grid != *grid {
		(*move_options_map)[name] = result_grid
	}
}

func main() {
	var grid [NUM_SPACES]int
	for {
		spawnNewNumber(&grid)
		printGrid(&grid)
		move_options_map := make(map[string]([NUM_SPACES]int))
		addMoveOptionIfValid(&move_options_map, &grid, "u", true, false)
		addMoveOptionIfValid(&move_options_map, &grid, "d", true, true)
		addMoveOptionIfValid(&move_options_map, &grid, "l", false, false)
		addMoveOptionIfValid(&move_options_map, &grid, "r", false, true)

		move_options := ""
		for move := range move_options_map {
			if move_options != "" {
				move_options += "/"
			}
			move_options += move
		}
		if move_options == "" {
			fmt.Println("Game over")
			break
		}
		prompt := "Enter a move (" + move_options + "): "
		var move string
		for {
			fmt.Print(prompt)
			fmt.Scanln(&move)
			result_grid, valid_move := move_options_map[move]
			if !valid_move {
				fmt.Println("Bad move!")
			} else {
				grid = result_grid
			}
		}
	}
}
