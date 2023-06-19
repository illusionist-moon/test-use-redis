package util

import (
	"reflect"
	"unsafe"
)

func Str2Byte(str string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&str))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

func Byte2Str(slice []byte) string {
	return *(*string)(unsafe.Pointer(&slice))
}
