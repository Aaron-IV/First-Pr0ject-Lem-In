package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// Graph — карта: имя комнаты -> список соседей
type Graph map[string][]string

type roomCoord struct {
	X int
	Y int
}

func (g Graph) addRoom(name string) {
	if _, exists := g[name]; !exists {
		g[name] = []string{}
	}
}

func (g Graph) addLink(from, to string) {
	g[from] = append(g[from], to)
	g[to] = append(g[to], from)
}

// parseG — читает и валидирует файл, возвращает граф, n, start, end и координаты
func parseG(fileName string) (Graph, int, string, string, map[string]roomCoord, error) {
	g := Graph{}
	coords := make(map[string]roomCoord)

	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, 0, "", "", nil, fmt.Errorf("read file: %w", err)
	}
	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	lines := strings.Split(text, "\n")
	if len(lines) == 0 {
		return nil, 0, "", "", nil, fmt.Errorf("empty file")
	}

	n, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil || n <= 0 {
		return nil, 0, "", "", nil, fmt.Errorf("invalid ants count in first line")
	}

	start, end := "", ""
	flagMode := "room"

	for _, raw := range lines[1:] {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "##") {
			continue
		}
		if strings.HasPrefix(line, "##") {
			switch strings.TrimSpace(line) {
			case "##start":
				flagMode = "start"
			case "##end":
				flagMode = "end"
			}
			continue
		}

		parts := strings.Fields(line)
		if len(parts) == 3 { // room line
			roomName := parts[0]
			x, errX := strconv.Atoi(parts[1])
			y, errY := strconv.Atoi(parts[2])
			if roomName == "" || errX != nil || errY != nil {
				return nil, 0, "", "", nil, fmt.Errorf("invalid room line: %s", line)
			}
			g.addRoom(roomName)
			coords[roomName] = roomCoord{X: x, Y: y}
			if flagMode == "start" {
				start = roomName
				flagMode = "room"
			} else if flagMode == "end" {
				end = roomName
				flagMode = "room"
			}
			continue
		}

		if strings.Contains(line, "-") && !strings.Contains(line, " ") { // link line
			linkParts := strings.Split(line, "-")
			if len(linkParts) != 2 {
				return nil, 0, "", "", nil, fmt.Errorf("invalid link: %s", line)
			}
			a := strings.TrimSpace(linkParts[0])
			b := strings.TrimSpace(linkParts[1])
			if a == b || a == "" || b == "" {
				return nil, 0, "", "", nil, fmt.Errorf("invalid link endpoints: %s", line)
			}
			if _, ok := g[a]; !ok {
				return nil, 0, "", "", nil, fmt.Errorf("link references undefined room: %s", a)
			}
			if _, ok := g[b]; !ok {
				return nil, 0, "", "", nil, fmt.Errorf("link references undefined room: %s", b)
			}
			// prevent duplicates (simple check)
			exists := false
			for _, v := range g[a] {
				if v == b {
					exists = true
					break
				}
			}
			if exists {
				return nil, 0, "", "", nil, fmt.Errorf("duplicate link: %s-%s", a, b)
			}
			g.addLink(a, b)
			continue
		}

		return nil, 0, "", "", nil, fmt.Errorf("wrong format near: %s", line)
	}

	if start == "" || end == "" {
		return nil, 0, "", "", nil, fmt.Errorf("start or end not defined")
	}
	return g, n, start, end, coords, nil
}

// getBestGroup — все простые пути start->end
func (g Graph) getBestGroup(start, end string) [][]string {
	ans := [][]string{}
	path := []string{start}
	visited := map[string]bool{start: true}
	var dfs func(string)
	dfs = func(cur string) {
		if cur == end {
			cp := append([]string{}, path...)
			ans = append(ans, cp)
			return
		}
		for _, nb := range g[cur] {
			if !visited[nb] {
				visited[nb] = true
				path = append(path, nb)
				dfs(nb)
				path = path[:len(path)-1]
				delete(visited, nb)
			}
		}
	}
	dfs(start)
	return ans
}

