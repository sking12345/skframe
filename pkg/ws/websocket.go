package ws

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"skframe/pkg/logger"
	"sync"
)

type Header struct {
	Auth   string `json:"auth"`
	Action string `json:"action"`
}

type Context struct {
	Header        Header                 `json:"header"`
	Body          map[string]interface{} `json:"body"`
	ConnectStatus bool                   `json:"connect_status"`
	Conn          *websocket.Conn
	NextStatus    bool
	UserData      map[string]interface{}
}

type Engine struct {
	allConnects         []*Context
	newConnectHandler   func(ctx *Context)
	closeConnectHandler func(ctx *Context)
	mutex               sync.Mutex
}

type HandlerFun func(*Context)

var handlersFunc = map[string][]HandlerFun{}

func (ctx *Engine) GET(relativePath string, handlers ...HandlerFun) {
	for _, funName := range handlers {
		handlersFunc[relativePath] = append(handlersFunc[relativePath], funName)
	}
}
func (ctx *Engine) POST(relativePath string, handlers ...HandlerFun) {
	for _, funName := range handlers {
		handlersFunc[relativePath] = append(handlersFunc[relativePath], funName)
	}
}

func (ctx *Engine) NewConnect(conn *websocket.Conn) *Context {
	context := &Context{
		Conn:          conn,
		ConnectStatus: true,
	}
	ctx.mutex.Lock()
	ctx.allConnects = append(ctx.allConnects, context)
	ctx.mutex.Unlock()
	return context
}

func (ctx *Engine) CloseConnect(clientCtx *Context) {
	clientCtx.ConnectStatus = false
	if ctx.closeConnectHandler != nil {
		ctx.closeConnectHandler(clientCtx)
	}
}

func (ctx *Engine) NewMessage(data []byte, clientCtx *Context) {
	err := json.Unmarshal(data, &clientCtx)
	if err != nil {
		logger.LogInfoIf(err)
		return
	}
	if len(handlersFunc[clientCtx.Header.Action]) <= 0 {
		logger.InfoString("ws", clientCtx.Header.Action, "call function not exits")
		return
	}
	clientCtx.NextStatus = true
	for _, funItem := range handlersFunc[clientCtx.Header.Action] {
		if clientCtx.NextStatus == false {
			break
		}
		clientCtx.NextStatus = false
		funItem(clientCtx)
	}
}
func (ctx *Context) Next() {
	ctx.NextStatus = true
}

func (ctx *Context) SetUserData(key string, val interface{}) bool { //?????????????????????????????????????????????
	ctx.UserData[key] = val
	return true
}

func (ctx *Context) GetUserData(key string) interface{} { //?????????????????????????????????
	return ctx.UserData[key]
}

func (ctx *Engine) SetNewHandler(handler func(ctx *Context)) {
	ctx.newConnectHandler = handler
}

func (ctx *Engine) SetCloseHandler(handler func(ctx *Context)) {
	ctx.closeConnectHandler = handler
}

func (ctx *Engine) AllConnect(handler func(ctx *Context) (rm bool)) { //???????????????????????????????????????????????????
	for index, item := range ctx.allConnects {
		if handler != nil {
			if handler(item) == true {
				ctx.mutex.Lock()
				ctx.allConnects = append(ctx.allConnects[:index], ctx.allConnects[(index+1):]...)
				ctx.mutex.Unlock()
			}
		}
	}
}
