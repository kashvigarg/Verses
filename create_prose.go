package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	auth "github.com/jaydee029/Verses/internal/auth"
	"github.com/jaydee029/Verses/internal/database"
)

type Post struct {
	ID         pgtype.UUID      `json:"id"`
	Userid     pgtype.UUID      `json:"userid"`
	Body       string           `json:"body"`
	User       *User            `json:"user,omitempty"`
	Created_at pgtype.Timestamp `json:"created_at"`
	Updated_at pgtype.Timestamp `json:"Updated_at"`
	Mine       bool             `json:"mine"`
}

type body struct {
	Body string `json:"body"`
}

func (cfg *apiconfig) postProse(w http.ResponseWriter, r *http.Request) {

	token, err := auth.BearerHeader(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	authorid, err := auth.ValidateToken(token, cfg.jwtsecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	decoder := json.NewDecoder(r.Body)
	params := body{}
	err = decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "parameters couldn't be decoded")
		return
	}

	cleanText := strings.TrimSpace(params.Body)

	if len([]rune(cleanText)) > 280 || cleanText == "" {
		respondWithError(w, http.StatusBadRequest, "Prose is invalid")
		return
	}

	var pgUUID pgtype.UUID
	err = pgUUID.Scan(authorid)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	uuids := uuid.New().String()
	var post_pgUUID pgtype.UUID

	err = post_pgUUID.Scan(uuids)
	if err != nil {
		fmt.Println("Error setting UUID:", err)
	}

	//total, _ := cfg.DB.Countprose(r.Context(), pgUUID)
	content := profane(cleanText)

	var pgtime pgtype.Timestamp
	err = pgtime.Scan(time.Now().UTC())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	tx, err := cfg.DBpool.Begin(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	qtx := cfg.DB.WithTx(tx)

	defer tx.Rollback(r.Context())

	prose, err := qtx.Createprose(r.Context(), database.CreateproseParams{
		ID:        post_pgUUID,
		AuthorID:  pgUUID,
		Body:      content,
		CreatedAt: pgtime,
		UpdatedAt: pgtime,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create Prose")
		return
	}

	respondWithJson(w, http.StatusCreated, Post{
		ID:         post_pgUUID,
		Body:       prose.Body,
		Created_at: prose.CreatedAt,
		Updated_at: prose.UpdatedAt,
		Mine:       true,
	})

	/*timelineuuid := uuid.New().String()
	var timelineId pgtype.UUID

	err = timelineId.Scan(timelineuuid)
	if err != nil {
		fmt.Println("Error setting timline id:", err)
	}*/

	err = qtx.InserinTimeline(r.Context(), database.InserinTimelineParams{
		ProseID: post_pgUUID,
		UserID:  pgUUID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't insert Prose in the timeline")
		return
	}

	var tl timeline_item

	tl.Userid = pgUUID
	tl.Post.ID = post_pgUUID
	tl.Post.Mine = true
	tl.Post.Userid = pgUUID
	tl.Post.Body = cleanText

	err = tx.Commit(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't commit the transaction")
	}

	go func(p Post) {

		u, err := cfg.DB.GetUserbyId(r.Context(), p.Userid)

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "couldn't fetch the user")
			return
		}

		p.User.Email = u.Email
		p.User.ID = u.ID
		p.User.Name = u.Name
		p.User.Username = u.Username
		p.Mine = false

		tl, err := cfg.fanoutprose(r.Context(), p)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "couldn't fanout the post")
			return
		}

		for _, i := range tl {
			fmt.Println(i)
			//TODO: Broadcast
		}

	}(tl.Post)

}

func (cfg *apiconfig) fanoutprose(ctx context.Context, p Post) ([]database.FetchTimelineItemsRow, error) {

	items, err := cfg.DB.FetchTimelineItems(ctx, database.FetchTimelineItemsParams{
		ProseID:    p.ID,
		FolloweeID: p.Userid,
	})

	if err != nil {
		return []database.FetchTimelineItemsRow{}, err
	}

	return items, nil

}

func profane(content string) string {
	contentslice := strings.Split(content, " ")

	for i, word := range contentslice {
		wordl := strings.ToLower(word)
		if wordl == "fuck" || wordl == "shit" || wordl == "fornax" {
			contentslice[i] = "****"
		}
	}

	return strings.Join(contentslice, " ")
}
