package chunker

import "strings"

func Split(markdown string) []string {
	markdown = stripFrontMatter(markdown)
	markdown = strings.TrimSpace(markdown)
	if markdown == "" {
		return nil
	}

	lines := strings.Split(markdown, "\n")
	var chunks []string
	var current strings.Builder

	flush := func() {
		chunk := strings.TrimSpace(current.String())
		if chunk != "" {
			chunks = append(chunks, chunk)
		}
		current.Reset()
	}

	for _, line := range lines {
		if strings.HasPrefix(line, "## ") && current.Len() > 0 {
			flush()
		}
		if current.Len() > 0 {
			current.WriteString("\n")
		}
		current.WriteString(line)
	}
	flush()

	return chunks
}

func stripFrontMatter(markdown string) string {
	trimmed := strings.TrimSpace(markdown)
	if !strings.HasPrefix(trimmed, "---") {
		return markdown
	}

	rest := trimmed[3:]
	end := strings.Index(rest, "\n---")
	if end < 0 {
		return markdown
	}

	return strings.TrimSpace(rest[end+4:])
}
