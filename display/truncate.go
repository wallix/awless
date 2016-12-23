package display

const truncateChars = "..."

func truncateLeft(str string, maxSize int) string {
	ltc := len(truncateChars)
	lstr := len(str)
	if maxSize < ltc {
		return str[lstr-maxSize : lstr]
	}
	if lstr > maxSize {
		return truncateChars + str[lstr-maxSize+ltc:lstr]
	}
	return str
}

func truncateRight(str string, maxSize int) string {
	ltc := len(truncateChars)
	lstr := len(str)
	if maxSize < ltc {
		return str[:maxSize]
	}
	if lstr > maxSize {
		return str[:maxSize-ltc] + truncateChars
	}
	return str
}
