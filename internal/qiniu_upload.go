package internal

import (
	"context"
	"time"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

// Qiniu SDK doc: https://developer.qiniu.com/kodo/1238/go
type QiniuUploadApi struct {
	ak              string
	sk              string
	bucket          string
	domains         []string
	token           string
	tokenExpireTime time.Time
}

func NewQuniuUploadApi(ak string, sk string, bucket string) *QiniuUploadApi {
	api := &QiniuUploadApi{
		ak:     ak,
		sk:     sk,
		bucket: bucket,
	}
	api.refreshToken()
	return api
}

func (api *QiniuUploadApi) refreshToken() {
	if api.tokenExpireTime.Add(-120 * time.Second).After(time.Now()) {
		return
	}
	// Get token by sdk.
	putPolicy := storage.PutPolicy{
		Scope: api.bucket,
	}
	mac := qbox.NewMac(api.ak, api.sk)
	api.token = putPolicy.UploadToken(mac)
}

// Upload file
func (api *QiniuUploadApi) Upload(remoteFilepath, localFilepath string) error {
	api.refreshToken()
	cfg := storage.Config{}
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	return formUploader.PutFile(context.Background(), &ret, api.token, remoteFilepath, localFilepath, nil)
}

func (api *QiniuUploadApi) GetDomains() ([]string, error) {
	mac := qbox.NewMac(api.ak, api.sk)
	cfg := storage.Config{
		UseHTTPS: true,
	}

	bucketManager := storage.NewBucketManager(mac, &cfg)
	var domainInfos []storage.DomainInfo
	var err error
	if domainInfos, err = bucketManager.ListBucketDomains(api.bucket); err != nil {
		return nil, err
	}

	domains := make([]string, len(domainInfos))
	for _, di := range domainInfos {
		if di.Domain != "" {
			continue
		}
		domains = append(domains, di.Domain)
	}
	api.domains = domains
	return api.domains, nil
}
