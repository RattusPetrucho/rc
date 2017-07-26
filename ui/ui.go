package ui

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jroimartin/gocui"
)

var ErrQuit = gocui.ErrQuit

type Window struct {
	gui     *gocui.Gui // Основной объект gui
	l_panel *panel     // Левая панель
	r_panel *panel     // Правая панель

	p *panel // Активная панель

	cur_panel_name string // Текущая активная панель
}

type panel struct {
	v         *gocui.View    // Отображение панели
	path      string         // Текущий пути панели
	elems     []element      // Элементы внутри панели
	show_hide bool           // Отображать ли скрытые файлы
	line      int            // Текущая выделенная строка в панели
	hystory   map[string]int // История выделенных строк по предыдущим путям панели
	search    string         // Последняя маска поиска
}

type element struct {
	name    string    // Имя файла/директории
	path    string    // Абсолютный путь
	is_dir  bool      // Это директория или файл
	size    int64     // Размер
	modtime time.Time // Время последней модификации
}

// Создание нового окна с двумя панелями
func New() (w *Window, err error) {
	w = new(Window)
	w.l_panel = new(panel)
	w.r_panel = new(panel)
	w.p = w.l_panel
	w.l_panel.hystory = make(map[string]int)
	w.r_panel.hystory = make(map[string]int)
	w.cur_panel_name = "l_panel"

	w.gui, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return
	}

	w.gui.FgColor = gocui.ColorDefault
	w.gui.BgColor = gocui.ColorDefault

	w.gui.SetManagerFunc(w.layout)

	if err = w.keybindings(); err != nil {
		return
	}

	return
}

func (w *Window) MainLoop() error {
	return w.gui.MainLoop()
}

func (w *Window) Close() {
	w.gui.Close()
}

// Инициализация панелей
func (w *Window) layout(g *gocui.Gui) error {
	var err error
	maxX, maxY := g.Size()
	home := os.Getenv("HOME")
	if w.l_panel.v, err = g.SetView("l_panel", 0, 0, (maxX/2)-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		w.l_panel.v.Highlight = true
		w.l_panel.v.SelBgColor = gocui.ColorGreen
		w.l_panel.v.SelFgColor = gocui.ColorBlack

		w.l_panel.path = home
		w.l_panel.v.Title = home
		w.l_panel.openDir(home)

		if _, err := g.SetCurrentView("l_panel"); err != nil {
			return err
		}
	}

	if w.r_panel.v, err = g.SetView("r_panel", maxX/2, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		w.r_panel.v.Highlight = false
		w.r_panel.v.SelBgColor = gocui.ColorGreen
		w.r_panel.v.SelFgColor = gocui.ColorBlack

		w.r_panel.path = home
		w.r_panel.v.Title = home
		w.r_panel.openDir(home)
	}

	return nil
}

