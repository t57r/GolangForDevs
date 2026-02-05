package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"lesson12/internal/documentstore"
	"lesson12/internal/handlers"
	"lesson12/internal/model"
)

type Server struct {
	store *documentstore.Store
}

func NewServer(store *documentstore.Store) *Server {
	return &Server{store: store}
}

func (s *Server) GetDocumentStore() *documentstore.Store {
	return s.store
}

func (s *Server) GetCollectionAsCollection(name string) (*documentstore.Collection, error) {
	c, err := s.store.GetCollection(name)

	if err != nil {
		return nil, err
	}
	// need concrete Collection for indexing/query methods.
	col, ok := c.(*documentstore.Collection)
	if !ok {
		return nil, fmt.Errorf("collection %q is not *documentstore.Collection", name)
	}
	return col, nil
}

func main() {

	s := NewServer(documentstore.NewStore())

	ln, err := net.Listen("tcp", model.Addr)
	if err != nil {
		log.Fatalf("listen error: %v", err)
	}
	log.Printf("documentstore server listening on %s", model.Addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept error: %v", err)
			continue
		}

		// each client in separate goroutine
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	_ = conn.SetDeadline(time.Time{}) // no deadline by default
	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)

	for {
		line, err := r.ReadBytes('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			log.Printf("read error: %v", err)
			return
		}

		var req model.Request
		if err := json.Unmarshal(model.TrimLine(line), &req); err != nil {
			_ = model.WriteJSON(w, model.Response{ID: "", Ok: false, Error: "invalid json: " + err.Error()})
			continue
		}

		resp := s.dispatch(req)
		if err := model.WriteJSON(w, resp); err != nil {
			log.Printf("write error: %v", err)
			return
		}
	}
}

func (s *Server) dispatch(req model.Request) model.Response {
	switch req.Op {

	// --- collection ops ---
	case model.RequestOperationCreateCollection:
		return handlers.HandleCreateCollection(s, req)

	case model.RequestOperationGetCollection:
		return handlers.HandleGetCollection(s, req)

	case model.RequestOperationDeleteCollection:
		return handlers.HandleDeleteCollection(s, req)

	// --- documents ---
	case model.RequestOperationPutDocument:
		return handlers.HandlePutDocument(s, req)

	case model.RequestOperationGetDocument:
		return handlers.HandleGetDocument(s, req)

	case model.RequestOperationDeleteDocument:
		return handlers.HandleDeleteDocument(s, req)

	case model.RequestOperationListDocuments:
		return handlers.HandleListDocuments(s, req)

	// --- indexes ---
	case model.RequestOperationCreateIndex:
		return handlers.HandleCreateIndex(s, req)

	case model.RequestOperationDeleteIndex:
		return handlers.HandleDeleteIndex(s, req)

	case model.RequestOperationQuery:
		return handlers.HandleQuery(s, req)

	default:
		return model.Response{ID: req.ID, Ok: false, Error: "unknown op: " + string(req.Op)}
	}
}
