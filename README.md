# Memory collector

1.) Spin up Vagrant instance.
```bash
cd vagrant
vagrant up 
```
1.1) fill in all necessary variables in golang or python scripts.
 * Python
  ```bash
  vim /python/main.py
  hosts = ["10.1.8.19","10.1.9.30",]
  username="" # dkorol
  keyPath="" # /Users/dmitriy_korol/.ssh/id_rsa

  ```
  * Go
```go
vim /src/main.go

hostmap["1"]="10.1.8.19"
hostmap["2"]="10.1.9.30"
...
//set username for ssh connection
username:="" //dkorol
rsaPath:="" // Users/dmitriy_korol/.ssh/id_rsa

```

2.) run python or golang client 

  2.1) Python
  ```bash
cd python 
python main.py
```
2.2) Golang
 - Install golang to your system. <br>
 - copy src folder to your GOPATH.
 - run goclient
 ```bash
 export GOPATH=/home/user/Mygoworkspace/
 cd go && cp -r ./src $GOPATH
  cd $GOPATH/src && go run main.go
  ```
             
3.) Observe metrics from your infrastructure in Graphite Browser

[https://127.0.0.1:8443/](https://127.0.0.1:8443/)