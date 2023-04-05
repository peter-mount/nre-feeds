package darwintty

import "strings"

// plainTextAgents contains signatures of the plain-text agents
var plainTextAgents = []string{
	"curl",
	"httpie",
	"lwp-request",
	"wget",
	"python-httpx",
	"python-requests",
	"openbsd ftp",
	"powershell",
	"fetch",
	"aiohttp",
	"http_get",
	"xh",
}

// IsPlainTextAgent returns true if userAgent is a plain-text agent
func IsPlainTextAgent(userAgent string) bool {
	userAgentLower := strings.ToLower(userAgent)
	for _, signature := range plainTextAgents {
		if strings.Contains(userAgentLower, signature) {
			return true
		}
	}
	return false
}
