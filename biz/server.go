package main

import (
	biz "biz/proto"
	"context"
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type BizServer struct {
	biz.UnimplementedBizServerServer

	db *sql.DB

	redisCli *redis.Client
}

func (b *BizServer) getUserFromQuery(rows *sql.Rows) (*biz.User, error) {
	var name, family, sex, createdAt string
	var id, age int32

	err := rows.Scan(&name, &family, &id, &age, &sex, &createdAt)
	if err != nil {
		return nil, err
	}
	return &biz.User{Name: name, Family: family, Id: id, Age: age, Sex: sex, CreatedAt: createdAt}, nil
}

func (b *BizServer) checkAuth(c context.Context, auth int32) error {
	_, err := b.redisCli.Get(c, fmt.Sprintf("authKey:%d", auth)).Result()
	return err
}

func (b *BizServer) GetUsers(c context.Context, user *biz.UserAuth) (*biz.UsersList, error) {
	err := b.checkAuth(c, user.AuthKey)
	if err != nil {
		return nil, err
	}
	var usersListQuery *sql.Rows
	if id, parseErr := strconv.ParseInt(user.UserId, 10, 64); parseErr != nil && id != 0 {
		usersListQuery, err = b.db.Query(fmt.Sprintf("SELECT * FROM USERS WHERE id=%s", user.UserId))
	} else {
		usersListQuery, err = b.db.Query("SELECT * FROM USERS ORDER BY id LIMIT 100")
	}
	if err != nil {
		return nil, err
	}
	defer usersListQuery.Close()
	var usersList biz.UsersList
	for usersListQuery.Next() {
		currUser, err := b.getUserFromQuery(usersListQuery)
		if err != nil {
			return nil, err
		}
		usersList.Users = append(usersList.Users, currUser)
	}
	return &usersList, nil
}

func (b *BizServer) GetUsersWithSqlInject(c context.Context, user *biz.UserAuth) (*biz.UsersList, error) {
	err := b.checkAuth(c, user.AuthKey)
	if err != nil {
		return nil, err
	}
	var usersListQuery *sql.Rows
	if user.UserId != "0" {
		usersListQuery, err = b.db.Query("SELECT * FROM USERS WHERE id=" + user.UserId)
	} else {
		usersListQuery, err = b.db.Query("SELECT * FROM USERS ORDER BY id LIMIT 100")
	}
	if err != nil {
		return nil, err
	}
	defer usersListQuery.Close()
	var usersList biz.UsersList
	for usersListQuery.Next() {
		currUser, err := b.getUserFromQuery(usersListQuery)
		if err != nil {
			return nil, err
		}
		usersList.Users = append(usersList.Users, currUser)
	}
	return &usersList, nil
}

func NewBizServer() (*BizServer, error) {
	connStr := "postgres://postgres:postgres@localhost/web-hw1?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	redisCli := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	if err != nil {
		return nil, err
	} else if err = db.Ping(); err != nil {
		return nil, err
	} else {
		bizServer := &BizServer{db: db, redisCli: redisCli}
		return bizServer, nil
	}
}

func main() {
	bizServer, err := NewBizServer()
	if err != nil {
		panic(err)
	}
	fmt.Print(bizServer.GetUsers(context.Background(), &biz.UserAuth{UserId: "0"}))
}
