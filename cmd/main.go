package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	lib "lem-in/helpers"
	// vis "lem-in/visual"
)

type Graph map[string][]string

func (g Graph) addRoom(name string) {
	if _, exists := g[name]; !exists {
		g[name] = []string{}
	} else {
		fmt.Printf("Room %s already exists\n", name)
		fmt.Printf("Wrong format in graph.txt\n")
		os.Exit(1)
	}
}

func (g Graph) addLink(from, to string) {
	g[from] = append(g[from], to)
	g[to] = append(g[to], from) // Assuming undirected graph
}

func parseG(fileName string) (Graph, int, string, string) {
	g := Graph{}
	data, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(strings.ReplaceAll(string(data), "\r\n", "\n")), "\n")
	if len(lines) == 0 {
		fmt.Println("No data in graph.txt")
		os.Exit(1)
	}
	if lines[0][len(lines[0])-1] == '\r' {
		lines[0] = lines[0][:len(lines[0])-1]
	}
	n, err := strconv.Atoi(lines[0])
	if err != nil {
		fmt.Println("Invalid number of ants:", lines[0])
	} else {
		if n == 0 {
			fmt.Println("Number of ants cannot be zero")
			os.Exit(1)
		}
		fmt.Println("############### ", n, " ants #################")
	}
	start, end := "", ""
	flag := "room"
	for _, line := range lines[1:] {
		if line == "" {
			continue
		}
		name := ""
		parts := strings.Split(line, " ")
		if parts[0] == "##start" {
			flag = "start"
			continue
		} else if parts[0] == "##end" {
			flag = "end"
			continue
		} else if len(parts) == 3 {
			copy := []string{}
			for ii, k := range parts {
				if ii == 0 {
					lib.RoomNames = append(lib.RoomNames, k)
					// if name == "" && flag == "start" {
					// fmt.Println("Room name is empty for start room")
					// if name == "" {
					name = k
					// }
				} else {
					copy = append(copy, k)
				}
			}
			lib.Matrix = append(lib.Matrix, copy)
		}
		part2 := strings.Split(line, "-")

		if len(part2) == 2 {
			room1, room2 := lib.ParseLink(line)
			g.addLink(room1, room2)
		} else if flag == "start" {
			// if start == "" {
			start = name
			// }
			flag = "room"
		} else if flag == "end" {
			end = name

			flag = "room"
		} else {
			g.addRoom(name)
		}

	}
	// graphVisualization(roomNames, matrix, g)
	// fmt.Println("Room Names:", roomNames)
	// fmt.Println("Matrix:", matrix)
	// fmt.Println("Flag:", flag)
	// fmt.Println("graph:", g)
	return g, n, start, end
}

func (g Graph) getBestGroup(start, end string) [][]string {
	// Implement your logic to find the best group of paths for the ants
	ans := [][]string{}

	path := []string{start}

	var dfs func(string)
	dfs = func(current string) {
		if current == end {
			ans = append(ans, append([]string{}, path...))
			return
		}
		for _, neighbor := range g[current] { // only visit once
			if !lib.Contains(path, neighbor) {
				path = append(path, neighbor)
				dfs(neighbor)
				path = path[:len(path)-1]
			}
		}
	}
	dfs(start)
	return ans
}

func (g Graph) sendTheAnts(n int, start, end string) [][]string {
	paths := g.getBestGroup(start, end)
	// fmt.Println("Best Group:", paths)

	if len(paths) == 0 {
		fmt.Println("No paths found from start to end.")
	}

	for i, p := range paths {
		// paths[i] = p[1 : len(p)-1]
		paths[i] = p[1 : len(p)-1]
	}

	for i := 0; i < len(paths); i++ {
		for ii := i + 1; ii < len(paths); ii++ {
			if len(paths[ii-1]) > len(paths[ii]) {
				paths[ii-1], paths[ii] = paths[ii], paths[ii-1]
			}
		}
	}

	// fmt.Println("Sorted Paths:", paths)
	ans := [][][]string{}
	for i := 0; i < len(paths); i++ {
		path2 := [][]string{paths[i]}
		for ii := i + 1; ii < len(paths); ii++ {
			if !lib.Conflicts(path2, paths[ii]) {
				path2 = append(path2, paths[ii])
			}
		}
		ans = append(ans, path2)
	}

	for i, group := range ans {
		for j, p := range group {
			ans[i][j] = append(p, end)
			// ans[i][j] = append([]string{start}, ans[i][j]...)
		}
	}

	return bestGroupByHeight(ans, n)
}

func getCombinedHeights(group [][]string, n int) []int {
	heights := make([]int, len(group))

	for i, k := range group {
		heights[i] = len(k)
	}
	for i := 0; i < n; i++ {
		midIndex := lib.IndexOfMin(heights)
		heights[midIndex]++
	}

	return heights
}

