package utils

import (
	"fmt"
	"regexp"
	"strings"
)

func GetChunksFromString(startTag string, endTag string, content string) []string {
	startTagLength := len(startTag)
	endTagLength := len(endTag)
	startTagReg := regexp.MustCompile(fmt.Sprintf(`%s(.*?)`, regexp.QuoteMeta(startTag)))
	num_snippets := len(startTagReg.FindAllString(content, -1))
	chunks := []string{}
	j := 0
	for i := 0; i < num_snippets; i++ {
		rawChunk := content[j : len(content)-1]
		start_idx := strings.Index(rawChunk, startTag)
		end_idx := strings.Index(rawChunk, endTag)
		if start_idx == -1 || end_idx == -1 {
			continue
		}
		chunk := rawChunk[start_idx+startTagLength : end_idx]
		chunks = append(chunks, chunk)
		j += end_idx + endTagLength
	}
	return chunks
}

func ReplaceAllChars(olds []string, new string, content string) string {
	for _, old := range olds {
		content = strings.ReplaceAll(content, old, new)
	}
	return content
}

func RemoveSpecialChars(raw string) string {
	return regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(raw, "")
}
