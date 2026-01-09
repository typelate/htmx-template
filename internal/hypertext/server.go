package hypertext

import (
	"context"
	"embed"
	"html/template"
	"net/http"
)

//go:embed templates/*.gohtml
var templateFiles embed.FS

var templates = template.Must(template.ParseFS(templateFiles, "templates/*.gohtml"))

type Server struct{}

type IndexData struct {
	Name string
}

func (srv *Server) Index(_ context.Context) IndexData {
	return IndexData{
		Name: "friend",
	}
}

// ContactForm represents form data with validation derived from HTML attributes.
// The HTML template defines: name (minlength=2, maxlength=50),
// email (pattern for email format), message (minlength=10, maxlength=500).
type ContactForm struct {
	Name    string `name:"name"`
	Email   string `name:"email"`
	Message string `name:"message"`
}

// ContactResult holds the result of a contact form submission.
type ContactResult struct {
	Form    ContactForm
	Success bool
}

// GetContact returns an empty form for the contact page.
func (srv *Server) GetContact(_ context.Context) ContactForm {
	return ContactForm{}
}

// Contact handles contact form submissions.
// Returns the form values and success status for template rendering.
func (srv *Server) Contact(_ context.Context, form ContactForm) ContactResult {
	return ContactResult{
		Form:    form,
		Success: true,
	}
}

// NotFoundError represents a resource not found error.
// Implements StatusCode() to return HTTP 404.
type NotFoundError struct {
	Resource string
}

func (e NotFoundError) Error() string {
	return e.Resource + " not found"
}

// StatusCode returns HTTP 404 Not Found.
// muxt uses this method to set the response status code.
func (e NotFoundError) StatusCode() int {
	return http.StatusNotFound
}

// PageData is returned by the Page handler to demonstrate NotFoundError.
type PageData struct {
	Slug    string
	Content string
}

// Page demonstrates returning a NotFoundError for unknown pages.
func (srv *Server) Page(_ context.Context, slug string) (PageData, error) {
	// Only "welcome" page exists - all others return NotFoundError
	if slug == "welcome" {
		return PageData{
			Slug:    slug,
			Content: "Welcome to the example page!",
		}, nil
	}
	return PageData{}, NotFoundError{Resource: "page"}
}
