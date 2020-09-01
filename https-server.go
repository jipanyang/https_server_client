// test run:
// sudo /usr/local/go/bin/go run ./https-server.go -alsologtostderr
//
// To force client cert verification:
// sudo /usr/local/go/bin/go run ./https-server.go -alsologtostderr -clientcert
//
// Path to x509-key-pair.sh should be specified if not in current folder
// sudo /usr/local/go/bin/go run https_server_client/https-server.go  -alsologtostderr -x509script https_server_client/
package main

import (
    "bytes"
    "crypto/tls"
    "crypto/x509"
    "flag"
    "fmt"
    "html"
    "io"
    "io/ioutil"
    "net/http"
    "os/exec"

    log "github.com/golang/glog"
)

var (
    SERVER_CERT = "/tmp/cert/server.crt"
    SERVER_KEY  = "/tmp/cert/server.key"
    CLIENT_CERT = "/tmp/cert/client.crt"

    requireClientCert = flag.Bool("clientcert", false, "When set, RequireAndVerifyClientCert")
    x509ScriptPath    = flag.String("x509script", ".", "Path to the x509-key-pair.sh script")
)

func generateX509Cert() {
    script := *x509ScriptPath + "/" + "x509-key-pair.sh"
    cmd := exec.Command("/bin/bash", script)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    err := cmd.Run()
    if err != nil {
        log.Fatal(err)
    }
    log.V(1).Infof("x509-key-pair generated: %V\n", out.String())
}

func main() {
    flag.Parse()
    generateX509Cert()

    // Accept client certificate as signed by known authority
    caCert, err := ioutil.ReadFile(CLIENT_CERT)
    if err != nil {
        log.Fatal(err)
    }
    caCertPool := x509.NewCertPool()
    caCertPool.AppendCertsFromPEM(caCert)

    cfg := &tls.Config{
        // RequestClientCert will ask client for a certificate but won't
        // require it to proceed. If certificate is provided, it will be
        // verified.
        ClientAuth: tls.RequestClientCert,
        ClientCAs:  caCertPool,
    }
    if *requireClientCert {
        cfg.ClientAuth = tls.RequireAndVerifyClientCert
    }
    srvHttps := &http.Server{
        Addr:      ":443",
        Handler:   &handler{},
        TLSConfig: cfg,
    }

    srvHttp := &http.Server{
        Addr:    ":8080",
        Handler: &handler{},
    }

    // server the same content on http
    go srvHttp.ListenAndServe()
    err = srvHttps.ListenAndServeTLS(SERVER_CERT, SERVER_KEY)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}

type handler struct{}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

    path := html.EscapeString(r.URL.Path)
    log.Infof("Extracted path: %v from request:\n\t %+v\n", path, *r)

    switch path {
    case "/":
        w.Header().Add("Content-Type", "application/json")
        io.WriteString(w, `{"status":"ok"}`)
    case "/hello":
        w.Header().Set("Content-Type", "text/plain")
        w.Write([]byte("This is an example server.\n"))
        fmt.Fprintf(w, "This is an example server ---.\n")
        io.WriteString(w, "This is an example server ^^^^.\n")
    default:
        fmt.Fprintf(w, "Hello, %q\n", html.EscapeString(r.URL.Path))
    }
}