// Установка сочетаний клавишь
func (w *Window) keybindings() error {

	// Перемещение курсора
	if err := w.gui.SetKeybinding("l_panel", gocui.KeyArrowDown, gocui.ModNone, w.cursorDown); err != nil {
		return err
	}
	if err := w.gui.SetKeybinding("l_panel", gocui.KeyArrowUp, gocui.ModNone, w.cursorUp); err != nil {
		return err
	}
	if err := w.gui.SetKeybinding("r_panel", gocui.KeyArrowDown, gocui.ModNone, w.cursorDown); err != nil {
		return err
	}
	if err := w.gui.SetKeybinding("r_panel", gocui.KeyArrowUp, gocui.ModNone, w.cursorUp); err != nil {
		return err
	}
	if err := w.gui.SetKeybinding("l_panel", gocui.KeyPgup, gocui.ModNone, w.cursorPgUp); err != nil {
		return err
	}
	if err := w.gui.SetKeybinding("r_panel", gocui.KeyPgup, gocui.ModNone, w.cursorPgUp); err != nil {
		return err
	}
	if err := w.gui.SetKeybinding("l_panel", gocui.KeyPgdn, gocui.ModNone, w.cursorPgDown); err != nil {
		return err
	}
	if err := w.gui.SetKeybinding("r_panel", gocui.KeyPgdn, gocui.ModNone, w.cursorPgDown); err != nil {
		return err
	}

	// Открытие директории/файла
	if err := w.gui.SetKeybinding("l_panel", gocui.KeyEnter, gocui.ModNone, w.enterLine); err != nil {
		return err
	}
	if err := w.gui.SetKeybinding("r_panel", gocui.KeyEnter, gocui.ModNone, w.enterLine); err != nil {
		return err
	}

	// Переоткрытие директории
	if err := w.gui.SetKeybinding("l_panel", gocui.KeyCtrlR, gocui.ModNone, w.reopen); err != nil {
		return err
	}
	if err := w.gui.SetKeybinding("r_panel", gocui.KeyCtrlR, gocui.ModNone, w.reopen); err != nil {
		return err
	}

	// Поиск в списке элементов
	if err := w.gui.SetKeybinding("l_panel", gocui.KeySpace, gocui.ModNone, w.showFindElementView); err != nil {
		return err
	}
	if err := w.gui.SetKeybinding("r_panel", gocui.KeySpace, gocui.ModNone, w.showFindElementView); err != nil {
		return err
	}
	if err := w.gui.SetKeybinding("find_elem", gocui.KeyCtrlC, gocui.ModNone, w.hideFindElementView); err != nil {
		return err
	}
	if err := w.gui.SetKeybinding("find_elem", gocui.KeyEnter, gocui.ModNone, w.startFindElementView); err != nil {
		return err
	}
	if err := w.gui.SetKeybinding("l_panel", gocui.KeyArrowRight, gocui.ModNone, w.l_panel.nextSearch); err != nil {
		return err
	}
	if err := w.gui.SetKeybinding("r_panel", gocui.KeyArrowRight, gocui.ModNone, w.r_panel.nextSearch); err != nil {
		return err
	}
	if err := w.gui.SetKeybinding("l_panel", gocui.KeyArrowLeft, gocui.ModNone, w.l_panel.prevSearch); err != nil {
		return err
	}
	if err := w.gui.SetKeybinding("r_panel", gocui.KeyArrowLeft, gocui.ModNone, w.r_panel.prevSearch); err != nil {
		return err
	}

	// Быстрое перемещение
	if err := w.gui.SetKeybinding("l_panel", gocui.KeyCtrlZ, gocui.ModNone, w.showJumpView); err != nil {
		return err
	}
	if err := w.gui.SetKeybinding("r_panel", gocui.KeyCtrlZ, gocui.ModNone, w.showJumpView); err != nil {
		return err
	}
	if err := w.gui.SetKeybinding("jump", gocui.KeyCtrlC, gocui.ModNone, w.hideJumpView); err != nil {
		return err
	}
	if err := w.gui.SetKeybinding("jump", gocui.KeyEnter, gocui.ModNone, w.startJump); err != nil {
		return err
	}

	// Открытие информации по текущему элементу
	if err := w.gui.SetKeybinding("l_panel", gocui.KeyF2, gocui.ModNone, w.showInfo); err != nil {
		return err
	}
	if err := w.gui.SetKeybinding("r_panel", gocui.KeyF2, gocui.ModNone, w.showInfo); err != nil {
		return err
	}
	if err := w.gui.SetKeybinding("info", gocui.KeyCtrlC, gocui.ModNone, w.hideInfo); err != nil {
		return err
	}

	// Создание файла
	if err := w.gui.SetKeybinding("l_panel", gocui.KeyF3, gocui.ModNone, w.mkFileView); err != nil {
		return err
	}
	if err := w.gui.SetKeybinding("r_panel", gocui.KeyF3, gocui.ModNone, w.mkFileView); err != nil {
		return err
	}
	if err := w.gui.SetKeybinding("new_file", gocui.KeyCtrlC, gocui.ModNone, w.hideMkFileView); err != nil {
		return err
	}
	if err := w.gui.SetKeybinding("new_file", gocui.KeyEnter, gocui.ModNone, w.mkFile); err != nil {
		return err
	}

	// Создание директории
	if err := w.gui.SetKeybinding("l_panel", gocui.KeyF7, gocui.ModNone, w.mkDirView); err != nil {
		return err
	}
	if err := w.gui.SetKeybinding("r_panel", gocui.KeyF7, gocui.ModNone, w.mkDirView); err != nil {
		return err
	}
	if err := w.gui.SetKeybinding("new_dir", gocui.KeyCtrlC, gocui.ModNone, w.hideMkDirView); err != nil {
		return err
	}
	if err := w.gui.SetKeybinding("new_dir", gocui.KeyEnter, gocui.ModNone, w.mkDir); err != nil {
		return err
	}

	// Удаление элемента
	if err := w.gui.SetKeybinding("l_panel", gocui.KeyF8, gocui.ModNone, w.delElemView); err != nil {
		return err
	}
	if err := w.gui.SetKeybinding("r_panel", gocui.KeyF8, gocui.ModNone, w.delElemView); err != nil {
		return err
	}
	if err := w.gui.SetKeybinding("del_elem", gocui.KeyCtrlC, gocui.ModNone, w.hideDelElemView); err != nil {
		return err
	}
	if err := w.gui.SetKeybinding("del_elem", gocui.KeyEnter, gocui.ModNone, w.delElem); err != nil {
		return err
	}

	// Переименование элемента
	if err := w.gui.SetKeybinding("l_panel", gocui.KeyF6, gocui.ModNone, w.renameElemView); err != nil {
		return err
	}
	if err := w.gui.SetKeybinding("r_panel", gocui.KeyF6, gocui.ModNone, w.renameElemView); err != nil {
		return err
	}
	if err := w.gui.SetKeybinding("rename_elem", gocui.KeyCtrlC, gocui.ModNone, w.hideRenameElemView); err != nil {
		return err
	}
	if err := w.gui.SetKeybinding("rename_elem", gocui.KeyEnter, gocui.ModNone, w.renameElem); err != nil {
		return err
	}

	// Переключение между панелями
	if err := w.gui.SetKeybinding("l_panel", gocui.KeyTab, gocui.ModNone, w.switchPanel); err != nil {
		return err
	}
	if err := w.gui.SetKeybinding("r_panel", gocui.KeyTab, gocui.ModNone, w.switchPanel); err != nil {
		return err
	}

	// Окно с помощью по hotkeys
	if err := w.gui.SetKeybinding("", gocui.KeyF1, gocui.ModNone, showHelp); err != nil {
		return err
	}
	if err := w.gui.SetKeybinding("help", gocui.KeyCtrlC, gocui.ModNone, w.hideHelp); err != nil {
		return err
	}

	// Выход из приложения
	if err := w.gui.SetKeybinding("", gocui.KeyF10, gocui.ModNone, quit); err != nil {
		return err
	}

	return nil
}

