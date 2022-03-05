# POKT Lint
An open-source diagnostic tool for Pocket Network node runners.

## Using the public API

The public deployment of this tool is available at the following baseURL:
https://2eqrf8goof.execute-api.us-east-1.amazonaws.com/test

### 👉 [Interactive RPC Spec](https://editor.swagger.io/?url=https://raw.githubusercontent.com/itsnoproblem/pokt-lint/master/doc/node-checker-rpc.yml)

---

## Deploying to AWS Lambda

Build the tests as AWS Lambda functions. A public deployment is maintained on AWS.  To build executables that can be uploaded to
AWS Lambda, run the following command:
```
make build-lambda
```

This will create 2 archives that can be uploaded to their corresponding
Lambda functions:
- `build/LambdaPingTestHandler.zip`
- `build/LambdaRelayTestHandler.zip`

---

## Build the test commands locally
This package provides 2 commands that can be used to test the operation of Pocket Network nodes: 
- `pingtest` measures the latency between the test client and a node
- `relaytest` runs relay tests on a node 

The commands can be built and run directly on your host, or they can be built and run in a docker container.



### Option 1) Build directly on your host
_Requirements:_
- [Go 1.17](https://go.dev/doc/install)
- GNU Make

```
# clone the repository
git clone https://github.com/itsnoproblem/pokt-lint

# build the commands
cd pokt-lint
make build-commands
```

### Option 2) Build in a docker container
_Requirements:_
- Docker

```
# clone the repository
git clone https://github.com/itsnoproblem/pokt-lint

# build the commands
cd pokt-lint
docker build .
```
---

This will create 2 executable files:

**./build/pingtest**
```
Usage of ./build/pingtest:
  -num int
    	-num 10 (default 1)
  -url string
    	-url https://www.example.com
```

**./build/relaytest**
```
Usage of ./build/relaytest:
  -chains string
    	comma separated chains ids, eg: -chains=0001,0003,0005
  -id string
    	node id
  -url string
    	node url
```

