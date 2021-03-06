// +build linux

/*
 * Copyright (C) 2019 IBM, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy ofthe License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specificlanguage governing permissions and
 * limitations under the License.
 *
 */

package runc

import (
	"bufio"
	"net"
	"os"
	"strings"
)

type hosts struct {
	IP       string              `json:"IP,omitempty"`
	Hostname string              `json:"Hostname,omitempty"`
	ByIP     map[string][]string `json:"ByIP,omitempty"`
}

func newHosts() *hosts {
	hosts := new(hosts)
	hosts.ByIP = make(map[string][]string)
	return hosts
}

func readHosts(path string) (*hosts, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	hosts := newHosts()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if err := scanner.Err(); err != nil {
			return nil, err
		}

		if i := strings.IndexByte(line, '#'); i >= 0 {
			// Discard comment.
			line = line[0:i]
		}

		f := strings.Fields(line)
		if len(f) < 2 {
			continue
		}

		ip := net.ParseIP(f[0])
		if ip == nil {
			continue
		}

		if ip.IsLoopback() {
			continue
		}

		hosts.IP = ip.String()
		hosts.ByIP[hosts.IP] = make([]string, 0)
		for i := 1; i < len(f); i++ {
			hosts.Hostname = strings.ToLower(f[i])
			hosts.ByIP[hosts.IP] = append(hosts.ByIP[hosts.IP], hosts.Hostname)
		}
	}

	return hosts, nil
}
