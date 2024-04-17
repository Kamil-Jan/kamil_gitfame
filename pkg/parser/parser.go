package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"gitfame/configs"
	"gitfame/pkg/cli_scanner"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type AuthorStats struct {
	Commits    map[string]bool
	Files      map[string]bool
	LinesCount int
}

func NewAuthorStats() *AuthorStats {
	return &AuthorStats{
		Commits:    make(map[string]bool),
		Files:      make(map[string]bool),
		LinesCount: 0,
	}
}

type Parser struct {
	Settings    *cli_scanner.Settings
	AuthorStats map[string]*AuthorStats
}

func NewParser(settings *cli_scanner.Settings) *Parser {
	return &Parser{
		Settings:    settings,
		AuthorStats: make(map[string]*AuthorStats),
	}
}

func runCommand(name string, args ...string) (string, error) {
	var stdout, stderr bytes.Buffer

	cmd := exec.Command(name, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf(stderr.String())
	}

	return strings.TrimSpace(stdout.String()), nil
}

func getPatterns(rawString string) []string {
	if rawString == "" {
		return nil
	}
	return strings.Split(rawString, ",")
}

func getLanguagesExtensions(languagesStr string) ([]string, error) {
	languages := getPatterns(languagesStr)
	filteredLanguages := make(map[string]bool)
	for _, lang := range languages {
		filteredLanguages[strings.ToLower(lang)] = true
	}

	allLanguages, err := configs.GetLanguages()
	if err != nil {
		return nil, err
	}

	var extensions []string
	for _, lang := range allLanguages {
		langName := strings.ToLower(lang.Name)
		if filteredLanguages[langName] {
			extensions = append(extensions, lang.Extensions...)
		}
	}
	return extensions, nil
}

func (p *Parser) findFiles() ([]string, error) {
	output, err := runCommand("git", "ls-tree", "-r", p.Settings.Revision, "--name-only", "--full-name", ".")
	if err != nil {
		return nil, err
	}

	if len(output) == 0 {
		return nil, nil
	}

	files := strings.Split(output, "\n")
	filteredFiles := make([]string, 0)

	languagesExtensions, err := getLanguagesExtensions(p.Settings.Languages)
	if err != nil {
		return nil, err
	}

	extensions := getPatterns(p.Settings.Extensions)
	exclude := getPatterns(p.Settings.Exclude)
	restrictTo := getPatterns(p.Settings.RestrictTo)

	for _, file := range files {
		if !p.isExtensionIncluded(file, extensions) {
			continue
		}

		if !p.isExtensionIncluded(file, languagesExtensions) {
			continue
		}

		if exclude != nil && p.isFileIncluded(file, exclude) {
			continue
		}

		if restrictTo != nil && !p.isFileIncluded(file, restrictTo) {
			continue
		}

		filteredFiles = append(filteredFiles, file)
	}

	return filteredFiles, nil
}

func (p *Parser) isExtensionIncluded(file string, extensions []string) bool {
	if extensions == nil {
		return true
	}

	for _, extension := range extensions {
		if filepath.Ext(file) == extension {
			return true
		}
	}
	return false
}

func (p *Parser) isFileIncluded(file string, patterns []string) bool {
	for _, pattern := range patterns {
		matched, _ := filepath.Match(pattern, file)
		if matched {
			return true
		}
	}
	return false
}

func (p *Parser) parseLastCommiter(file string) error {
	output, err := runCommand("git", "log", p.Settings.Revision, "-s", "--format=%H%n%an", "--", file)
	if err != nil {
		return err
	}

	result := strings.Split(output, "\n")
	commit, author := result[0], result[1]
	if _, ok := p.AuthorStats[author]; !ok {
		p.AuthorStats[author] = NewAuthorStats()
	}

	p.AuthorStats[author].Files[file] = true
	p.AuthorStats[author].Commits[commit] = true
	return nil
}

func (p *Parser) parseFile(file string) error {
	output, err := runCommand("git", "blame", p.Settings.Revision, "--porcelain", file)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(strings.NewReader(output))

	commitAuthor := make(map[string]string)

	authorPrefix := "author "
	if p.Settings.UseCommitter {
		authorPrefix = "committer "
	}

	for scanner.Scan() {
		lineGroupInfo := strings.Split(scanner.Text(), " ")
		commit := lineGroupInfo[0]
		linesCountInGroup, err := strconv.Atoi(lineGroupInfo[len(lineGroupInfo)-1])
		if err != nil {
			return err
		}

		if author, ok := commitAuthor[commit]; ok {
			p.AuthorStats[author].LinesCount += linesCountInGroup
		}

		for i := 0; i < linesCountInGroup; {
			if !scanner.Scan() {
				break
			}

			line := scanner.Text()
			if strings.HasPrefix(line, "\t") {
				i++
				continue
			}

			if strings.HasPrefix(line, authorPrefix) {
				author := strings.TrimSpace(strings.TrimPrefix(line, authorPrefix))
				if _, ok := p.AuthorStats[author]; !ok {
					p.AuthorStats[author] = NewAuthorStats()
				}

				if _, ok := commitAuthor[commit]; !ok {
					commitAuthor[commit] = author
				}

				p.AuthorStats[author].Files[file] = true
				p.AuthorStats[author].Commits[commit] = true
				p.AuthorStats[author].LinesCount += linesCountInGroup
				continue
			}
		}
	}

	if len(commitAuthor) == 0 {
		return p.parseLastCommiter(file)
	}

	return nil
}

func (p *Parser) parseFiles(files []string) error {
	for _, file := range files {
		err := p.parseFile(file)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) Parse() error {
	err := os.Chdir(p.Settings.Repository)
	if err != nil {
		return err
	}

	files, err := p.findFiles()
	if err != nil {
		return err
	}

	err = p.parseFiles(files)
	if err != nil {
		return err
	}

	formatter, err := NewFormatter(p.Settings.Format, p.Settings.OrderBy)
	if err != nil {
		return err
	}

	err = formatter.Print(p.AuthorStats)
	if err != nil {
		return err
	}

	return nil
}
