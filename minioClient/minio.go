/**
 * Copyright (c) 2016 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package minioClient

import (
	"io"
	"github.com/minio/minio-go"
	"github.com/trustedanalytics/tap-go-common/logger"
	"errors"
)

var (
	bucketName string
	minioClient *minio.Client
	logger = logger_wrapper.InitLogger("minio")
	err error
)


const (
	//TODO: those information should be stored in kubernetes secrets ?
        endpoint = "127.0.0.1:9000"
        accessKeyID = "<accessKeyID>"
        secretAccessKey = "<secretAccessKey>"
        ssl = false
)

const (
	ErrBucketNotExist = "The specified bucket does not exist."
	ErrKeyNotExist = "The specified key does not exist."
	ErrKeyAlreadyInUse = "The specified key already exists."
)

func CreateMinioClient(_bucketName string) error {
	bucketName = _bucketName

	logger.Info("Initializing Minio Client in bucket -", bucketName)

        minioClient, err = minio.New(endpoint, accessKeyID, secretAccessKey, ssl)
        if err != nil {
		return err
        }
	logger.Info("Successfully started Minio Client. Connected to Minio Server -", endpoint)

        err = minioClient.BucketExists(bucketName)
        if err != nil {
		switch err.Error() {
		case ErrBucketNotExist:
			err = minioClient.MakeBucket(bucketName, "")
			if err != nil {
				return err
			}
			logger.Info("Successfully created bucket -", bucketName)

		default:
			return err
		}
        }

	logger.Info("Minio Client is ready in bucket -", bucketName)
	return nil
}


func StoreInMinio(objectName string, object io.Reader) error {
	logger.Info("Trying to store object -", objectName, "- in bucket -", bucketName)

	_, err := minioClient.StatObject(bucketName, objectName)
	if(err == nil) {
		return errors.New(ErrKeyAlreadyInUse)
	}
	if err.Error() != ErrKeyNotExist {
		return err
	}

        _, err = minioClient.PutObject(bucketName, objectName, object, "application/octet-stream")
	if err != nil {
                return err
        }

	logger.Info("Object successfully stored -", objectName)
	return nil
}

func RemoveFromMinio(objectName string) error {
	logger.Info("Trying to remove object -", objectName, "- in bucket -", bucketName)

	_, err := minioClient.StatObject(bucketName, objectName)
	if err != nil {
		return err
	}

        err = minioClient.RemoveObject(bucketName, objectName)
        if err != nil {
                return err
        }

	logger.Info("Object successfully removed -" + objectName)
	return nil
}

func RetrieveFromMinio(objectName string) (*minio.Object, error) {
	logger.Info("Trying to retrieve object -", objectName, "- from bucket -", bucketName)

        blob, err := minioClient.GetObject(bucketName, objectName)
        if err != nil {
                return nil, err
        }

	logger.Info("Object successfully retrieved -" + objectName)
        return blob, nil
}
