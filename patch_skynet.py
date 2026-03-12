import re

with open('edunexus/skynet_ac.go', 'r') as f:
    content = f.read()

replacement = """		for i := 0; i <= totalWords-ngramSize; i++ {
			window := strings.Join(targetWords[i:i+ngramSize], " ")

			latestMatchStr := ""
			for _, pattern := range sourceNgrams {
				if window == pattern {
					matchesFound++
					foundMatches = append(foundMatches, Match{Index: i, Pattern: pattern})
					runtime.EventsEmit(a.ctx, "log", fmt.Sprintf("[Skynet_Core] MATCH DETECTED at word index %d: '%s'", i, pattern))
					latestMatchStr = window
					break
				}
			}

			// Calculate progress based on sliding window position
			progress := float64(i) / float64(totalWords-ngramSize) * 100.0

			runtime.EventsEmit(a.ctx, "skynet_update", map[string]interface{}{
				"progress_percent": progress,
				"matches_found":    matchesFound,
				"scan_line":        progress, // Visually link progress to top-down scan line
				"latest_match":     latestMatchStr,
			})"""

pattern = r"""		for i := 0; i <= totalWords-ngramSize; i\+\+ \{\n\t\t\twindow := strings\.Join\(targetWords\[i:i\+ngramSize\], " "\)\n\n\t\t\tfor _, pattern := range sourceNgrams \{\n\t\t\t\tif window == pattern \{\n\t\t\t\t\tmatchesFound\+\+\n\t\t\t\t\tfoundMatches = append\(foundMatches, Match\{Index: i, Pattern: pattern\}\)\n\t\t\t\t\truntime\.EventsEmit\(a\.ctx, "log", fmt\.Sprintf\("\[Skynet_Core\] MATCH DETECTED at word index %d: '%s'", i, pattern\)\)\n\t\t\t\t\tbreak\n\t\t\t\t\}\n\t\t\t\}\n\n\t\t\t// Calculate progress based on sliding window position\n\t\t\tprogress := float64\(i\) / float64\(totalWords-ngramSize\) \* 100\.0\n\n\t\t\truntime\.EventsEmit\(a\.ctx, "skynet_update", map\[string\]interface\{\}\{\n\t\t\t\t"progress_percent": progress,\n\t\t\t\t"matches_found":    matchesFound,\n\t\t\t\t"scan_line":        progress, // Visually link progress to top-down scan line\n\t\t\t\t"latest_match":     window,\n\t\t\t\}\)"""

new_content = re.sub(pattern, replacement, content, flags=re.DOTALL)

with open('edunexus/skynet_ac.go', 'w') as f:
    f.write(new_content)
