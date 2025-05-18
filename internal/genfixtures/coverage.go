package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Block struct {
	File      string
	StartLine int
	StartCol  int
	EndLine   int
	EndCol    int
}

type BlockInfo struct {
	Block Block
	Stmts int
}

type CoverSet map[Block]bool

var profileLine = regexp.MustCompile(`^(.+):([0-9]+)\.([0-9]+),([0-9]+)\.([0-9]+) ([0-9]+) ([0-9]+)$`)

func readCoverage(covDir string) (CoverSet, map[Block]BlockInfo, error) {
	outFile, err := os.CreateTemp("", "dsa-cover-*.out")
	if err != nil {
		return nil, nil, err
	}
	outName := outFile.Name()
	outFile.Close()
	defer os.Remove(outName)

	cmd := goTool("tool", "covdata", "textfmt", "-i="+covDir, "-o="+outName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, nil, fmt.Errorf("covdata textfmt: %w: %s", err, strings.TrimSpace(string(output)))
	}

	data, err := os.ReadFile(outName)
	if err != nil {
		return nil, nil, err
	}

	covered := make(CoverSet)
	infos := make(map[Block]BlockInfo)
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "mode:") {
			continue
		}
		match := profileLine.FindStringSubmatch(line)
		if match == nil {
			continue
		}
		startLine, _ := strconv.Atoi(match[2])
		startCol, _ := strconv.Atoi(match[3])
		endLine, _ := strconv.Atoi(match[4])
		endCol, _ := strconv.Atoi(match[5])
		stmts, _ := strconv.Atoi(match[6])
		count, _ := strconv.Atoi(match[7])
		block := Block{
			File:      match[1],
			StartLine: startLine,
			StartCol:  startCol,
			EndLine:   endLine,
			EndCol:    endCol,
		}
		infos[block] = BlockInfo{Block: block, Stmts: stmts}
		if count > 0 {
			covered[block] = true
		}
	}
	if len(infos) == 0 {
		return nil, nil, fmt.Errorf("empty coverage profile in %s", covDir)
	}
	return covered, infos, nil
}

func mergeInfos(dst map[Block]BlockInfo, src map[Block]BlockInfo) {
	for block, info := range src {
		dst[block] = info
	}
}

type coverageStat struct {
	CoveredStmts int
	TotalStmts   int
	Covered      int
	Total        int
	Uncovered    []Block
}

func (s coverageStat) Percent() float64 {
	if s.TotalStmts == 0 {
		return 0
	}
	return 100 * float64(s.CoveredStmts) / float64(s.TotalStmts)
}

func statsFor(set CoverSet, infos map[Block]BlockInfo, target func(Block) bool) coverageStat {
	var stat coverageStat
	for block, info := range infos {
		if target != nil && !target(block) {
			continue
		}
		stat.Total++
		stat.TotalStmts += info.Stmts
		if set[block] {
			stat.Covered++
			stat.CoveredStmts += info.Stmts
		} else {
			stat.Uncovered = append(stat.Uncovered, block)
		}
	}
	sort.Slice(stat.Uncovered, func(i, j int) bool {
		if stat.Uncovered[i].StartLine != stat.Uncovered[j].StartLine {
			return stat.Uncovered[i].StartLine < stat.Uncovered[j].StartLine
		}
		return stat.Uncovered[i].StartCol < stat.Uncovered[j].StartCol
	})
	return stat
}

func unionCoverage(sets ...CoverSet) CoverSet {
	result := make(CoverSet)
	for _, set := range sets {
		for block := range set {
			result[block] = true
		}
	}
	return result
}

func addCoverage(dst CoverSet, src CoverSet) int {
	added := 0
	for block := range src {
		if !dst[block] {
			dst[block] = true
			added++
		}
	}
	return added
}

func newCoverage(set CoverSet, covered CoverSet) CoverSet {
	result := make(CoverSet)
	for block := range set {
		if !covered[block] {
			result[block] = true
		}
	}
	return result
}

func (b Block) String() string {
	return fmt.Sprintf("%s:%d.%d", filepath.Base(b.File), b.StartLine, b.StartCol)
}
