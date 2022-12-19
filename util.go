package main

func truncate(str string, maxwidth int) string {
	if len(str) <= maxwidth {
		return str
	}

	return str[0:maxwidth-1] + "…"
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
