package converter

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
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

	fmt.Fprintf(&b, "🎉 Article processing done in %s. converted: %d", r.Duration, r.ConvertedCount)

	if !r.reconvert {
		fmt.Fprintf(&b, ", up-to-date: %d", r.UpToDateCount)
	}

	fmt.Fprintf(&b, ", deleted: %d\n", r.DeletedCount)

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

		source, err := os.ReadFile(in)

		if err != nil {
			return res, fmt.Errorf("error reading %s: %w", file.Name(), err)
		}

		sourceHash := hashSource(source)
		needsConvert, err := shouldConvert(out, sourceHash)

		if err != nil {
			return res, fmt.Errorf("error checking %s: %w", file.Name(), err)
		}

		if !needsConvert && !reconvert {
			res.UpToDateCount += 1
			continue
		}

		article, err := parseArticleFromData(string(source))

		if err != nil {
			return res, fmt.Errorf(`error parsing %s: %v`, file.Name(), err)
		}

		article.SourceHash = sourceHash

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

	files, err = os.ReadDir(outputPath)

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

// hashSource returns the hex-encoded SHA-256 of an article's markdown source.
func hashSource(source []byte) string {
	sum := sha256.Sum256(source)
	return hex.EncodeToString(sum[:])
}

// shouldConvert reports whether destFn is absent, unreadable as JSON, or was
// generated from markdown whose hash differs from sourceHash.
func shouldConvert(destFn, sourceHash string) (bool, error) {
	data, err := os.ReadFile(destFn)

	if err != nil {
		if os.IsNotExist(err) {
			return true, nil
		}
		return false, err
	}

	var dest struct {
		SourceHash string `json:"sourceHash"`
	}

	if err := json.Unmarshal(data, &dest); err != nil {
		return true, nil
	}

	return dest.SourceHash != sourceHash, nil
}
