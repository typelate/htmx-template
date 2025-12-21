# A Simple MUXT Project Template using HTMX

This is an example [muxt](http://github.com/typelate/muxt) hypertext server.
See this article for how to use this template: https://go.dev/blog/gonew.

You may also want to see:
- [HTMX](https://htmx.org): configured for site interactivity
- [Pico CSS](https://picocss.com/docs): configured to make the site look a bit better by default
- [counterfeiter](https://github.com/maxbrunsfeld/counterfeiter): used to generate fake Server implementations for testing templates
- [unpkg.com](https://unpkg.com): use the meta query param to get new version integrity values

## Developing

Add new HTML in `./internal/hypertext/templates`.
Use inline template declarations in the template files using muxt template name semantics
then run Go generate to generate a function that maps hypertext concepts to helpful handler methods.

## Testing

After developing routes, run `go generate ./...`.

Write handler tests in `./internal/hypertext/server_test.go`.

Write template tests in `./internal/hypertext/template_routes_test.go`.

Run tests with: `go test ./...`

## Licensing

Please consider removing or replacing the LICENSE file.