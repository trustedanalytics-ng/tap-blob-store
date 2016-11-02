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

package util

import (
	"os"
	"fmt"
	"errors"
)

func GetConnectionAddressFromEnvs(componentName string) (string, error) {
	hostEnvName := componentName+"_HOST"
	portEnvName := componentName+"_PORT"

	var errorArray []error

	host, err := GetEnvOrError(hostEnvName)
	if err != nil {
		errorArray = append(errorArray, err)
	}

	port, err := GetEnvOrError(portEnvName)
	if err != nil {
		errorArray = append(errorArray, err)
	}

	address := fmt.Sprintf("%s:%s", host, port)

	if len(errorArray) > 0 {
		return address, sumErrors(errorArray)
	}
	return address, nil
}

func GetConnectionCredentialsFromEnvs(componentName string) (string, string, error) {
	userEnvName := componentName+"_USER"
	passEnvName := componentName+"_PASS"

	var errorArray []error

	username, err := GetEnvOrError(userEnvName)
	if err != nil {
		errorArray = append(errorArray, err)
	}

	password, err := GetEnvOrError(passEnvName)
	if err != nil {
		errorArray = append(errorArray, err)
	}

	if len(errorArray) > 0 {
		return username, password, sumErrors(errorArray)
	}
	return username, password, nil
}


func GetEnvOrError(envName string) (string, error){
	value := os.Getenv(envName)
	if value == "" {
		return value, errors.New(envName+" not set!")
	}
	return value, nil
}

func sumErrors(errorsArray []error) error {
	finalByteMessage := []byte{}
	byteCounter := 0

	for _, err := range errorsArray {
		byteCounter += copy(finalByteMessage[byteCounter:], fmt.Sprintf("%s, ", err.Error()))
	}

	return errors.New(string(finalByteMessage))
}