package aws

import (
	"bufio"
	"bytes"
	"regexp"
	"sort"
	"strings"
)

var (
	profileHeaderPattern = regexp.MustCompile(`^profile\s+(\S+)\]`)
	defaultHeaderPattern = regexp.MustCompile(`^default\]`)
	attrPattern          = regexp.MustCompile(`^([a-zA-Z_]+)\s*=\s*(.+)$`)
)

type Profile struct {
	Name       string
	Attributes map[string]string
}

func profileFromBytes(b []byte) *Profile {
	scanner := bufio.NewScanner(bytes.NewReader(b))
	if !scanner.Scan() {
		return nil
	}
	header := strings.TrimSpace(scanner.Text())

	var name string
	if m := profileHeaderPattern.FindStringSubmatch(header); m != nil {
		name = m[1]
	} else if defaultHeaderPattern.MatchString(header) {
		name = "default"
	} else {
		return nil
	}

	p := &Profile{
		Name:       name,
		Attributes: map[string]string{},
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if m := attrPattern.FindStringSubmatch(line); m != nil {
			key := strings.TrimSpace(m[1])
			value := strings.TrimSpace(m[2])
			p.Attributes[key] = value
		}
	}

	return p
}

func (p *Profile) ToString() string {
	keys := make([]string, 0, len(p.Attributes))
	for k := range p.Attributes {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var b strings.Builder
	b.WriteString("\n")
	for _, k := range keys {
		b.WriteString(k + ": " + p.Attributes[k] + "\n")
	}
	return b.String()
}
