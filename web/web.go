package web

import (
	"fmt"
	"hash/fnv"
	"io"
	"kv-go/db"
	"net/http"
)

type Server struct {
	db         *db.Database
	shardIndex int
	shardCount int
	addresses  map[int]string
}

func NewServer(db *db.Database, shardIndex, shardCount int, addresses map[int]string) *Server {
	return &Server{
		db:         db,
		shardIndex: shardIndex,
		shardCount: shardCount,
		addresses:  addresses,
	}
}

func (s *Server) getShard(key string) int {
	h := fnv.New64()
	h.Write([]byte(key))
	return int(h.Sum64() % uint64(s.shardCount))
}

func (s *Server) redirect(shard int, w http.ResponseWriter, r *http.Request) {

	url := fmt.Sprintf("http://%s%s", s.addresses[shard], r.RequestURI)
	fmt.Fprintf(w, "redirecting from shard %d to shard %d (%q) \n", s.shardIndex, shard, url)

	resp, err := http.Get(url)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error redirecting the request %v", err)
		return
	}

	defer resp.Body.Close()
	io.Copy(w, resp.Body)
}

func (s *Server) GetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	shard := s.getShard(key)
	value, err := s.db.GetKey(key)

	if shard != s.shardIndex {
		s.redirect(shard, w, r)
		return
	}

	fmt.Fprintf(w, "Shard = %d , current shardIndex = %d, address = %q , Value = %q, error = %v \n", shard, s.shardIndex, s.addresses[shard], value, err)
}

func (s *Server) SetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	value := r.Form.Get("value")

	shard := s.getShard(key)
	if shard != s.shardIndex {
		s.redirect(shard, w, r)
		return
	}

	err := s.db.SetKey(key, []byte(value))
	fmt.Fprintf(w, "Error = %v, shardIndex  = %d, current shardIndex = %d \n", err, shard, s.shardIndex)
}

func (s *Server) ListenAndServe(addr *string) error {
	return http.ListenAndServe(*addr, nil)
}
