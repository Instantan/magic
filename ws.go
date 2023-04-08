package magic

import (
	"log"
	"net"
	"net/http"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type wss struct {
	conn net.Conn
}

func establishWSConnection(w http.ResponseWriter, r *http.Request) *wss {
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		return nil
	}
	return &wss{
		conn: conn,
	}
}

func (w *wss) Write(p []byte) (int, error) {
	err := wsutil.WriteServerText(w.conn, p)
	return 0, err
}

func (w *wss) Read(ctx *pageContext) {
	log.Println("Started WS")
	ctx.epb.setWriter(w)
	defer func() {
		w.conn.Close()
		ctx.epb.setWriter(nil)
		log.Println("Ended WS")
	}()

	for {
		msg, op, err := wsutil.ReadClientData(w.conn)
		if err != nil {
			// handle error
			break
		}
		err = wsutil.WriteServerMessage(w.conn, op, msg)
		if err != nil {
			// handle error
			break
		}
	}
}
