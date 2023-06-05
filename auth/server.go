package main

import (
	auth "auth/proto"

	"bytes"
	"context"
	"crypto/sha1"
	"flag"
	"fmt"
	"math"
	"math/rand"
	"net"
	"time"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

var primes = [20]int{461, 151, 197, 79, 239, 263, 137, 127, 139, 113, 101, 277, 83, 479, 397, 233, 23, 449, 223, 251}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var redisCli = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

type Server struct {
	auth.UnimplementedAuthGeneratorServer
}

func (s *Server) ReqPq(ctx context.Context, in *auth.PGRequest) (*auth.PGResponse, error) {
	fmt.Print("HDHDHDHD")
	serverNonce := nonceGen()
	p := primeNumberGen()
	g := rand.Intn(50)
	redisCli.Set(ctx, getKey(in.Nonce, serverNonce), [2]int{p, g}, 20*time.Minute)

	return &auth.PGResponse{
		Nonce:       in.Nonce,
		ServerNonce: serverNonce,
		MessageId:   in.MessageId + 1,
		P:           int32(p),
		G:           int32(g),
	}, nil
}

func (s *Server) Req_DHParams(ctx context.Context, in *auth.DHParamsRequest) (*auth.DHParamsResponse, error) {
	val, err := redisCli.Get(ctx, getKey(in.Nonce, in.ServerNonce)).Result()
	if err != nil {
		return nil, err
	}
	p, g := val[0], val[1]
	b := rand.Intn(50)
	pubB := uint8(math.Pow(float64(g), float64(b))) % p
	x := uint8(math.Pow(float64(in.A), float64(b)))
	key := p % x

	redisCli.Set(ctx, "authKey:"+string(key), key, 0)
	redisCli.Del(ctx, getKey(in.Nonce, in.ServerNonce))
	return &auth.DHParamsResponse{
		Nonce:       in.Nonce,
		ServerNonce: in.ServerNonce,
		MessageId:   in.MessageId + 1,
		B:           int32(pubB),
	}, nil
}

func getKey(s1 string, s2 string) string {
	enc := sha1.Sum([]byte(s1 + s2))
	return bytes.NewBuffer(enc[:]).String()
}

func primeNumberGen() int {
	return primes[rand.Intn(20)]
}

func nonceGen() string {
	b := make([]rune, 20)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func main() {
	port := flag.Int("port", 5052, "Port number")
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		fmt.Printf("failed to listen: %v\n", err)
	}

	s := grpc.NewServer()
	serv := Server{}

	auth.RegisterAuthGeneratorServer(s, &serv)
	fmt.Printf("server listening at %d\n", *port)
	if err := s.Serve(lis); err != nil {
		fmt.Printf("failed to serve: %v\n", err)
	}
}
