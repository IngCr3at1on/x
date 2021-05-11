package pjs

import (
	"bufio"
	"strings"

	"github.com/jackc/pgx/v4"
)

const (
	funcStart = `do $$`
	funcEnd   = `end $$`
)

func scanQueries(sql string) (pgx.Batch, []string, error) {
	scanner := bufio.NewScanner(strings.NewReader(sql))
	scanner.Split(scanPG)

	var queries []string
	var batch pgx.Batch
	for scanner.Scan() {
		query := strings.TrimSpace(scanner.Text())
		if query == `` {
			break
		}

		queries = append(queries, query)
		batch.Queue(query)
	}

	return batch, queries, nil
}

// This is probably overly simple and/or naive.. Test cases likely need expanding.
func scanPG(data []byte, atEOF bool) (int, []byte, error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	start := strings.Index(string(data), funcStart)
	idx := strings.Index(string(data), ";")
	if idx > -1 {
		if start == -1 || idx < start {
			data = data[:idx+1]
			return len(data), data, nil
		}
	}

	if strings.HasPrefix(string(data), funcStart) && strings.Contains(string(data), funcEnd) {
		data = data[:strings.Index(string(data), funcEnd)+len(funcEnd)]
		return len(data), data, nil
	}

	if atEOF {
		return len(data), data, nil
	}

	return 0, nil, nil
}
