package fileservice

import (
	"context"
	"time"

	FS "github.com/TOMMy-Net/tages/pkg/protobuf/file_service/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)
func (fs *FileServer) ListFiles(ctx context.Context, _ *FS.ListRequest) (*FS.ListResponse, error) {
	clientIP := getClientIP(ctx)
	
	limits := fs.limiter.GetLimits(clientIP)
	if !limits.ListFilesConn.Add() {
		return &FS.ListResponse{}, status.Error(codes.PermissionDenied, "many connections")
	}
	defer limits.ListFilesConn.Release()


	fs.metaDataMutex.RLock()
	defer fs.metaDataMutex.RUnlock()

	files := make([]*FS.FileInfo, 0, len(fs.metaData))
	for name, meta := range fs.metaData {
		files = append(files, &FS.FileInfo{
			Filename:  name,
			CreatedAt: meta.CreatedAt.Format(time.RFC3339),
			UpdatedAt: meta.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &FS.ListResponse{Files: files}, nil
}