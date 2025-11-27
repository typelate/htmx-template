package hypertext_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/crhntr/dom/domtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/crhntr/muxt-template-module-htmx/internal/fake"
	"github.com/crhntr/muxt-template-module-htmx/internal/hypertext"
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
	} {
		t.Run(tc.Name, run(tc))
	}
}
