# Https server and client in golang

Simple https server and client testing code in golang with self signed certificate

## Running server

```sh
sudo /usr/local/go/bin/go run https_server_client/https-server.go  \
                          -alsologtostderr -x509script https_server_client/
```

`-x509script` provide the path to x509-key-pair.sh script which generates self signed certificate for both server and client. We reuse the same certicate just for testing purpose.

`sudo` is used here for official https port `443` is being tested against.
 
## Running client

```sh
/usr/local/go/bin/go run https-client.go
```

Client will access server in three different versions: http/22, http/1.1 and http/1.1 without TLS.
```
http/2 response:
200 OK
This is an example server.
This is an example server ---.
This is an example server ^^^^.

http/1.1 response:
200 OK
This is an example server.
This is an example server ---.
This is an example server ^^^^.

http/1.1 response without TLS:
200 OK
This is an example server.
This is an example server ---.
This is an example server ^^^^.
```

if `-alsologtostderr` is provided for running server, server side output may be seen as well. 
Check the example out for the difference
```
I0901 07:28:34.173480    3564 https-server.go:94] Extracted path: /hello from request:
	 {Method:GET URL:/hello Proto:HTTP/2.0 ProtoMajor:2 ProtoMinor:0 Header:map[Accept-Encoding:[gzip] User-Agent:[Go-http-client/2.0]] Body:0xc000304900 GetBody:<nil> ContentLength:0 TransferEncoding:[] Close:false Host:localhost:443 Form:map[] PostForm:map[] MultipartForm:<nil> Trailer:map[] RemoteAddr:127.0.0.1:42202 RequestURI:/hello TLS:0xc0003740b0 Cancel:<nil> Response:<nil> ctx:0xc0003681c0}
```
```	 
I0901 07:28:34.194044    3564 https-server.go:94] Extracted path: /hello from request:
	 {Method:GET URL:/hello Proto:HTTP/1.1 ProtoMajor:1 ProtoMinor:1 Header:map[Accept-Encoding:[gzip] User-Agent:[Go-http-client/1.1]] Body:{} GetBody:<nil> ContentLength:0 TransferEncoding:[] Close:false Host:localhost:443 Form:map[] PostForm:map[] MultipartForm:<nil> Trailer:map[] RemoteAddr:127.0.0.1:42204 RequestURI:/hello TLS:0xc00011a2c0 Cancel:<nil> Response:<nil> ctx:0xc0000127c0}
```	 
```
I0901 07:28:34.194521    3564 https-server.go:94] Extracted path: /hello from request:
	 {Method:GET URL:/hello Proto:HTTP/1.1 ProtoMajor:1 ProtoMinor:1 Header:map[Accept-Encoding:[gzip] User-Agent:[Go-http-client/1.1]] Body:{} GetBody:<nil> ContentLength:0 TransferEncoding:[] Close:false Host:localhost:8080 Form:map[] PostForm:map[] MultipartForm:<nil> Trailer:map[] RemoteAddr:127.0.0.1:56284 RequestURI:/hello TLS:<nil> Cancel:<nil> Response:<nil> ctx:0xc000012840}
```

# Access server from web browser

Use address: `https://localhost/`

chrome browser  will warn you with message like `Your connection is not private` that is because the certificate is self signed and not from known authority. Select 'Advance' then 'Proceed to localhost (unsafe)' (Don't worry, this is just the test server program), you should see:

```
{"status":"ok"}
```
Try different address and check the server debug output to understand the message format.

Note: if `-clientcert` is also set when running server, web browser will be rejected and get error 
```
localhost unexpectedly closed the connection.
```

Again, check the server debug output:
```
2020/09/01 07:53:23 http: TLS handshake error from [::1]:40018: remote error: tls: unknown certificate
2020/09/01 07:53:27 http: TLS handshake error from [::1]:40022: remote error: tls: unknown certificate
2020/09/01 07:53:27 http: TLS handshake error from [::1]:40024: tls: client didn't provide a certificate
```
# Reference
Largely based on
https://github.com/jcbsmpsn/golang-https-example


