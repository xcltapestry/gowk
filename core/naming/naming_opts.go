package naming

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
	"time"
)

type NamingOption func(*Naming)

func NewNaming(options ...func(*Naming)) *Naming {

	cli := &Naming{}
	cli.config.Endpoints = []string{_DefaultEndpoints}
	cli.config.DialTimeout = 2 * time.Second
	cli.prefix = _DefaultEtcdPrefix

	for _, f := range options {
		f(cli)
	}

	return cli
}

func WithAddress(addrs []string) NamingOption {
	return func(c *Naming) {
		c.config.Endpoints = addrs
	}
}

func WithDialTimeout(dialTimeout time.Duration) NamingOption {
	return func(c *Naming) {
		c.config.DialTimeout = dialTimeout
	}
}
