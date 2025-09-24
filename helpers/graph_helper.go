package helpers

import (
	"image"
	"log"
	"os"
	"strconv"

	// cmd "lem-in/cmd"
	// hp "lem-in/helpers"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type Room struct {
	Name string
	X, Y int
}

var maxLenName int

type CustomParagraph struct {
	name string
	*widgets.Paragraph
	ShowTop    bool
	ShowLeft   bool
	ShowRight  bool
	ShowBottom bool
}

func NewP(name string, Current string, x, y int) *CustomParagraph {
	p := NewCustomParagraph(1)
	p.Text = name + " " + Current // + " " + strconv.Itoa(x) + " " + strconv.Itoa(y)
	p.SetRect(x, y, x+10+(maxLenName), y+3)
	p.Border = true
	return p
}

func GraphVisualization(roomNames []string, matrix [][]string, g []VizAnt, start, end string) {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	// Example room data: name and coordinates
	rooms := []Room{}
	roomMap := map[string]Room{}
	minX, maxX := 5, 5
	minY, maxY := 5, 5
	for i, k := range roomNames {
		co1, _ := strconv.Atoi(matrix[i][0])
		if co1 < minX {
			minX = co1
		}
		if co1 > maxX {
			maxX = co1
		}
		co2, _ := strconv.Atoi(matrix[i][1])
		if co2 < minY {
			minY = co2
		}
		if co2 > maxY {
			maxY = co2
		}
		room := Room{k, co1, co2}
		rooms = append(rooms, room)
		// Map name to Room struct
	}

	var roomWidgets []ui.Drawable
	roomNum := 0
	// targetW, targetH := 77, 77
	termW, termH := ui.TerminalDimensions()
	paraW, paraH := 15, 3
	targetW := termW - paraW
	targetH := termH - paraH
	margin := 5
	maxLenName := 0
	for range rooms {
		if len(rooms[roomNum].Name) > maxLenName {
			maxLenName = len(rooms[roomNum].Name)
		}
	}
	for i, k := range rooms {
		// Avoid division by zero
		if maxX != minX {
			rooms[i].X = margin + (rooms[i].X-minX)*(targetW-margin*2)/(maxX-minX) // was targetW-1
		} else {
			rooms[i].X = targetW / 2
		}
		if maxY != minY {
			rooms[i].Y = margin + (rooms[i].Y-minY)*(targetH-margin*2)/(maxY-minY) // was targetH-1
		} else {
			rooms[i].Y = targetH / 2
		}

		roomMap[k.Name] = Room{k.Name, rooms[i].X, rooms[i].Y}
		// Create a new paragraph for the room
		x := rooms[i].X
		y := rooms[i].Y
		j := widgets.NewParagraph()
		j.Text = rooms[roomNum].Name
		// if len(rooms[roomNum].Name) > maxLenName {
		// 	maxLenName = len(rooms[roomNum].Name)
		// }
		j.SetRect(x, y, x+10+(maxLenName), y+3)
		roomWidgets = append(roomWidgets, j)
		roomNum++

	}
	widgetsList := [][]ui.Drawable{}

	for _, k := range g {
		s := len(k.So)
		linesDraw := []ui.Drawable{}
		Current := []ui.Drawable{}
		for _, kk := range k.So { // need Path[] for loop only ? outside of g[]
			str2 := Room{}
			if kk.Current == 0 {
				if tempStr2, ok := roomMap[start]; ok {
					str2 = tempStr2
				} else {
					log.Println("Start room not found in roomMap", roomMap[start])
					os.Exit(1)
				}
			} else {
				if tempStr2, ok := roomMap[kk.Path[kk.Current-1]]; ok {
					str2 = tempStr2
				} else {
					log.Println("Room not found in roomMap:", kk.Path[kk.Current-1])
					os.Exit(1)
				}
			}
			getXY := 0
			if str1, ok := roomMap[kk.Path[kk.Current]]; ok {
				lastCenterX := str1.X + 15/2
				lastCenterY := str1.Y + 3/2 + 2
				// Center of Current paragraph
				currCenterX := str2.X + 15/2
				currCenterY := str2.Y + 3/2
				if lastCenterY == currCenterY { // || currCenterY-lastCenterY < 5 || lastCenterX-currCenterX < 5 {
					linesDraw = append(linesDraw, DrawHorizontalLine(
						min(lastCenterX, currCenterX),
						lastCenterY,
						max(lastCenterX, currCenterX),
					))
				} else if lastCenterX == currCenterX { // || currCenterY-lastCenterY < 5 || lastCenterX-currCenterX < 5 {
					linesDraw = append(linesDraw, DrawVerticalLine(
						lastCenterX,
						min(lastCenterY, currCenterY),
						max(lastCenterY, currCenterY),
					))
				} else {
					midX := (lastCenterX + currCenterX) / 2
					midY := (lastCenterY + currCenterY) / 2

					// Horizontal from lastCenterX to midX at lastCenterY
					linesDraw = append(linesDraw, DrawHorizontalLine(
						min(lastCenterX, midX),
						lastCenterY,
						max(lastCenterX, midX),
					))
					// Vertical from lastCenterY to midY at midX
					linesDraw = append(linesDraw, DrawVerticalLine(
						midX,
						min(lastCenterY, midY),
						max(lastCenterY, midY),
					))
					// Horizontal from midX to currCenterX at midY
					linesDraw = append(linesDraw, DrawHorizontalLine(
						min(midX, currCenterX),
						midY,
						max(midX, currCenterX),
					))
					// Vertical from midY to currCenterY at currCenterX
					linesDraw = append(linesDraw, DrawVerticalLine(
						currCenterX,
						min(midY, currCenterY),
						max(midY, currCenterY),
					))
					// // Draw horizontal then vertical (L-shape)
					// linesDraw = append(linesDraw, DrawHorizontalLine(
					// 	min(lastCenterX, currCenterX),
					// 	lastCenterY,
					// 	max(lastCenterX, currCenterX),
					// ))
					// linesDraw = append(linesDraw, DrawVerticalLine(
					// 	currCenterX,
					// 	min(lastCenterY, currCenterY),
					// 	max(lastCenterY, currCenterY),
					// ))
				}
			}
			getXY++

		}
		k.Index++
		num := 0
		lastX, lastY := -10, -10
		for num < s {
			if str, ok := roomMap[k.So[num].Path[k.So[num].Current]]; ok {

				if lastX == str.X && lastY == str.Y {
					Current[len(Current)-1].(*CustomParagraph).Text += "\n" + k.So[num].Name + " " + k.So[num].Path[k.So[num].Current]
					currXY := Current[len(Current)-1].(*CustomParagraph).GetRect()
					Current[len(Current)-1].(*CustomParagraph).SetRect(currXY.Min.X, currXY.Min.Y, currXY.Max.X, currXY.Max.Y+1)
					Current[len(Current)-1].(*CustomParagraph).name = k.So[num].Name
					num++
					continue
				}
				lastX, lastY = str.X, str.Y
				if k.So[num].Name != "" && str.X > 0 && str.Y > 0 {
					x := str.X
					y := str.Y
					Current = append(Current, NewP(k.So[num].Name, k.So[num].Path[k.So[num].Current], x, y))
					Current[len(Current)-1].(*CustomParagraph).name = k.So[num].Name
				}
			}
			num++
		}
		widgetsList = append(widgetsList, linesDraw)
		widgetsList = append(widgetsList, Current)
	}
	selected := 0
	for {
		ui.Clear()
		for _, group := range roomWidgets {
			ui.Render(group)
		}

		for i, group := range widgetsList {

			for _, w := range group {
				if p, ok := w.(*CustomParagraph); ok {
					if i == selected {
						p.TextStyle.Bg = ui.ColorYellow // Highlight selected
					} else {
						p.TextStyle.Bg = ui.ColorBlack
					}
				} else if p, ok := w.(*widgets.Paragraph); ok {
					if i == selected {
						p.TextStyle.Bg = ui.ColorYellow
					} else {
						p.TextStyle.Bg = ui.ColorWhite
					}
				}
			}

			if i <= selected {
				ui.Render(group...)
			} else {
			}
		}
		e := <-ui.PollEvents()
		if e.Type == ui.KeyboardEvent {
			switch e.ID {
			case "<Right>", "j":
				if selected == len(widgetsList)-1 {
					selected = 0
				}
				selected = (selected + 1) % len(widgetsList)
			case "<Left>", "k":
				selected = (selected - 1 + len(widgetsList)) % len(widgetsList)
			case "q", "<C-c>":
				return
			}
		}
	}
}

