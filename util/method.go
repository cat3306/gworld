package util

import "hash/crc32"

func MethodHash(method string) uint32 {
	return crc32.ChecksumIEEE([]byte(method))
}
