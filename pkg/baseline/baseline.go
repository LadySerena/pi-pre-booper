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

package baseline

// todo extra baselines to add
// set motd
// disable unattended updates
// disable general ubuntu things (flatpack, snap, cloud-init)
// install my stuff (not sure yet ?????)

import (
	"fmt"
	"io/ioutil"
	"strings"

	"golang.org/x/crypto/ssh"
)

func CreateUser(name string, pubKeyPath string, sshSession *ssh.Session) error {
	sshDirPath := fmt.Sprintf("/home/%s/.ssh", name)
	authKeyPath := fmt.Sprintf("%s/authorized_keys", sshDirPath)
	key, keyReadErr := AddAuthorizedKey(pubKeyPath)
	if keyReadErr != nil {
		return keyReadErr
	}
	addUser := fmt.Sprintf("sudo adduser --disabled-password %s", name)
	createSSHDir := fmt.Sprintf("sudo mkdir %s", sshDirPath)
	setSSHDirOwners := fmt.Sprintf("sudo chown -R %s:%s %s", name, name, sshDirPath)
	setSSHDirPermissions := fmt.Sprintf("sudo chmod -R 0700 %s", sshDirPath)
	createAuthorizedKeyFile := fmt.Sprintf("sudo touch %s", authKeyPath)
	setAuthorizedKeyFileOwners := fmt.Sprintf("sudo chown -R %s:%s %s", name, name, authKeyPath)
	setAuthorizedKeyPermissions := fmt.Sprintf("sudo chmod 0600 %s", authKeyPath)
	addKey := fmt.Sprintf("echo \"%s\" | sudo tee %s", key, authKeyPath)

	// todo add to sudoers file as well and then we can really rock

	commands := []string{
		addUser,
		createSSHDir,
		setSSHDirOwners,
		setSSHDirPermissions,
		createAuthorizedKeyFile,
		setAuthorizedKeyFileOwners,
		setAuthorizedKeyPermissions,
		addKey,
	}
	commandString := strings.Join(commands, "; ")
	_, runErr := sshSession.CombinedOutput(commandString)
	if runErr != nil {
		return runErr
	}
	return nil
}

func AddAuthorizedKey(path string) (string, error) {
	keyData, readErr := ioutil.ReadFile(path)
	if readErr != nil {
		return "", readErr
	}
	return string(keyData), nil
}
