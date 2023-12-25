package util

import (
	"log"
	"os"
	"strings"
)

func LoadFromFile(filepath string) map[string]string {
	b, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalln(err)
	}
	lines := strings.Split(strings.ReplaceAll(string(b), "\r\n", "\n"), "\n")
	linesWithoutBlankLines := make([]string, 0, len(lines))
	for _, line := range lines {
		if len(strings.TrimSpace(line)) > 0 {
			linesWithoutBlankLines = append(linesWithoutBlankLines, line)
		}
	}

	s := strings.Split(strings.Join(linesWithoutBlankLines, "\n"), ";")
	units := s[:len(s)-1]
	q := make(map[string]string, 0)
	for _, u := range units {
		if !strings.Contains(u, "name:") {
			continue
		}

		s := strings.SplitN(strings.TrimSpace(u), "\n", 2)
		nameComment, query := s[0], s[1]

		name := strings.TrimSpace(strings.SplitN(nameComment, ":", 2)[1])
		query = strings.TrimSpace(strings.ReplaceAll(query, "\n", " ")) + ";"

		_, ok := q[name]
		if ok {
			log.Fatalln(name, "already in queries")
		}

		q[name] = query
	}

	return q
}
