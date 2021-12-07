package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

type Counter struct {
	Word  string
	Count int
}

func Top10(s string) []string {
	if s == "" {
		return []string{}
	}

	ss := make([]Counter, 0, len(s))

	input := strings.Fields(s)
	counter := make(map[string]int)
	sorted := make([]string, 0, 10)

	for _, word := range input {
		_, matched := counter[word]
		if matched {
			counter[word]++
		} else {
			counter[word] = 1
		}
	}

	for k, v := range counter {
		ss = append(ss, Counter{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		switch {
		case ss[i].Count == ss[j].Count:
			return ss[i].Word < ss[j].Word
		default:
			return ss[i].Count > ss[j].Count
		}
	})

	for _, v := range ss {
		sorted = append(sorted, v.Word)
	}

	if len(sorted) > 10 {
		return sorted[:10]
	}

	return sorted
}

var splitRegExp = regexp.MustCompile(`\s-\s|[!"#$%&'()*+,\\./:;<=>?@[\]^_\x60{|}~\s]+`)

func Top10Asterisk(s string) []string {
	if s == "" {
		return []string{}
	}

	ssAsterisk := make([]Counter, 0, len(s))

	input := strings.ToLower(s)
	counter := make(map[string]int)
	sorted := make([]string, 0, 10)

	for _, word := range splitRegExp.Split(input, -1) {
		if word != "" {
			counter[word]++
		}
	}

	for k, v := range counter {
		ssAsterisk = append(ssAsterisk, Counter{k, v})
	}

	sort.Slice(ssAsterisk, func(i, j int) bool {
		switch {
		case ssAsterisk[i].Count == ssAsterisk[j].Count:
			return ssAsterisk[i].Word < ssAsterisk[j].Word
		default:
			return ssAsterisk[i].Count > ssAsterisk[j].Count
		}
	})

	for _, v := range ssAsterisk {
		sorted = append(sorted, v.Word)
	}

	if len(sorted) > 10 {
		return sorted[:10]
	}

	return sorted
}
