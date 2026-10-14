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

	"github.com/tenthirtyam/go-vnc"
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
	serverMessages := make(chan vnc.ServerMessage, 8)

	slog.Debug("connecting vnc client")
	vncConn, err := vnc.ClientWithOptions(ctx, ws,
		// Auth is managed at the HTTP/websocket level,therefore no password is
		// needed for the VNC client.
		vnc.WithAuth(&vnc.ClientAuthNone{}),
		vnc.WithServerMessageChannel(serverMessages),
	)
	if err != nil {
		return err
	}
	defer vncConn.Close()

	width, height := vncConn.GetFrameBufferSize()

	slog.Debug("requesting framebuffer update")
	if err := vncConn.FramebufferUpdateRequest(false, 0, 0, width, height); err != nil {
		return err
	}

	slog.Debug("waiting for server message")
	var framebufferUpdate *vnc.FramebufferUpdateMessage
	for framebufferUpdate == nil {
		select {
		case msg, ok := <-serverMessages:
			if !ok {
				return fmt.Errorf("vnc connection closed before receiving framebuffer update")
			}
			update, ok := msg.(*vnc.FramebufferUpdateMessage)
			if !ok {
				slog.Debug("ignoring unexpected server message", "type", msg.Type())
				continue
			}
			framebufferUpdate = update
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	// Sanity checks
	if len(framebufferUpdate.Rectangles) != 1 ||
		framebufferUpdate.Rectangles[0].X != 0 ||
		framebufferUpdate.Rectangles[0].Y != 0 ||
		framebufferUpdate.Rectangles[0].Width != width ||
		framebufferUpdate.Rectangles[0].Height != height {
		return fmt.Errorf("received invalid frame buffer update")
	}

	slog.Debug("composing image from framebuffer")
	img := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
	{
		rect := framebufferUpdate.Rectangles[0]
		enc := rect.Enc.(*vnc.RawEncoding)
		pf := vncConn.GetPixelFormat()
		for index, clr := range enc.Colors {
			x, y := index%int(width), index/int(width)
			img.Set(int(rect.X)+x, int(rect.Y)+y, color.RGBA{
				R: scaleColor(clr.R, pf.RedMax),
				G: scaleColor(clr.G, pf.GreenMax),
				B: scaleColor(clr.B, pf.BlueMax),
				A: 255,
			})
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

func scaleColor(value, maxValue uint16) uint8 {
	if maxValue == 0 {
		return 0
	}
	return uint8(uint32(value) * 255 / uint32(maxValue))
}
