package handlers

type Server struct {
	repo Repository
}

func NewServer(repo Repository) *Server {
	return &Server{repo: repo}
}

type OkResponse struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
}
