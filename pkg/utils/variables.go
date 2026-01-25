package utils

import "regexp"

var doubleMustacheRegex = regexp.MustCompile(`\{\{(\w+)\}\}`)

func ReplaceVariables(str string, variables map[string]string) string {
	if len(variables) == 0 {
		return str
	}

	return doubleMustacheRegex.ReplaceAllStringFunc(str, func(match string) string {
		groups := doubleMustacheRegex.FindStringSubmatch(match)
		if len(groups) > 1 {
			if val, ok := variables[groups[1]]; ok {
				return val
			}
		}
		return match
	})
}
