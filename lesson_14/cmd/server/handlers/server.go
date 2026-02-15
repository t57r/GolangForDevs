package handlers

import "go.mongodb.org/mongo-driver/mongo"

type Server struct {
	db *mongo.Database
}

func NewServer(db *mongo.Database) *Server {
	return &Server{db: db}
}

type OkResponse struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
}
