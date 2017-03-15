/*
 * controller.go - goControlTor
 * Copyright (C) 2014  Yawning Angel, David Stainton
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 * Modified by twitchyliquid64 for inclusion in subshard.
 */

package main

import (
	"bufio"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net"
	"net/textproto"
	"path"
	"strings"
)

const (
	cmdOK            = 250
	cmdAuthenticate  = "AUTHENTICATE"
	cmdAuthChallenge = "AUTHCHALLENGE"
	//	authMethodCookie     = "COOKIE"
	//	authMethodNull       = "NULL"

	respAuthChallenge = "AUTHCHALLENGE "

	argServerHash  = "SERVERHASH="
	argServerNonce = "SERVERNONCE="

	authMethodSafeCookie = "SAFECOOKIE"
	authNonceLength      = 32

	authServerHashKey = "Tor safe cookie authentication server-to-controller hash"
	authClientHashKey = "Tor safe cookie authentication controller-to-server hash"
)

type TorControl struct {
	controlConn            net.Conn
	textprotoReader        *textproto.Reader
	authenticationPassword string
}

func (t *TorControl) Close() error {
	return t.controlConn.Close()
}

// Dial handles unix domain sockets and tcp!
func (t *TorControl) Dial(network, addr string) error {
	var err error = nil
	t.controlConn, err = net.Dial(network, addr)
	if err != nil {
		return err
	}
	reader := bufio.NewReader(t.controlConn)
	t.textprotoReader = textproto.NewReader(reader)
	return nil
}

func (t *TorControl) TorVersion() (string, error) {
	code, out, err := t.SendMultiCommand("GETINFO version\n")
	if err != nil {
		return "", err
	}
	if code != 250 {
		return "", fmt.Errorf("Getting version failed: %d", code)
	}
	firstline := strings.Split(out, "\n")[0]
	return strings.TrimPrefix(firstline, "version="), nil
}

func (t *TorControl) SocksListenersAreLocal() (bool, error) {
	code, out, err := t.SendMultiCommand("GETINFO net/listeners/socks\n")
	if err != nil {
		return false, err
	}
	if code != 250 {
		return false, fmt.Errorf("Getting socks listeners failed: %d", code)
	}
	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, ":") {
			if !(strings.Contains(line, "localhost") || strings.Contains(line, "127.0.0.1")) {
				return false, nil
			}
		}
	}
	return true, nil
}

func (t *TorControl) ControlListenersAreLocal() (bool, error) {
	code, out, err := t.SendMultiCommand("GETINFO net/listeners/control\n")
	if err != nil {
		return false, err
	}
	if code != 250 {
		return false, fmt.Errorf("Getting control listeners failed: %d", code)
	}
	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, ":") {
			if !(strings.Contains(line, "localhost") || strings.Contains(line, "127.0.0.1")) {
				return false, nil
			}
		}
	}
	return true, nil
}

func (t *TorControl) IsDormant() (bool, error) {
	code, out, err := t.SendMultiCommand("GETINFO dormant\n")
	if err != nil {
		return false, err
	}
	if code != 250 {
		return false, fmt.Errorf("Getting dormancy failed: %d", code)
	}
	firstline := strings.Split(out, "\n")[0]
	return firstline == "dormant=1", nil
}

func (t *TorControl) CircuitsEstablished() (bool, error) {
	code, out, err := t.SendMultiCommand("GETINFO status/circuit-established\n")
	if err != nil {
		return false, err
	}
	if code != 250 {
		return false, fmt.Errorf("Getting circuit status failed: %d", code)
	}
	firstline := strings.Split(out, "\n")[0]
	return firstline == "status/circuit-established=1", nil
}

func (t *TorControl) EnoughDirInfo() (bool, error) {
	code, out, err := t.SendMultiCommand("GETINFO status/enough-dir-info\n")
	if err != nil {
		return false, err
	}
	if code != 250 {
		return false, fmt.Errorf("Getting dir status failed: %d", code)
	}
	firstline := strings.Split(out, "\n")[0]
	return firstline == "status/enough-dir-info=1", nil
}

func (t *TorControl) SendCommand(command string) (int, string, error) {
	var code int
	var message string
	var err error

	_, err = t.controlConn.Write([]byte(command))
	if err != nil {
		return 0, "", fmt.Errorf("writing to tor control port: %s", err)
	}
	code, message, err = t.textprotoReader.ReadCodeLine(cmdOK)
	if err != nil {
		return code, message, fmt.Errorf("reading tor control port command status: %s", err)
	}
	return code, message, nil
}

func (t *TorControl) SendMultiCommand(command string) (int, string, error) {
	var code int
	var message string
	var err error

	_, err = t.controlConn.Write([]byte(command))
	if err != nil {
		return 0, "", fmt.Errorf("writing to tor control port: %s", err)
	}
	code, message, err = t.textprotoReader.ReadResponse(cmdOK)
	if err != nil {
		return code, message, fmt.Errorf("reading tor control port command status: %s", err)
	}
	return code, message, nil
}

