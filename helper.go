package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	m, err := ParseYaml("materials.yaml")
	if err != nil {
		fmt.Println("error:", err)
	}
	for _, r := range FindRepetitions(m) {
		fmt.Println(r)
	}
}

type Material struct {
	Name        string
	Recommended bool
	Why         string
	Where       string
	Price       string
	Duration    string
	Reference   string

	Line int
}

func ParseYaml(filename string) ([]Material, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\r\n")
	var materials []Material
	var m *Material
	_ = m
	for i, line := range lines {
		t, k, v := parseYamlLine(line)
		_ = k
		_ = v
		switch t {
		case YamlUnrecognized:
			continue
		case YamlList:
			if m != nil {
				materials = append(materials, *m)
				m = nil
			}
		case YamlKeyVal:
			if m == nil {
				m = new(Material)
				m.Line = i
			}
			switch k {
			case "name":
				m.Name = v
			case "recommended":
				v = strings.TrimSpace(strings.ToLower(v))
				if v == "true" {
					m.Recommended = true
				} else {
					m.Recommended = false
				}
			case "why":
				m.Why = v
			case "where":
				m.Where = v
			case "price":
				m.Price = v
			case "duration":
				m.Duration = v
			case "reference":
				m.Reference = v
			}
		}
	}
	if m != nil {
		materials = append(materials, *m)
		m = nil
	}
	return materials, nil
}

type YamlToken int

const (
	YamlUnrecognized = iota
	YamlList         = iota
	YamlKeyVal       = iota
)

func parseYamlLine(l string) (YamlToken, string, string) {
	var k, v string
	if l == "-" {
		return YamlList, "", ""
	}
	l = strings.TrimSpace(l)
	delimiter := strings.Index(l, ":")
	if delimiter < 0 {
		return YamlUnrecognized, "", ""
	}
	k = l[:delimiter]
	v = strings.TrimSpace(l[delimiter+1:])
	return YamlKeyVal, k, v
}

func FindRepetitions(m []Material) []string {
	var r []string
	for i := 0; i < len(m); i++ {
		for j := i + 1; j < len(m); j++ {
			if i == j {
				continue
			}
			m1 := &m[i]
			m2 := &m[j]
			lines := fmt.Sprintf("(%d, %d)", m1.Line, m2.Line)
			if m1.Name == m2.Name {
				r = append(r, fmt.Sprintf("%s name repititions: %s", lines, m1.Name))
			}
			if m1.Reference == m2.Reference {
				r = append(r, fmt.Sprintf("%s reference repititions: %s", lines, m1.Reference))
			}
		}
	}
	return r
}
