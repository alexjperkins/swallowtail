package handler

import (
	"github.com/monzo/slog"
	"github.com/monzo/typhon"
)

func POSTToChannel(req typhon.Request) typhon.Response {
	slog.Info(req, "RECEIVED: %v", req)
	return req.Response(nil)
}
