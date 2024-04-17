package parser

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type AuthorSummary struct {
	Name    string `json:"name"`
	Lines   int    `json:"lines"`
	Commits int    `json:"commits"`
	Files   int    `json:"files"`
}

func getSortFunction(summaries []AuthorSummary, sortOrder []string) func(i, j int) bool {
	sortByCriteria := map[string]func(i, j int) int{
		"Lines":   compareInt(summaries, func(a AuthorSummary) int { return a.Lines }),
		"Commits": compareInt(summaries, func(a AuthorSummary) int { return a.Commits }),
		"Files":   compareInt(summaries, func(a AuthorSummary) int { return a.Files }),
		"Name":    compareStr(summaries, func(a AuthorSummary) string { return a.Name }),
	}

	var sortFunctions []func(i, j int) int
	for _, criteria := range sortOrder {
		if sortFunc, ok := sortByCriteria[criteria]; ok {
			sortFunctions = append(sortFunctions, sortFunc)
		}
	}

	return func(i, j int) bool {
		for _, sortFunc := range sortFunctions {
			result := sortFunc(i, j)
			if result != 0 {
				return result > 0
			}
		}
		return false
	}
}

func compareInt(summaries []AuthorSummary, getValue func(AuthorSummary) int) func(int, int) int {
	return func(i, j int) int {
		a, b := getValue(summaries[i]), getValue(summaries[j])
		if a == b {
			return 0
		} else if a > b {
			return 1
		}
		return -1
	}
}

func compareStr(summaries []AuthorSummary, getValue func(AuthorSummary) string) func(int, int) int {
	return func(i, j int) int {
		return -strings.Compare(getValue(summaries[i]), getValue(summaries[j]))
	}
}

func getSummaries(statsMap map[string]*AuthorStats, sortOrder []string) []AuthorSummary {
	var summaries []AuthorSummary

	for author, stats := range statsMap {
		summary := AuthorSummary{
			Name:    author,
			Lines:   stats.LinesCount,
			Commits: len(stats.Commits),
			Files:   len(stats.Files),
		}
		summaries = append(summaries, summary)
	}

	sort.SliceStable(summaries, getSortFunction(summaries, sortOrder))
	return summaries
}

type Formatter interface {
	Print(map[string]*AuthorStats) error
}

func NewFormatter(format string, orderBy string) (Formatter, error) {
	sortOrders := map[string][]string{
		"lines":   {"Lines", "Commits", "Files", "Name"},
		"commits": {"Commits", "Lines", "Files", "Name"},
		"files":   {"Files", "Lines", "Commits", "Name"},
	}

	sortOrder, ok := sortOrders[orderBy]
	if !ok {
		return nil, fmt.Errorf("invalid order")
	}

	switch format {
	case "tabular":
		return &TabularFormatter{SortOrder: sortOrder}, nil
	case "csv":
		return &CsvFormatter{SortOrder: sortOrder}, nil
	case "json":
		return &JsonFormatter{SortOrder: sortOrder}, nil
	case "json-lines":
		return &JsonLinesFormatter{SortOrder: sortOrder}, nil
	default:
		return nil, fmt.Errorf("invalid format")
	}
}

type TabularFormatter struct {
	SortOrder []string
}

func (tf *TabularFormatter) Print(statsMap map[string]*AuthorStats) error {
	header := []string{"Name", "Lines", "Commits", "Files"}
	summaries := getSummaries(statsMap, tf.SortOrder)

	colWidths := make([]int, len(header))
	for i, col := range header {
		colWidths[i] = len(col)
	}
	for _, summary := range summaries {
		for i, field := range []string{summary.Name, strconv.Itoa(summary.Lines), strconv.Itoa(summary.Commits), strconv.Itoa(summary.Files)} {
			if len(field) > colWidths[i] {
				colWidths[i] = len(field)
			}
		}
	}

	for i, field := range header {
		if i != len(header)-1 {
			fmt.Printf("%-*s ", colWidths[i], field)
		} else {
			fmt.Println(field)
		}
	}

	for _, summary := range summaries {
		row := []string{summary.Name, strconv.Itoa(summary.Lines), strconv.Itoa(summary.Commits), strconv.Itoa(summary.Files)}
		for i, field := range row {
			if i != len(header)-1 {
				fmt.Printf("%-*s ", colWidths[i], field)
			} else {
				fmt.Println(field)
			}
		}
	}

	return nil
}

type CsvFormatter struct {
	SortOrder []string
}

func (cf *CsvFormatter) Print(statsMap map[string]*AuthorStats) error {
	header := []string{"Name", "Lines", "Commits", "Files"}
	summaries := getSummaries(statsMap, cf.SortOrder)

	writer := csv.NewWriter(os.Stdout)
	err := writer.Write(header)
	if err != nil {
		return err
	}

	for _, summary := range summaries {
		row := []string{summary.Name,
			strconv.Itoa(summary.Lines),
			strconv.Itoa(summary.Commits),
			strconv.Itoa(summary.Files)}
		err := writer.Write(row)
		if err != nil {
			return err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return err
	}

	return nil
}

type JsonFormatter struct {
	SortOrder []string
}

func (jf *JsonFormatter) Print(statsMap map[string]*AuthorStats) error {
	summaries := getSummaries(statsMap, jf.SortOrder)

	jsonData, err := json.Marshal(summaries)
	if err != nil {
		return err
	}

	fmt.Println(string(jsonData))
	return nil
}

type JsonLinesFormatter struct {
	SortOrder []string
}

func (jlf *JsonLinesFormatter) Print(statsMap map[string]*AuthorStats) error {
	summaries := getSummaries(statsMap, jlf.SortOrder)

	for _, summary := range summaries {
		jsonData, err := json.Marshal(summary)
		if err != nil {
			return err
		}
		fmt.Println(string(jsonData))
	}

	return nil
}
