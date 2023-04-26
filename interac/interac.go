package interac

import (
	"bufio"
	"context"
	"github.com/xpwu/go-cmd/exe"
	"github.com/xpwu/go-log/log"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func unixSocketFile() string {
	return filepath.Join(exe.Exe.AbsDir, "." + exe.Exe.Name + ".unix_socket_for_client_cmd")
}

type Request struct {
	Content  string
	Response chan<- string
}

// ChanFromClient don't call twice, before ctx.Done
func ChanFromClient(ctx context.Context) (ch <-chan Request, err error) {
	ctx, logger := log.WithCtx(ctx)

	uf := unixSocketFile()
	_ = os.Remove(uf)

	unixAddr, _ := net.ResolveUnixAddr("unix", uf)
	unixListener, err := net.ListenUnix("unix", unixAddr)

	ret := make(chan Request, 10)

	if err != nil {
		logger.Error("ListenUnix error! CAN NOT USE interactive, error: ", err)
		return nil, err
	}
	logger.Debug("ListenUnix ok.")

	go func() {
		for {
			unixConn, err := unixListener.AcceptUnix()
			if err != nil {
				logger.Error(err)
				return
			}

			go doConn(ctx, unixConn, ret)
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				_ = unixListener.Close()
				logger.Debug("ListenUnix closed.")
				return
			}
		}
	}()

	return ret, nil
}

func addBlankLine(request string) string {
	for len(request) > 0 && request[len(request)-1] == '\n' {
		request = strings.TrimSuffix(request, "\n")
	}
	if len(request) == 0 {
		request += "\n"
	} else {
		request += "\n\n"
	}

	return request
}

func readAll(reader *bufio.Reader) (ret string, err error) {
	allMsg := ""
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}

		// blank line
		if msg == "\n" {
			break
		}

		allMsg += msg
	}

	return allMsg, nil
}

// long connection. one-request-one-response
// end：blank line, \n

// ChanFromServer close(ch.Response): read/writer error
func ChanFromServer(ctx context.Context) (ch chan<- Request, err error) {
	ctx, logger := log.WithCtx(ctx)

	ret := make(chan Request)

	unixAddr, err := net.ResolveUnixAddr("unix", unixSocketFile())
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUnix("unix", nil, unixAddr)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				_ = conn.Close()
			}
		}
	}()

	go func() {

	forLabel:
		for {
			select {
			case req, ok := <-ret:
				if !ok {
					break forLabel
				}

				request := addBlankLine(req.Content)

				_, err = conn.Write([]byte(request))
				if err != nil {
					logger.Error("Write error! ", err)
					close(req.Response)
					break forLabel
				}

				reader := bufio.NewReader(conn)
				allMsg, err := readAll(reader)
				if err != nil {
					logger.Error("Read error! error: ", err)
					close(req.Response)
					break forLabel
				}

				req.Response <- allMsg
			}
		}
	}()

	return ret, nil
}

func doConn(ctx context.Context, unixConn *net.UnixConn, ret chan<- Request) {
	ctx, logger := log.WithCtx(ctx)

	defer func() {
		if r := recover(); r != nil {
			logger.Error(r)
		}
		_ = unixConn.Close()
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				_ = unixConn.Close()
			}
		}
	}()

	reader := bufio.NewReader(unixConn)

	for {
		allMsg, err := readAll(reader)
		if err != nil {
			logger.Warning(err)
			break
		}

		response := make(chan string, 1)
		ret <- Request{
			Content:  allMsg,
			Response: response,
		}
		res := <-response

		res = addBlankLine(res)

		_, err = unixConn.Write([]byte(res))
		if err != nil {
			logger.Warning(err)
			break
		}
	}
}