// Переключение панелей
func (w *Window) switchPanel(g *gocui.Gui, v *gocui.View) error {
	switch w.cur_panel_name {
	case "l_panel":
		if _, err := g.SetCurrentView("r_panel"); err != nil {
			return err
		}

		w.r_panel.v.Highlight = true
		w.l_panel.v.Highlight = false
		w.cur_panel_name = "r_panel"
		w.p = w.r_panel
	case "r_panel":
		if _, err := g.SetCurrentView("l_panel"); err != nil {
			return err
		}

		w.r_panel.v.Highlight = false
		w.l_panel.v.Highlight = true
		w.cur_panel_name = "l_panel"
		w.p = w.l_panel
	}

	return nil
}

// Выход из приложения
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

// Отображение окна помощи
func showHelp(g *gocui.Gui, v *gocui.View) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("help", maxX/2-30, maxY/12, maxX/2+30, maxY/12+17); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		fmt.Fprintln(v, " * F1 - Help. Это окно.")
		fmt.Fprintln(v, " * F2 - Выводит информацию о выбранном элементе.")
		fmt.Fprintln(v, " * F3 - Создание нового файла.")
		fmt.Fprintln(v, " * F4 - Редактирование файла в vim.")
		fmt.Fprintln(v, " * F5 - Копирование.")
		fmt.Fprintln(v, " * F5 - Переименование текущего элемента.")
		fmt.Fprintln(v, " * F7 - Создание директории.")
		fmt.Fprintln(v, " * F8 - Удаление текущего элемента.")
		fmt.Fprintln(v, " * F10 - Выход из rc.")
		fmt.Fprintln(v, " * Space - Поиск элемента в текущей панели.")
		fmt.Fprintln(v, " * Enter - Входит в папку либо открывает файл.")
		fmt.Fprintln(v, "    - Видео файлы открываются с помощью mplayer.")
		fmt.Fprintln(v, "    - Текстовые файлы открываются с помощью Sublime Text.")
		fmt.Fprintln(v, "    - Pdf и изображения открываются с помощью preview.")
		fmt.Fprintln(v, "")
		fmt.Fprintln(v, "        Чтобы закрыть модальное окно нажмите Ctrl+C.")

		if _, err := g.SetCurrentView("help"); err != nil {
			return err
		}
		v.Title = "Help"
	}

	return nil
}

