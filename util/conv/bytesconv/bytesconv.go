package bytesconv

import (
	"reflect"
	"unsafe"
)

// StringToBytes はメモリー割り当てを行わずに文字列型からバイト型配列へ変換します
func StringToBytes(s string) (b []byte) {
	sh := *(*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bh.Data, bh.Len, bh.Cap = sh.Data, sh.Len, sh.Len
	return b
}

// BytesToString はメモリー割り当てを行わずにバイト型配列から文字列型へ変換します
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
