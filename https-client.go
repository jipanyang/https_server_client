// test run:
// /usr/local/go/bin/go run https-client.go
package main

import (
    "crypto/tls"
    "crypto/x509"
    "flag"
    "fmt"
    "golang.org/x/net/http2"
    "io/ioutil"
    "net/http"

    log "github.com/golang/glog"
)

var (
    SERVER_CERT = "/tmp/cert/server.crt"
    CLIENT_CERT = "/tmp/cert/client.crt"
    CLIENT_KEY  = "/tmp/cert/client.key"
)

func printHtml(resp *http.Response) error {
    htmlData, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println(err)
        return err
    }
    defer resp.Body.Close()
    fmt.Printf("%v\n", resp.Status)
    fmt.Printf(string(htmlData))
    return nil
}

func main() {

    flag.Parse()
    caCert, err := ioutil.ReadFile(SERVER_CERT)
    if err != nil {
        log.Fatal(err)
    }
    caCertPool := x509.NewCertPool()
    caCertPool.AppendCertsFromPEM(caCert)

    cert, err := tls.LoadX509KeyPair(CLIENT_CERT, CLIENT_KEY)
    if err != nil {
        log.Fatal(err)
    }

    // http2
    client2 := &http.Client{
        Transport: &http2.Transport{
            TLSClientConfig: &tls.Config{
                RootCAs:      caCertPool,
                Certificates: []tls.Certificate{cert},
                // InsecureSkipVerify: true,
            },
        },
    }

    resp, err := client2.Get("https://localhost:443/hello")
    if err != nil {
        log.Errorf("%v", err)
        return
    }
    fmt.Printf("http/2 response:\n")
    printHtml(resp)

    // http 1.1
    client1_1 := &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{
                RootCAs:      caCertPool,
                Certificates: []tls.Certificate{cert},
                // InsecureSkipVerify: true,
            },
        },
    }
    resp, err = client1_1.Get("https://localhost:443/hello")
    if err != nil {
        log.Errorf("%v", err)
        return
    }
    fmt.Printf("\nhttp/1.1 response:\n")
    printHtml(resp)

    // http 1.1 no TLS
    client1_1 = &http.Client{}
    resp, err = client1_1.Get("http://localhost:8080/hello")
    if err != nil {
        log.Errorf("%v", err)
        return
    }
    fmt.Printf("\nhttp/1.1 response without TLS:\n")
    printHtml(resp)
}
