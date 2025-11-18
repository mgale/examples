# Example on using docker-compose from within Golang

```
go run .
....
Example docker-compose management
==================================

Docker-compose body:

services:
  testservice:
    image: ubuntu:latest
    environment:
      - COMPOSE_TEST=TestMe123
    command: sleep infinity
    volumes:
      - /tmp:/tmp
    init: true
  testservice2:
    build: .


==================================
Docker compose project created, list of known services:
Service: testservice
  Image: ubuntu:latest
Service: testservice2
  Image:
Docker service up...
[+] Running 3/3
 ✔ Network testproject_default           Created                                                                                                                                                                                                 0.0s
 ✔ Container testproject-testservice2-1  Started                                                                                                                                                                                                 0.6s
 ✔ Container testproject-testservice-1   Started                                                                                                                                                                                                 0.6s
hello world
2025/11/18 18:59:39 Command result: 0  and err: <nil>
bin  boot  dev  etc  home  lib  lib64  media  mnt  opt  proc  root  run  sbin  srv  sys  tmp  usr  var
2025/11/18 18:59:39 Command result: 0  and err: <nil>
Docker service down...
[+] Running 3/3
 ✔ Container testproject-testservice-1   Removed                                                                                                                                                                                                 0.1s
 ✔ Container testproject-testservice2-1  Removed                                                                                                                                                                                                10.1s
 ✔ Network testproject_default           Removed
```
