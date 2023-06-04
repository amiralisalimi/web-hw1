package AuthClient

import (
	__ "AuthService/Authenticate"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"math"
	"math/rand"
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

func sendPGRequest(messageID int) (int, error) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("did not connect: %v", err)
	}

	defer conn.Close()

	c := __.NewAuthGeneratorClient(conn)

	r, err := c.ReqPq(context.Background(), &__.PGRequest{
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

func sendDHParamsRequest(messageID int) (int, int, error) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("did not connect: %v", err)
	}

	defer conn.Close()

	c := __.NewAuthGeneratorClient(conn)

	r, err := c.Req_DHParams(context.Background(), &__.DHParamsRequest{
		Nonce:       nonce,
		ServerNonce: serverNonce,
		MessageId:   uint32(messageID),
		A:           int32(aPub),
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
	xPub := int(uint8(math.Pow(float64(g), float64(x))) % uint8(p))
	return x, xPub
}

func getKey(p, bPub, a int) int {
	return int(uint8(p) % uint8(math.Pow(float64(bPub), float64(a))))
}
