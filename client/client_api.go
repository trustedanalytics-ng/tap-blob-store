package client

import (
	"fmt"
	"net/http"

	"bytes"
	brokerHttp "github.com/trustedanalytics/tapng-go-common/http"
	"github.com/trustedanalytics/tapng-go-common/logger"
	"io"
	"mime/multipart"
)

var (
	logger = logger_wrapper.InitLogger("client")
)

type TapBlobStoreApi interface {
	StoreBlob(blob_id string, file multipart.File) error
	GetBlob(blob_id string, dest io.Writer) error
	DeleteBlob(blob_id string) (int, error)
}

func NewTapBlobStoreApiWithBasicAuth(address, username, password string) (*TapBlobStoreApiConnector, error) {
	client, _, err := brokerHttp.GetHttpClientWithBasicAuth()
	if err != nil {
		return nil, err
	}
	return &TapBlobStoreApiConnector{address, username, password, client}, nil
}

func NewTapBlobStoreApiWithSSLAndBasicAuth(address, username, password, certPemFile, keyPemFile, caPemFile string) (*TapBlobStoreApiConnector, error) {
	client, _, err := brokerHttp.GetHttpClientWithCertAndCaFromFile(certPemFile, keyPemFile, caPemFile)
	if err != nil {
		return nil, err
	}
	return &TapBlobStoreApiConnector{address, username, password, client}, nil
}

type TapBlobStoreApiConnector struct {
	Address  string
	Username string
	Password string
	Client   *http.Client
}

func (c *TapBlobStoreApiConnector) getApiConnector(url string) brokerHttp.ApiConnector {
	return brokerHttp.ApiConnector{
		BasicAuth: &brokerHttp.BasicAuth{c.Username, c.Password},
		Client:    c.Client,
		Url:       url,
	}
}

func (c *TapBlobStoreApiConnector) StoreBlob(blob_id string, file multipart.File) error {
	connector := c.getApiConnector(fmt.Sprintf("%s/api/v1/blobs", c.Address))

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	err := bodyWriter.WriteField("blob_id", blob_id)
	if err != nil {
		return err
	}

	fileWriter, err := bodyWriter.CreateFormFile("uploadfile", "blob.tar.gz")
	if err != nil {

		fmt.Println("error writing to buffer")
		return err
	}

	size, err := io.Copy(fileWriter, file)
	if err != nil {
		return err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	logger.Infof("Doing: POST %v Sending %v bytes", connector.Url, size)
	_, err = http.Post(connector.Url, contentType, bodyBuf)
	return err
}

func (c *TapBlobStoreApiConnector) GetBlob(blob_id string, dest io.Writer) error {
	connector := c.getApiConnector(fmt.Sprintf("%s/api/v1/blobs/%s", c.Address, blob_id))
	size, err := brokerHttp.DownloadBinary(connector.Url, connector.BasicAuth, connector.Client, dest)
	if err != nil {
		return err
	}
	logger.Infof("Written %v bytes of binary data to destination", size)
	return err
}

func (c *TapBlobStoreApiConnector) DeleteBlob(blob_id string) (int, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/api/v1/blobs/%s", c.Address, blob_id))
	status, _, err := brokerHttp.RestDELETE(connector.Url, "", connector.BasicAuth, connector.Client)
	if err != nil {
		return status, err
	}

	return status, err
}