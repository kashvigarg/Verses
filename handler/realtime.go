package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jaydee029/Verses/pubsub"
)

func (cfg *handler) subscribeTotimeline(w http.ResponseWriter, ctx context.Context, userid pgtype.UUID) {

	f, ok := w.(http.Flusher)

	if !ok {
		respondWithError(w, http.StatusBadRequest, "streaming unsupported")
	}

	subch, err := pubsub.Consume[timeline_item](cfg.pubsub, "timeline_direct", "timeline_queue", "timeline_item."+uuid.UUID(userid.Bytes).String())
	if err != nil {
		log.Printf("error consuming items: %v", err)
	}
	/*
		cl := &timelineclient{
			timeline: ti,
			Userid:   userid,
		}

		cfg.Clients.timelineClients.Store(cl, nil)*/

	go func() {
		<-ctx.Done()
	}()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for item := range subch {
		writesse(w, item)
		f.Flush()
	}

}

func (cfg *handler) subscribeTocomments(w http.ResponseWriter, ctx context.Context, _ pgtype.UUID, proseid pgtype.UUID) {

	f, ok := w.(http.Flusher)

	if !ok {
		respondWithError(w, http.StatusBadRequest, "streaming unsupported")
	}

	subch, err := pubsub.Consume[Comment](cfg.pubsub, "comment_direct", "comment_queue", "comment_item."+uuid.UUID(proseid.Bytes).String())
	if err != nil {
		log.Printf("error consuming items: %v", err)
	}
	/*
		c := make(chan Comment)

		cl := &commentclient{
			comments: c,
			Userid:   userid,
			Proseid:  proseid,
		}

		cfg.Clients.commentClients.Store(cl, nil)
	*/
	go func() {
		<-ctx.Done()
	}()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for item := range subch {
		writesse(w, item)
		f.Flush()
	}

}

func (cfg *handler) subscribeTonotifications(w http.ResponseWriter, ctx context.Context, userid pgtype.UUID) {

	f, ok := w.(http.Flusher)

	if !ok {
		respondWithError(w, http.StatusBadRequest, "streaming unsupported")
	}
	subch, err := pubsub.Consume[Notification](cfg.pubsub, "notification_direct", "notification_queue", "notification_item."+uuid.UUID(userid.Bytes).String())
	if err != nil {
		log.Printf("error consuming items: %v", err)
	}
	/*
		n := make(chan Notification)

		cl := &notificationclient{
			notifications: n,
			Userid:        userid,
		}

		cfg.Clients.notificationClients.Store(cl, nil)
	*/
	go func() {
		<-ctx.Done()

	}()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for item := range subch {
		writesse(w, item)
		f.Flush()
	}

}
