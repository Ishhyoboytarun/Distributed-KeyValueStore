package main

import (
	"bytes"
	"net"
	"sync"
)

type Server struct {
	id      string
	kvs     *KVStore
	ln      net.Listener
	clients map[string]net.Conn
	mu      sync.Mutex
	data    map[interface{}]interface{}
}

func NewServer(id string, endpoints []string) (*Server, error) {
	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return nil, err
	}
	kvs, err := NewKVStore(endpoints)
	if err != nil {
		return nil, err
	}
	return &Server{
		id:      id,
		kvs:     kvs,
		ln:      ln,
		clients: make(map[string]net.Conn),
	}, nil
}

func (s *Server) Run() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			return err
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	s.mu.Lock()
	s.clients[conn.RemoteAddr().String()] = conn
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.clients, conn.RemoteAddr().String())
		s.mu.Unlock()
		conn.Close()
	}()

	var buf [4096]byte
	for {
		n, err := conn.Read(buf[:])
		if err != nil {
			return
		}

		parts := bytes.SplitN(buf[:n], []byte{' '}, 2)
		if len(parts) != 2 {
			conn.Write([]byte("Invalid request\n"))
			continue
		}

		key := string(parts[0])
		value := string(parts[1])
		switch string(key) {
		case "GET":
			val, ok := s.data[value]
			if !ok {
				conn.Write([]byte("Key not found\n"))
				continue
			}
			conn.Write([]byte(val.(string) + "\n"))
		case "SET":
			s.mu.Lock()
			s.data[value] = string(parts[1])
			s.mu.Unlock()
			conn.Write([]byte("OK\n"))
		case "DELETE":
			s.mu.Lock()
			delete(s.data, value)
			s.mu.Unlock()
			conn.Write([]byte("OK\n"))
		default:
			conn.Write([]byte("Invalid request\n"))
		}
	}
}
