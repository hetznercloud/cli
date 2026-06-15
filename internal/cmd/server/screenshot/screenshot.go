package screenshot

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log/slog"
	"net/url"
	"os"

	"github.com/alexsnet/go-vnc"
	"github.com/alexsnet/go-vnc/rfbflags"
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

func dialWebsocket(ctx context.Context, wsURL string) (*websocket.Conn, error) {
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

func TakeScreenshot(ctx context.Context, wsURL string, filename string) error {
	slog.Debug("dialing websocket")
	ws, err := dialWebsocket(ctx, wsURL)
	if err != nil {
		return err
	}
	defer ws.Close()

	slog.Debug("creating vnc client")
	vncCfg := &vnc.ClientConfig{
		Auth: []vnc.ClientAuth{
			// Auth is managed at the HTTP/websocket level,therefore no password is
			// needed for the VNC client.
			&vnc.ClientAuthNone{},
		},
		ServerMessages: []vnc.ServerMessage{
			&vnc.FramebufferUpdate{},
		},
		ServerMessageCh: make(chan vnc.ServerMessage, 1),
	}

	slog.Debug("connecting vnc client")
	vncConn, err := vnc.Connect(ctx, ws, vncCfg)
	if err != nil {
		return err
	}
	defer vncConn.Close()

	go func() {
		slog.Debug("listening for vnc server messages")
		if err := vncConn.ListenAndHandle(); err != nil {
			slog.Error("failed to listen and handle incoming vnc server messages", "err", err)
		}
	}()

	slog.Debug("requesting framebuffer update")
	if err := vncConn.FramebufferUpdateRequest(
		rfbflags.RFBFalse,
		0, 0,
		vncConn.FramebufferWidth(), vncConn.FramebufferHeight(),
	); err != nil {
		return err
	}

	slog.Debug("waiting for server message")
	msg := <-vncCfg.ServerMessageCh

	framebufferUpdate, ok := msg.(*vnc.FramebufferUpdate)
	if !ok {
		return fmt.Errorf("received unexpected server message: %T", msg)
	}

	// Sanity checks
	if len(framebufferUpdate.Rects) != 1 ||
		framebufferUpdate.Rects[0].X != 0 ||
		framebufferUpdate.Rects[0].Y != 0 ||
		framebufferUpdate.Rects[0].Width != vncConn.FramebufferWidth() ||
		framebufferUpdate.Rects[0].Height != vncConn.FramebufferHeight() {
		return fmt.Errorf("received invalid frame buffer update")
	}

	width := int(vncConn.FramebufferWidth())
	height := int(vncConn.FramebufferHeight())

	slog.Debug("composing image from framebuffer")
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	{
		rect := framebufferUpdate.Rects[0]
		enc := rect.Enc.(*vnc.RawEncoding)
		for index, clr := range enc.Colors {
			x, y := index%width, index/width
			img.Set(int(rect.X)+x, int(rect.Y)+y, color.RGBA{uint8(clr.R), uint8(clr.G), uint8(clr.B), 255})
		}
	}

	slog.Debug("encoding image")
	imgData := bytes.NewBuffer(nil)
	{
		pngEncoder := png.Encoder{CompressionLevel: png.DefaultCompression}
		if err := pngEncoder.Encode(io.Writer(imgData), img); err != nil {
			return err
		}
	}

	slog.Debug("writing image to file")
	if err := os.WriteFile(filename, imgData.Bytes(), 0600); err != nil {
		return err
	}

	return nil
}
