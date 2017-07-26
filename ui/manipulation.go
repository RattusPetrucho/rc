package ui

import (
	"fmt"
	"strings"

	"github.com/RattusPetrucho/rc/iolib"
	"github.com/jroimartin/gocui"
)

// Окно ввода имени файла
func (w *Window) mkFileView(g *gocui.Gui, v *gocui.View) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("new_file", maxX/2-15, maxY/12, maxX/2+15, maxY/12+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = true
		v.Wrap = true
		g.Cursor = true

		if _, err := g.SetCurrentView("new_file"); err != nil {
			return err
		}

		v.Title = "New file"
	}

	return nil
}

// Отмена создани файла
func (w *Window) hideMkFileView(g *gocui.Gui, v *gocui.View) error {
	g.Cursor = false
	if err := g.DeleteView("new_file"); err != nil {
		return err
	}
	if _, err := g.SetCurrentView(w.cur_panel_name); err != nil {
		return err
	}
	return nil
}

// Подтверждение и создание файла
func (w *Window) mkFile(g *gocui.Gui, v *gocui.View) error {
	path := w.p.path

	name := v.Buffer()
	name = strings.TrimSpace(name)

	if err := iolib.CreateFile(path + "/" + name); err != nil {
		v.Clear()
		fmt.Fprintln(v, err)
		return nil
	}

	w.p.v.Clear()
	if err := w.p.openDir(path); err != nil {
		return err
	}
	w.p.line = 0
	w.p.findElem(name)

	if err := w.p.setCursor(); err != nil {
		return err
	}

	g.Cursor = false
	if err := g.DeleteView("new_file"); err != nil {
		return err
	}
	if _, err := g.SetCurrentView(w.cur_panel_name); err != nil {
		return err
	}
	return nil
}

// Окно ввода имени директории
func (w *Window) mkDirView(g *gocui.Gui, v *gocui.View) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("new_dir", maxX/2-15, maxY/12, maxX/2+15, maxY/12+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = true
		v.Wrap = true
		g.Cursor = true

		if _, err := g.SetCurrentView("new_dir"); err != nil {
			return err
		}

		v.Title = "New directory"
	}

	return nil
}

// Отмена создани директории
func (w *Window) hideMkDirView(g *gocui.Gui, v *gocui.View) error {
	g.Cursor = false
	if err := g.DeleteView("new_dir"); err != nil {
		return err
	}
	if _, err := g.SetCurrentView(w.cur_panel_name); err != nil {
		return err
	}
	return nil
}

// Подтверждение и создание директории
func (w *Window) mkDir(g *gocui.Gui, v *gocui.View) error {
	path := w.p.path

	name := v.Buffer()
	name = strings.TrimSpace(name)

	if err := iolib.MkDir(path + "/" + name); err != nil {
		v.Clear()
		fmt.Fprintln(v, err)
		return nil
	}

	w.p.v.Clear()
	if err := w.p.openDir(path); err != nil {
		return err
	}
	w.p.line = 0
	w.p.findElem(name)

	if err := w.p.setCursor(); err != nil {
		return err
	}

	g.Cursor = false
	if err := g.DeleteView("new_dir"); err != nil {
		return err
	}
	if _, err := g.SetCurrentView(w.cur_panel_name); err != nil {
		return err
	}
	return nil
}

// Окно подтверждения удаления элемента
func (w *Window) delElemView(g *gocui.Gui, v *gocui.View) error {
	if w.p.elems[w.p.line].name == ".." {
		return nil
	}
	maxX, maxY := g.Size()

	mess := " Do you really want to delete " + w.p.elems[w.p.line].name + "? "
	l := len([]rune(mess))

	if v, err := g.SetView("del_elem", maxX/2-(l/2), maxY/12, maxX/2+(l/2)+2, maxY/12+4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		fmt.Fprintln(v, "")
		fmt.Fprintln(v, "")
		fmt.Fprintln(v, mess)
		fmt.Fprintln(v, "")

		if _, err := g.SetCurrentView("del_elem"); err != nil {
			return err
		}

		v.Title = "Delete directory"
	}

	return nil
}

// Отмена удаления элемента
func (w *Window) hideDelElemView(g *gocui.Gui, v *gocui.View) error {
	g.Cursor = false
	if err := g.DeleteView("del_elem"); err != nil {
		return err
	}
	if _, err := g.SetCurrentView(w.cur_panel_name); err != nil {
		return err
	}
	return nil
}

// Подтверждение удаления элемента
func (w *Window) delElem(g *gocui.Gui, v *gocui.View) error {
	if err := iolib.Delete(w.p.elems[w.p.line].path); err != nil {
		v.Clear()
		fmt.Fprintln(v, err)
		return nil
	}

	w.p.v.Clear()
	if err := w.p.openDir(w.p.path); err != nil {
		return err
	}

	g.Cursor = false
	if err := g.DeleteView("del_elem"); err != nil {
		return err
	}
	if _, err := g.SetCurrentView(w.cur_panel_name); err != nil {
		return err
	}
	return nil
}

// Окно переименования элемента
func (w *Window) renameElemView(g *gocui.Gui, v *gocui.View) error {
	if w.p.elems[w.p.line].name == ".." {
		return nil
	}
	maxX, maxY := g.Size()

	if v, err := g.SetView("rename_elem", maxX/2-15, maxY/12, maxX/2+15, maxY/12+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = true
		v.Wrap = true
		g.Cursor = true

		fmt.Fprintln(v, w.p.elems[w.p.line].name)

		if _, err := g.SetCurrentView("rename_elem"); err != nil {
			return err
		}

		v.Title = "Rename"
	}

	return nil
}

// Отмена переименования элемента
func (w *Window) hideRenameElemView(g *gocui.Gui, v *gocui.View) error {
	g.Cursor = false
	if err := g.DeleteView("rename_elem"); err != nil {
		return err
	}
	if _, err := g.SetCurrentView(w.cur_panel_name); err != nil {
		return err
	}
	return nil
}

// Подтверждение и переименования элемента
func (w *Window) renameElem(g *gocui.Gui, v *gocui.View) error {
	name := v.Buffer()
	name = strings.TrimSpace(name)

	if err := iolib.Rename(w.p.elems[w.p.line].path, w.p.path+"/"+name); err != nil {
		v.Clear()
		fmt.Fprintln(v, err)
		return nil
	}

	if err := w.p.openDir(w.p.path); err != nil {
		return err
	}

	g.Cursor = false
	if err := g.DeleteView("rename_elem"); err != nil {
		return err
	}
	if _, err := g.SetCurrentView(w.cur_panel_name); err != nil {
		return err
	}
	return nil
}

// Поиск элемента в панели
func (p *panel) findElem(name string) {
	if p.line == len(p.elems)-1 {
		return
	}

	for i := p.line + 1; i < len(p.elems); i++ {
		if strings.Contains(strings.ToLower(p.elems[i].name), strings.ToLower(name)) {
			p.line = i
			return
		}
	}
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
