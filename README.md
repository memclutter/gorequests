# gorequests

[![Go](https://github.com/memclutter/gorequests/actions/workflows/go.yml/badge.svg)](https://github.com/memclutter/gorequests/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/memclutter/gorequests/branch/main/graph/badge.svg?token=1IWTNCLCAQ)](https://codecov.io/gh/memclutter/gorequests)

http requests wrapper for go.

## Motivation

The motivation for this project was my feeling of disgust at the "convenience" of working with the golang http client. 
Let's say I want to request an ip address from the ipify service (this is not advertising!). 
Well, how do we do it on go? Dude, it’s easier than a steamed turnip, relax...

```go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
)

func GetIP() (net.IP, error) {
	res, err := http.Get("https://api.ipify.org?format=json")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(res.Status)
	}
	defer res.Body.Close()
	resData := make(map[string]string)
	if err := json.Unmarshal(body, &resData); err != nil {
		return nil, err
	}
	ipStr, ok := resData["ip"]
	if !ok {
		return nil, fmt.Errorf("not found ip")
	}
	ip := net.ParseIP(ipStr)
	if len(ip) == 0 {
		return nil, fmt.Errorf("invalid ip address '%s'", ipStr)
	}
	return ip, nil
}

// ...
```

cool right? I looked for alternatives, but did not find what would suit me. That's why I decided to make this project.

## Idea

If we consider the previous option, the idea is as follows

```go
package main

import (
	"net"
	"net/http"
	"github.com/memclutter/gorequests"
)

func GetIPEasy() (ip net.IP, err error) {
    err = gorequests.Request().
        Method(http.MethodGet).
        Url("https://api.ipify.org?format=json").
		ResponseCodeOk(http.StatusOK).
        ResponseJson(&ip, ".ip").
		Exec()
    return
}
// ...
```

or more short

```go
package main

import (
	"net"
	"net/http"
	"github.com/memclutter/gorequests"
)

func GetIPEasy() (ip net.IP, err error) {
	err = gorequests.Get("https://api.ipify.org?format=json").
		ResponseCodeOk(http.StatusOK).
		ResponseJson(&ip, ".ip").
		Exec()
	return
}
// ...
```

Wow! Now I can focus on the business logic of the application, and not the details of decoding the server response.