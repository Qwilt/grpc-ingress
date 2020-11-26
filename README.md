# grpc-ingress
Hello world example of GRPC exposed in GKE over ingress

### prerequisites
+ linux or mac development environment 
+ go 1.15 (though 1.14 should work)
+ docker
+ GKE cluster ( this poc used v1.17.13-gke.1400 and it seems that other versions may behave very differently) 
+ kubectl configured to a namespace you can work with
+ access to a dns server. GCP CloudDNS is good but anyone would do (even `/etc/hosts`)

###setup
+ replace all instances of <my domain> with your real domain name e.g. grpcserver.mycompany.com  
files: Makefile,deployment.yaml,main/cli.go:
+ replace all instances of <my docker image> with your docker image name e.g. mycompany/grpcserver  
files: Makefile,deployment.yaml
+ install protoc and protoc-gen-go make sure that protoc-gen-go is in your PATH
+ run `make generate` which will create `chat/chat.pb.go` and `chat/chat_grpc.pb.go`
+ run `make cert` to generate a self signed certificate and private key in the `cert` directory
###test locally
+ run `go run main/server.go` verify that the server is starting
+ install `evans` from https://github.com/ktr0731/evans
+ run `echo '{"body":"my grpc service"}'| evans --host 127.0.0.1  -p 8443 -t -r --servername <my domain>
   --cacert cert/cert.pem cli call chat.ChatService.SayHello`
+ verify that you get a response.
+ recommended: play with evans in REPL interactive mode (very cool): `evans --host 127.0.0.1  -p 8443 -t -r --servername <my domain>
                                                        --cacert cert/cert.pem` 
### build and deploy
+ run `make build` - verify an executable was created in the bin directory `ls -l  bin`
+ run `make docker` - verify with `docker images` that a docker image was created 
+ run `make push` - very that push succeeded (login before if you need to)
+ create a `tls` secret based on newly generated cert and key `make secret` - verify with `kubectl get secrets` that the secret was created
+ deploy to k8s `kubectl apply -f deployment.yaml`
+ wait a while - can take 20 mins
+ find your ingress external IP `kubectl get ing grpcserver -ojsonpath='{.status.loadBalancer.ingress[0].ip}'` 
and update your dns (gcp https://console.cloud.google.com/net-services/dns). map an A record from <your domain> to the 
ingress external IP
+ wait until the dns record propagates. You can check it with `dig
+ run `go run main/cli.go` and verify that the response is "Hello I'm the client" 
+ you can also try evans `evans --host <my domain> -p 443 -t -r --cacert cert/cert.pem`

##Open Issues
+ GRPC [Health Checking Protocol](https://github.com/grpc/grpc/blob/master/doc/health-checking.md) works for kubernetes container probes as shown in https://cloud.google.com/blog/topics/developers-practitioners/health-checking-your-grpc-servers-gke, 
however it doesn't work for GCP External HTTP(S) LoadBalancer created from Ingress.  
[BackendConfig](https://cloud.google.com/kubernetes-engine/docs/concepts/ingress#health_checks) in GRPC protocol as specified in https://cloud.google.com/load-balancing/docs/health-check-concepts#category_and_protocol does not work in the GKE version specified above
+ What are the requirements for self signed certificates? Some certificates did not work for me. Will an official CA certificate work?
+ Why do [managed certificates](https://cloud.google.com/kubernetes-engine/docs/how-to/managed-certs) not work? Using it seemed to keep port 443 closed at the LB Listener 