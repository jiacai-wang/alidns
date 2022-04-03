# alidns
update domain record using AliCloud Go SDK. ip is fetched from https://ifconfig.me.

更新阿里云域名的解析IP地址。IP地址从https://ifconfig.me 获取


### Usage
- 在[此处](https://usercenter.console.aliyun.com/)申请AccessKey。
- 在 `config.json`文件中填入AccessKey和AccessSecret，以及需要更新的域名、主机记录。
- 运行`alidns`程序。如果config.json配置文件信息正确，会输出success信息。
- (默认配置文件为`./config.json`，如需自定义配置文件，执行`alidns -config /path/to/config.json`)


- First, create an AccessKey at https://usercenter.console.aliyun.com/.
- Then fill in `config.json` the AccessKey, AccessSecret, and the domain you want to setup.
- Run `alidns`. If provided correct config file, things should work fine.
- (config file defaults to `./config.json`. To specify another config file, run `alidns -config /path/to/config.json`)


> To compile from source
> - Install Aliyun Go SDK by running `go get -u github.com/aliyun/alibaba-cloud-sdk-go/sdk`.
> - Run `go build alidns.go`. *Run `go mod init alidns` if needed.*

### Release

Windows
> [alidns.zip (to be updated.)](https://github.com/jiacai-wang/alidns/releases/download/v0.1/alidns.x86.zip)

Linux
> [alidns.tar.gz (to be updated.)](https://github.com/jiacai-wang/alidns/releases/download/v0.1/alidns.x86.tar.gz)
