package ui

import "github.com/jroimartin/gocui"

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

	// Быстрое перемещение в GOPATH
	if err := w.gui.SetKeybinding("l_panel", gocui.KeyCtrlG, gocui.ModNone, w.l_panel.goToGopath); err != nil {
		return err
	}
	if err := w.gui.SetKeybinding("r_panel", gocui.KeyCtrlG, gocui.ModNone, w.r_panel.goToGopath); err != nil {
		return err
	}

	// Быстрое перемещение в HOME
	if err := w.gui.SetKeybinding("l_panel", gocui.KeyCtrlH, gocui.ModNone, w.l_panel.goToHome); err != nil {
		return err
	}
	if err := w.gui.SetKeybinding("r_panel", gocui.KeyCtrlH, gocui.ModNone, w.r_panel.goToHome); err != nil {
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
	if err := w.gui.SetKeybinding("l_panel", gocui.KeyF1, gocui.ModNone, showHelp); err != nil {
		return err
	}
	if err := w.gui.SetKeybinding("r_panel", gocui.KeyF1, gocui.ModNone, showHelp); err != nil {
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
