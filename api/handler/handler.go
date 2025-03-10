package handler

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jaydee029/Verses/internal/database"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

/*
type Clients struct {
	timelineClients     sync.Map
	notificationClients sync.Map
	commentClients      sync.Map
}*/

type handler struct {
	fileservercounts int
	Jwtsecret        string
	apiKey           string
	DB               *database.Queries
	DBpool           *pgxpool.Pool
	pubsub           *amqp.Connection
	logger           *zap.Logger
	//Clients          *Clients
}

func New(fscounts int, jwt, apikey string, DBQueries *database.Queries, DBPool *pgxpool.Pool, pubsubconn *amqp.Connection, logger *zap.Logger) *handler {
	return &handler{
		fileservercounts: fscounts,
		Jwtsecret:        jwt,
		apiKey:           apikey,
		DB:               DBQueries,
		DBpool:           DBPool,
		pubsub:           pubsubconn,
		logger:           logger,
	}
}
