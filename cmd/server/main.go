package main

import (
	"database/sql"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	appGrpc "github.com/ManuelJ0aquim/go-grpc-inventory-service/internal/infrastructure/grpc"
	inventoryv1 "github.com/ManuelJ0aquim/go-grpc-inventory-service/internal/infrastructure/grpc/pb/inventory/v1"
	"github.com/ManuelJ0aquim/go-grpc-inventory-service/internal/infrastructure/repository"
	"github.com/ManuelJ0aquim/go-grpc-inventory-service/internal/service"
)

func main() {
	dbConnStr := os.Getenv("DATABASE_URL")
	if dbConnStr == "" {
		dbConnStr = "postgres://postgres:postgres@localhost:5432/inventory_db?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatalf("Falha ao abrir conexao com PostgreSQL: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatalf("Falha ao conectar ao PostgreSQL: %v", err)
	}
	log.Println("Conectado ao PostgreSQL com sucesso!")

	repo := repository.NewPostgresRepository(db)
	svc := service.NewInventoryService(repo)
	handler := appGrpc.NewHandler(svc)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Falha ao abrir porta: %v", err)
	}

	grpcServer := grpc.NewServer()
	inventoryv1.RegisterInventoryServiceServer(grpcServer, handler)
	reflection.Register(grpcServer)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("Servidor gRPC rodando na porta :50051...")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Erro ao executar servidor gRPC: %v", err)
		}
	}()

	<-stop
	log.Println("Encerrando servidor graciosamente...")
	grpcServer.GracefulStop()
	log.Println("Servidor finalizado.")
}