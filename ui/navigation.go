package ui

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/RattusPetrucho/rc/iolib"

	"github.com/jroimartin/gocui"
)

// Окно поиска элемента в панели
func (w *Window) showFindElementView(g *gocui.Gui, v *gocui.View) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("find_elem", maxX/2-15, maxY/12, maxX/2+15, maxY/12+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = true
		v.Wrap = true
		g.Cursor = true

		if _, err := g.SetCurrentView("find_elem"); err != nil {
			return err
		}

		v.Title = "Find"
	}

	return nil
}

// Отмена поиска элемента
func (w *Window) hideFindElementView(g *gocui.Gui, v *gocui.View) error {
	g.Cursor = false
	if err := g.DeleteView("find_elem"); err != nil {
		return err
	}
	if _, err := g.SetCurrentView(w.cur_panel_name); err != nil {
		return err
	}
	return nil
}

// Запуск поиска элемента
func (w *Window) startFindElementView(g *gocui.Gui, v *gocui.View) error {
	name := v.Buffer()
	name = strings.TrimSpace(name)

	w.p.findElem(name)

	if err := w.p.setCursor(); err != nil {
		return err
	}

	g.Cursor = false
	if err := g.DeleteView("find_elem"); err != nil {
		return err
	}
	if _, err := g.SetCurrentView(w.cur_panel_name); err != nil {
		return err
	}

	w.p.search = name

	return nil
}

// Продолжение поиска элемента в следующих элементах панели
func (p *panel) nextSearch(g *gocui.Gui, v *gocui.View) error {
	if p.search == "" {
		return nil
	}

	for i := p.line + 1; i < len(p.elems); i++ {
		if strings.Contains(strings.ToLower(p.elems[i].name), strings.ToLower(p.search)) {
			p.line = i
			if err := p.setCursor(); err != nil {
				return err
			}
			return nil
		}
	}

	return nil
}

// Продолжение поиска элемента в предыдущих элементах панели
func (p *panel) prevSearch(g *gocui.Gui, v *gocui.View) error {
	if p.search == "" {
		return nil
	}

	for i := p.line - 1; i > 0; i-- {
		if strings.Contains(strings.ToLower(p.elems[i].name), strings.ToLower(p.search)) {
			p.line = i
			if err := p.setCursor(); err != nil {
				return err
			}
			return nil
		}
	}

	return nil
}

// Окно быстрого перехода
func (w *Window) showJumpView(g *gocui.Gui, v *gocui.View) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("jump", maxX/2-15, maxY/12, maxX/2+15, maxY/12+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = true
		v.Wrap = true
		g.Cursor = true

		if _, err := g.SetCurrentView("jump"); err != nil {
			return err
		}

		v.Title = "Jump"
	}

	return nil
}

// Отмена быстрого перехода
func (w *Window) hideJumpView(g *gocui.Gui, v *gocui.View) error {
	g.Cursor = false
	if err := g.DeleteView("jump"); err != nil {
		return err
	}
	if _, err := g.SetCurrentView(w.cur_panel_name); err != nil {
		return err
	}
	return nil
}

// Запуск быстрого перехода
func (w *Window) startJump(g *gocui.Gui, v *gocui.View) error {
	name := v.Buffer()
	name = strings.TrimSpace(name)

	if err := w.p.openDir(name); err != nil {
		return err
	}

	w.p.line = 0
	w.p.path = name
	w.p.v.Title = name
	for key := range w.p.hystory {
		delete(w.p.hystory, key)
	}

	if err := w.p.setCursor(); err != nil {
		return err
	}

	g.Cursor = false
	if err := g.DeleteView("jump"); err != nil {
		return err
	}
	if _, err := g.SetCurrentView(w.cur_panel_name); err != nil {
		return err
	}
	return nil
}

// Быстрое перемещение в GOPATH
func (p *panel) goToGopath(g *gocui.Gui, v *gocui.View) error {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		return nil
	}
	if err := p.openDir(gopath); err != nil {
		return err
	}
	p.line = 0
	p.path = gopath
	p.v.Title = gopath
	for key := range p.hystory {
		delete(p.hystory, key)
	}

	if err := p.setCursor(); err != nil {
		return err
	}

	return nil
}

// Быстрое перемещение в HOME
func (p *panel) goToHome(g *gocui.Gui, v *gocui.View) error {
	home := os.Getenv("HOME")
	if home == "" {
		return nil
	}
	if err := p.openDir(home); err != nil {
		return err
	}
	p.line = 0
	p.path = home
	p.v.Title = home
	for key := range p.hystory {
		delete(p.hystory, key)
	}

	if err := p.setCursor(); err != nil {
		return err
	}

	return nil
}

// Выделение следующего элемента внизу
func (w *Window) cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		if w.p.line < len(w.p.elems)-1 {
			cx, cy := w.p.v.Cursor()
			if err := w.p.v.SetCursor(cx, cy+1); err != nil {
				ox, oy := w.p.v.Origin()
				if err := w.p.v.SetOrigin(ox, oy+1); err != nil {
					return err
				}
			}
		}
	}
	if w.p.line < len(w.p.elems)-1 {
		w.p.line++
	}
	return nil
}

