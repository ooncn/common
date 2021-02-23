package util

func ArrayStringIndex(s []string, str string) (i int) {
	i = -1
	if len(s) < 1 {
		return
	}
	for j, t := range s {
		if t == str {
			i = j
			return
		}
	}
	return
}