func getHeight(group [][]string, n int) (int, []int) {
	heights := getCombinedHeights(group, n)
	if len(heights) == 0 {
		return 0, nil
	}

	return heights[0], heights
}

func bestGroupByHeight(groups [][][]string, n int) [][]string {
	heights := make([]int, len(groups))
	for i, group := range groups {
		heights[i], _ = getHeight(group, n)
	}

	minIndex := lib.IndexOfMin(heights)
	if minIndex == -1 {
		fmt.Println("\n\t\tERROR: No valid input data found")
		os.Exit(1)
	}
	return groups[minIndex]
}

func getPathHeights(group [][]string) []int {
	heights := make([]int, len(group))
	for i, path := range group {
		heights[i] = len(path)
	}
	return heights
}

func getAntHeights(combined, heights []int) []int {
	for i := range combined {
		combined[i] -= heights[i]
	}
	return combined
}

func main() {
	fileName := ""
	x := os.Args
	if len(x) == 1 {
		fmt.Println("No input file specified.")
		return
	} else if len(x) > 2 {
		fmt.Println("Too many arguments. Please provide only the input file.")
		return
	}
	if !strings.HasSuffix(x[1], ".txt") {
		fmt.Println("Input file must have a .txt extension.")
		return
	} else {
		fileName = x[1]
	}
	g, n, start, end := parseG(fileName)
	bestGroup := g.sendTheAnts(n, start, end)
	// fmt.Println("Best Valid Group:", bestGroup)

	height, combined := getHeight(bestGroup, n)
	// fmt.Println("Height:", height)

	heights := getPathHeights(bestGroup)
	// fmt.Println("Heights:", heights)

	antProxy := getAntHeights(combined, heights)
	// fmt.Println("Ant Proxy Heights:", antProxy)

	counter := 0
	ants := make([]lib.Ant, 0)
	visAnts := []lib.VizAnt{}

	// Initial ant placement - only place one ant per path initially
	for i := range bestGroup {
		if n > 0 && i < len(antProxy) && antProxy[i] > 0 {
			counter++
			ant := lib.Ant{
				Name:    fmt.Sprintf("L%d", counter),
				Current: 0,
				Path:    bestGroup[i],
			}
			ants = append(ants, ant)
			n--
			antProxy[i]--
		}
	}

	for height > 1 {
		height--
		structAnt := []lib.Ant{}

		// Move existing ants
		for i := range ants {
			fmt.Print(ants[i].Name, "->", ants[i].Path[ants[i].Current], " ")
			structAnt = append(structAnt, ants[i])
			ants[i].Current++
		}

		if len(structAnt) > 0 {
			visAnts = append(visAnts, lib.VizAnt{So: structAnt})
		}

		// Remove ants that reached the end
		for i := len(ants) - 1; i >= 0; i-- {
			if ants[i].Current >= len(ants[i].Path) {
				ants = append(ants[:i], ants[i+1:]...)
			}
		}

		// Add new ants if we still have ants to send and available paths
		for i := 0; i < len(antProxy) && n > 0; i++ {
			if antProxy[i] > 0 {
				counter++
				ant := lib.Ant{
					Name:    fmt.Sprintf("L%d", counter),
					Current: 0,
					Path:    bestGroup[i],
				}
				ants = append(ants, ant)
				n--
				antProxy[i]--
			}
		}

		fmt.Println()
	}

	// Continue until all ants are sent and reach the end
	for len(ants) > 0 || n > 0 {
		structAnt := []lib.Ant{}

		// Move existing ants
		for i := range ants {
			if ants[i].Current < len(ants[i].Path) {
				fmt.Print(ants[i].Name, "->", ants[i].Path[ants[i].Current], " ")
				structAnt = append(structAnt, ants[i])
				ants[i].Current++
			}
		}

		if len(structAnt) > 0 {
			visAnts = append(visAnts, lib.VizAnt{So: structAnt})
		}

		// Remove ants that reached the end
		for i := len(ants) - 1; i >= 0; i-- {
			if ants[i].Current >= len(ants[i].Path) {
				ants = append(ants[:i], ants[i+1:]...)
			}
		}

		// Add new ants
		for i := 0; i < len(antProxy) && n > 0; i++ {
			if antProxy[i] > 0 {
				counter++
				ant := lib.Ant{
					Name:    fmt.Sprintf("L%d", counter),
					Current: 0,
					Path:    bestGroup[i],
				}
				ants = append(ants, ant)
				n--
				antProxy[i]--
			}
		}

		fmt.Println()
		if n > 0 {
			fmt.Printf("Remaining ants to send: %d\n", n)
		}

		// Safety break to prevent infinite loop
		if len(ants) == 0 && n > 0 {
			fmt.Printf("Warning: %d ants couldn't be sent (no available paths)\n", n)
			break
		}
	}

	lib.GraphVisualization(lib.RoomNames, lib.Matrix, visAnts, start, end)

	// fmt.Println("Final Ants:", visAnts)
}