// Выделение следующего элемента вверху
func (w *Window) cursorUp(g *gocui.Gui, v *gocui.View) error {
	if w.p.line > 0 {
		ox, oy := w.p.v.Origin()
		cx, cy := w.p.v.Cursor()
		if err := w.p.v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := w.p.v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	if w.p.line > 0 {
		w.p.line--
	}
	return nil
}

// Выделение нулевого элемента
func (w *Window) cursorPgUp(g *gocui.Gui, v *gocui.View) (err error) {
	if w.p.line == 0 {
		return nil
	}
	_, y := w.p.v.Size()
	cx, cy := w.p.v.Cursor()

	if len(w.p.elems)-1 <= y {
		if err = w.p.v.SetCursor(cx, 0); err != nil {
			return
		}
		w.p.line = 0
		return
	}

	if cy > 0 {
		if err = w.p.v.SetCursor(cx, 0); err != nil {
			return
		}
		w.p.line -= cy
		return
	}

	page := w.p.line / y
	ox, _ := w.p.v.Origin()
	if err = w.p.v.SetOrigin(ox, (page-1)*y); err != nil {
		return
	}
	w.p.line = (page - 1) * y

	return
}

// Выделение последнего элемента
func (w *Window) cursorPgDown(g *gocui.Gui, v *gocui.View) (err error) {
	if w.p.line == len(w.p.elems)-1 {
		return nil
	}
	_, y := w.p.v.Size()
	cx, _ := w.p.v.Cursor()

	if len(w.p.elems)-1 <= y {
		if err = w.p.v.SetCursor(cx, len(w.p.elems)-1); err != nil {
			return
		}
		w.p.line = len(w.p.elems) - 1
	} else {
		if w.p.line < y-1 {
			if err = w.p.v.SetCursor(cx, y-1); err != nil {
				return
			}
			w.p.line = y - 1
		} else {
			ox, oy := w.p.v.Origin()
			oy = ((oy / y) * y) + y
			if len(w.p.elems)-1 <= oy {
				if err = w.p.v.SetCursor(ox, y-(oy-(len(w.p.elems)-1))); err != nil {
					return
				}
				w.p.line = len(w.p.elems) - 1
			} else {
				if err = w.p.v.SetOrigin(ox, oy); err != nil {
					return
				}
				if len(w.p.elems)-1 < (oy+y)-1 {
					if err = w.p.v.SetCursor(ox, (len(w.p.elems)-oy)-1); err != nil {
						return
					}
					w.p.line = len(w.p.elems) - 1
				} else {
					if err = w.p.v.SetCursor(ox, y-1); err != nil {
						return
					}
					w.p.line = (oy + y) - 1
				}
			}
		}
	}

	return
}

// Открытие директории/файла
func (w *Window) enterLine(g *gocui.Gui, v *gocui.View) (err error) {
	if w.p.elems[w.p.line].is_dir {
		e_path := w.p.elems[w.p.line].path
		if err = w.p.openDir(e_path); err != nil {
			return
		}
	} else {
		if err = iolib.OpenFile(w.p.elems[w.p.line].path); err != nil {
			return
		}
	}

	return
}

// Переоткрытие директории/файла
func (w *Window) reopen(g *gocui.Gui, v *gocui.View) (err error) {
	e_path := w.p.path
	if err = w.p.openDir(e_path); err != nil {
		return
	}

	return
}

// Вывод элементов директории в панель с запоминанием текущего положения курсора
func (p *panel) openDir(dir string) error {
	cur_len := len(p.elems)
	var e_path, name string
	if cur_len > 0 {
		p.v.Clear()

		if dir != p.path {
			name = p.elems[p.line].name
			e_path = p.elems[p.line].path

			if name == ".." {
				line, ok := p.hystory[e_path]
				if ok {
					p.line = line
					delete(p.hystory, p.path)
				} else {
					p.line = 0
				}
			} else {
				p.hystory[p.path] = p.line
				p.line = 0
			}
		} else {
			e_path = p.path
		}
	}

	// Читаем данные папки
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	p.elems = nil

	if dir != "/" {
		elem := element{
			name:   "..",
			path:   path.Dir(dir),
			is_dir: true,
		}
		p.elems = append(p.elems, elem)
	}

	// Выбираем директории
	var dirs []element
	for _, val := range entries {
		if !p.show_hide {
			if strings.HasPrefix(val.Name(), ".") {
				continue
			}
		}
		if val.IsDir() {
			elem := element{
				name:   val.Name(),
				path:   dir + "/" + val.Name(),
				is_dir: val.IsDir(),
			}

			dirs = append(dirs, elem)
		}
	}
	sort.Sort(elements(dirs))

	// Выбираем файлы
	var files []element
	for _, val := range entries {
		if !p.show_hide {
			if strings.HasPrefix(val.Name(), ".") {
				continue
			}
		}
		if !val.IsDir() {
			elem := element{
				name:   val.Name(),
				path:   dir + "/" + val.Name(),
				is_dir: val.IsDir(),
			}

			files = append(files, elem)
		}
	}

	sort.Sort(elements(files))

	p.elems = append(p.elems, dirs...)
	p.elems = append(p.elems, files...)

	// Заполнение view
	x, y := p.v.Size()
	for _, val := range p.elems {
		l := x - len([]rune(val.name))

		value := ""
		if l <= 0 {
			value = string([]rune(val.name)[:x-3]) + "..."
		} else {
			value = val.name + strings.Repeat(" ", l)
		}

		fmt.Fprintln(p.v, value)
	}

	if dir != "/Users/Rattus" {
		if len(p.elems)-1 < y {
			for i := 0; i < y-len(p.elems)-1; i++ {
				fmt.Fprintln(p.v, strings.Repeat(" ", x))
			}
		}
	}

	if cur_len > 0 {
		if e_path == p.path {
			if len(p.elems) > p.line {
				if err = p.setCursor(); err != nil {
					return err
				}
			} else {
				p.line = len(p.elems) - 1

				if err = p.setCursor(); err != nil {
					return err
				}
			}
		} else {
			p.path = e_path
			p.v.Title = e_path

			if err = p.setCursor(); err != nil {
				return err
			}
		}
	}

	return nil
}