// ...existing code...

func DrawHorizontalLine(x0, y, x1 int) *widgets.Paragraph {
	p := widgets.NewParagraph()
	line := ""
	for i := 0; i < x1-x0; i++ {
		line += "─"
	}
	p.Text = line
	p.Border = false
	p.SetRect(x0, y, x1, y+1)
	return p
}

// Vertical line between two paragraphs
func DrawVerticalLine(x, y0, y1 int) *widgets.Paragraph {
	p := widgets.NewParagraph()
	line := ""
	for i := 0; i < y1-y0; i++ {
		line += "│\n"
	}
	p.Text = line
	p.Border = false
	p.SetRect(x, y0, x+1, y1)
	return p
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func NewCustomParagraph(num int) *CustomParagraph {
	p := widgets.NewParagraph()
	if num == 1 {
		return &CustomParagraph{
			Paragraph:  p,
			ShowTop:    false,
			ShowLeft:   false,
			ShowRight:  false,
			ShowBottom: false,
		}
	}
	return &CustomParagraph{}
}

func (cp *CustomParagraph) Draw(buf *ui.Buffer) {
	cp.Paragraph.Block.Border = false // Disable default border

	// Draw the paragraph content
	cp.Paragraph.Draw(buf)

	// Get rectangle
	r := cp.Paragraph.Inner
	x0, y0, x1, y1 := r.Min.X-1, r.Min.Y-1, r.Max.X, r.Max.Y

	// Draw custom borders
	if cp.ShowTop {
		for x := x0; x < x1; x++ {
			buf.SetCell(ui.NewCell('─', cp.Paragraph.BorderStyle), image.Pt(x, y0))
		}
	}
	if cp.ShowLeft {
		for y := y0; y < y1; y++ {
			buf.SetCell(ui.NewCell('│', cp.Paragraph.BorderStyle), image.Pt(x0, y))
		}
	}
	if cp.ShowRight {
		for y := y0; y < y1; y++ {
			buf.SetCell(ui.NewCell('│', cp.Paragraph.BorderStyle), image.Pt(x1-1, y))
		}
	}
	if cp.ShowBottom {
		for x := x0; x < x1; x++ {
			buf.SetCell(ui.NewCell('─', cp.Paragraph.BorderStyle), image.Pt(x, y1-1))
		}
	}
	// Draw corners if needed (optional)
	if cp.ShowTop && cp.ShowLeft {
		buf.SetCell(ui.NewCell('┌', cp.Paragraph.BorderStyle), image.Pt(x0, y0))
	}
	if cp.ShowTop && cp.ShowRight {
		buf.SetCell(ui.NewCell('┐', cp.Paragraph.BorderStyle), image.Pt(x1-1, y0))
	}
	if cp.ShowBottom && cp.ShowLeft {
		buf.SetCell(ui.NewCell('└', cp.Paragraph.BorderStyle), image.Pt(x0, y1-1))
	}
	if cp.ShowBottom && cp.ShowRight {
		buf.SetCell(ui.NewCell('┘', cp.Paragraph.BorderStyle), image.Pt(x1-1, y1-1))
	}
}
