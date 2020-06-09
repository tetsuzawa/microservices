package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"

	"github.com/tetsuzawa/microservices/backend/pkg/api"
)

func main() {
	// get configuration
	address := flag.String("server", "", "gRPC server in format host:port")
	flag.Parse()

	// Set up a connection to the server.
	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := api.NewPostServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Call Create
	req1 := api.CreatePostRequest{
		UserId: "11111111-1111-1111-1111-111111111111",
		Text:   "test post 1",
	}
	res1, err := c.CreatePost(ctx, &req1)
	if err != nil {
		log.Fatalf("CreatePost failed: %v", err)
	}
	log.Printf("CreatePost result: <%+v>\n\n", res1)

	// Read
	req2 := api.GetPostRequest{
		Id: res1.Post.Id,
	}
	res2, err := c.GetPost(ctx, &req2)
	if err != nil {
		log.Fatalf("GetPost failed: %v", err)
	}
	log.Printf("GetPost result: <%+v>\n\n", res2)

	// Update
	req3 := api.UpdatePostRequest{
		Id:     res1.Post.Id,
		UserId: res1.Post.UserId,
		Text:   "updated post 1",
	}
	res3, err := c.UpdatePost(ctx, &req3)
	if err != nil {
		log.Fatalf("UpdatePost failed: %v", err)
	}
	log.Printf("UpdatePost result: <%+v>\n\n", res3)

	// Call ReadAll
	req4 := api.ListPostsRequest{}
	res4, err := c.ListPosts(ctx, &req4)
	if err != nil {
		log.Fatalf("ListPosts failed: %v", err)
	}
	log.Printf("ListPosts result: <%+v>\n\n", res4)

	// Delete
	req5 := api.DeletePostRequest{
		Id:     res1.Post.Id,
		UserId: res1.Post.UserId,
	}
	res5, err := c.DeletePost(ctx, &req5)
	if err != nil {
		log.Fatalf("Delete failed: %v", err)
	}
	log.Printf("Delete result: <%+v>\n\n", res5)
}
