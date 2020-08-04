package datetimeconv

import (
	"time"
	"unsafe"
)

// StringToDateTime 書式に合わせた文字列型から`*time.Time`へ変換
func StringToDateTime(layout *DateTimeFormat, dateTime string) *time.Time {
	t, err := time.Parse(*(*string)(unsafe.Pointer(layout)), dateTime)
	if err != nil {
		return nil
	}
	return &t
}

// DateTimeToString `*time.Time`を書式に合わせた文字列型へ変換
func DateTimeToString(layout *DateTimeFormat, dateTime *time.Time) string {
	return dateTime.Format(*(*string)(unsafe.Pointer(layout)))
}
