package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Simplified Rabin-Karp / Substring sliding window for plagiarism detection
// In a true AC Automaton we would build a Trie of multiple patterns.
// Here we are comparing a Source document to a Target document by looking for
// shared phrases of a certain length (N-grams).
func extractNgrams(text string, n int) []string {
	words := strings.Fields(text)
	var ngrams []string
	if len(words) < n {
		return ngrams
	}
	for i := 0; i <= len(words)-n; i++ {
		ngrams = append(ngrams, strings.Join(words[i:i+n], " "))
	}
	return ngrams
}

type Match struct {
	Index   int    `json:"index"`
	Pattern string `json:"pattern"`
}

// 4. Skynet Plagiarism Matrix
func (a *Backend) RunSkynetPlagiarism(sourceText string, targetText string) {
	if len(sourceText) == 0 || len(targetText) == 0 {
		runtime.EventsEmit(a.ctx, "log", "[Skynet_Core] Error: Source or Target text is empty.")
		return
	}
	taskCtx := a.startNewTask()
	go func() {
		runtime.EventsEmit(a.ctx, "log", "[Skynet_Core] Initializing text analysis matrix...")

		ngramSize := 3 // Match 3 consecutive words
		sourceNgrams := extractNgrams(sourceText, ngramSize)

		if len(sourceNgrams) == 0 {
			runtime.EventsEmit(a.ctx, "log", "[Skynet_Core] Error: Source text too short for N-gram extraction.")
			return
		}

		runtime.EventsEmit(a.ctx, "log", fmt.Sprintf("[Skynet_Core] Extracted %d patterns from Source. Scanning Target...", len(sourceNgrams)))

		targetWords := strings.Fields(targetText)
		totalWords := len(targetWords)
		matchesFound := 0
		var foundMatches []Match

		// Hash set optimization
		sourceNgramSet := make(map[string]struct{})
		for _, ngram := range sourceNgrams {
			sourceNgramSet[ngram] = struct{}{}
		}

		for i := 0; i <= totalWords-ngramSize; i++ {
			select {
			case <-taskCtx.Done():
				runtime.EventsEmit(a.ctx, "log", "[Skynet_Core] Scan cancelled.")
				return
			default:
			}

			window := strings.Join(targetWords[i:i+ngramSize], " ")
			var matchedPhrase string
			if _, exists := sourceNgramSet[window]; exists {
				matchesFound++
				foundMatches = append(foundMatches, Match{Index: i, Pattern: window})
				matchedPhrase = window
				runtime.EventsEmit(a.ctx, "log", fmt.Sprintf("[Skynet_Core] MATCH DETECTED at word index %d: '%s'", i, window))
			}

			// Guard against division by zero
			progress := 100.0
			if totalWords > ngramSize {
				progress = float64(i) / float64(totalWords-ngramSize) * 100.0
			}

			runtime.EventsEmit(a.ctx, "skynet_update", map[string]interface{}{
				"progress_percent": progress,
				"matches_found":    matchesFound,
				"scan_line":        progress, // Visually link progress to top-down scan line
				"latest_match":     matchedPhrase,
			})

			select {
			case <-taskCtx.Done():
				runtime.EventsEmit(a.ctx, "log", "[Skynet_Core] Scan cancelled.")
				return
			case <-time.After(50 * time.Millisecond):
			}
		}

		// Ensure 100% at end
		runtime.EventsEmit(a.ctx, "skynet_update", map[string]interface{}{
			"progress_percent": 100.0,
			"matches_found":    matchesFound,
			"scan_line":        100.0,
			"latest_match":     "",
		})

		matchRate := 0.0
		if len(sourceNgrams) > 0 {
			matchRate = float64(matchesFound) / float64(len(sourceNgrams)) * 100.0
			if matchRate > 100.0 {
				matchRate = 100.0
			}
		}

		runtime.EventsEmit(a.ctx, "skynet_complete", map[string]interface{}{
			"total_matches": matchesFound,
			"match_rate":    matchRate,
			"all_matches":   foundMatches,
		})
		runtime.EventsEmit(a.ctx, "log", fmt.Sprintf("[Skynet_Core] Scan complete. Final Match Rate: %.1f%%", matchRate))
	}()
}
