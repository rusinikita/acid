package code

import (
	"github.com/rusinikita/acid/ui/theme"
	"regexp"
	"strings"
	"unicode"
)

// Highlight takes SQL code as input and returns the code with:
// - SQL keywords wrapped in <b></b> tags
// - Non-keywords, non-alphanumeric content wrapped in <i></i> tags
// It handles cases where keywords have punctuation attached.
func Highlight(code string) string {
	// Get a map of all SQL keywords for faster lookup
	keywordMap := GetKeywordMap()

	// Split the code into tokens including spaces, punctuation, identifiers and literals
	tokenRegex := `(\s+|[(),;.=<>!%^&*+\-\[\]{}]|"[^"]*"|'[^']*'|\b\w+\b)`
	re := regexp.MustCompile(tokenRegex)
	tokens := re.FindAllString(code, -1)

	var result strings.Builder

	// Helper function to check if a string contains only alphanumeric characters
	isAlphanumeric := func(s string) bool {
		for _, r := range s {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
				return false
			}
		}
		return true
	}

	for _, token := range tokens {
		// Trim any leading/trailing punctuation to check if the word is a keyword
		trimmedToken := strings.Trim(token, " \t\n\r(),;.=<>!%^&*+[]{}-")

		// If we're dealing with a quoted string, don't apply any highlighting
		if (strings.HasPrefix(token, "'") && strings.HasSuffix(token, "'")) ||
			(strings.HasPrefix(token, "\"") && strings.HasSuffix(token, "\"")) {
			result.WriteString(token)
			continue
		}

		// Check if the trimmed token is a SQL keyword (case insensitive)
		_, isKeyword := keywordMap[strings.ToUpper(trimmedToken)]

		// If it's a keyword, wrap it in <b> tags
		if isKeyword && len(trimmedToken) > 0 {
			// Find the position of the keyword within the token
			upperToken := strings.ToUpper(token)
			upperTrimmed := strings.ToUpper(trimmedToken)
			keywordPos := strings.Index(upperToken, upperTrimmed)

			// Write any characters before the keyword (with <i> if not alphanumeric)
			prefix := token[:keywordPos]
			if prefix != "" && !isAlphanumeric(prefix) {
				result.WriteString(prefix)
			} else {
				result.WriteString(theme.SQLWordsStyle.Render(prefix))
			}

			// Write the keyword wrapped in <b> tags
			keyword := token[keywordPos : keywordPos+len(trimmedToken)]
			result.WriteString(theme.SQLKeywordStyle.Render(keyword))

			// Write any characters after the keyword (with <i> if not alphanumeric)
			suffix := token[keywordPos+len(trimmedToken):]
			if suffix != "" && !isAlphanumeric(suffix) {
				result.WriteString(suffix)
			} else {
				result.WriteString(theme.SQLWordsStyle.Render(suffix))
			}
		} else {
			// Not a keyword, check if it needs <i> tags
			if strings.TrimSpace(token) != "" && !isAlphanumeric(token) {
				result.WriteString(token)
			} else {
				result.WriteString(theme.SQLWordsStyle.Render(token))
			}
		}
	}

	return result.String()
}
