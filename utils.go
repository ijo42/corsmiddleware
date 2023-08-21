package corsmiddleware

import (
	"fmt"
	"regexp"
	"strings"
)

func RemoveDuplicates[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	result := []T{}

	for _, val := range slice {
		if _, ok := seen[val]; !ok {
			seen[val] = true
			result = append(result, val)
		}
	}
	return result
}

func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

func MergeAndUniques[T comparable](slices ...[]T) []T {
	var input []T

	for _, chunk := range slices {
		input = append(input, chunk...)
	}

	return RemoveDuplicates(input)
}

func CompileOrigins(origin []string) ([]*regexp.Regexp, error) {
	regexReplacer := strings.NewReplacer(".", "\\.", "*", ".*")

	var originsRegex []*regexp.Regexp
	originsRegex = make([]*regexp.Regexp, len(origin))

	for i, v := range origin {
		var err error
		vFix := regexReplacer.Replace(v)
		originsRegex[i], err = regexp.Compile(vFix)

		if err != nil {
			return nil, fmt.Errorf("error compiling origin '%s': %s", v, err)
		}

	}

	return originsRegex, nil
}

func AllowOrigin(origins []*regexp.Regexp, origin string) bool {
	for _, reg := range origins {
		if reg.MatchString(origin) {
			return true
		}
	}

	return false
}
