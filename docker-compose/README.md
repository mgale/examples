# Example on using docker-compose from within Golang

```
go run .
....
Example docker-compose management
==================================

Docker-compose body:

version: "3.8"
services:
  testservice:
    image: ubuntu:latest
    environment:
      - COMPOSE_TEST=TestMe123
    command: sleep infinity
    init: true


==================================
Docker service up...
[+] Running 2/2
 ✔ Network testproject_default          Created                                                                                                                                                                  0.1s
 ✔ Container testproject-testservice-1  Started                                                                                                                                                                  0.1s
hello world
2023/12/02 19:59:18 Command result: 0  and err: <nil>
bin  boot  dev  etc  home  lib  lib32  lib64  libx32  media  mnt  opt  proc  root  run  sbin  srv  sys  tmp  usr  var
2023/12/02 19:59:18 Command result: 0  and err: <nil>
Docker service down...
[+] Running 2/2
 ✔ Container testproject-testservice-1  Removed                                                                                                                                                                  0.2s
 ✔ Network testproject_default          Removed
```
