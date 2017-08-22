package ui

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/RattusPetrucho/rc/iolib"

	"github.com/jroimartin/gocui"
)

// Отображение информации по выбранному элементу
func (w *Window) showInfo(g *gocui.Gui, v *gocui.View) error {
	if w.p.elems[w.p.line].name == ".." {
		return nil
	}

	maxX, maxY := g.Size()

	if v, err := g.SetView("info", maxX/2-30, maxY/12, maxX/2+30, maxY/12+6); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		inf, err := os.Stat(w.p.elems[w.p.line].path)
		if err != nil {
			return err
		}

		if w.p.elems[w.p.line].is_dir {
			PrintInfo(v, w.p.elems[w.p.line].is_dir, "Calculate size...", w.p.elems[w.p.line].name, inf.Mode().String(), inf.ModTime())

			w.info_done = make(chan struct{})
			go w.collectInfo(v, w.p.elems[w.p.line].is_dir, w.p.elems[w.p.line].path, w.p.elems[w.p.line].name, inf.Mode().String(), inf.ModTime())
		} else {
			PrintInfo(v, w.p.elems[w.p.line].is_dir, iolib.SizeIntToString(inf.Size()), w.p.elems[w.p.line].name, inf.Mode().String(), inf.ModTime())
		}

		if _, err := g.SetCurrentView("info"); err != nil {
			return err
		}
		v.Title = "Info"
	}

	return nil
}

// Скрытие окна информации по выбранному элементу
func (w *Window) hideInfo(g *gocui.Gui, v *gocui.View) error {
	if !w.cancelled_info() {
		close(w.info_done)
	}
	if err := g.DeleteView("info"); err != nil {
		return err
	}
	if _, err := g.SetCurrentView(w.cur_panel_name); err != nil {
		return err
	}

	runtime.GC()

	return nil
}

// Проверка отмены обхода каталогов для получения размера содержимого
func (w *Window) cancelled_info() bool {
	select {
	case <-w.info_done:
		return true
	default:
		return false
	}
}

// Вывод информации об элементе
func PrintInfo(v *gocui.View, is_dir bool, size, name, permissions string, date time.Time) {
	fmt.Fprintln(v, " * Имя:              "+name)

	if is_dir {
		fmt.Fprintln(v, " * Тип:              Директория")
	} else {
		fmt.Fprintln(v, " * Тип:              Файл")
	}

	fmt.Fprintln(v, " * Размер:           "+size)

	fmt.Fprintln(v, " * Permissions:     ", permissions)

	fmt.Fprintln(v, " * Дата модификации:", date.Format("15:04 02.01.2006 MST"))
}

// Сбор информации о селержимом каталога
func (w *Window) collectInfo(v *gocui.View, is_dir bool, path, name, permissions string, date time.Time) {
	var total_size int64
	filesize := make(chan int64)

	var wg sync.WaitGroup
	wg.Add(1)
	go iolib.GetFilesSizeInDir(path, &wg, filesize, w.info_done)

	go func(w *Window) {
		wg.Wait()
		close(filesize)
	}(w)

	tick := time.NewTicker(500 * time.Millisecond)
	defer tick.Stop()

loop:
	for {
		select {
		case <-w.info_done:
			for range filesize {
			}
			fmt.Fprintln(w.r_panel.v, runtime.NumGoroutine())
			return
		case size, ok := <-filesize:
			if !ok {
				if !w.cancelled_info() {
					close(w.info_done)
				}
				break loop
			}
			total_size += size
		case <-tick.C:
			go w.gui.Execute(func(g *gocui.Gui) error {
				v.Clear()
				PrintInfo(v, is_dir, "..."+iolib.SizeIntToString(total_size), name, permissions, date)
				return nil
			})

		}
	}

	go w.gui.Execute(func(g *gocui.Gui) error {
		v.Clear()
		PrintInfo(v, is_dir, iolib.SizeIntToString(total_size), name, permissions, date)
		return nil
	})
}
