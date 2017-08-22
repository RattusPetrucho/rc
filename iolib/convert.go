package iolib

import "strconv"

// Конвертация размера
func SizeIntToString(size int64) string {
	if size < 1024 {
		return strconv.FormatFloat(float64(size), 'f', 2, 64) + " B"
	} else if size/1024 < 1024 {
		return strconv.FormatFloat(float64(size)/1024, 'f', 2, 64) + " KB"
	} else if size/(1024*1024) < 1024 {
		return strconv.FormatFloat(float64(size)/(1024*1024), 'f', 2, 64) + " MB"
	} else if size/(1024*1024*1024) < 1024 {
		return strconv.FormatFloat(float64(size)/(1024*1024*1024), 'f', 2, 64) + " GB"
	}

	return strconv.FormatFloat(float64(size)/(1024*1024*1024*1024), 'f', 2, 64) + " TB"
}
