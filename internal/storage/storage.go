package storage

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"time"

	pb "github.com/krekio/TagesTest/protos"
)

type FileStorage struct {
	storagePath string
}

func NewFileStorage(path string) (*FileStorage, error) {
	if _, err := os.Stat(path); err != nil {
		if err = os.MkdirAll(path, os.ModePerm); err != nil {
			return nil, err
		}
	}

	return &FileStorage{storagePath: path}, nil
}

func (s *FileStorage) Upload(stream pb.FileService_UploadFileServer) error {
	var filename string
	var file *os.File

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if file == nil {
			filename = filepath.Join(s.storagePath, req.GetFilename())
			file, err = os.Create(filename)
			if err != nil {
				return err
			}
			defer file.Close()
		}

		if _, err := file.Write(req.GetData()); err != nil {
			return err
		}
	}

	return stream.SendAndClose(&pb.UploadFileResponse{Message: "Файл успешно загружен"})
}

func (s *FileStorage) List() (*pb.ListFilesResponse, error) {
	files, err := os.ReadDir(s.storagePath)
	if err != nil {
		return nil, err
	}

	var fileInfos []*pb.FileInfo
	for _, f := range files {
		if f.IsDir() {
			continue
		}

		filePath := filepath.Join(s.storagePath, f.Name())
		fileStat, err := os.Stat(filePath)
		if err != nil {
			return nil, err
		}

		fileInfos = append(fileInfos, &pb.FileInfo{
			Filename:  f.Name(),
			CreatedAt: fileStat.ModTime().Format(time.RFC3339),
			UpdatedAt: fileStat.ModTime().Format(time.RFC3339),
		})
	}

	return &pb.ListFilesResponse{Files: fileInfos}, nil
}

func (s *FileStorage) Download(filename string, stream pb.FileService_DownloadFileServer) error {
	filePath := filepath.Join(s.storagePath, filename)
	file, err := os.Open(filePath)
	if err != nil {
		return errors.New("Ошибка скачивания файла")
	}
	defer file.Close()

	buf := make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if err := stream.Send(&pb.DownloadFileResponse{Data: buf[:n]}); err != nil {
			return err
		}
	}
	return nil
}
