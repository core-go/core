package upload

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
)

type UploadService interface {
	Upload(ctx context.Context, id string, data Request) (*Upload, error)
	Delete(ctx context.Context, id string, url string) (int64, error)
}

type StoragePort interface {
	Upload(ctx context.Context, directory string, filename string, data []byte, contentType string) (string, error)
	Delete(ctx context.Context, id string) (bool, error)
}

type UploadUseCase struct {
	repository       StorageRepository
	Service          StoragePort
	Provider         string
	GeneralDirectory string
	Directory        string
	KeyFile          string
}

func NewUploadService(
	repository StorageRepository,
	service StoragePort, provider string, generalDirectory string,
	keyFile string, directory string) UploadService {

	return &UploadUseCase{Service: service, Provider: provider, GeneralDirectory: generalDirectory,
		KeyFile: keyFile, Directory: directory, repository: repository}
}

func (u *UploadUseCase) Upload(ctx context.Context, id string, req Request) (*Upload, error) {
	_, err := u.repository.Load(ctx, id)
	if err != nil {
		return nil, err
	}
	url, err := u.uploadFileOnServer(ctx, req.Filename, req.Type, req.Size, req.Data)
	if err != nil {
		return nil, err
	}
	attachment := Upload{
		OriginalFileName: req.OriginalFileName,
		FileName:         req.Filename,
		Type:             req.Type,
		Size:             req.Size,
		Url:              url,
	}
	rows, err := u.repository.Update(ctx, id, attachment)
	if err != nil {
		return nil, err
	}
	if rows > 0 {
		return &attachment, nil
	}
	return nil, nil
}

func (u *UploadUseCase) Delete(ctx context.Context, id string, url string) (int64, error) {
	attachment, err := u.repository.Load(ctx, id)
	if err != nil {
		return 0, err
	}
	if attachment == nil {
		return -1, errors.New("not found item")
	}
	exist := false
	if attachment.Url == url {
		_, err2 := u.deleteFile(url, ctx)
		if err2 != nil {
			return 0, err2
		}
		exist = true
	}
	if exist == false {
		return -1, errors.New("no exist file " + url)
	}
	_, err2 := u.repository.Update(ctx, id, *attachment)
	if err2 != nil {
		return 0, err2
	}
	return 1, nil
}

func (u *UploadUseCase) uploadFileOnServer(ctx context.Context, fileName string, contentType string, size int64, data []byte) (rs string, errorRespone error) {
	directory := u.Directory
	rs, err2 := u.Service.Upload(ctx, directory, fileName, data, contentType)
	if err2 != nil {
		return rs, err2
	}
	return
}

func (u *UploadUseCase) deleteFile(url string, ctx context.Context) (bool, error) {
	arrOrigin := strings.Split(url, "/")
	delOriginUrl := arrOrigin[len(arrOrigin)-2] + "/" + arrOrigin[len(arrOrigin)-1]
	rs, err := u.Service.Delete(ctx, delOriginUrl)
	return rs, err
}

func getExt(file string) string {
	ext := filepath.Ext(file)
	if strings.HasPrefix(ext, ":") {
		ext = ext[1:]
		return ext
	}
	return ext
}
