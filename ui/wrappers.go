package ui

import (
	"fmt"
	"strings"
)

// Устанавливает активным элемент согласно заданной линии панели
func (p *panel) setCursor() (err error) {
	_, y := p.v.Size()
	cx, _ := p.v.Cursor()
	ox, _ := p.v.Origin()

	page := p.line / y

	if err = p.v.SetOrigin(ox, page*y); err != nil {
		fmt.Println("yopt0", page)
		return
	}
	if err = p.v.SetCursor(cx, p.line-(page*y)); err != nil {
		fmt.Println("yopt1", page)
		return
	}
	return
}

// Тип и функции для сортировки элементов панели
type elements []element

func (e elements) Len() int { return len(e) }

func (e elements) Swap(i, j int) { e[i], e[j] = e[j], e[i] }

func (e elements) Less(i, j int) bool {
	return strings.ToLower(e[i].name)[0] < strings.ToLower(e[j].name)[0]
}
