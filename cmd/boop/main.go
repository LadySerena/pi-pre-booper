/*
 * Copyright (c) 2021. Serena Tiede <serena.tiede@gmail.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"log"

	"github.com/LadySerena/pi-pre-booper/pkg/baseline"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
)

const (
	usernameEnv      = "remote_user"
	passwordEnv      = "remote_password"
	addressEnv       = "remote_address"
	publicKeyPathEnv = "ssh_pub_key_path"
)

func main() {

	viper.AutomaticEnv()
	username := viper.GetString(usernameEnv)
	if username == "" {
		log.Println("need to specify username")
		return
	}

	password := viper.GetString(passwordEnv)
	if password == "" {
		log.Println("need to specify password")
		return
	}

	address := viper.GetString(addressEnv)
	if address == "" {
		log.Println("need to specify address")
		return
	}

	pubKeyPath := viper.GetString(publicKeyPathEnv)
	if pubKeyPath == "" {
		log.Println("need to specify pubkey path")
		return
	}

	sshConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // todo do some kind of validation ?
	}

	client, err := ssh.Dial("tcp", address, sshConfig) // do some kind of inventories file
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}

	defer func() {
		closeErr := client.Close()
		if closeErr != nil {
			log.Fatalf("failed to close client: %v", closeErr)
		}
	}()

	session, sessionErr := client.NewSession()
	if sessionErr != nil {
		log.Printf("failed to create session: %v", sessionErr)
		return
	}

	defer func() {
		sessionCloseErr := session.Close()
		log.Fatalf("failed to close session: %v", sessionCloseErr)
	}()

	createUserErr := baseline.CreateUser("serena", pubKeyPath, session)
	if createUserErr != nil {
		log.Printf("could not create user: %v", createUserErr)
		return
	}

}
