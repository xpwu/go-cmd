package interac

import (
  "bufio"
  "github.com/xpwu/go-cmd/pid"
  "github.com/xpwu/go-log/log"
  "io"
  "net"
  "os"
  "path/filepath"
  "strings"
)

//  客户端以 \n 为标记，作为一次发送；服务器以io.EOF的接收作为一次发送

func getFilePath (key string) (file string,err error)  {
  key, err = filepath.Abs(key)
  if err != nil {
    return "", err
  }

  if err = os.MkdirAll(key, 0777); err != nil {
    log.Error("mkdir unix-socket error! CAN NOT USE interactive, error: ", err)
    return
  }

  file = filepath.Join(key, ".unix_socket")

  return
}

type RequestChan struct {
  Request  string
  Response chan<- string
}

func ChanFromClient(key string) <-chan RequestChan {
  file,_ := getFilePath(key)
  unixAddr, _ := net.ResolveUnixAddr("unix", file)
  unixListener, err := net.ListenUnix("unix", unixAddr)

  ret := make(chan RequestChan)

  if err != nil {
    log.Error("ListenUnix error! CAN NOT USE interactive, error: ", err)
    return ret
  }

  go func() {
    for {
      // 不让退出
      doConn(unixListener, ret)
    }
  }()

  return ret
}

func ChanFromClientByPID() <-chan RequestChan {
  err := pid.Write(os.Getpid())
  pd, err := pid.Read()
  if err != nil {
    log.Error("write or read pid error ", err)
    return nil
  }
  return ChanFromClient(pd)
}

func ChanFromServer(key string) chan<- RequestChan {

  ret := make(chan RequestChan)

  file,_ := getFilePath(key)

  go func() {
  forLabel:
    for {
      select {
      case req,ok := <-ret:
        if !ok {
          break forLabel
        }
        unixAddr, err := net.ResolveUnixAddr("unix", file)
        if err != nil {
          log.Error(err)
          close(req.Response)
          break
        }
        conn,err := net.DialUnix("unix", nil, unixAddr)
        if err != nil {
          log.Error("DialUnix error! error: ", err)
          close(req.Response)
          break
        }
        _,err = conn.Write([]byte(req.Request))
        if err != nil {
          log.Error("Write error! error: ", err)
          close(req.Response)
          break
        }

        reader := bufio.NewReader(conn)
        allMsg := ""
        for {
          msg, err := reader.ReadString('\n')
          if err == nil {
            allMsg += msg
            continue
          }
          if err != nil && err != io.EOF {
            log.Error("Read error! error: ", err)
            close(req.Response)
            break forLabel
          }
          break
        }

        req.Response <- strings.TrimRight(allMsg, "\n")
      }
    }
  }()

  return ret
}

func doConn(unixListener *net.UnixListener, ret chan RequestChan) {
  defer func() {
    if r := recover(); r != nil {
      // nothing to do
    }
  }()

  for {
    unixConn, err := unixListener.AcceptUnix()
    if err != nil {
      log.Warning(err)
      continue
    }

    reader := bufio.NewReader(unixConn)
    message, err := reader.ReadString('\n')
    if err != nil {
      log.Warning(err)
      continue
    }

    response := make(chan string)
    ret <- RequestChan{
      Request:  message,
      Response: response,
    }
    res := <- response
    _,err = unixConn.Write([]byte(res + "\n"))
    if err != nil {
      log.Warning(err)
    }
    close(response)
    _ = unixConn.Close()
  }
}
