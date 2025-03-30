package server

import (
	"context"
	"github.com/krekio/TagesTest/internal/storage"
	pb "github.com/krekio/TagesTest/protos"
)

type fileServiceServer struct {
	pb.UnimplementedFileServiceServer
	fileStorage           *storage.FileStorage
	uploadDownloadLimiter chan struct{}
	listLimiter           chan struct{}
}

func NewFileServiceServer(storage *storage.FileStorage) pb.FileServiceServer {
	return &fileServiceServer{
		fileStorage:           storage,
		uploadDownloadLimiter: make(chan struct{}, 10),
		listLimiter:           make(chan struct{}, 100),
	}
}

func (s *fileServiceServer) UploadFile(stream pb.FileService_UploadFileServer) error {

	s.uploadDownloadLimiter <- struct{}{}
	defer func() { <-s.uploadDownloadLimiter }()

	return s.fileStorage.Upload(stream)
}

func (s *fileServiceServer) ListFiles(ctx context.Context, req *pb.ListFilesRequest) (*pb.ListFilesResponse, error) {
	s.listLimiter <- struct{}{}
	defer func() { <-s.listLimiter }()

	return s.fileStorage.List()
}

func (s *fileServiceServer) DownloadFile(req *pb.DownloadFileRequest, stream pb.FileService_DownloadFileServer) error {

	s.uploadDownloadLimiter <- struct{}{}
	defer func() { <-s.uploadDownloadLimiter }()

	return s.fileStorage.Download(req.Filename, stream)
}
