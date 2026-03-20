package hypertext

//go:generate go tool muxt generate --use-receiver-type=Server
//go:generate go tool counterfeiter -generate

//counterfeiter:generate -o=../fake/server.go --fake-name=Server . RoutesReceiver
