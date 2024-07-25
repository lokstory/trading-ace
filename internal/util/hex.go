package util

import (
	"github.com/ethereum/go-ethereum/common"
	"strings"
)

func BytesToHexWithPrefix(b []byte) string {
	return AddHexPrefix(BytesToHexWithoutPrefix(b))
}

func BytesToHexWithoutPrefix(b []byte) string {
	return common.Bytes2Hex(b)
}

func HexToBytes(hex string) []byte {
	return common.Hex2Bytes(RemoveHexPrefix(hex))
}

func AddHexPrefix(hex string) string {
	if !strings.HasPrefix(hex, "0x") {
		hex = "0x" + hex
	}

	return hex
}

func RemoveHexPrefix(hex string) string {
	if strings.HasPrefix(hex, "0x") {
		hex = hex[2:]
	}

	return hex
}
