package shared

import "regexp"

// Color constants for terminal output
const (
	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorMagenta = "\033[35m"
	ColorCyan    = "\033[36m"
	ColorWhite   = "\033[37m"
	ColorReset   = "\033[0m"
)

type JSExtractionLog struct {
	ParentURL string   `json:"parent_url,omitempty"`
	JSURL     []string `json:"js_url,omitempty"`
}

// Result represents the structured output for extracted data
type Result struct {
	URL               string   `json:"url"`
	Endpoints         []string `json:"endpoints,omitempty"`
	Paths             []string `json:"paths,omitempty"`
	Info              []string `json:"info,omitempty"`
	CriticalPaths     []string `json:"critical_paths,omitempty"`
	SensitiveKeywords []string `json:"sensitive,omitempty"`
	Scripts           []string `json:"scripts,omitempty"`
	// Add other fields as needed
}
type RegexPatterns struct {
	HTMLTagFilter     *regexp.Regexp
	Endpoint          *regexp.Regexp
	Path              *regexp.Regexp
	Info              *regexp.Regexp
	CriticalPath      *regexp.Regexp
	SensitiveKeywords *regexp.Regexp
	Script            *regexp.Regexp
}

var Patterns = RegexPatterns{
	HTMLTagFilter: regexp.MustCompile(`^/(/|html|body|tbody|thead|head|style|title|header|css|h1|h2|h3|h4|h5|h6|section|footer|section|div|nav|ul|li|a|u|table|tr|td|span|button|i|th|select|option|g|p|b|form|font|textarea|label|javascript|script|main|strong)\b`),
	Endpoint:      regexp.MustCompile(`(?i)(https?://[^\s"'<>]+|/v\d+(/[a-zA-Z0-9_\-./]*)*|/api(/[a-zA-Z0-9_\-./]*)*|/[a-zA-Z0-9_\-./]+(\?[a-zA-Z0-9_\-=&%\.]*)?)`),
	Path:          regexp.MustCompile(`(?i)(/[a-zA-Z0-9_\-./]+(\.[a-zA-Z0-9]+)?(/[a-zA-Z0-9_\-./]*)*(\?[a-zA-Z0-9_\-=&%\.]*)?)`),
	Info: regexp.MustCompile(`(?i)(
    api[_-]?key\s*[:=]\s*[A-Za-z0-9\-_\.]+|
    secret[_-]?key\s*[:=]\s*[A-Za-z0-9\-_\.]+|
    access[_-]?token\s*[:=]\s*[A-Za-z0-9\-_\.]+|
    refresh[_-]?token\s*[:=]\s*[A-Za-z0-9\-_\.]+|
    bearer[_-]?token\s*[:=]\s*[A-Za-z0-9\-_\.]+|
    bearer\s+[A-Za-z0-9\-_\.]+|
    client[_-]?secret\s*[:=]\s*[A-Za-z0-9\-_\.]+|
    private[_-]?key\s*[:=]\s*[A-Za-z0-9\-_\.]+|
    public[_-]?key\s*[:=]\s*[A-Za-z0-9\-_\.]+|
    session[_-]?id\s*[:=]\s*[A-Za-z0-9\-_\.]+|
    sessionid\s*[:=]\s*[A-Za-z0-9\-_\.]+|
    jwt\s*[:=]\s*[A-Za-z0-9\-_\.]+|
    eyJ[A-Za-z0-9\-_\.]{10,}\.[A-Za-z0-9\-_\.]{10,}\.[A-Za-z0-9\-_\.]{10,}|  # JWT tokens
    [A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}|                          # Emails
    password\s*[:=]\s*[A-Za-z0-9\-_\.!@#$%^&*]+|
    passphrase\s*[:=]\s*[A-Za-z0-9\-_\.!@#$%^&*]+|
    credentials\s*[:=]\s*[A-Za-z0-9\-_\.]+|
    pin\s*[:=]\s*\d{4,}|
    otp\s*[:=]\s*\d{4,}|
    mfa\s*[:=]\s*\d{4,}|
    2fa\s*[:=]\s*\d{4,}|
    ssn\s*[:=]\s*\d{3}-?\d{2}-?\d{4}|
    credit[_-]?card\s*[:=]\s*\d{12,19}|
    card[_-]?number\s*[:=]\s*\d{12,19}|
    iban\s*[:=]\s*[A-Z0-9]{15,34}|
    swift\s*[:=]\s*[A-Z0-9]{8,11}|
    tax[_-]?id\s*[:=]\s*[A-Za-z0-9\-_\.]+|
    license\s*[:=]\s*[A-Za-z0-9\-_\.]+|
    recovery[_-]?code\s*[:=]\s*[A-Za-z0-9\-_\.]+|
    recovery[_-]?token\s*[:=]\s*[A-Za-z0-9\-_\.]+|
    reset[_-]?token\s*[:=]\s*[A-Za-z0-9\-_\.]+
)`),
	CriticalPath: regexp.MustCompile(`(?i)(
    /api/[\w\-/\.]+|
    /admin[\w\-/\.]*|
    /adminpanel[\w\-/\.]*|
    /superadmin[\w\-/\.]*|
    /root[\w\-/\.]*|
    /sys[\w\-/\.]*|
    /system[\w\-/\.]*|
    /auth[\w\-/\.]*|
    /login[\w\-/\.]*|
    /logout[\w\-/\.]*|
    /register[\w\-/\.]*|
    /signup[\w\-/\.]*|
    /signin[\w\-/\.]*|
    /forgot[\w\-/\.]*|
    /reset[\w\-/\.]*|
    /verify[\w\-/\.]*|
    /confirm[\w\-/\.]*|
    /session[\w\-/\.]*|
    /token[\w\-/\.]*|
    /key[\w\-/\.]*|
    /debug[\w\-/\.]*|
    /internal[\w\-/\.]*|
    /config[\w\-/\.]*|
    /secret[\w\-/\.]*|
    /password[\w\-/\.]*|
    /credential[\w\-/\.]*|
    /user[\w\-/\.]*|
    /account[\w\-/\.]*|
    /profile[\w\-/\.]*|
    /private[\w\-/\.]*|
    /secure[\w\-/\.]*|
    /test[\w\-/\.]*|
    /dev[\w\-/\.]*|
    /staging[\w\-/\.]*|
    /sandbox[\w\-/\.]*|
    /backup[\w\-/\.]*|
    /dump[\w\-/\.]*|
    /export[\w\-/\.]*|
    /import[\w\-/\.]*|
    /logs?[\w\-/\.]*|
    /monitor[\w\-/\.]*|
    /status[\w\-/\.]*|
    /health[\w\-/\.]*|
    /billing[\w\-/\.]*|
    /invoice[\w\-/\.]*|
    /payment[\w\-/\.]*|
    /checkout[\w\-/\.]*|
    /cart[\w\-/\.]*|
    /order[\w\-/\.]*|
    /purchase[\w\-/\.]*|
    /report[\w\-/\.]*|
    /graphql[\w\-/\.]*|
    /ws[\w\-/\.]*|
    /socket[\w\-/\.]*|
    /upload[\w\-/\.]*|
    /download[\w\-/\.]*|
    /cmd[\w\-/\.]*|
    /exec[\w\-/\.]*|
    /run[\w\-/\.]*|
    /eval[\w\-/\.]*|
    /robots\.txt|
    /sitemap\.xml|
    /api-docs[\w\-/\.]*|
    /swagger[\w\-/\.]*|
    /openapi[\w\-/\.]*
)`),
	SensitiveKeywords: regexp.MustCompile(`(?i)(
    api[_-]?key|
    secret[_-]?key|
    access[_-]?token|
    refresh[_-]?token|
    bearer[_-]?token|
    bearer\s+[A-Za-z0-9._-]+|
    password|
    passphrase|
    credentials|
    auth[_-]?token|
    client[_-]?secret|
    private[_-]?key|
    public[_-]?key|
    session[_-]?id|
    sessionid|
    jwt|
    cookie|
    csrf|
    xsrf|
    pin|
    otp|
    mfa|
    2fa|
    ssn|
    credit[_-]?card|
    card[_-]?number|
    iban|
    swift|
    tax[_-]?id|
    license|
    secret|
    security[_-]?answer|
    recovery[_-]?code|
    recovery[_-]?token|
    reset[_-]?token|
    auth[_-]?cookie|
    refresh[_-]?token|
    bearer[_-]?token|
    access[_-]?key|
    secret[_-]?access[_-]?key
)`),
	Script: regexp.MustCompile(`<script[^>]+src=["']([^"']+)["']`),
}
