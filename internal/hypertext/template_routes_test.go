package hypertext_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/typelate/dom/domtest"

	"github.com/typelate/htmx-template/internal/fake"
	"github.com/typelate/htmx-template/internal/hypertext"
)

func TestTemplates(t *testing.T) {
	type (
		Given struct {
			srv *fake.Server
		}
		When struct{}
		Then struct {
			res *http.Response
			srv *fake.Server
		}
		Case struct {
			Name  string
			Given func(*testing.T, *Given)
			When  func(*testing.T, *When) *http.Request
			Then  func(*testing.T, *Then)
		}
	)

	run := func(tc Case) func(t *testing.T) {
		return func(t *testing.T) {
			srv := new(fake.Server)
			mux := http.NewServeMux()
			hypertext.TemplateRoutes(mux, srv)
			rec := httptest.NewRecorder()

			if tc.Given != nil {
				tc.Given(t, &Given{srv})
			}
			req := tc.When(t, &When{})
			mux.ServeHTTP(rec, req)
			tc.Then(t, &Then{
				res: rec.Result(),
				srv: srv,
			})
		}
	}

	for _, tc := range []Case{
		{
			Name: "the header has the name",
			Given: func(t *testing.T, g *Given) {
				g.srv.IndexReturns(hypertext.IndexData{
					Name: "somebody",
				})
			},
			When: func(t *testing.T, _ *When) *http.Request {
				return httptest.NewRequest(http.MethodGet, "/", nil)
			},
			Then: func(t *testing.T, then *Then) {
				if assert.Equal(t, 1, then.srv.IndexCallCount()) {
					ctx := then.srv.IndexArgsForCall(0)
					require.NotNil(t, ctx)
				}
				assert.Equal(t, http.StatusOK, then.res.StatusCode)
				doc := domtest.ParseResponseDocument(t, then.res)
				if el := doc.QuerySelector(`h1`); assert.NotNil(t, el) {
					assert.Equal(t, "Hello, somebody!", strings.TrimSpace(el.TextContent()))
				}
			},
		},
		{
			Name: "the about page is routable",
			When: func(t *testing.T, _ *When) *http.Request {
				return httptest.NewRequest(http.MethodGet, "/about", nil)
			},
			Then: func(t *testing.T, then *Then) {
				assert.Equal(t, http.StatusOK, then.res.StatusCode)
				doc := domtest.ParseResponseDocument(t, then.res)
				if el := doc.QuerySelector(`h1`); assert.NotNil(t, el) {
					assert.Equal(t, "About", strings.TrimSpace(el.TextContent()))
				}
			},
		},
		{
			Name: "contact form renders with validation attributes",
			When: func(t *testing.T, _ *When) *http.Request {
				return httptest.NewRequest(http.MethodGet, "/contact", nil)
			},
			Then: func(t *testing.T, then *Then) {
				assert.Equal(t, http.StatusOK, then.res.StatusCode)
				doc := domtest.ParseResponseDocument(t, then.res)

				// Use QuerySelectorSequence to iterate form inputs
				inputCount := 0
				for el := range doc.QuerySelectorSequence("form input, form textarea") {
					inputCount++
					// Use Matches() to verify specific input attributes
					if el.Matches("[name='name']") {
						assert.Equal(t, "2", el.GetAttribute("minlength"))
						assert.Equal(t, "50", el.GetAttribute("maxlength"))
						assert.True(t, el.HasAttribute("required"))
					}
					if el.Matches("[name='email']") {
						assert.True(t, el.HasAttribute("pattern"))
						assert.True(t, el.HasAttribute("required"))
					}
					if el.Matches("[name='message']") {
						assert.Equal(t, "10", el.GetAttribute("minlength"))
						assert.Equal(t, "500", el.GetAttribute("maxlength"))
					}
				}
				assert.Equal(t, 3, inputCount)
			},
		},
		{
			Name: "contact form submission sets HX-Trigger",
			Given: func(t *testing.T, g *Given) {
				g.srv.ContactReturns(hypertext.ContactResult{
					Form:    hypertext.ContactForm{Name: "Test User"},
					Success: true,
				})
			},
			When: func(t *testing.T, _ *When) *http.Request {
				form := url.Values{}
				form.Set("name", "Test User")
				form.Set("email", "test@example.com")
				form.Set("message", "Hello world!")
				req := httptest.NewRequest(http.MethodPost, "/contact", strings.NewReader(form.Encode()))
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				return req
			},
			Then: func(t *testing.T, then *Then) {
				assert.Equal(t, http.StatusOK, then.res.StatusCode)
				assert.Equal(t, "contact-sent", then.res.Header.Get("HX-Trigger"))
				doc := domtest.ParseResponseDocument(t, then.res)
				successDiv := doc.QuerySelector(".success")
				require.NotNil(t, successDiv)
				assert.Contains(t, successDiv.TextContent(), "Test User")
			},
		},
		{
			Name: "page handler returns 404 for unknown pages",
			Given: func(t *testing.T, g *Given) {
				g.srv.PageReturns(hypertext.PageData{}, hypertext.NotFoundError{Resource: "page"})
			},
			When: func(t *testing.T, _ *When) *http.Request {
				return httptest.NewRequest(http.MethodGet, "/page/unknown", nil)
			},
			Then: func(t *testing.T, then *Then) {
				assert.Equal(t, http.StatusNotFound, then.res.StatusCode)
				doc := domtest.ParseResponseDocument(t, then.res)
				errorEl := doc.QuerySelector(".error")
				require.NotNil(t, errorEl)
				assert.Contains(t, errorEl.TextContent(), "not found")
			},
		},
		{
			Name: "page handler returns content for known pages",
			Given: func(t *testing.T, g *Given) {
				g.srv.PageReturns(hypertext.PageData{
					Slug:    "welcome",
					Content: "Welcome to the example page!",
				}, nil)
			},
			When: func(t *testing.T, _ *When) *http.Request {
				return httptest.NewRequest(http.MethodGet, "/page/welcome", nil)
			},
			Then: func(t *testing.T, then *Then) {
				assert.Equal(t, http.StatusOK, then.res.StatusCode)
				doc := domtest.ParseResponseDocument(t, then.res)
				h1 := doc.QuerySelector("h1")
				require.NotNil(t, h1)
				assert.Equal(t, "welcome", strings.TrimSpace(h1.TextContent()))
			},
		},
	} {
		t.Run(tc.Name, run(tc))
	}
}
