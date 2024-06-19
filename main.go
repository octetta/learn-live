package main

import (
	"bytes"
	"context"
	"fmt"
	"text/template"
  "io"
	"net/http"
	"github.com/jfyne/live"
)

type CounterModel struct {
	Count int
}

func NewCounterModel(s live.Socket) *CounterModel {
  fmt.Println("-> NewCounterModel", s)
  m, ok := s.Assigns().(*CounterModel)
  if !ok {
    m = &CounterModel{
      Count: 0,
    }
  }
  fmt.Println("NewCounterModel ->", m)
  return m
}

func CounterMount(ctx context.Context, s live.Socket) (interface{}, error) {
  fmt.Println("-> CounterMount", s)
  return NewCounterModel(s), nil
}

func CounterAdd(ctx context.Context, s live.Socket, p live.Params) (interface{}, error) {
  fmt.Println("-> CounterAdd", s, p)
  model := NewCounterModel(s)
  model.Count += 1
  return model, nil
}

func CounterSub(ctx context.Context, s live.Socket, p live.Params) (interface{}, error) {
  fmt.Println("-> CounterSub", s, p)
  model := NewCounterModel(s)
  model.Count -= 1
  return model, nil
}

func CounterRender(ctx context.Context, data *live.RenderContext) (io.Reader, error) {
  fmt.Println("-> CounterRender")
  t, err := template.New("counter").Parse(`
<html>
<head>
</head>
<body>
<button live-click="sub">-</button>
<div>{{.Assigns.Count}}</div>
<button live-click="add">+</button>
<script src="/live.js"></script>
</body>
</html>
`)
  if err != nil {
    return nil, err
  }
  var buf bytes.Buffer
  if err := t.Execute(&buf, data); err != nil {
    return nil, err
  }
  return &buf, nil
}

func main() {
  fmt.Println("before live.NewHandler")
	h := live.NewHandler()

  fmt.Println("before h.HandleMount")
	h.HandleMount(CounterMount)

  fmt.Println("before h.HandleRender")
  h.HandleRender(CounterRender)

  fmt.Println("before h.HandleEvent")
  h.HandleEvent("add", CounterAdd)
  h.HandleEvent("sub", CounterSub)

  fmt.Println("before http.Handle /counter")
  http.Handle("/counter", live.NewHttpHandler(live.NewCookieStore("session-name", []byte("weak-secret")), h))

  fmt.Println("before http.Handle /live.js")
  http.Handle("/live.js", live.Javascript{})

  fmt.Println("before http.Handle /auto.js.map")
  http.Handle("/auto.js.map", live.Javascript{})

  fmt.Println("before http.ListenAndServe")
  http.ListenAndServe(":8080", nil)
}
