package lfile

import "regexp"

type SearchResult struct {
	// File where the matches were found
	File string
	// Matches locations within file
	// TODO: this currently encapsulates the FindAllStringSubmatchIndex
	// when I probably just want to find the whole regex. For more generic
	// searches, I'll just want the [term][start/end]; or I can squash them down
	// so they're all the start/end in order
	Matches [][][]int
	// Counts for each search term
	Counts []int
	// Total number of counts
	Sum int
}

// RegexWords creates a slice of case insensitive Regexp for each word.
func RegexWords(words []string) []*regexp.Regexp {
	wordREs := make([]*regexp.Regexp, len(words))
	for i, w := range words {
		wordREs[i] = regexp.MustCompile("(?i)" + w)
	}
	return wordREs
}

type RegexSearch struct {
	Terms   []*regexp.Regexp
	Results []SearchResult
}

func (rs *RegexSearch) HandleIter(i Iterator) {
	if i.Stat().IsDir() {
		return
	}
	foundAll, m, c, sum := CountAll(string(i.Data()), rs.Terms)
	if foundAll {
		rs.Results = append(rs.Results, SearchResult{
			File:    i.Path(),
			Matches: m,
			Counts:  c,
			Sum:     sum,
		})
	}
}

// CountAll searches the string for the given words. If all words are not
// present, all return values will be their zero value. If all words are found,
// bool will be true, the second value indicates the locations of every word,
// the third return is the count of each word and the fourth return is the total
// number of matches.
func CountAll(s string, words []*regexp.Regexp) (foundAll bool, matches [][][]int, counts []int, sum int) {
	matches = make([][][]int, len(words))
	counts = make([]int, len(words))
	for i, re := range words {
		c := re.FindAllStringSubmatchIndex(s, -1)
		ln := len(c)
		if ln == 0 {
			return false, nil, nil, 0
		}
		matches[i] = c
		counts[i] = ln
		sum += ln
	}
	foundAll = true
	return
}
