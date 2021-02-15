package services

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
	"net"
	"net/http"
	"time"
	"strings"
	"context"

	muxCtx "github.com/gorilla/context"
	"github.com/gorilla/mux"
)

type RouterFunc func(m *mux.Router)

type BaseService struct {
	Name string
}

//
type HTTPService struct {
	BaseService

	Addr  string
	svc   *http.Server
	route *mux.Router

	ConnectTimeout time.Duration
	WriteTimeout   time.Duration
	ReadTimeout    time.Duration
	MaxHeaderBytes int

	HTTPSsl
}

func NewHTTPService() *HTTPService {
	svc := &HTTPService{}
	svc.Initialize()
	return svc
}

func (s *HTTPService) Initialize() error {
	s.ConnectTimeout, s.WriteTimeout, s.ReadTimeout = 60*time.Second, 60*time.Second, 60*time.Second
	s.MaxHeaderBytes = 1 << 20
	s.Addr = ":8000"
	return nil
}

func (s *HTTPService) SetHTTPAddr(addr string) {
	s.Addr = addr
}

func (s *HTTPService) Listen(addrs... string) {
	var addr string
	for _, a := range addrs {
		addr = a
	}

	if strings.TrimSpace(addr) != "" {
		s.Addr = addr
	}
}


func (s *HTTPService) SetHTTPTimeout(
	connectTimeout, writeTimeout, readTimeout time.Duration) {
	s.ConnectTimeout, s.WriteTimeout, s.ReadTimeout = connectTimeout, writeTimeout, readTimeout
}

func (s *HTTPService) SetMaxHeaderBytes(maxHeaderBytes int) {
	if maxHeaderBytes > 0 {
		s.MaxHeaderBytes = maxHeaderBytes
	}
}

// Router  :	websvc.Router(svc.Router(svc.RegisterHandlers))
func (s *HTTPService) Router(f RouterFunc) *HTTPService {
	f(s.mux())
	return s
}

// mux 使用mux来作为路由
func (s *HTTPService) mux() *mux.Router {
	if s.route != nil {
		return s.route
	}
	s.route = mux.NewRouter()
	s.route.StrictSlash(true)
	return s.route
}

func (s *HTTPService) init() error {
	if s.route == nil {
		return fmt.Errorf("%s", "route is null!")
	}

	http.Handle("/", s.route)
	s.svc = &http.Server{
		TLSConfig:      s.TLSConfig,
		Addr:           s.Addr,
		Handler:        muxCtx.ClearHandler(http.DefaultServeMux),
		ReadTimeout:    s.ReadTimeout,
		WriteTimeout:   s.WriteTimeout,
		MaxHeaderBytes: s.MaxHeaderBytes,
	}
	return nil
}

func (s *HTTPService) Run() error {
	if s.svc == nil {
		s.init()
	}

	//listen
	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return fmt.Errorf("err:%s", err)
	}

	//https
	if s.TLSConfig != nil {
		if err = s.svc.Serve(ln); err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("HTTPS Addr:%s err:%s", s.Addr, err)
		}
		return nil
	}

	//http
	if err = s.svc.Serve(ln); err != nil {
		return fmt.Errorf("HTTP Addr:%s err:%s", s.Addr, err)
	}

	return nil
}

func (s *HTTPService) Stop(ctx context.Context) {
	fmt.Println(" Serve stop. addr:", s.Addr)
}
