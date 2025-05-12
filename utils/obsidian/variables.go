package obsidian

import (
	"fmt"
	"github.com/maycon-jesus/mj-cli/utils"
	"regexp"
	"strings"
	"time"
)

func DateReplacer(str string, t time.Time) string {
	regex := "{{DATE:([^\\s}]+)}}"
	a, _ := regexp.Compile(regex)
	matchs := a.FindAllStringSubmatch(str, -1)
	for _, match := range matchs {

		match[1] = strings.ReplaceAll(match[1], "MMMM", utils.GetMonthName(t.Month()))

		match[1] = strings.ReplaceAll(match[1], "YYYY", "2006")
		match[1] = strings.ReplaceAll(match[1], "MM", "01")
		match[1] = strings.ReplaceAll(match[1], "DD", "02")
		match[1] = t.Format(match[1])

		_, weekN := t.ISOWeek()
		match[1] = strings.ReplaceAll(match[1], "w", fmt.Sprintf("%d", weekN))
		str = strings.ReplaceAll(str, match[0], match[1])
	}

	return str
}
