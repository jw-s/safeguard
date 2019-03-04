# safeguard

Safeguard is a custom admission controller for kubernetes used to enforce protection on kubernetes resources.

## Usage

1.  Configure [safeguard.yml](safeguard.yml) with your own ca bundle base64 encoded.
2.  Modify [secret.yml](contrib/secret.yml) with your own `tls.crt` and `tls.key`

    NOTE: the certificates have to be signed by the same CA as your api server!
3. Run the following;
```
kubectl create -f safeguard.yml -n NAMESPACE
kubectl create -f contrib/secret.yml -n NAMESPACE
kubectl create -f contrib/deployment.yml -n NAMESPACE # this should be in the same namespace as the secret
kubectl create -f contrib/service.yml -n NAMESPACE # this should be in the same namespace as the deployment
```

4. Decorate your resources with the following annotation to protect them.
```
---
apiVersion: v1
kind: Service
metadata:
  name: example
  annotations:
    safeguard.jw-s.com/protected: 'true'
...

```

5. Try to delete the protected resource!


## Development
### Prerequistities
* Go 1.12.x
* Make

```bash
go get -d github.com/jw-s/safeguard
cd $GOPATH/src/github.com/jw-s/safeguard
make build
```