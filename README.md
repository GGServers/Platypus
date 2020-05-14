# Platypus
## Simple realtime server monitoring

[![CircleCI](https://circleci.com/gh/gmemstr/Platypus.svg?style=svg)](https://circleci.com/gh/gmemstr/Platypus)
[![https://img.shields.io/badge/demo-status.gmem.ca-black.svg](https://img.shields.io/badge/demo-status.gmem.ca-black.svg?style=for-the-badge)](https://status.gmem.ca) 
[![https://img.shields.io/badge/frontend-gmemstr%2Fplatypus--react-blue.svg](https://img.shields.io/badge/frontend-gmemstr%2Fplatypus--react-blue.svg?style=for-the-badge)](https://github.com/gmemstr/platypus-react)

### Dependencies

```bash
go get github.com/gorilla/mux
go get github.com/gorilla/websocket
go get github.com/go-yaml/yaml
go get github.com/shirou/gopsutil
```

### Usage

Master server:
```bash
go build -o platypus main.go
chmod +x platypus
./platypus
```

Client servers:
```bash
go build -o platypus_client client/client.go
chmod +x platypus_client
nano config.yml
# Input your secret key and master server IP here, secret key found on master server in .secret
# master: example.com
# secret: s3cr3tk3y
# End config
./platypus_client
```

Navigate to your master server and check out the stats.

## Rewrite 

Rewriting this from the ground up. Why did I do this in Python.

The goal of the rewrite is to move away from using Python for the entire stack 
and instead break things up into smaller chunks, maybe moving this to it's own
GitHub / Gitlab org, which will allow it to be much more modular and open ended
when it comes to what kind of information you want to monitor and how.

### Polling

There is an optional route that works well with plugins - the `/poll` endpoint.

Requires a few headers
