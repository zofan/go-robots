package robots

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"testing"
)

var (
	robotsFile, _   = os.Open(`./robots.txt`)
	robotsConfig, _ = ParseStream(robotsFile)

	baseHTTPResponse = http.Response{
		StatusCode: 200,
		Header:     map[string][]string{`Content-Type`: {`text/plain`}},
		Body:       ioutil.NopCloser(bytes.NewBufferString(``)),
	}
)

func TestOptions(t *testing.T) {
	if len(robotsConfig.groups) == 0 {
		t.Error(`Groups is empty`)
	}

	if robotsConfig.Host == nil || robotsConfig.Host.Host != `example.com` {
		t.Error(`Invalid host`)
	}

	if len(robotsConfig.SiteMaps) != 2 {
		t.Error(`Incorrect sitemaps count`)
	}

	if _, ok := robotsConfig.SiteMaps[`https://example.com/sitemap.xml`]; !ok {
		t.Error(`Invalid sitemaps`)
	}
}

func TestMatchingOther(t *testing.T) {
	group := robotsConfig.MatchGroup(`YahooBot/2.0`)
	if group == nil {
		t.Error(`No match group`)
		return
	}

	if group.CrawlDelay != 1234567890 {
		t.Error(`Invalid crawl delay`)
	}

	u, _ := url.Parse(`/`)
	if group.IsAllowed(u) {
		t.Error(`Bad access, expected disallow`)
	}

	if group.IsAllowedString(`/posts`) {
		t.Error(`Bad access, expected disallow`)
	}

	if group.IsAllowedString(`/?my=1&param=2`) {
		t.Error(`Bad access with params, expected disallow`)
	}
}

func TestMatchingGoogle(t *testing.T) {
	group := robotsConfig.MatchGroup(`GoogleBot/1.0`)
	if group == nil {
		t.Error(`Invalid group`)
		return
	}

	if group.CrawlDelay != 0.5 {
		t.Error(`Invalid crawl delay`)
	}

	if group.VisitTime.From.Format(`15:04`) != `06:00` {
		t.Error(`Invalid VisitTime.From`)
	}

	if group.VisitTime.To.Format(`15:04`) != `08:45` {
		t.Error(`Invalid VisitTime.To`)
	}
}

func TestCleanParam(t *testing.T) {
	group := robotsConfig.MatchGroup(`GoogleBot/1.0`)
	if group == nil {
		t.Error(`Invalid group`)
		return
	}

	u, _ := url.Parse(`https://www.example.com/posts/toys?sid=1&sort=asc&param=2&page=3`)
	group.CleanParam(u)

	if u.String() != `https://www.example.com/posts/toys?page=3&param=2` {
		t.Error(`Invalid url after clean`, u.String())
	}
}

func TestAllowCase1(t *testing.T) {
	group := robotsConfig.MatchGroup(`case-1`)
	if group == nil {
		t.Error(`Invalid group`)
		return
	}

	if !group.IsAllowedString(`/page`) {
		t.Error(`Bad access, expected allow`)
	}
}

func TestAllowCase2(t *testing.T) {
	group := robotsConfig.MatchGroup(`case-2`)
	if group == nil {
		t.Error(`Invalid group`)
		return
	}

	if !group.IsAllowedString(`/folder/page`) {
		t.Error(`Bad access, expected allow`)
	}
}

func TestAllowCase4(t *testing.T) {
	group := robotsConfig.MatchGroup(`case-4`)
	if group == nil {
		t.Error(`Invalid group`)
		return
	}

	if !group.IsAllowedString(`/`) {
		t.Error(`Bad access, expected allow`)
	}
}

func TestAllowCase5(t *testing.T) {
	group := robotsConfig.MatchGroup(`case-5`)
	if group == nil {
		t.Error(`Invalid group`)
		return
	}

	if group.IsAllowedString(`/page.htm`) {
		t.Error(`Bad access, expected disallow`)
	}
}

func TestParseResponseOK(t *testing.T) {
	resp := baseHTTPResponse

	config, err := ParseResponse(&resp)
	if err != nil {
		t.Error(err)
		return
	}
	if config == nil {
		t.Error(`Fail parse response`)
	}
}

func TestParseResponseUnavailable(t *testing.T) {
	resp := baseHTTPResponse
	resp.StatusCode = 500

	_, err := ParseResponse(&resp)
	if err != ErrUnavailable {
		t.Error(err)
	}
}

func TestParseResponseWrongType(t *testing.T) {
	resp := baseHTTPResponse
	resp.Header = map[string][]string{`Content-Type`: {`text/html`}}

	_, err := ParseResponse(&resp)
	if err != ErrWrongContentType {
		t.Error(err)
	}
}

func TestParseResponseInvalidContent(t *testing.T) {
	resp := baseHTTPResponse
	resp.Body = ioutil.NopCloser(bytes.NewBufferString(`Hello world!`))

	_, err := ParseResponse(&resp)
	if err != ErrInvalidContent {
		t.Error(err)
	}
}

func BenchmarkPatternContains(b *testing.B) {
	group := robotsConfig.MatchGroup(`case-1`)
	if group == nil {
		b.Error(`Invalid group`)
		return
	}

	for i := 0; i < b.N; i++ {
		group.IsAllowedString(`/page`)
	}
}

func BenchmarkPatternRegexp(b *testing.B) {
	group := robotsConfig.MatchGroup(`case-4`)
	if group == nil {
		b.Error(`Invalid group`)
		return
	}

	for i := 0; i < b.N; i++ {
		group.IsAllowedString(`/`)
	}
}
