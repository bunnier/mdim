package qiniu

import (
	"context"
	"errors"
	"path/filepath"
	"time"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

// Qiniu SDK doc: https://developer.qiniu.com/kodo/1238/go

// QiniuUploadApi is a wrapper upload api for Qiniu SDK
type QiniuUploadApi struct {
	ak              string
	sk              string
	bucket          string
	defaultFolder   string
	domain          string
	token           string
	tokenExpireTime time.Time

	mac           *qbox.Mac
	bucketManager *storage.BucketManager
}

type QiniuUploadApiOption func(api *QiniuUploadApi)

func QiniuUploadApiDomainOption(domain string) QiniuUploadApiOption {
	return func(api *QiniuUploadApi) {
		api.domain = domain
	}
}
func QiniuUploadApiDefaultFolderOption(defaultFolder string) QiniuUploadApiOption {
	return func(api *QiniuUploadApi) {
		api.defaultFolder = defaultFolder
	}
}

func NewQuniuUploadApi(ak, sk, bucket string, options ...QiniuUploadApiOption) *QiniuUploadApi {
	api := &QiniuUploadApi{
		ak:     ak,
		sk:     sk,
		bucket: bucket,
		mac:    qbox.NewMac(ak, sk),
	}

	cfg := storage.Config{
		UseHTTPS: true,
	}
	api.bucketManager = storage.NewBucketManager(api.mac, &cfg)

	for _, option := range options {
		option(api)
	}

	if api.domain == "" {
		if err := api.fetchDomain(); err != nil {
			panic(err)
		}
	}

	api.refreshToken()

	return api
}

// refreshToken
func (api *QiniuUploadApi) refreshToken() {
	const timeoutOffset time.Duration = 120 * time.Second
	if api.tokenExpireTime.Add(-timeoutOffset).After(time.Now()) {
		return
	}

	// Get token by sdk.
	putPolicy := storage.PutPolicy{
		Scope: api.bucket,
	}
	api.token = putPolicy.UploadToken(api.mac)
}

// fetchDomains of the specific bucket.
func (api *QiniuUploadApi) fetchDomain() error {
	var domainInfos []storage.DomainInfo
	var err error
	if domainInfos, err = api.bucketManager.ListBucketDomains(api.bucket); err != nil {
		return err
	}
	for _, di := range domainInfos {
		if di.Domain == "" {
			continue
		}
		api.domain = di.Domain
		return nil
	}
	return errors.New("no available domain")
}

// Upload file to Quniu cloud, and then return the public url.
func (api *QiniuUploadApi) Upload(remoteFilepath, localFilepath string) (string, error) {
	return api.UploadContext(context.Background(), remoteFilepath, localFilepath)
}

// Upload file to Quniu cloud, and then return the public url.
func (api *QiniuUploadApi) UploadContext(ctx context.Context, remoteFilepath, localFilepath string) (string, error) {
	if remoteFilepath == "" {
		remoteFilepath = filepath.Join(api.defaultFolder, filepath.Base(localFilepath))
	}
	api.refreshToken()
	cfg := storage.Config{}
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	if err := formUploader.PutFile(ctx, &ret, api.token, remoteFilepath, localFilepath, nil); err != nil {
		return "", err
	}
	return storage.MakePublicURL(api.domain, remoteFilepath), nil
}
