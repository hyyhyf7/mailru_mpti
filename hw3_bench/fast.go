package search

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

//easyjson:json
type User struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Browsers []string `json:"browsers"`
}

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	i := -1
	sc := bufio.NewScanner(file)
	browserSet := make(map[string]bool)

	_, _ = fmt.Fprintln(out, "found users:")
	for sc.Scan() {
		user := new(User)
		if err = user.UnmarshalJSON(sc.Bytes()); err != nil {
			panic(err)
		}

		i++
		var hasAndroid, hasMSIE bool
		for _, browser := range user.Browsers {
			containsAndroid := strings.Contains(browser, "Android")
			containsMSIE := strings.Contains(browser, "MSIE")

			if containsAndroid {
				hasAndroid = true
			}
			if containsMSIE {
				hasMSIE = true
			}

			if containsAndroid || containsMSIE {
				if !browserSet[browser] {
					browserSet[browser] = true
				}
			}
		}

		if hasAndroid && hasMSIE {
			email := strings.ReplaceAll(user.Email, "@", " [at] ")
			_, _ = fmt.Fprintf(out, "[%d] %s <%s>\n", i, user.Name, email)
		}
	}

	_, _ = fmt.Fprintln(out, "\nTotal unique browsers", len(browserSet))
}
