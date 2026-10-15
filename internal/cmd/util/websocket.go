package util

import (
	"context"
	"errors"
	"io"
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

	ws, err := cfg.DialContext(ctx)
	if err != nil {
		return nil, err
	}
	ws.PayloadType = websocket.BinaryFrame

	return ws, nil
}

func ProxyWebsocket(ctx context.Context, conn net.Conn, ws *websocket.Conn) error {
	closeAll := func() {
		_ = conn.Close()
		_ = ws.Close()
	}
	stop := context.AfterFunc(ctx, closeAll)
	defer stop()

	errs := make(chan error, 2)
	go func() {
		_, err := io.Copy(ws, conn)
		errs <- err
	}()
	go func() {
		_, err := io.Copy(conn, ws)
		errs <- err
	}()

	// close both connections if one side errors
	err := <-errs
	closeAll()
	<-errs

	if ctx.Err() != nil || errors.Is(err, net.ErrClosed) {
		return nil
	}
	return err
}
