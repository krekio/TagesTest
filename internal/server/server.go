package server

import (
	"context"

	"github.com/krekio/TagesTest/internal/storage"
	pb "github.com/krekio/TagesTest/protos"
	"golang.org/x/sync/semaphore"
)

type FileServiceServer struct {
	fileStorage           *storage.FileStorage
	uploadDownloadLimiter *semaphore.Weighted
	listLimiter           *semaphore.Weighted
}

func NewFileServiceServer(storage *storage.FileStorage) *FileServiceServer {
	return &FileServiceServer{
		fileStorage:           storage,
		uploadDownloadLimiter: semaphore.NewWeighted(10),
		listLimiter:           semaphore.NewWeighted(100),
	}
}

func (s *FileServiceServer) UploadFile(stream pb.FileService_UploadFileServer) error {
	if err := s.uploadDownloadLimiter.Acquire(context.Background(), 1); err != nil {
		return err
	}
	defer s.uploadDownloadLimiter.Release(1)

	return s.fileStorage.Upload(stream)
}

func (s *FileServiceServer) ListFiles(ctx context.Context, req *pb.ListFilesRequest) (*pb.ListFilesResponse, error) {
	if err := s.listLimiter.Acquire(context.Background(), 1); err != nil {
		return nil, err
	}
	defer s.listLimiter.Release(1)

	return s.fileStorage.List()
}

func (s *FileServiceServer) DownloadFile(req *pb.DownloadFileRequest, stream pb.FileService_DownloadFileServer) error {

	if err := s.uploadDownloadLimiter.Acquire(context.Background(), 1); err != nil {
		return err
	}
	defer s.uploadDownloadLimiter.Release(1)

	return s.fileStorage.Download(req.Filename, stream)
}
