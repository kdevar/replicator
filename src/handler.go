package main

import (
	"github.com/siddontang/go-mysql/canal"
	"log"
	"time"
)

type Handler struct {
	canal.DummyEventHandler
	Streamer Streamer
}

func (h *Handler) OnRow(e *canal.RowsEvent) error {
	m := make(map[string]interface{})
	//revisit possibly memoize or cache
	for _, r := range e.Rows {
		for _, c := range e.Table.Columns {
			val, err := e.Table.GetColumnValue(c.Name, r)

			if err != nil {

			}

			m[c.Name] = val
		}
	}

	log.Printf("Packaged item\n%+v", m)

	h.Streamer.Stream(&StreamItem{
		TableName: e.Table.Name,
		Data: m,
		ProcessTime: time.Now().UTC(),
	})

	return nil
}

func NewHandler(s Streamer) *Handler {
	return &Handler{
		Streamer: s,
	}
}