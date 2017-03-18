package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/elazarl/goproxy"
)

//Code mostly from github.com/elazar/goproxy - all credit his.

var unauthorizedMsg = []byte("407 Proxy Authentication Required")

func basicUnauthorizedResponse(req *http.Request, realm string) *http.Response {
	return &http.Response{
		StatusCode:    407,
		ProtoMajor:    1,
		ProtoMinor:    1,
		Request:       req,
		Header:        http.Header{"Proxy-Authenticate": []string{"Basic realm=" + realm}},
		Body:          ioutil.NopCloser(bytes.NewBuffer(unauthorizedMsg)),
		ContentLength: int64(len(unauthorizedMsg)),
	}
}

var proxyAuthorizationHeader = "Proxy-Authorization"

//Checks for password
func checkRequestAuthentication(req *http.Request, ctx *goproxy.ProxyCtx, f func(ctx *goproxy.ProxyCtx, user, passwd string) (bool, map[string]interface{})) bool {
	authheader := strings.SplitN(req.Header.Get(proxyAuthorizationHeader), " ", 2)
	req.Header.Del(proxyAuthorizationHeader)
	if len(authheader) != 2 || authheader[0] != "Basic" {
		return false
	}
	userpassraw, err := base64.StdEncoding.DecodeString(authheader[1])
	if err != nil {
		return false
	}
	userpass := strings.SplitN(string(userpassraw), ":", 2)
	if len(userpass) != 2 {
		return false
	}
	ok, mapData := f(ctx, userpass[0], userpass[1])
	if ok {
		ctx.UserData = map[string]interface{}{
			"auth":     true,
			"user":     userpass[0],
			"authkind": "password",
			"info":     mapData,
		}
	}
	return ok
}

// Basic returns a basic HTTP authentication handler for requests
//
// You probably want to use auth.ProxyBasic(proxy) to enable authentication for all proxy activities
func Basic(realm string, f func(ctx *goproxy.ProxyCtx, user, passwd string) (bool, map[string]interface{})) goproxy.ReqHandler {
	return goproxy.FuncReqHandler(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		if !checkRequestAuthentication(req, ctx, f) {
			return nil, basicUnauthorizedResponse(req, realm)
		}
		return req, nil
	})
}

// BasicConnect returns a basic HTTP authentication handler for CONNECT requests
//
// You probably want to use auth.ProxyBasic(proxy) to enable authentication for all proxy activities
func BasicConnect(realm string, f func(ctx *goproxy.ProxyCtx, user, passwd string) (bool, map[string]interface{})) goproxy.HttpsHandler {
	return goproxy.FuncHttpsHandler(func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
		if !checkRequestAuthentication(ctx.Req, ctx, f) {
			ctx.Resp = basicUnauthorizedResponse(ctx.Req, realm)
			return goproxy.RejectConnect, host
		}
		return goproxy.OkConnect, host
	})
}

// SetupProxyAuthentication will force HTTP authentication or valid TLS client certs before any request to the proxy is processed.
func SetupProxyAuthentication(proxy *goproxy.ProxyHttpServer, realm string, users []User) {
	f := func(ctx *goproxy.ProxyCtx, user, passwd string) (bool, map[string]interface{}) {
		sha1Hash := fmt.Sprintf("%x", sha256.Sum256([]byte(passwd)))
		for _, usr := range users {
			if usr.Username == user && usr.Password == sha1Hash {
				return true, map[string]interface{}{}
			}
		}
		return false, nil
	}

	proxy.OnRequest().Do(Basic(realm, f))
	proxy.OnRequest().HandleConnect(BasicConnect(realm, f))
}