func indexOfMin(arr []int) int {
	if len(arr) == 0 {
		return -1
	}
	min := arr[0]
	idx := 0
	for i := 1; i < len(arr); i++ {
		if arr[i] < min {
			min = arr[i]
			idx = i
		}
	}
	return idx
}

func getCombinedHeights(group [][]string, n int) []int {
	heights := make([]int, len(group))
	for i, k := range group {
		heights[i] = len(k)
	}
	for i := 0; i < n; i++ {
		midIndex := indexOfMin(heights)
		if midIndex >= 0 {
			heights[midIndex]++
		}
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
	minIndex := indexOfMin(heights)
	if minIndex == -1 {
		return [][]string{}
	}
	return groups[minIndex]
}

func (g Graph) sendTheAnts(n int, start, end string) [][]string {
	paths := g.getBestGroup(start, end)
	if len(paths) == 0 {
		return [][]string{}
	}
	sort.Slice(paths, func(i, j int) bool {
		if len(paths[i]) != len(paths[j]) {
			return len(paths[i]) < len(paths[j])
		}
		return strings.Join(paths[i], ",") < strings.Join(paths[j], ",")
	})
	for i, p := range paths {
		if len(p) >= 2 {
			paths[i] = p[1 : len(p)-1]
		} else {
			paths[i] = []string{}
		}
	}
	for i := 0; i < len(paths); i++ {
		for ii := i + 1; ii < len(paths); ii++ {
			if len(paths[ii-1]) > len(paths[ii]) {
				paths[ii-1], paths[ii] = paths[ii], paths[ii-1]
			}
		}
	}
	ans := [][][]string{}
	for i := 0; i < len(paths); i++ {
		path2 := [][]string{paths[i]}
		for ii := i + 1; ii < len(paths); ii++ {
			conflict := false
			for _, s1 := range paths[ii] {
				for _, arr1 := range path2 {
					for _, s2 := range arr1 {
						if s1 == s2 {
							conflict = true
							break
						}
					}
					if conflict {
						break
					}
				}
				if conflict {
					break
				}
			}
			if !conflict {
				path2 = append(path2, paths[ii])
			}
		}
		ans = append(ans, path2)
	}
	for i, group := range ans {
		for j, p := range group {
			ans[i][j] = append(p, end)
		}
	}
	return bestGroupByHeight(ans, n)
}

// generateMoves — формирует шаги вида "L1-roomA L2-roomB"
func generateMoves(g Graph, n int, start, end string) []string {
	bestGroup := g.sendTheAnts(n, start, end)
	height, combined := getHeight(bestGroup, n)
	heights := make([]int, len(bestGroup))
	for i, path := range bestGroup {
		heights[i] = len(path)
	}
	// combined[i] -= heights[i]
	for i := range combined {
		combined[i] -= heights[i]
	}

	counter := 0
	type Ant struct {
		Name    string
		Current int
		Path    []string
	}
	ants := make([]Ant, 0)

	for i := range bestGroup {
		if n > 0 && i < len(combined) && combined[i] > 0 {
			counter++
			ant := Ant{Name: fmt.Sprintf("L%d", counter), Current: 0, Path: bestGroup[i]}
			ants = append(ants, ant)
			n--
			combined[i]--
		}
	}

	var steps []string

	for height > 1 {
		height--
		lineParts := []string{}
		for i := range ants {
			lineParts = append(lineParts, fmt.Sprintf("%s-%s", ants[i].Name, ants[i].Path[ants[i].Current]))
			ants[i].Current++
		}
		// remove finished
		for i := len(ants) - 1; i >= 0; i-- {
			if ants[i].Current >= len(ants[i].Path) {
				ants = append(ants[:i], ants[i+1:]...)
			}
		}
		for i := 0; i < len(combined) && n > 0; i++ {
			if combined[i] > 0 {
				counter++
				ant := Ant{Name: fmt.Sprintf("L%d", counter), Current: 0, Path: bestGroup[i]}
				ants = append(ants, ant)
				n--
				combined[i]--
			}
		}
		steps = append(steps, strings.TrimSpace(strings.Join(lineParts, " ")))
	}

	for len(ants) > 0 || n > 0 {
		lineParts := []string{}
		for i := range ants {
			if ants[i].Current < len(ants[i].Path) {
				lineParts = append(lineParts, fmt.Sprintf("%s-%s", ants[i].Name, ants[i].Path[ants[i].Current]))
				ants[i].Current++
			}
		}
		for i := len(ants) - 1; i >= 0; i-- {
			if ants[i].Current >= len(ants[i].Path) {
				ants = append(ants[:i], ants[i+1:]...)
			}
		}
		for i := 0; i < len(combined) && n > 0; i++ {
			if combined[i] > 0 {
				counter++
				ant := Ant{Name: fmt.Sprintf("L%d", counter), Current: 0, Path: bestGroup[i]}
				ants = append(ants, ant)
				n--
				combined[i]--
			}
		}
		if len(lineParts) > 0 {
			steps = append(steps, strings.TrimSpace(strings.Join(lineParts, " ")))
		} else {
			break
		}
	}

	return steps
}

type roomJSON struct {
	Name    string   `json:"name"`
	X       int      `json:"x"`
	Y       int      `json:"y"`
	IsStart bool     `json:"isStart"`
	IsEnd   bool     `json:"isEnd"`
	Links   []string `json:"links"`
}

type dataJSON struct {
	Rooms []roomJSON `json:"rooms"`
	Moves []string   `json:"moves"`
}

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	file := flag.String("file", "examples/example05.txt", "input graph file")
	webDir := flag.String("web", "web", "web directory")
	flag.Parse()

	// Positional args support: server [file] [addr]
	// Example: go run ./cmd/server examples/example02.txt :8080
	args := flag.Args()
	if len(args) >= 1 {
		*file = args[0]
	}
	if len(args) >= 2 {
		*addr = args[1]
	}

	// Validate input file early and print CLI-like output before starting HTTP
	if *file == "" {
		fmt.Println("No input file specified.")
		os.Exit(1)
	}
	if !strings.HasSuffix(*file, ".txt") {
		fmt.Println("Input file must have a .txt extension.")
		os.Exit(1)
	}
	useFile := *file
	if !filepath.IsAbs(useFile) {
		useFile = filepath.Join(".", useFile)
	}

	g, n, start, end, _, err := parseG(useFile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Print original file contents like CLI does (first line, then the rest)
	if dataRaw, err := os.ReadFile(useFile); err == nil {
		dataStr := strings.ReplaceAll(string(dataRaw), "\r\n", "\n")
		lines := strings.Split(dataStr, "\n")
		if len(lines) > 0 && strings.TrimSpace(lines[0]) != "" {
			fmt.Println(lines[0])
		}
		if len(lines) > 1 {
			fmt.Print(strings.Join(lines[1:], "\n"))
		}
		fmt.Println()
		fmt.Println()
	}

	// Generate moves and print them line-by-line to stdout, identical formatting
	steps := generateMoves(g, n, start, end)
	for _, line := range steps {
		fmt.Println(line)
	}
	fmt.Println()

	// Print link to visualization for this file
	fmt.Printf("Open visualization: http://localhost%s/visual.html?file=%s\n", *addr, *file)

	// static files
	fs := http.FileServer(http.Dir(*webDir))
	http.Handle("/", fs)

	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		qFile := r.URL.Query().Get("file")
		useFile := *file
		if qFile != "" {
			useFile = qFile
		}
		// make absolute path if needed
		if !filepath.IsAbs(useFile) {
			useFile = filepath.Join(".", useFile)
		}
		g, n, start, end, coords, err := parseG(useFile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		steps := generateMoves(g, n, start, end)

		rooms := make([]roomJSON, 0, len(g))
		for name, links := range g {
			c := coords[name]
			rj := roomJSON{
				Name:    name,
				X:       c.X,
				Y:       c.Y,
				IsStart: name == start,
				IsEnd:   name == end,
				Links:   append([]string{}, links...),
			}
			rooms = append(rooms, rj)
		}

		resp := dataJSON{Rooms: rooms, Moves: steps}
		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		_ = enc.Encode(resp)
	})

	fmt.Printf("server listening on %s (serving %s)\n", *addr, *webDir)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		panic(err)
	}
}