func (t *TorControl) SafeCookieAuthenticate(cookiePath string) error {

	var code int
	var message string

	cookie, err := readAuthCookie(cookiePath)
	if err != nil {
		return err
	}

	cookie, err = t.authSafeCookie(cookie)
	if err != nil {
		return err
	}
	cookieStr := hex.EncodeToString(cookie)
	authReq := fmt.Sprintf("%s %s\n", cmdAuthenticate, cookieStr)

	code, message, err = t.SendCommand(authReq)
	if err != nil {
		return fmt.Errorf("Safe Cookie Authentication fail: %s %s %s", code, message, err)
	}

	return nil
}

func (t *TorControl) CookieAuthenticate(cookiePath string) error {

	var code int
	var message string

	cookie, err := readAuthCookie(cookiePath)
	if err != nil {
		return err
	}

	cookieStr := hex.EncodeToString(cookie)
	authReq := fmt.Sprintf("%s %s\n", cmdAuthenticate, cookieStr)

	code, message, err = t.SendCommand(authReq)
	if err != nil {
		return fmt.Errorf("Cookie Authentication fail: %s %s %s", code, message, err)
	}

	return nil
}

func readAuthCookie(path string) ([]byte, error) {
	// Read the cookie auth file.
	cookie, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading cookie auth file: %s", err)
	}
	return cookie, nil
}

func (t *TorControl) authSafeCookie(cookie []byte) ([]byte, error) {
	var code int
	var message string
	var err error

	clientNonce := make([]byte, authNonceLength)
	if _, err := rand.Read(clientNonce); err != nil {
		return nil, fmt.Errorf("generating AUTHCHALLENGE nonce: %s", err)
	}
	clientNonceStr := hex.EncodeToString(clientNonce)

	// Send and process the AUTHCHALLENGE.
	authChallengeReq := []byte(fmt.Sprintf("%s %s %s\n", cmdAuthChallenge, authMethodSafeCookie, clientNonceStr))
	if _, err := t.controlConn.Write(authChallengeReq); err != nil {
		return nil, fmt.Errorf("writing AUTHCHALLENGE request: %s", err)
	}

	code, message, err = t.textprotoReader.ReadCodeLine(cmdOK)
	if err != nil {
		return nil, fmt.Errorf("reading tor control port command status: %s %s %s", code, message, err)
	}

	lineStr := strings.TrimSpace(message)
	respStr := strings.TrimPrefix(lineStr, respAuthChallenge)
	if respStr == lineStr {
		return nil, fmt.Errorf("parsing AUTHCHALLENGE response")
	}
	splitResp := strings.SplitN(respStr, " ", 2)
	if len(splitResp) != 2 {
		return nil, fmt.Errorf("parsing AUTHCHALLENGE response")
	}
	hashStr := strings.TrimPrefix(splitResp[0], argServerHash)
	serverHash, err := hex.DecodeString(hashStr)
	if err != nil {
		return nil, fmt.Errorf("decoding AUTHCHALLENGE ServerHash: %s", err)
	}
	serverNonceStr := strings.TrimPrefix(splitResp[1], argServerNonce)
	serverNonce, err := hex.DecodeString(serverNonceStr)
	if err != nil {
		return nil, fmt.Errorf("decoding AUTHCHALLENGE ServerNonce: %s", err)
	}

	// Validate the ServerHash.
	m := hmac.New(sha256.New, []byte(authServerHashKey))
	m.Write([]byte(cookie))
	m.Write([]byte(clientNonce))
	m.Write([]byte(serverNonce))
	dervServerHash := m.Sum(nil)
	if !hmac.Equal(serverHash, dervServerHash) {
		return nil, fmt.Errorf("AUTHCHALLENGE ServerHash is invalid")
	}

	// Calculate the ClientHash.
	m = hmac.New(sha256.New, []byte(authClientHashKey))
	m.Write([]byte(cookie))
	m.Write([]byte(clientNonce))
	m.Write([]byte(serverNonce))

	return m.Sum(nil), nil
}

func (t *TorControl) PasswordAuthenticate(password string) error {
	authCmd := fmt.Sprintf("%s \"%s\"\n", cmdAuthenticate, password)
	_, _, err := t.SendCommand(authCmd)
	return err
}

// Creates a Tor Hidden Service with the HiddenServiceDirGroupReadable option
// set so that the service's hostname file will have group read permission set.
// At this time of writing this feature is only available in the alpha version
// of tor. See https://trac.torproject.org/projects/tor/ticket/11291
func (t *TorControl) CreateHiddenService(serviceDir string, listenAddrs map[int]string) error {
	var createCmd string = fmt.Sprintf("SETCONF hiddenservicedir=%s", serviceDir)
	for virtPort, listenAddr := range listenAddrs {
		createCmd += fmt.Sprintf(" hiddenserviceport=\"%d %s\"", virtPort, listenAddr)
	}
	createCmd += " HiddenServiceDirGroupReadable=1\n"
	_, _, err := t.SendCommand(createCmd)
	return err
}

func (t *TorControl) DeleteHiddenService(serviceDir string) error {
	var deleteCmd string = fmt.Sprintf("SETCONF hiddenservicedir=%s\n", serviceDir)
	_, _, err := t.SendCommand(deleteCmd)
	return err
}

func ReadOnion(serviceDir string) (string, error) {
	onion, err := ioutil.ReadFile(path.Join(serviceDir, "hostname"))
	if err != nil {
		return "", fmt.Errorf("reading Tor hidden service hostname file: %s", err)
	}
	return string(onion), nil
}
