// go run phone_book.go

// Implement a simple phone book manager.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type phoneQuery struct {
	kind   string
	number string
	name   string
}

func processPhoneQueries(queries []phoneQuery) []string {
	contacts := make(map[string]string)
	result := make([]string, 0)

	for _, query := range queries {
		switch query.kind {
		case "add":
			contacts[query.number] = query.name
		case "del":
			delete(contacts, query.number)
		case "find":
			if name, ok := contacts[query.number]; ok {
				result = append(result, name)
			} else {
				result = append(result, "not found")
			}
		}
	}

	return result
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	queries := make([]phoneQuery, 0, n)
	for i := 0; i < n; i++ {
		var query phoneQuery
		if _, err := fmt.Fscan(reader, &query.kind, &query.number); err != nil {
			return
		}
		if query.kind == "add" {
			if _, err := fmt.Fscan(reader, &query.name); err != nil {
				return
			}
		}
		queries = append(queries, query)
	}

	fmt.Println(strings.Join(processPhoneQueries(queries), "\n"))
}
