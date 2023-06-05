package auth

import (
	auth "auth/proto"

	"context"
	"fmt"
	"math"
	"math/rand"

	"google.golang.org/grpc"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var a = 0
var aPub = 0
var bPub = 0
var p = 0
var g = 0
var nonce = ""
var serverNonce = ""

func nonceGen() string {
	b := make([]rune, 20)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func SendPGRequest(messageID int) (int, error) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":5052", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("did not connect: %v", err)
	}

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(conn)

	c := auth.NewAuthGeneratorClient(conn)

	r, err := c.ReqPq(context.Background(), &auth.PGRequest{
		Nonce:     nonceGen(),
		MessageId: uint32(messageID),
	})
	if err != nil {
		return 0, err
	}

	p, g = int(r.P), int(r.G)
	a, aPub = getAPrivetAndAPub(g, p)
	nonce, serverNonce = r.Nonce, r.ServerNonce

	return int(r.MessageId), nil
}

func SendDHParamsRequest(messageID int) (int, int, error) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":5052", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("did not connect: %v", err)
	}

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(conn)

	c := auth.NewAuthGeneratorClient(conn)
	r, err := c.Req_DHParams(context.Background(), &auth.DHParamsRequest{
		Nonce:       nonce,
		ServerNonce: serverNonce,
		MessageId:   uint32(messageID),
		A:           uint64(aPub),
	})
	if err != nil {
		return 0, 0, err
	}

	bPub = int(r.B)
	key := getKey(p, bPub, a)

	return key, int(r.MessageId), nil
}

func getAPrivetAndAPub(g, p int) (int, int) {
	x := rand.Intn(50)
	xPub := 1.
	for i := 0; i < x; i++ {
		xPub = math.Mod(xPub*float64(g), float64(p))
	}
	return x, int(xPub)
}

func getKey(p, bPub, a int) int {
	key := 1.
	for i := 0; i < a; i++ {
		key = math.Mod(key*float64(bPub), float64(p))
	}
	return int(key)
}
