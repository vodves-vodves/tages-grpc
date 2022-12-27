package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"tages-grpc/client/menu"
	pb "tages-grpc/proto"

	ratelimit "github.com/tommy-sho/rate-limiter-grpc-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	menu.PrintMenu([]string{"Upload file to server", "Download file from server", "File list", "Exit"})
	choice := menu.UserInputInteger("Choose your option")
	switch choice {

	default:
		fmt.Println("Invalid choice")

	case 1:
		fileName := menu.UserInput("Enter file name")
		if fileName == "" || fileName == "\n" {
			log.Println("Invalid file name")
			os.Exit(0)
		}
		conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithStreamInterceptor(ratelimit.StreamClientInterceptor(ratelimit.NewLimiter(10))))
		if err != nil {
			log.Fatal("client can not connect to grpc service:", err)
		}
		c := pb.NewTagesServiceClient(conn)
		fUploadImage(c, fileName)

	case 2:
		fileName := menu.UserInput("Enter file name")
		if fileName == "" || fileName == "\n" {
			log.Println("Invalid file name")
			os.Exit(0)
		}
		conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithStreamInterceptor(ratelimit.StreamClientInterceptor(ratelimit.NewLimiter(10))))
		if err != nil {
			log.Fatal("client can not connect to grpc service:", err)
		}
		c := pb.NewTagesServiceClient(conn)
		fDownloadImage(c, fileName)

	case 3:
		conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithStreamInterceptor(ratelimit.StreamClientInterceptor(ratelimit.NewLimiter(100))))
		if err != nil {
			log.Fatal("client can not connect to grpc service:", err)
		}
		c := pb.NewTagesServiceClient(conn)
		fListImage(c)

	case 4:
		os.Exit(0)
	}
}

// download file from server
func fDownloadImage(c pb.TagesServiceClient, fileName string) {
	saveFileName := "client/files/" + fileName
	fileStreamResponse, err := c.DownloadImage(context.Background(), &pb.Messages{
		Mes: "server/files/" + fileName,
	})
	if err != nil {
		log.Println("error downloading:", err)
		return
	}
	file, err := os.Create(saveFileName)
	if err != nil {
		log.Println("error create file: ", err)
		return
	}
	defer file.Close()
	for {
		chunkResponse, err := fileStreamResponse.Recv()
		if err == io.EOF {
			log.Println("received all chunks")
			break
		}
		if err != nil {
			log.Println("error receiving chunk: ", err)
			break
		}
		log.Printf("got new chunk with data \n")
		_, err = file.Write(chunkResponse.Chunk)
		if err != nil {
			log.Println("error write file: ", err)
			break
		}
	}
}

// upload file to server
func fUploadImage(c pb.TagesServiceClient, fileName string) {
	bufferSize := 64 * 1024
	fileStream, err := c.UploadImage(context.Background())
	if err != nil {
		log.Println("error downloading: ", err)
		return
	}
	fileDir := "client/files/" + fileName
	file, err := os.Open(fileDir)
	if err != nil {
		log.Println("no such file: ", err)
		return
	}
	defer file.Close()
	buff := make([]byte, bufferSize)
	for {
		bytesRead, err := file.Read(buff)
		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}
			break
		}
		resp := &pb.UploadRequest{
			Mes:   fileName,
			Chunk: buff[:bytesRead],
		}
		err = fileStream.Send(resp)
		if err != nil {
			log.Println("error while sending chunk: ", err)
			return
		}
	}
	res, err := fileStream.CloseAndRecv()
	if err != nil {
		log.Println("error receive message", err)
		return
	}
	log.Println(res.GetMes())
}

// send list files
func fListImage(c pb.TagesServiceClient) {
	stream, err := c.ListImage(context.Background(), &pb.NoParam{})
	if err != nil {
		log.Println("error receiving:", err)
		return
	}

	for {
		message, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("error receiving files list: ", err)
			break
		}
		log.Println(message.GetMes())
	}
}
