package main

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"

	"github.com/TiagoBehencks/gRPC/api/db"
	"github.com/TiagoBehencks/gRPC/api/handlers"
	"github.com/TiagoBehencks/gRPC/api/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	if err := db.Connect(); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.InitSchema(); err != nil {
		log.Fatalf("failed to init schema: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/products", handleProducts)
	mux.HandleFunc("/api/products/", handleProductByID)

	log.Printf("HTTP server listening at http://localhost:8080")
	log.Printf("gRPC server listening at :50051")

	go http.ListenAndServe(":8080", mux)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	gs := grpc.NewServer()
	pb.RegisterProductServiceServer(gs, handlers.NewProductServer())
	reflection.Register(gs)
	gs.Serve(lis)
}

var clientConn *grpc.ClientConn

func getClient() (*grpc.ClientConn, error) {
	if clientConn == nil {
		conn, err := grpc.Dial(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, err
		}
		clientConn = conn
	}
	return clientConn, nil
}

func handleProducts(w http.ResponseWriter, r *http.Request) {
	conn, err := getClient()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	client := pb.NewProductServiceClient(conn)

	switch r.Method {
	case http.MethodGet:
		resp, err := client.List(r.Context(), &pb.ListProductsRequest{})
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		jsonResp, _ := json.Marshal(resp.Products)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResp)

	case http.MethodPost:
		body, _ := io.ReadAll(r.Body)
		var req pb.CreateProductRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		resp, err := client.Create(r.Context(), &req)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		jsonResp, _ := json.Marshal(resp)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResp)
	}
}

func handleProductByID(w http.ResponseWriter, r *http.Request) {
	conn, err := getClient()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	client := pb.NewProductServiceClient(conn)

	id := r.URL.Path[len("/api/products/"):]
	if id == "" {
		http.Error(w, "missing id", 400)
		return
	}

	switch r.Method {
	case http.MethodGet:
		resp, err := client.Get(r.Context(), &pb.GetProductRequest{Id: id})
		if err != nil {
			http.Error(w, err.Error(), 404)
			return
		}
		jsonResp, _ := json.Marshal(resp)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResp)

	case http.MethodPut:
		body, _ := io.ReadAll(r.Body)
		var req pb.UpdateProductRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		req.Id = id
		resp, err := client.Update(r.Context(), &req)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		jsonResp, _ := json.Marshal(resp)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResp)

	case http.MethodDelete:
		_, err := client.Delete(r.Context(), &pb.DeleteProductRequest{Id: id})
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(204)
	}
}