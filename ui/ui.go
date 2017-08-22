package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/jroimartin/gocui"
)

var ErrQuit = gocui.ErrQuit

type Window struct {
	gui     *gocui.Gui // Основной объект gui
	l_panel *panel     // Левая панель
	r_panel *panel     // Правая панель

	p *panel // Активная панель

	cur_panel_name string // Текущая активная панель

	info_done chan struct{}
}

type panel struct {
	show_hide bool           // Отображать ли скрытые файлы
	line      int            // Текущая выделенная строка в панели
	v         *gocui.View    // Отображение панели
	path      string         // Текущий пути панели
	search    string         // Последняя маска поиска
	elems     []element      // Элементы внутри панели
	hystory   map[string]int // История выделенных строк по предыдущим путям панели
}

type element struct {
	is_dir bool   // Это директория или файл
	name   string // Имя файла/директории
	path   string // Абсолютный путь
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

	if v, err := g.SetView("help", maxX/2-30, maxY/14, maxX/2+30, maxY/14+23); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		fmt.Fprintln(v, " * F1 - Help. Это окно.")
		fmt.Fprintln(v, " * F2 - Выводит информацию о выбранном элементе.")
		fmt.Fprintln(v, " * F3 - Создание нового файла.")
		fmt.Fprintln(v, " * F5 - Копирование.")
		fmt.Fprintln(v, " * F6 - Переименование текущего элемента.")
		fmt.Fprintln(v, " * F7 - Создание директории.")
		fmt.Fprintln(v, " * F8 - Удаление текущего элемента.")
		fmt.Fprintln(v, " * F10 - Выход из rc.")
		fmt.Fprintln(v, strings.Repeat("-", 60))
		fmt.Fprintln(v, " * Space - Поиск элемента в текущей панели.")
		fmt.Fprintln(v, "      -> - Next search.")
		fmt.Fprintln(v, "      <- - Prev search.")
		fmt.Fprintln(v, strings.Repeat("-", 60))
		fmt.Fprintln(v, " * Ctrl+G - Переход в GOPATH.")
		fmt.Fprintln(v, " * Ctrl+H - Переход в домашнюю директорию.")
		fmt.Fprintln(v, strings.Repeat("-", 60))
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
