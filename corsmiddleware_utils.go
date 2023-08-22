package corsmiddleware

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// RemoveDuplicates create a new slice omitting duplicate values.
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

// Contains find if array contains the received value.
func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

// MergeAndUniques merge slices whit uniques values.
func MergeAndUniques(slices ...[]string) []string {
	var input []string

	for _, chunk := range slices {
		input = append(input, chunk...)
	}

	return RemoveDuplicates(input)
}

// CompileOrigins transform raw origins to Regexp.
func CompileOrigins(origin []string) ([]*regexp.Regexp, error) {
	regexReplacer := strings.NewReplacer(".", "\\.", "*", ".*")
	originsRegex := make([]*regexp.Regexp, len(origin))

	for i, v := range origin {
		var err error
		vFix := regexReplacer.Replace(v)
		originsRegex[i], err = regexp.Compile(vFix)

		if err != nil {
			return nil, fmt.Errorf("error compiling origin '%s': %w", v, err)
		}
	}
	return originsRegex, nil
}

// AllowOrigin check if the array of origins match with of received origin request.
func AllowOrigin(origins []*regexp.Regexp, origin string) bool {
	for _, reg := range origins {
		if reg.MatchString(origin) {
			return true
		}
	}
	return false
}

// WriteLogLine write plugin log
func WriteLogLine(name string, message string) {
	if strings.HasSuffix(message, "\n") {
		message = strings.Trim(message, "\n")
	}
	_, _ = os.Stdout.WriteString(fmt.Sprintf("corsmiddleware:[%v]> %v\n", name, message))
}
