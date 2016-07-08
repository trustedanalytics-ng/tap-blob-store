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
package main

import (
	"net/http"

	"github.com/gocraft/web"

	"github.com/trustedanalytics/tapng-image-factory/logger"
	"os"
)

var (
	logger = logger_wrapper.InitLogger("main")
	port   = os.Getenv("IMAGE_FACTORY_PORT")
)

func main() {
	c := Context{}
	c.SetupContext()
	r := web.New(c)
	r.Middleware(web.LoggerMiddleware)
	r.Post("/api/v1/image", c.BuildImage)

	err := http.ListenAndServe("localhost:"+port, r)
	if err != nil {
		logger.Critical("Couldn't serve app on port ", port, " Application will be closed now.")
	}
}
