package hypertext

//go:generate go run github.com/typelate/muxt/cmd/muxt generate --use-receiver-type=Server
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate -o=../fake/server.go --fake-name=Server . RoutesReceiver
