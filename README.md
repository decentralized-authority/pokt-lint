# POKT Lint
An open-source diagnostic tool for Pocket Network node runners.

- [Using the public API](#using-the-public-api)
- [Deploying to AWS Lambda](#deploying-to-aws-lambda)
- [Build the test commands locally](#build-the-test-commands-locally)
- [Dev tools (Makefile) reference](#dev-toolchain-reference)
---
## Using the public API

The public deployment of this tool is available at the following baseURLs:
- staging: https://2eqrf8goof.execute-api.us-east-1.amazonaws.com/test
- production: https://2eqrf8goof.execute-api.us-east-1.amazonaws.com/prod

### 👉 [Interactive OpenAPI Spec](https://editor.swagger.io/?url=https://raw.githubusercontent.com/itsnoproblem/pokt-lint/master/doc/node-checker-rpc.yml)

### Run the OpenAPI docs locally
```bash
make docserver
```

---

## Deploying to AWS Lambda

A public deployment is maintained on AWS. To build executables that can be uploaded to
AWS Lambda, run the following command:
```
make build-lambda
```

This will create 2 archives that can be uploaded to their corresponding
Lambda functions:
- `build/LambdaPingTestHandler.zip`
- `build/LambdaRelayTestHandler.zip`

You can also build *and* deploy to AWS using the aws-cli:
```
make deploy-lambda-all-qa 
# or
make deploy-lambda-all
```
This uses the aws-cli docker container and expects `~/.aws/config` and `~/.aws/credentials` to exist. 


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
---
## Dev Toolchain Reference
```
pokt-lint % make
usage: make [target] ...

targets:
-------
help                            Show this help message.
docserver                       Run an interactive OpenAPI spec on port 3333
docserver-stop                  Stop the interactive spec
build-commands                  <-- compiles executables to ${BUILD_DIR}
build-lambda                    <-- builds lambda function bundles in ${BUILD_DIR}
build-lambda-pingtest           builds the pingtest lambda function
build-lambda-relaytest          builds the relaytest lambda function
build-lambda-cors               builds the cors handler (to return access-control-* headers)
deploy-lambda-all-qa            <-- deploys all lambda QA functions
deploy-lambda-all               <-- deploys all lambda functions
deploy-lambda-pingtest-qa       builds and deploys the pingtest QA function
deploy-lambda-relaytest-qa      builds and deploys the relaytest QA function
deploy-lambda-cors-qa           builds and deploys the CORS QA function
deploy-lambda-pingtest          builds and deploys the pingtest function
deploy-lambda-relaytest         builds and deploys the relaytest function
deploy-lambda-cors              builds and deploys the CORS function
test                            runs the unit tests
clean                           deletes build artifacts
```