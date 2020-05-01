package main

import (
	"context"
	"fmt"
)

func main() {
	//lambda.Start(DeleteVideoHandler)
	fmt.Println(DeleteVideoHandler(context.TODO(), Data{UserId: 420, VideoId: "420/VMc1MnmUZL4WYVKkgYzWW.mp4", Password: "vHQHmwysxamWwkXHffxy"}))
}
