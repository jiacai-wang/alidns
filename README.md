# alidns
update domain record using AliCloud Go SDK.

### Usage
- First, install Aliyun Go SDK by running `go get github.com/aliyun/alibaba-cloud-sdk-go/services/alidns`.
- Create a AccessKey at https://usercenter.console.aliyun.com/.
- Fill in `config.json` the AccessKey, secret, and the domain you'd setup.
- To compile, just run `go build alidns.go`. *Run `go mod init alidns` first if needed.*
- If provided correct config file, things should work fine.

ip is fetched from http://ip-api.com/json
