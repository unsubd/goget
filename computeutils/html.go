package computeutils

import (
	bytes2 "bytes"
	"golang.org/x/net/html"
	"io"
	"regexp"
)

func ParseHtml(reader io.Reader) (*html.Node, error) {
	return html.Parse(reader)
}

func ExtractFileNames(node *html.Node) []string {
	data := node.Data
	var results []string

	if node.Type == html.ElementNode && data == "a" {
		url := extractContent(node)
		if url != "" {
			results = []string{url}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		results = append(results, ExtractFileNames(child)...)
	}
	return results
}

func extractContent(node *html.Node) string {
	data := parseHtmlNode(node)
	regex := regexp.MustCompile("<a.*>(.*)</a>")
	submatch := regex.FindStringSubmatch(data)

	if len(submatch) < 2 {
		return ""
	}

	return submatch[1]
}

func parseHtmlNode(node *html.Node) string {
	var bytes bytes2.Buffer
	writer := io.Writer(&bytes)
	err := html.Render(writer, node)
	if err != nil {
		return ""
	}

	return bytes.String()
}
