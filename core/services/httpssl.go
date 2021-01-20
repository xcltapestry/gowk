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
 * @Project gowk
 * @Description go framework
 * @author XiongChuanLiang<br/>(xcl_168@aliyun.com)
 * @license http://www.apache.org/licenses/  Apache v2 License
 * @version 1.0
 */

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
)

type HTTPSsl struct {
	//https : "server.pem", "server.key","client.pem"
	PublicKey, PrivateKey, clientCert string
	TLSConfig                         *tls.Config
}

//ServeCrt  服务端加载证书
func (s *HTTPSsl) ServeCrt(publicKey, privateKey string) error {
	s.PublicKey, s.PrivateKey = publicKey, privateKey

	var err error

	cert := make([]tls.Certificate, 1)
	cert[0], err = tls.LoadX509KeyPair(s.PublicKey, s.PrivateKey)
	if err != nil {
		return fmt.Errorf("LoadX509KeyPair err:%s", err)
	}

	s.TLSConfig = &tls.Config{
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
		},
		PreferServerCipherSuites: true,
		MinVersion:               tls.VersionTLS12,
		NextProtos:               []string{"http/1.1"},
		// CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		Certificates: cert,
	}

	return nil
}

//ClientCrt 客户端加载证书并做双向认证
func (s *HTTPSsl) ClientCrt(publicKey, privateKey, clientCert string) error {
	s.PublicKey, s.PrivateKey = publicKey, privateKey
	s.clientCert = clientCert

	var err error

	cert := make([]tls.Certificate, 1)
	cert[0], err = tls.LoadX509KeyPair(s.PublicKey, s.PrivateKey)
	if err != nil {
		return fmt.Errorf("LoadX509KeyPair err:%s", err)
	}

	certBytes, err := ioutil.ReadFile(s.clientCert)
	if err != nil {
		return errors.New("(" + s.clientCert + ") err:" + err.Error())
	}
	clientCertPool := x509.NewCertPool()
	ok := clientCertPool.AppendCertsFromPEM(certBytes)
	if !ok {
		return errors.New("root certificate.")
	}

	s.TLSConfig = &tls.Config{
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		},
		PreferServerCipherSuites: true,
		ClientAuth:               tls.RequireAndVerifyClientCert,
		ClientCAs:                clientCertPool,
		MinVersion:               tls.VersionTLS12,
		NextProtos:               []string{"http/1.1"},
		Certificates:             cert,
	}

	return nil
}

//ClientTLS 用于单向认证
func (s *HTTPSsl) ClientAuth() {
	s.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
}
