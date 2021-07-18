package appcast

import (
	"regexp"
	"strings"

	glob "github.com/gobwas/glob"
)

type MatchFunc func(string) bool

func RegexMatch(regex string) MatchFunc {
	r := regexp.MustCompile(regex)
	return func(s string) bool { return r.MatchString(s) }
}

func GlobMatch(globs string) MatchFunc {
	g := glob.MustCompile(globs)
	return func(s string) bool { return g.Match(s) }
}

func ChainMatch(fn ...MatchFunc) MatchFunc {
	return func(s string) bool {
		for _, match := range fn {
			if match(s) {
				return true
			}
		}
		return false
	}
}

func matchFallback(fn1, fn2 MatchFunc) MatchFunc {
	return func(s string) bool {
		if fn1 != nil {
			return fn1(s)
		}
		return fn2(s)
	}
}

func isMacOS(url string) bool {
	return strings.HasSuffix(url, ".dmg")
}

func isWindows64(url string) bool {
	return strings.HasSuffix(url, "64-bit.msi")
}

func isWindows32(url string) bool {
	return strings.HasSuffix(url, "32-bit.msi")
}
