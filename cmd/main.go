package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	lib "lem-in/helpers"
)

// Graph — карта: имя комнаты -> список соседей
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
	g[to] = append(g[to], from)
}

// parseG — читает и валидирует файл, возвращает граф, n, start, end.
func parseG(fileName string) (Graph, int, string, string) {
	g := Graph{}

	data, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println("Ошибка чтения файла:", err)
		os.Exit(1)
	}
	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	lines := strings.Split(text, "\n")

	if len(lines) == 0 {
		fmt.Println("No data in graph.txt")
		os.Exit(1)
	}

	// защититься от возможного '\r' в конце первой строки
	if len(lines[0]) > 0 && lines[0][len(lines[0])-1] == '\r' {
		lines[0] = lines[0][:len(lines[0])-1]
	}

	// первая строка — количество муравьёв
	n, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		fmt.Println("Invalid number of ants:", lines[0])
		os.Exit(1)
	}
	// теперь требуем строго положительное число (>0)
	if n <= 0 {
		fmt.Println("Number of ants must be a positive integer")
		os.Exit(1)
	}

	start, end := "", ""
	flag := "room"

	// парсим остальные строки
	for _, raw := range lines[1:] {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}

		// пропускаем комментарии одной решётки: "#something"
		// но обрабатываем директивы "##start" и "##end"
		if strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "##") {
			continue
		}

		// директивы типа "##start" / "##end"
		if strings.HasPrefix(line, "##") {
			switch strings.TrimSpace(line) {
			case "##start":
				flag = "start"
				continue
			case "##end":
				flag = "end"
				continue
			default:
				continue
			}
		}

		// описание комнаты: name X Y
		parts := strings.Fields(line)
		if len(parts) == 3 {
			roomName := parts[0]
			if roomName == "" {
				fmt.Println("Invalid room name in line:", line)
				os.Exit(1)
			}
			if _, err := strconv.Atoi(parts[1]); err != nil {
				fmt.Println("Invalid X coordinate for room", roomName)
				os.Exit(1)
			}
			if _, err := strconv.Atoi(parts[2]); err != nil {
				fmt.Println("Invalid Y coordinate for room", roomName)
				os.Exit(1)
			}

			lib.RoomNames = append(lib.RoomNames, roomName)
			lib.Matrix = append(lib.Matrix, []string{parts[1], parts[2]})

			g.addRoom(roomName)

			if flag == "start" {
				start = roomName
				flag = "room"
			} else if flag == "end" {
				end = roomName
				flag = "room"
			}
			continue
		}

		// описание связи: A-B
		if strings.Contains(line, "-") && !strings.Contains(line, " ") {
			room1, room2 := lib.ParseLink(line)

			if room1 == "" || room2 == "" {
				fmt.Printf("Invalid link format: %s\n", line)
				os.Exit(1)
			}
			if room1 == room2 {
				fmt.Printf("Invalid link: room %s cannot be linked to itself\n", room1)
				os.Exit(1)
			}
			if _, ok := g[room1]; !ok {
				fmt.Printf("Invalid link: room %s is not defined\n", room1)
				os.Exit(1)
			}
			if _, ok := g[room2]; !ok {
				fmt.Printf("Invalid link: room %s is not defined\n", room2)
				os.Exit(1)
			}
			if lib.Contains(g[room1], room2) {
				fmt.Printf("Invalid link: duplicate link %s-%s\n", room1, room2)
				os.Exit(1)
			}

			g.addLink(room1, room2)
			continue
		}

		// иначе неверный формат
		fmt.Printf("Wrong format in graph.txt near line: %s\n", line)
		os.Exit(1)
	}

	// start и end обязаны быть заданы
	if start == "" {
		fmt.Println("Start room not defined")
		os.Exit(1)
	}
	if end == "" {
		fmt.Println("End room not defined")
		os.Exit(1)
	}

	return g, n, start, end
}

// простой DFS, возвращает все простые пути от start до end
func (g Graph) getBestGroup(start, end string) [][]string {
	ans := [][]string{}
	path := []string{start}

	var dfs func(string)
	dfs = func(current string) {
		if current == end {
			ans = append(ans, append([]string{}, path...))
			return
		}
		for _, neighbor := range g[current] {
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

// распределяем муравьёв: создаём набор групп путей и выбираем лучший по высоте.
// Добавлена сортировка путей для детерминизма.
func (g Graph) sendTheAnts(n int, start, end string) [][]string {
	paths := g.getBestGroup(start, end)
	if len(paths) == 0 {
		fmt.Println("No paths found from start to end.")
	}

	// сортируем пути: сначала по длине (короче лучше), потом лексикографически
	sort.Slice(paths, func(i, j int) bool {
		if len(paths[i]) != len(paths[j]) {
			return len(paths[i]) < len(paths[j])
		}
		return strings.Join(paths[i], ",") < strings.Join(paths[j], ",")
	})

	// теперь убираем старт и конец (как у вас было)
	for i, p := range paths {
		if len(p) >= 2 {
			paths[i] = p[1 : len(p)-1]
		} else {
			paths[i] = []string{}
		}
	}

	// ваша исходная логика объединения путей по конфликтам
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
			if !lib.Conflicts(path2, paths[ii]) {
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

// вспомогательные функции
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
	if len(os.Args) < 2 {
		fmt.Println("No input file specified.")
		return
	}
	fileName := os.Args[1]
	if !strings.HasSuffix(fileName, ".txt") {
		fmt.Println("Input file must have a .txt extension.")
		return
	}

	// сначала парсим и валидируем — если есть ошибка, parseG сделает os.Exit(1)
	g, n, start, end := parseG(fileName)

	// после успешного парсинга — выводим исходный файл (количество муравьёв + остальное)
	dataRaw, err := os.ReadFile(fileName)
	if err == nil {
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

	// симуляция и печать шагов (ваша логика сохранена)
	bestGroup := g.sendTheAnts(n, start, end)
	height, combined := getHeight(bestGroup, n)
	heights := getPathHeights(bestGroup)
	antProxy := getAntHeights(combined, heights)

	counter := 0
	ants := make([]lib.Ant, 0)

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
		for i := range ants {
			fmt.Printf("%s-%s ", ants[i].Name, ants[i].Path[ants[i].Current])
			structAnt = append(structAnt, ants[i])
			ants[i].Current++
		}
		for i := len(ants) - 1; i >= 0; i-- {
			if ants[i].Current >= len(ants[i].Path) {
				ants = append(ants[:i], ants[i+1:]...)
			}
		}
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

	for len(ants) > 0 || n > 0 {
		structAnt := []lib.Ant{}
		for i := range ants {
			if ants[i].Current < len(ants[i].Path) {
				fmt.Printf("%s-%s ", ants[i].Name, ants[i].Path[ants[i].Current])
				structAnt = append(structAnt, ants[i])
				ants[i].Current++
			}
		}
		for i := len(ants) - 1; i >= 0; i-- {
			if ants[i].Current >= len(ants[i].Path) {
				ants = append(ants[:i], ants[i+1:]...)
			}
		}
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
		if len(ants) == 0 && n > 0 {
			fmt.Printf("Warning: %d ants couldn't be sent (no available paths)\n", n)
			break
		}
	}
}
