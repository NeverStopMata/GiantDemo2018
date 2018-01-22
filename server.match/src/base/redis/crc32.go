package redis

import (
	"hash/crc32"
	//"strings"
)

func Crc32Slot(key *string) uint32 {
	//s := strings.Index(*key, "(")
	//e := strings.Index(*key, ")")
	//if s != -1 && e != -1 {
	//	return crc32.ChecksumIEEE([]byte((*key)[s-1:e])) % MAXSLOTNUM
	//}
	return crc32.ChecksumIEEE([]byte(*key)) % MAXSLOTNUM
}
