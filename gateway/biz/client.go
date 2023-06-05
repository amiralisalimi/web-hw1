package biz

import (
	biz "biz/proto"
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var grpcPort = ":5062"

var conn *grpc.ClientConn
var client biz.BizServerClient

func Init() {
	conn, err := grpc.Dial(grpcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("did not connect: %v", err)
	}
	client = biz.NewBizServerClient(conn)
}

func Close() {
	conn.Close()
}

func GetUsers(authKey, userId string, withSqlInject bool) (any, error) {
	var userList *biz.UsersList
	var err error
	if withSqlInject {
		userList, err = client.GetUsersWithSqlInject(context.Background(), &biz.UserAuth{
			UserId:  string(userId),
			AuthKey: authKey,
		})
	} else {
		userList, err = client.GetUsers(context.Background(), &biz.UserAuth{
			UserId:  string(userId),
			AuthKey: authKey,
		})
	}
	if err != nil {
		return nil, err
	}
	return userList, nil
}