// Скрытие окна помощи
func (w *Window) hideHelp(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("help"); err != nil {
		return err
	}
	if _, err := g.SetCurrentView(w.cur_panel_name); err != nil {
		return err
	}
	return nil
}

// Отображение информации по выбранному элементу
func (w *Window) showInfo(g *gocui.Gui, v *gocui.View) error {
	if w.p.elems[w.p.line].name == ".." {
		return nil
	}
	maxX, maxY := g.Size()

	if v, err := g.SetView("info", maxX/2-30, maxY/12, maxX/2+30, maxY/12+5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		fmt.Fprintln(v, " * Имя - "+w.p.elems[w.p.line].name)
		size := w.p.elems[w.p.line].size
		if size < 1024 {
			fmt.Fprintln(v, " * Размер - "+strconv.FormatInt(size, 10)+"B")
		} else if size/1024 < 1024 {
			fmt.Fprintln(v, " * Размер - "+strconv.FormatInt(size/1024, 10)+"KB")
		} else if size/(1024*1024) < 1024 {
			fmt.Fprintln(v, " * Размер - "+strconv.FormatInt(size/(1024*1024), 10)+"MB")
		} else if size/(1024*1024*1024) < 1024 {
			fmt.Fprintln(v, " * Размер - "+strconv.FormatInt(size/(1024*1024*1024), 10)+"GB")
		} else if size/(1024*1024*1024*1024) < 1024 {
			fmt.Fprintln(v, " * Размер - "+strconv.FormatInt(size/(1024*1024*1024*1024), 10)+"TB")
		}
		if w.p.elems[w.p.line].is_dir {
			fmt.Fprintln(v, " * Тип - директория")
		} else {
			fmt.Fprintln(v, " * Тип - файл")
		}
		fmt.Fprintln(v, " * Дата модификации -", w.p.elems[w.p.line].modtime.Format("15:04 02.01.2006 MST"))

		if _, err := g.SetCurrentView("info"); err != nil {
			return err
		}
		v.Title = "Info"
	}

	return nil
}

// Скрытие окна информации по выбранному элементу
func (w *Window) hideInfo(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("info"); err != nil {
		return err
	}
	if _, err := g.SetCurrentView(w.cur_panel_name); err != nil {
		return err
	}
	return nil
}
