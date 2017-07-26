package iolib

import (
	"os"
	"os/exec"
	"path"
	"strings"
)

var video = map[string]struct{}{
	".avi":  struct{}{},
	".mkv":  struct{}{},
	".vob":  struct{}{},
	".mp4":  struct{}{},
	".m4v":  struct{}{},
	".wmv":  struct{}{},
	".flv":  struct{}{},
	".mpg":  struct{}{},
	".mpeg": struct{}{},
}
var source = map[string]struct{}{
	".pl":   struct{}{},
	".pm":   struct{}{},
	".py":   struct{}{},
	".sh":   struct{}{},
	".js":   struct{}{},
	".html": struct{}{},
	".cpp":  struct{}{},
	".hpp":  struct{}{},
	".c":    struct{}{},
	".h":    struct{}{},
	".rb":   struct{}{},
	".erb":  struct{}{},
	".txt":  struct{}{},
	".go":   struct{}{},
	".sql":  struct{}{},
	".xml":  struct{}{},
	".json": struct{}{},
}
var img = map[string]struct{}{
	".jpg":  struct{}{},
	".jpeg": struct{}{},
	".png":  struct{}{},
	".gif":  struct{}{},
}

// Открывает файл
func OpenFile(file_path string) (err error) {
	ext := strings.ToLower(path.Ext(file_path))

	if _, ok := video[ext]; ok {
		err = exec.Command("mplayer", file_path).Start()
	} else if _, ok := source[ext]; ok {
		err = exec.Command("open", "-a", "Sublime Text", file_path).Start()
	} else if _, ok := img[ext]; ok {
		err = exec.Command("open", "-a", "preview", file_path).Start()
	} else if ext == ".pdf" {
		err = exec.Command("open", "-a", "preview", file_path).Start()
	}
	return
}

// Создание файла
func CreateFile(name string) error {
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	file.Close()
	return nil
}

// Создание директории
func MkDir(name string) error {
	return os.Mkdir(name, 0666)
}

// Переименование
func Rename(oldpath, newpath string) error {
	return os.Rename(oldpath, newpath)
}

// Удаление
func Delete(path string) error {
	return os.RemoveAll(path)
}
