package logger

import "regexp"

const redacted = "***"

var (
	// "password": "x", "confirm_password":"x", "token": "x"
	jsonSecretPattern = regexp.MustCompile(`(?i)("[a-z_]*(?:password|token|secret)[a-z_]*"\s*:\s*)"(?:[^"\\]|\\.)*"`)
	// password=x&token=y (form-encoded bodies)
	formSecretPattern = regexp.MustCompile(`(?i)((?:^|[&?])[a-z_]*(?:password|token|secret)[a-z_]*=)[^&\s]*`)
)

// redactSecrets masks values of password, token and secret fields so request
// bodies can be logged without leaking credentials or session tokens.
func redactSecrets(s string) string {
	s = jsonSecretPattern.ReplaceAllString(s, `${1}"`+redacted+`"`)
	return formSecretPattern.ReplaceAllString(s, `${1}`+redacted)
}
