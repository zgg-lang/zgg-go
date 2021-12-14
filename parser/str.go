package parser

import (
	"strings"
)

const (
	parseStrNormal = iota
	parseStrMetSlash
	parseStrUnicode
	parseStrOctByte
)

func parseStr(raw string) string {
	return parseStrContent(raw, 1, -1)
}

func parseStrContent(raw string, begin, end int) string {
	if end <= 0 {
		end += len(raw)
	}
	rawChars := []rune(raw[begin:end])
	var v strings.Builder
	status := parseStrNormal
	digitsLen := 0
	var tmpUnicode rune
	var tmpByte byte
	for _, ch := range rawChars {
		switch status {
		case parseStrNormal:
			if ch == '\\' {
				status = parseStrMetSlash
			} else {
				v.WriteRune(ch)
			}
		case parseStrMetSlash:
			switch ch {
			case '\\':
				v.WriteRune('\\')
				status = parseStrNormal
			case 'b':
				v.WriteRune('\b')
				status = parseStrNormal
			case 'f':
				v.WriteRune('\f')
				status = parseStrNormal
			case 'n':
				v.WriteRune('\n')
				status = parseStrNormal
			case 'r':
				v.WriteRune('\r')
				status = parseStrNormal
			case 't':
				v.WriteRune('\t')
				status = parseStrNormal
			case 'v':
				v.WriteRune('\v')
				status = parseStrNormal
			case '\'':
				v.WriteRune('\'')
				status = parseStrNormal
			case '$':
				v.WriteRune('$')
				status = parseStrNormal
			case '"':
				v.WriteRune('"')
				status = parseStrNormal
			case 'u', 'U':
				status = parseStrUnicode
				digitsLen = 4
				tmpUnicode = 0
			case 'x', 'X':
				status = parseStrOctByte
				digitsLen = 2
				tmpByte = 0
			}
		case parseStrUnicode:
			switch ch {
			case '0':
				tmpUnicode = (tmpUnicode << 4)
			case '1':
				tmpUnicode = (tmpUnicode << 4) | 1
			case '2':
				tmpUnicode = (tmpUnicode << 4) | 2
			case '3':
				tmpUnicode = (tmpUnicode << 4) | 3
			case '4':
				tmpUnicode = (tmpUnicode << 4) | 4
			case '5':
				tmpUnicode = (tmpUnicode << 4) | 5
			case '6':
				tmpUnicode = (tmpUnicode << 4) | 6
			case '7':
				tmpUnicode = (tmpUnicode << 4) | 7
			case '8':
				tmpUnicode = (tmpUnicode << 4) | 8
			case '9':
				tmpUnicode = (tmpUnicode << 4) | 9
			case 'a', 'A':
				tmpUnicode = (tmpUnicode << 4) | 0x0a
			case 'b', 'B':
				tmpUnicode = (tmpUnicode << 4) | 0x0b
			case 'c', 'C':
				tmpUnicode = (tmpUnicode << 4) | 0x0c
			case 'd', 'D':
				tmpUnicode = (tmpUnicode << 4) | 0x0d
			case 'e', 'E':
				tmpUnicode = (tmpUnicode << 4) | 0x0e
			case 'f', 'F':
				tmpUnicode = (tmpUnicode << 4) | 0x0f
			}
			digitsLen--
			if digitsLen == 0 {
				v.WriteRune(tmpUnicode)
				status = parseStrNormal
			}
		case parseStrOctByte:
			switch ch {
			case '0':
				tmpByte = (tmpByte << 4)
			case '1':
				tmpByte = (tmpByte << 4) | 1
			case '2':
				tmpByte = (tmpByte << 4) | 2
			case '3':
				tmpByte = (tmpByte << 4) | 3
			case '4':
				tmpByte = (tmpByte << 4) | 4
			case '5':
				tmpByte = (tmpByte << 4) | 5
			case '6':
				tmpByte = (tmpByte << 4) | 6
			case '7':
				tmpByte = (tmpByte << 4) | 7
			case '8':
				tmpByte = (tmpByte << 4) | 8
			case '9':
				tmpByte = (tmpByte << 4) | 9
			case 'a', 'A':
				tmpByte = (tmpByte << 4) | 0x0a
			case 'b', 'B':
				tmpByte = (tmpByte << 4) | 0x0b
			case 'c', 'C':
				tmpByte = (tmpByte << 4) | 0x0c
			case 'd', 'D':
				tmpByte = (tmpByte << 4) | 0x0d
			case 'e', 'E':
				tmpByte = (tmpByte << 4) | 0x0e
			case 'f', 'F':
				tmpByte = (tmpByte << 4) | 0x0f
			}
			digitsLen--
			if digitsLen == 0 {
				v.WriteByte(tmpByte)
				status = parseStrNormal
			}
		}
	}
	return v.String()
}
