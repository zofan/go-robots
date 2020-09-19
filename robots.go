package robots

import (
	"bufio"
	"errors"
	"github.com/zofan/go-qexp"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	commentRemover = regexp.MustCompile(`\s*#.*`)

	ErrUnavailable      = errors.New(`robots: resource temporary unavailable`)
	ErrWrongContentType = errors.New(`robots: wrong content type`)
	ErrInvalidContent   = errors.New(`robots: invalid content`)
)

type Group struct {
	allows    []qexp.Matcher
	disallows []qexp.Matcher

	cleanParams []cleanParam

	VisitTime  VisitTime
	CrawlDelay float64
}

type VisitTime struct {
	From *time.Time
	To   *time.Time
}

type cleanParam struct {
	pattern qexp.Matcher
	params  []string
}

type Config struct {
	lastGroup *Group
	groups    map[string]*Group
	groupKeys []string

	SiteMaps map[string]*url.URL
	Host     *url.URL
}

func ParseResponse(resp *http.Response) (*Config, error) {
	if resp == nil {
		return &Config{}, nil
	}

	if !strings.Contains(resp.Header.Get(`Content-Type`), `text/plain`) {
		return nil, ErrWrongContentType
	}

	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		return &Config{}, nil
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return ParseStream(resp.Body)
	}

	return nil, ErrUnavailable
}

func ParseStream(stream io.Reader) (*Config, error) {
	config := &Config{
		groups:   map[string]*Group{},
		SiteMaps: map[string]*url.URL{},
	}

	scanner := bufio.NewScanner(stream)
	for scanner.Scan() {
		line := commentRemover.ReplaceAllString(scanner.Text(), ``)
		line = strings.Trim(line, " \t\n\v\f\r\xef\xbb\xbf")

		if len(line) == 0 {
			continue
		}

		param := strings.SplitN(line, `:`, 2)
		if len(param) != 2 {
			return nil, ErrInvalidContent
		}
		key := strings.ToLower(strings.TrimSpace(param[0]))
		value := strings.TrimSpace(param[1])

		if config.parseMainParam(key, value) {
			continue
		}

		config.parseGroupParam(key, value)
	}

	sort.Slice(config.groupKeys, func(i, j int) bool {
		return len(config.groupKeys[i]) > len(config.groupKeys[j])
	})

	return config, nil
}

func (config *Config) parseMainParam(key string, value string) bool {
	switch key {
	case `sitemap`:
		u, err := url.Parse(value)
		if err == nil {
			config.SiteMaps[value] = u
		}
	case `host`:
		u, err := url.Parse(value)
		if err == nil {
			config.Host = u
		}
	default:
		return false
	}

	return true
}

func (config *Config) parseGroupParam(key string, value string) {
	switch key {
	case `user-agent`, `useragent`:
		value = strings.ToLower(value)
		if _, ok := config.groups[value]; !ok {
			config.groups[value] = &Group{}
			config.groupKeys = append(config.groupKeys, value)
		}
		config.lastGroup = config.groups[value]
	case `crawl-delay`, `crawldelay`:
		config.lastGroup.CrawlDelay, _ = strconv.ParseFloat(value, 64)
	case `request-rate`, `requestrate`:
		split := strings.SplitN(value, `/`, 2)
		period, _ := strconv.Atoi(split[1])
		count, _ := strconv.Atoi(split[0])
		if config.lastGroup.CrawlDelay == 0 {
			config.lastGroup.CrawlDelay = float64(period / count)
		}
	case `visit-time`, `visittime`:
		split := strings.SplitN(value, `-`, 2)
		from, _ := time.Parse(`1504`, split[0])
		to, _ := time.Parse(`1504`, split[1])

		config.lastGroup.VisitTime = VisitTime{&from, &to}
	case `clean-param`, `cleanparam`:
		split := strings.SplitN(value, ` `, 2)
		if len(split) == 1 {
			split = append(split, `/`)
		}

		params := cleanParam{
			pattern: patternCompile(split[1]),
			params:  strings.Split(split[0], `&`),
		}
		config.lastGroup.cleanParams = append(config.lastGroup.cleanParams, params)
	case `allow`:
		config.lastGroup.allows = append(config.lastGroup.allows, patternCompile(value))
	case `disallow`:
		config.lastGroup.disallows = append(config.lastGroup.disallows, patternCompile(value))
	}
}

func (config *Config) MatchGroup(userAgent string) *Group {
	userAgent = strings.ToLower(userAgent)

	for _, name := range config.groupKeys {
		if name != `*` && strings.Contains(userAgent, name) {
			return config.groups[name]
		}
	}

	if defaultGroup, ok := config.groups[`*`]; ok {
		return defaultGroup
	}

	return &Group{}
}

func (group *Group) IsAllowedString(url string) bool {
	allowed := true

	if len(group.disallows) == 0 {
		return true
	}

	for _, disallow := range group.disallows {
		if disallow.MatchString(url) {
			allowed = false
			break
		}
	}

	for _, allow := range group.allows {
		if allow.MatchString(url) {
			return true
		}
	}

	return allowed
}

func (group *Group) IsAllowed(u *url.URL) bool {
	return group.IsAllowedString(u.RequestURI())
}

func (group *Group) CleanParam(u *url.URL) {
	values := u.Query()
	ru := u.RequestURI()
	for _, rule := range group.cleanParams {
		if rule.pattern.MatchString(ru) {
			for _, k := range rule.params {
				values.Del(k)
			}
		}
	}
	u.RawQuery = values.Encode()
}

func patternCompile(s string) qexp.Matcher {
	s = strings.TrimRight(s, `*`)

	s = regexp.QuoteMeta(s)
	s = strings.ReplaceAll(s, `\*`, `.*`)
	s = strings.ReplaceAll(s, `\$`, `$`)

	return qexp.MustCompile(`^` + s)
}
