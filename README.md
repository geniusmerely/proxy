# Simple proxy server
## configure proxy.toml
### bind_ip
your ip for outgoing calls. Default empty (default ip of your machine)
```toml
#example
bind_ip = "ffff:ffff:ffff:2" #supported ipv6 or ipv4
# if empty will be used default local ip
bind_ip = ""
```
### listen_addr 
which ip or port to bind the server to. Default ":8888"
```toml
#example
listen_addr = ":8080"
```
### users
Do not specify [[users]] settings to make the proxy open
```toml
#example
[[users]]
username = "user1"
password = "userpassword"
[[users]]
username = "user2"
password = "userpassword"
```
## how to use
```bash
./proxy -conf ./proxy.toml
```
where proxy.toml is the path to the config file
