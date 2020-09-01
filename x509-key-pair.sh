#!/usr/bin/env bash

# for debugging
# openssl x509 -in local_ca.crt -text -noout

# These certs file is only for localhost testing.
HOST='localhost'
CERT_DIR='/tmp/cert'

mkdir -p $CERT_DIR
pushd $CERT_DIR
# Create CA certificate
openssl req \
    -newkey rsa:4096 -nodes -sha256 -keyout local_ca.key \
    -x509 -days 365 -out local_ca.crt -subj '/C=US/ST=CA/L=San Jose/O=JPY/CN=*/emailAddress=jipan.yang@gmail.com'

# Generate a Certificate Signing Request
 # Key considerations for algorithm "RSA" ≥ 4096-bit
openssl req \
    -newkey rsa:4096 -nodes -sha256 -keyout $HOST.key \
    -out $HOST.csr -subj '/C=US/ST=CA/L=San Jose/O=JPY/CN=*/emailAddress=jipan.yang@gmail.com'

# Generate the certificate of local registry host
# echo subjectAltName = IP:$IP > extfile.cnf
echo subjectAltName = DNS:$HOST > extfile.cnf
openssl x509 -req -days 3650 -in $HOST.csr -CA local_ca.crt \
	-CAkey local_ca.key -CAcreateserial -extfile extfile.cnf -out $HOST.crt	
	
# Copy to cert default location
cp $HOST.crt $CERT_DIR/server.crt
cp $HOST.key $CERT_DIR/server.key

# Also use as client key pair
cp $HOST.crt $CERT_DIR/client.crt
cp $HOST.key $CERT_DIR/client.key
popd