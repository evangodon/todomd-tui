package main

func truncate(str string, maxwidth int) string {
	if len(str) <= maxwidth {
		return str
	}

	maxwidth = max(1, maxwidth)
	return str[0:maxwidth-1] + "…"
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func clamp(minVal int, val int, maxVal int) int {
	return max(minVal, min(val, maxVal))
}
