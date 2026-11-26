package vncutil

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/url"

	"golang.org/x/net/websocket"
)

func websocketOrigin(wsURL string) (string, error) {
	u, err := url.Parse(wsURL)
	if err != nil {
		return "", err
	}
	u.Scheme = "https"
	return u.String(), nil
}

func DialWebsocket(ctx context.Context, wsURL string) (*websocket.Conn, error) {
	origin, err := websocketOrigin(wsURL)
	if err != nil {
		return nil, err
	}

	cfg, err := websocket.NewConfig(wsURL, origin)
	if err != nil {
		return nil, err
	}

	slog.Debug("dialing vnc websocket", "host", cfg.Location.Host)
	ws, err := cfg.DialContext(ctx)
	if err != nil {
		// unwrap *websocket.DialError to not leak connection token to stdout
		if dialErr, ok := errors.AsType[*websocket.DialError](err); ok {
			err = dialErr.Err
		}
		slog.Debug("vnc websocket dial failed", "host", cfg.Location.Host, "err", err)
		return nil, fmt.Errorf("could not connect to VNC console: %w", err)
	}
	ws.PayloadType = websocket.BinaryFrame
	slog.Debug("vnc websocket connected", "host", cfg.Location.Host)

	return ws, nil
}

type proxyResult struct {
	direction string
	bytes     int64
	err       error
}

func ProxyWebsocket(ctx context.Context, conn net.Conn, ws *websocket.Conn) error {
	closeAll := func() {
		_ = conn.Close()
		_ = ws.Close()
	}
	stop := context.AfterFunc(ctx, closeAll)
	defer stop()

	slog.Debug("proxying vnc connection", "client", conn.RemoteAddr())

	results := make(chan proxyResult, 2)
	go func() {
		n, err := io.Copy(ws, conn)
		results <- proxyResult{"client->websocket", n, err}
	}()
	go func() {
		n, err := io.Copy(conn, ws)
		results <- proxyResult{"websocket->client", n, err}
	}()

	// close both connections if one side errors
	first := <-results
	slog.Debug("vnc proxy stream ended, closing connection",
		"client", conn.RemoteAddr(), "direction", first.direction, "bytes", first.bytes, "err", first.err)
	closeAll()
	second := <-results
	slog.Debug("vnc proxy stream ended",
		"client", conn.RemoteAddr(), "direction", second.direction, "bytes", second.bytes, "err", second.err)

	err := first.err
	if ctx.Err() != nil || errors.Is(err, net.ErrClosed) {
		return nil
	}
	return err
}
