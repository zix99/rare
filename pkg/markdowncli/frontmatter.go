package markdowncli

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

type Frontmatter map[string]string

func (s Frontmatter) Description() string {
	return s["description"]
}

func (s Frontmatter) Order() int {
	v, _ := strconv.Atoi(s["order"])
	return v
}

func (s Frontmatter) Depth() int {
	v, _ := strconv.Atoi(s["depth"])
	return v
}

func ExtractFrontmatter(r io.Reader) Frontmatter {
	ret := make(Frontmatter)

	scanner := bufio.NewScanner(r)

	scanner.Scan()
	if scanner.Text() != "---" {
		return ret
	}

	for scanner.Scan() {
		line := scanner.Text()
		if line == "---" {
			break
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			ret[key] = value
		}
	}

	return ret
}
