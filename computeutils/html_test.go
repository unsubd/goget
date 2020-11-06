package computeutils

import (
	"sort"
	"strings"
	"testing"
)

func TestExtractFileNames(t *testing.T) {
	html := `
	<html>
	<body>
	<a href="four/">four/</a>
	<a href="five/">five/</a>
	<a href="six/">six/</a>
	<a href="seven/">seven/</a>
</body>
</html>
`
	parseHtml, err := ParseHtml(strings.NewReader(html))

	if err != nil {
		t.Error("Invalid Html")
	}

	fileNames := ExtractFileNames(parseHtml)

	if len(fileNames) != 4 {
		t.Error("Incorrect Result", fileNames)
	}

	expected := []string{"four/", "five/", "six/", "seven/"}

	sort.Strings(expected)
	sort.Strings(fileNames)

	for i, s := range expected {
		if fileNames[i] != s {
			t.Error("Incorrect Result", fileNames)
		}
	}
}
