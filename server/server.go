package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"

	pb "tages-grpc/proto"
	myfile "tages-grpc/server/file"

	"github.com/djherbis/times"
	"google.golang.org/grpc"
)

type myserver struct {
	pb.TagesServiceServer
}

func main() {
	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("Fail start server: %v", err)
	}
	defer lis.Close()
	grpcServer := grpc.NewServer()
	pb.RegisterTagesServiceServer(grpcServer, &myserver{})
	log.Printf("Server started: %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Fail start: %v", err)
	}

}

func (srv *myserver) DownloadImage(req *pb.Messages, responseStream pb.TagesService_DownloadImageServer) error {
	log.Println("Started download file")
	bufferSize := 64 * 1024
	file, err := os.Open(req.GetMes())
	if err != nil {
		log.Println(err)
		return err
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
		resp := &pb.SendChunks{
			Chunk: buff[:bytesRead],
		}
		err = responseStream.Send(resp)
		if err != nil {
			log.Println("error while sending chunk:", err)
			return err
		}
	}
	log.Println("Finished download file")
	return nil
}

func (srv *myserver) UploadImage(stream pb.TagesService_UploadImageServer) error {
	log.Println("Started download file")
	file := myfile.NewFile()
	for {
		chunkResponse, err := stream.Recv()
		if file.GetName() == "" {
			file.EditName(chunkResponse.GetMes())
		}
		if err == io.EOF {
			errFile := myfile.Store(file)
			if errFile != nil {
				log.Println("error create file: ", errFile)
				return errFile
			}
			log.Println("received all chunks")
			break
		}
		if err != nil {
			log.Println("error receiving chunk:", err)
			break
		}
		log.Println("got new chunk")
		errFile := file.Write(chunkResponse.Chunk)
		if errFile != nil {
			log.Println("error write file: ", errFile)
			break
		}
	}
	log.Println("Finished download file")
	return stream.SendAndClose(&pb.Messages{Mes: "done upload"})
}

func (srv *myserver) ListImage(req *pb.NoParam, stream pb.TagesService_ListImageServer) error {
	log.Println("Started send list files")
	dir := "server/files/"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if !file.IsDir() {
			t, _ := times.Stat(dir + file.Name())
			mes := &pb.Messages{
				Mes: fmt.Sprintf("%s | %v | %v \n", file.Name(), t.BirthTime(), t.ModTime()),
			}
			err := stream.Send(mes)
			if err != nil {
				log.Println("error send mes: ", err)
				return err
			}
		}
	}
	log.Println("Finished send list files")
	return nil
}
