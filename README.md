# Generating certificates
## Create folder
* `mkdir secrets`
* `cd secrets`
## Create root CA
* Create a Self-Signed Root CA - `openssl req -x509 -sha256 -days 1825 -newkey rsa:2048 -keyout ca.key -out ca.crt`
* Create configuration text file 'v3.ext'
  
    ```
    authorityKeyIdentifier=keyid,issuer
    basicConstraints=CA:FALSE
    keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
    subjectAltName = @alt_names
    [alt_names]
    DNS.1 = localhost
    ```
## Create server certificate
* Generate certificate signing request (CSR) for server cert - `openssl req -newkey rsa:2048 -nodes -keyout server.key -out server.csr`
* Sign the csr with ca cert generated above - `openssl x509 -req -CA ca.crt -CAkey ca.key -in server.csr -out server.crt -days 365 -CAcreateserial -extfile v3.ext`
* Verify the certificate - `openssl x509 -text -noout -in server.crt`

## Add root cert to truststore
* Create trust store with ca certificate for client `keytool -keystore kafka.client.truststore.jks -alias CARoot -importcert -file ca.crt -storepass changeit` and for server `keytool -keystore kafka.server.truststore.jks -alias CARoot -importcert -file ca.crt -storepass changeit`
  
## Create keystore

* Export server cert to p12 format - `openssl pkcs12 -export -in server.crt -inkey server.key -out server.p12 -name localhost -CAfile ca.crt -caname root`

* Create keystore
    ```
    keytool -importkeystore \
            -deststorepass changeit -destkeypass changeit -destkeystore kafka.server.keystore.jks \
            -srckeystore server.p12 -srcstoretype PKCS12 -srcstorepass changeit \
            -alias localhost
    ```
* Verify cert - `keytool -list -v -keystore kafka.server.keystore.jks`

## Create client certificate
* Generate certificate signing request (CSR) for client cert - `openssl req -newkey rsa:2048 -nodes -keyout client.key -out client.csr`
* Sign the csr with ca cert generated above - `openssl x509 -req -CA ca.crt -CAkey ca.key -in client.csr -out client.crt -days 365 -CAcreateserial -extfile v3.ext`
* Verify the certificate - `openssl x509 -text -noout -in client.crt`