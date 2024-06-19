package main

import (
	"bytes"
	"context"
	"log"
	"text/template"
  "io"
	"net/http"
	"github.com/jfyne/live"
)

type CounterModel struct {
	Count int
}

func NewCounterModel(s live.Socket) *CounterModel {
  log.Println("-> NewCounderModel", s)
  m, ok := s.Assigns().(*CounterModel)
  if !ok {
    m = &CounterModel{
      Count: 0,
    }
  }
  log.Println("NewCounderModel ->", m)
  return m
}

func CounterMount(ctx context.Context, s live.Socket) (interface{}, error) {
  log.Println("-> CounterMount", s)
  return NewCounterModel(s), nil
}

func CounterRender(ctx context.Context, data *live.RenderContext) (io.Reader, error) {
  log.Println("-> CounterRender")
  t, err := template.New("counter").Parse(`
<html>
<head>
</head>
<body>
<div>{{.Assigns.Count}}</div>
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
  log.Println("before live.NewHandler")
	h := live.NewHandler()
  log.Println("before h.HandleMount")
	h.HandleMount(CounterMount)
  log.Println("before h.HandleRender")
  h.HandleRender(CounterRender)
  log.Println("before http.Handle /counter")
  http.Handle("/counter", live.NewHttpHandler(live.NewCookieStore("session-name", []byte("weak-secret")), h))
  log.Println("before http.Handle /live.js")
  http.Handle("/live.js", live.Javascript{})
  log.Println("before http.Handle /auto.js.map")
  http.Handle("/auto.js.map", live.Javascript{})
  log.Println("before http.ListenAndServe")
  http.ListenAndServe(":8080", nil)
}
