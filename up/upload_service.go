package upload

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

type UploadService interface {
	UploadGallery(data Upload, r *http.Request) ([]UploadInfo, error)
	DeleteGalleryFile(id string, url string, r *http.Request) (int64, error)
	UploadCover(id string, data []UploadData, contentType string, r *http.Request) (string, error)
	UploadImage(id string, data []UploadData, contentType string, r *http.Request) (string, error)
	UpdateGallery(data []UploadInfo, id string, r *http.Request) (int64, error)
	GetGallery(id string, r *http.Request) ([]UploadInfo, error)
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
	SizesImage       []int
	SizesCover       []int
}

func NewUploadService(
	repository StorageRepository,
	service StoragePort, provider string, generalDirectory string,
	keyFile string, directory string,
	sizesCover []int,
	sizesImage []int) UploadService {

	var sizesI = []int{40, 400}
	var sizesC = []int{576, 768}

	if len(sizesCover) != 0 {
		sizesC = sizesCover
	}
	if len(sizesImage) != 0 {
		sizesI = sizesCover
	}
	return &UploadUseCase{Service: service, Provider: provider, GeneralDirectory: generalDirectory,
		KeyFile: keyFile, Directory: directory,
		SizesImage: sizesI, SizesCover: sizesC, repository: repository}
}
func (u *UploadUseCase) UploadCover(id string, data []UploadData, contentType string, r *http.Request) (string, error) {
	//delete
	result, err := u.repository.Load(r.Context(), id)
	if err != nil {
		return "", err
	}
	if result == nil {
		return "", nil
	}

	if result.CoverURL != nil {
		_, err := u.DeleteFileUpload(u.SizesCover, *result.CoverURL, r)
		if err != nil {
			return "", errors.New("internal server error")
		}
	}
	//upload
	var newUrl string
	for i := range data {
		file := data[i]
		rs, errorRespone := u.UploadFile(file.Name, contentType, file.Data, r)
		if errorRespone != nil {
			return "", errorRespone
		}
		if i != 0 {
			continue
		}
		newUrl = rs
	}
	user := UploadModel{Id: id, CoverURL: &newUrl}

	_, err1 := u.Update(r.Context(), user)
	if err1 != nil {
		return "", err1
	}
	return newUrl, nil
}

func (u *UploadUseCase) UploadImage(id string, data []UploadData, contentType string, r *http.Request) (string, error) {
	//delete
	result, err := u.repository.Load(r.Context(), id)
	if err != nil {
		return "", err
	}
	if result == nil {
		return "", nil
	}

	if result.ImageURL != nil {
		_, err := u.DeleteFileUpload(u.SizesImage, *result.ImageURL, r)
		if err != nil {
			return "", errors.New("internal server error")
		}
	}
	//upload
	var newUrl string
	for i := range data {
		file := data[i]
		rs, errorRespone := u.UploadFile(file.Name, contentType, file.Data, r)
		if errorRespone != nil {
			return "", errorRespone
		}
		if i != 0 {
			continue
		}
		newUrl = rs
	}
	user := UploadModel{Id: id, ImageURL: &newUrl}

	_, err1 := u.Update(r.Context(), user)
	if err1 != nil {
		return "", err1
	}
	return newUrl, nil
}

func (u *UploadUseCase) UploadGallery(data Upload, r *http.Request) ([]UploadInfo, error) {
	result, err := u.repository.Load(r.Context(), data.Id)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	var gallery []UploadInfo
	if result.Gallery != nil {
		gallery = result.Gallery
	}

	rs, errorRespone := u.UploadFile(data.Name, data.Type, data.Data, r)
	if errorRespone != nil {
		return nil, errorRespone
	}

	gallery = append(gallery, UploadInfo{Source: data.Source, Type: strings.Split(data.Type, "/")[0], Url: rs})
	user := UploadModel{Id: data.Id, Gallery: gallery}

	_, err = u.Update(r.Context(), user)
	if err != nil {
		return nil, err
	}
	return gallery, nil
}

func (u *UploadUseCase) UploadFile(fileName string, contentType string, data []byte, r *http.Request) (rs string, errorRespone error) {
	directory := u.Directory
	rs, err2 := u.Service.Upload(r.Context(), directory, fileName, data, contentType)
	if err2 != nil {
		errorRespone = err2
		return
	}
	return
}

func (u *UploadUseCase) DeleteGalleryFile(id string, url string, r *http.Request) (int64, error) {
	rs, err := u.repository.Load(r.Context(), id)
	if err != nil {
		return 0, err
	}
	if rs == nil {
		return -1, nil
	}
	gallery := rs.Gallery
	// if find url in gallery
	idx := -1
	for i := range gallery {
		if gallery[i].Url == url {
			idx = i
		}
	}
	if idx == -1 {
		return 0, nil
	}
	_, err2 := u.DeleteFile(url, r)
	if err2 != nil {
		return 0, err2
	}
	gallery = append(gallery[:idx], gallery[idx+1:]...)
	user := UploadModel{Id: id, Gallery: gallery}
	_, err3 := u.Update(r.Context(), user)
	if err3 != nil {
		return 0, err3
	}
	return 1, nil
}

func (u *UploadUseCase) GetGallery(id string, r *http.Request) ([]UploadInfo, error) {
	rs, err := u.repository.Load(r.Context(), id)
	if err != nil {
		return nil, err
	}
	if rs == nil {
		return nil, nil
	}
	return rs.Gallery, err
}

func (u *UploadUseCase) UpdateGallery(data []UploadInfo, id string, r *http.Request) (int64, error) {
	user := UploadModel{Id: id, Gallery: data}
	_, err2 := u.Update(r.Context(), user)
	if err2 != nil {
		return 0, err2
	}
	return 1, err2
}
func (u *UploadUseCase) DeleteFileUpload(sizes []int, url string, r *http.Request) (bool, error) {
	rs, err := u.DeleteFile(url, r)
	fmt.Print(rs, err)
	// if err != nil {
	// 	return false, errors.New("internal server error")
	// }
	for i := range sizes {
		oldUrl := removeExt(url) + "_" + strconv.Itoa(sizes[i]) + getExt(url)
		arr := strings.Split(oldUrl, "/")
		delUrl := arr[len(arr)-2] + "/" + arr[len(arr)-1]
		rss, err := u.DeleteFile(delUrl, r)
		fmt.Print(rss, err)
	}
	return true, nil
}

func (u *UploadUseCase) DeleteFile(url string, r *http.Request) (bool, error) {
	arrOrigin := strings.Split(url, "/")
	delOriginUrl := arrOrigin[len(arrOrigin)-2] + "/" + arrOrigin[len(arrOrigin)-1]
	rs, err := u.Service.Delete(r.Context(), delOriginUrl)
	return rs, err
}

func (u *UploadUseCase) Update(ctx context.Context, user UploadModel) (bool, error) {
	_, err := u.repository.Update(ctx, user)
	if err == nil {
		return false, nil
	}
	return false, err
}

func getExt(file string) string {
	ext := filepath.Ext(file)
	if strings.HasPrefix(ext, ":") {
		ext = ext[1:]
		return ext
	}
	return ext
}

func removeExt(file string) string {
	return file[:len(file)-len(filepath.Ext(file))]
}
