package server

import (
	"context"
	"grpc-file-service/internal/storage"
	pb "grpc-file-service/proto"
)

type fileServiceServer struct {
	pb.UnimplementedFileServiceServer
	fileStorage           *storage.FileStorage
	uploadDownloadLimiter chan struct{}
	listLimiter           chan struct{}
}

// NewFileServiceServer Creating a new FileService Server
//
// This function creates a new FileServiceServer object.
//
// It takes a single argument, storage, which is a pointer to a
// storage.FileStorage object. This object is used to store and retrieve
// files.
//
// The function returns a pointer to a fileServiceServer object, which
// implements the FileServiceServer interface.
func NewFileServiceServer(storage *storage.FileStorage) pb.FileServiceServer {
	return &fileServiceServer{
		fileStorage:           storage,
		uploadDownloadLimiter: make(chan struct{}, 10),  // Ограничение на 10 одновременных запросов загрузки/скачивания
		listLimiter:           make(chan struct{}, 100), // Ограничение на 100 одновременных запросов списка
	}
}

// UploadFile handles the uploading of a file through a stream.
//
// It uses a limiter to restrict the number of simultaneous uploads/downloads.
// After acquiring the limiter, it delegates the upload process to the
// FileStorage's Upload method.
func (s *fileServiceServer) UploadFile(stream pb.FileService_UploadFileServer) error {
	// Acquire limiter slot to restrict simultaneous uploads/downloads
	s.uploadDownloadLimiter <- struct{}{}
	defer func() { <-s.uploadDownloadLimiter }() // Release limiter slot on function exit

	// Delegate the upload process to FileStorage
	return s.fileStorage.Upload(stream)
}

// ListFiles Viewing a list of files
//
// This function returns a list of files stored in the FileStorage.
//
// It uses a limiter to restrict the number of simultaneous requests to
// retrieve the list of files.
func (s *fileServiceServer) ListFiles(ctx context.Context, req *pb.ListFilesRequest) (*pb.ListFilesResponse, error) {
	// Acquire limiter slot to restrict simultaneous list requests
	s.listLimiter <- struct{}{}
	defer func() { <-s.listLimiter }() // Release limiter slot on function exit

	// Delegate the list retrieval process to FileStorage
	return s.fileStorage.List()
}

// DownloadFile Downloading a file
//
// It takes a single argument, req, which is a pointer to a
// pb.DownloadFileRequest object. This object contains the filename
// to be downloaded.
//
// It uses a limiter to restrict the number of simultaneous uploads/downloads.
// After acquiring the limiter, it delegates the download process to the
// FileStorage's Download method.
func (s *fileServiceServer) DownloadFile(req *pb.DownloadFileRequest, stream pb.FileService_DownloadFileServer) error {
	// Acquire limiter slot to restrict simultaneous uploads/downloads
	s.uploadDownloadLimiter <- struct{}{}
	defer func() { <-s.uploadDownloadLimiter }() // Release limiter slot on function exit

	// Delegate the download process to FileStorage
	return s.fileStorage.Download(req.Filename, stream)
}
