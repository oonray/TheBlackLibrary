## TCP Forwarder in go

Forwards from one TCP port on host to one on target.

```
local==>host==>target
        8081   80
```

### Usage 
```
Usage of tcpforwarder:
  -H string
        The host to conenct to
  -P string
        The port to connect to
  -lp string
        The port to listen on
```

### Dependency 
1. Logging: [github.com/sirupsen/logrus](http://github.com/sirupsen/logrus)
