package main

import (
	auth "auth/proto"

	"bytes"
	"context"
	"crypto/sha1"
	"strconv"
	"fmt"
	"math"
	"math/rand"
	"net"
	"time"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

var primes = [20]int{2, 13, 37, 53, 17, 29, 3, 41, 43, 31, 7, 5, 23, 11, 19, 83, 101, 71, 97, 223}
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
	serverNonce := nonceGen()
	p := primeNumberGen()
	g := rand.Intn(50)

	redisCli.Set(ctx, "p:"+getKey(in.Nonce, serverNonce), p, 20*time.Minute)
	redisCli.Set(ctx, "g:"+getKey(in.Nonce, serverNonce), g, 20*time.Minute)

	return &auth.PGResponse{
		Nonce:       in.Nonce,
		ServerNonce: serverNonce,
		MessageId:   in.MessageId + 1,
		P:           int32(p),
		G:           int32(g),
	}, nil
}

func (s *Server) Req_DHParams(ctx context.Context, in *auth.DHParamsRequest) (*auth.DHParamsResponse, error) {
	pStr, err := redisCli.Get(ctx, "p:"+getKey(in.Nonce, in.ServerNonce)).Result()
	if err != nil {
		return nil, err
	}
	gStr, err := redisCli.Get(ctx, "g:"+getKey(in.Nonce, in.ServerNonce)).Result()
	if err != nil {
		return nil, err
	}
	p, err := strconv.Atoi(pStr)
	g, err := strconv.Atoi(gStr)
	b := rand.Intn(50)
	pubB := 1.
	key := 1.
	for i := 0; i < b; i++ {
		pubB = math.Mod(pubB*float64(g), float64(p))
		key = math.Mod(key*float64(in.A), float64(p))
	}

	redisCli.Set(ctx, "authKey:"+fmt.Sprintf("%d", int(key)), key, 0)
	redisCli.Del(ctx, getKey(in.Nonce, in.ServerNonce))

	return &auth.DHParamsResponse{
		Nonce:       in.Nonce,
		ServerNonce: in.ServerNonce,
		MessageId:   in.MessageId + 1,
		B:           uint64(pubB),
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
	lis, err := net.Listen("tcp", ":5052")
	if err != nil {
		fmt.Printf("failed to listen: %v\n", err)
	}

	s := grpc.NewServer()
	serv := Server{}

	auth.RegisterAuthGeneratorServer(s, &serv)
	fmt.Println("server listening at 5052")
	if err := s.Serve(lis); err != nil {
		fmt.Printf("failed to serve: %v\n", err)
	}
}
