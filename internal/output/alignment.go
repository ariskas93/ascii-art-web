package output

func visibleWidth(s string) int {
	width := 0
	for i := 0; i < len(s); i++ {
		if i+1 < len(s) {
			if s[i] == '\033' && s[i+1] == '[' {
				for i < len(s) && s[i] != 'm' {
					i++
				}
				continue
			}
		}
		width++
	}
	return width
}
