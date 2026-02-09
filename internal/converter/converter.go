package converter

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"
)

type ConvertResult struct {
	ConvertedPaths []string
	Duration       time.Duration
	ConvertedCount int
	UpToDateCount  int
	DeletedCount   int
	reconvert      bool
}

func (r ConvertResult) String() string {
	var b strings.Builder

	b.WriteString(
		fmt.Sprintf("🎉 Article processing done in %s. converted: %d", r.Duration, r.ConvertedCount),
	)

	if !r.reconvert {
		b.WriteString(fmt.Sprintf(", up-to-date: %d", r.UpToDateCount))
	}

	b.WriteString(fmt.Sprintf(", deleted: %d\n", r.DeletedCount))

	return b.String()
}

func Convert(inputPath, outputPath string, reconvert bool) (ConvertResult, error) {
	res := ConvertResult{reconvert: reconvert}
	files, err := os.ReadDir(inputPath)

	if err != nil {
		return ConvertResult{}, err
	}

	res.ConvertedPaths = make([]string, 0, len(files))
	begin := time.Now()
	validOutputFiles := make([]string, 0, len(files))

	for _, file := range files {
		ext := filepath.Ext(file.Name())

		if ext != ".md" {
			continue
		}

		in := path.Join(inputPath, file.Name())
		out := path.Join(outputPath, strings.TrimSuffix(file.Name(), ext)+".json")
		validOutputFiles = append(validOutputFiles, out)

		shouldConvert, err := shouldConvert(in, out)

		if !shouldConvert && !reconvert {
			res.UpToDateCount += 1
			continue
		}

		article, err := Parse(in)

		if err != nil {
			return res, fmt.Errorf(`error parsing %s: %v`, file.Name(), err)
		}

		data, err := json.MarshalIndent(article, "", "  ")

		if err != nil {
			return res, err
		}

		err = os.WriteFile(out, data, os.ModePerm)

		if err != nil {
			return res, err
		}

		res.ConvertedPaths = append(res.ConvertedPaths, in)
		res.ConvertedCount += 1
	}

	files, err = os.ReadDir(inputPath)

	if err != nil {
		return res, err
	}

	for _, file := range files {
		ext := filepath.Ext(file.Name())

		if ext != ".json" {
			continue
		}

		out := path.Join(outputPath, file.Name())

		if !slices.Contains(validOutputFiles, out) {
			fmt.Printf("⚠️ deleted %s...\n", out)
			err := os.Remove(out)

			if err != nil {
				return res, err
			}

			res.DeletedCount += 1
		}
	}

	res.Duration = time.Since(begin)

	return res, nil
}

func isGitDirty(filePath string) (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain", "--", filePath)
	output, err := cmd.Output()

	if err != nil {
		return false, err
	}

	return strings.TrimSpace(string(output)) != "", nil
}

func getGitModTime(filePath string) (time.Time, error) {
	cmd := exec.Command("git", "log", "-1", "--format=%ct", "--", filePath)
	output, err := cmd.Output()

	if err != nil {
		return time.Time{}, err
	}

	trimmed := strings.TrimSpace(string(output))

	if trimmed == "" {
		return time.Time{}, nil
	}

	timestamp, err := strconv.ParseInt(trimmed, 10, 64)

	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(timestamp, 0), nil
}

func shouldConvert(sourceFn, destFn string) (bool, error) {
	if _, err := os.Stat(destFn); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return true, nil
		}
		return false, err
	}

	dirty, err := isGitDirty(sourceFn)

	if err != nil {
		return false, err
	}

	if dirty {
		return true, nil
	}

	inTime, err := getGitModTime(sourceFn)

	if err != nil {
		return false, err
	}

	outTime, err := getGitModTime(destFn)

	if err != nil {
		return false, err
	}

	if inTime.IsZero() || outTime.IsZero() {
		return true, nil
	}

	return inTime.After(outTime), nil
}
