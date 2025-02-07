package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"

	desc "github.com/AndreiAvdko/auth/pkg/auth_v1"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	grpcPort = 50051
	dbDSN    = "host=localhost port=54321 dbname=auth user=authuser password=authpassword sslmode=disable"
	DB_NAME  = "auth_users"
)

type server struct {
	desc.UnimplementedAuthV1Server
	Pool *pgxpool.Pool
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {

	// Делаем запрос на выборку записей из таблицы note
	builderSelect := sq.Select("id", "name", "email", "password", "password_confirm", "is_admin", "create_time", "update_time").
		From(DB_NAME).
		PlaceholderFormat(sq.Dollar).
		OrderBy("id ASC").
		Where(sq.Eq{"id": req.Id})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	rows, err := s.Pool.Query(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to find user: %v", err)
	}

	var id int
	var name, email, password, passwordConfirm string
	var isAdmin bool
	var createTime time.Time
	var updateTime sql.NullTime

	for rows.Next() {
		err = rows.Scan(&id, &name, &email, &password, &passwordConfirm, &isAdmin, &createTime, &updateTime)
		if err != nil {
			log.Fatalf("failed to scan user: %v", err)
		}

		log.Printf("id: %d, name: %s, email: %s, password: %s, passwordConfirm: %s, isAdmin: %t, createTime: %v, updateTime: %v\n", id, name, email, password, passwordConfirm, isAdmin, createTime, updateTime)
	}

	log.Println("==============================")
	log.Printf("User id: %d", req.GetId())

	return &desc.GetResponse{
		User: &desc.User{
			Id:             (int64(id)),
			Name:           name,
			Email:          email,
			Password:       password,
			PassworConfirm: passwordConfirm,
			IsAdmin:        isAdmin,
			CreatedAt:      (timestamppb.New(createTime)),
			UpdatedAt:      nil,
		},
	}, nil
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	// Делаем запрос на вставку записи в таблицу note
	builderInsert := sq.Insert(DB_NAME).
		PlaceholderFormat(sq.Dollar).
		Columns("name", "email", "password", "password_confirm", "is_admin").
		Values(req.GetName(), req.GetEmail(), req.GetPassword(), req.GetPassworConfirm(), req.GetIsAdmin()).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	var userID int
	err = s.Pool.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		log.Fatalf("failed to insert note: %v", err)
	}

	// Выводим данные созданного пользователя

	log.Println("===========================================")
	log.Printf("inserted note with id: %d", userID)
	log.Printf("User name: %s", req.GetName())
	log.Printf("User email: %s", req.GetEmail())
	log.Printf("User password: %s", req.GetPassword())
	log.Printf("User password_confirm: %s", req.GetPassworConfirm())
	log.Printf("User is_admin: %t", req.GetIsAdmin())

	return &desc.CreateResponse{
		Id: int64(userID),
	}, nil
}

func main() {
	ctx := context.Background()

	// Создаем пул соединений с базой данных
	pool, err := pgxpool.Connect(ctx, dbDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	// // Делаем запрос на обновление записи в таблице note
	// builderUpdate := sq.Update("note").
	// 	PlaceholderFormat(sq.Dollar).
	// 	Set("title", gofakeit.City()).
	// 	Set("body", gofakeit.Address().Street).
	// 	Set("updated_at", time.Now()).
	// 	Where(sq.Eq{"id": noteID})

	// query, args, err = builderUpdate.ToSql()
	// if err != nil {
	// 	log.Fatalf("failed to build query: %v", err)
	// }

	// res, err := pool.Exec(ctx, query, args...)
	// if err != nil {
	// 	log.Fatalf("failed to update note: %v", err)
	// }

	// log.Printf("updated %d rows", res.RowsAffected())

	// // Делаем запрос на получение измененной записи из таблицы note
	// builderSelectOne := sq.Select("id", "title", "body", "created_at", "updated_at").
	// 	From("note").
	// 	PlaceholderFormat(sq.Dollar).
	// 	Where(sq.Eq{"id": noteID}).
	// 	Limit(1)

	// query, args, err = builderSelectOne.ToSql()
	// if err != nil {
	// 	log.Fatalf("failed to build query: %v", err)
	// }

	// err = pool.QueryRow(ctx, query, args...).Scan(&id, &title, &body, &createdAt, &updatedAt)
	// if err != nil {
	// 	log.Fatalf("failed to select notes: %v", err)
	// }

	// log.Printf("id: %d, title: %s, body: %s, created_at: %v, updated_at: %v\n", id, title, body, createdAt, updatedAt)

	// ///////////////////////////////////////////////////

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterAuthV1Server(s, &server{Pool: pool})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
