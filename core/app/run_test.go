package app

/**
 * Copyright 2021  gowrk Author. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/xcltapestry/gowk/core/services"
)

func TestRun(t *testing.T) {

	glog.Info("golog ------------- ")

	fmt.Println(" ok ")
	New()
	httpsvc2 := services.NewHTTPService()
	httpsvc2.Router(RegisterHandlers)
	httpsvc2.SetHTTPAddr(":8003")
	App.Serve(httpsvc2)
	fmt.Println("Serve --> :8003")
	Run()

	defer Flush()

}

//RegisterHandlers 路由
func RegisterHandlers(m *mux.Router) {
	m.HandleFunc("/v1/ping", ArticlesCategoryHandler)
	m.HandleFunc("/v1/health", HealthHandler)
	m.Handle("/", http.NotFoundHandler())
}

func ArticlesCategoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category: %v\n", vars["category"])
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Health: %v\n", time.Now().Unix())
}
