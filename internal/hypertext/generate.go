package hypertext

//go:generate go run github.com/typelate/muxt/cmd/muxt generate --receiver-type=Server --routes-func=TemplateRoutes
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate -o=../fake/server.go --fake-name=Server . RoutesReceiver
