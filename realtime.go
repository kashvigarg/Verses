package main

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
)

func (cfg *apiconfig) subscribeTotimeline(w http.ResponseWriter, ctx context.Context, userid pgtype.UUID) {

	f, ok := w.(http.Flusher)

	if !ok {
		respondWithError(w, http.StatusBadRequest, "streaming unsupported")
	}

	ti := make(chan timeline_item)

	cl := &timelineclient{
		timeline: ti,
		Userid:   userid,
	}

	cfg.Clients.timelineClients.Store(cl, nil)

	go func() {
		<-ctx.Done()
		cfg.Clients.timelineClients.Delete(cl)
		close(ti)
	}()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for item := range ti {
		writesse(w, item)
		f.Flush()
	}

}

func (cfg *apiconfig) subscribeTocomments(w http.ResponseWriter, ctx context.Context, userid pgtype.UUID, proseid pgtype.UUID) {

	f, ok := w.(http.Flusher)

	if !ok {
		respondWithError(w, http.StatusBadRequest, "streaming unsupported")
	}

	c := make(chan Comment)

	cl := &commentclient{
		comments: c,
		Userid:   userid,
		Proseid:  proseid,
	}

	cfg.Clients.commentClients.Store(cl, nil)

	go func() {
		<-ctx.Done()
		cfg.Clients.timelineClients.Delete(cl)
		close(c)
	}()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for item := range c {
		writesse(w, item)
		f.Flush()
	}

}

func (cfg *apiconfig) subscribeTonotifications(w http.ResponseWriter, ctx context.Context, userid pgtype.UUID) {

	f, ok := w.(http.Flusher)

	if !ok {
		respondWithError(w, http.StatusBadRequest, "streaming unsupported")
	}

	n := make(chan Notification)

	cl := &notificationclient{
		notifications: n,
		Userid:        userid,
	}

	cfg.Clients.notificationClients.Store(cl, nil)

	go func() {
		<-ctx.Done()
		cfg.Clients.timelineClients.Delete(cl)
		close(n)
	}()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for item := range n {
		writesse(w, item)
		f.Flush()
	}

}
