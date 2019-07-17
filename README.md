# sk

![usage](./image/sk.gif)

# Examlpe:
vim ~/.sk.yml
```yaml
- name: Aliyun Server
  children:
  - { name: http server, user: ctg, host: 192.168.11.11, port: 22, password: hellohello }
  - { name: redis server, user: root, host: 192.168.11.12, port: 22, password: hellohello }
  - { name: nginx server, user: root, host: 192.168.11.13, port: 22, password: hellohello }
  - { name: mysql server, user: root, host: 192.168.11.14, port: 22, password: hellohello }

- name: AWS Server
  children:
  - { name: Jenkins , user: ec2-user, host: 192.168.128.16, port: 22, keypath: ./Ops-key.pem}
  - { name: Jumpserver, user: root, host: 192.168.32.13, port: 15622, password: 123456 }
  - { name: zookeeper, user: root, host: 192.168.32.14, port: 15622, password: 123456 }
  - { name: jenkins-old, user: ec2-user, host: 192.168.32.16, port: 22, keypath: ~/.ssh/id_rsa }


- { name: Host1, user: root, host: 10.0.123.12, port: 22, keypath: ./key/proc-wall-jmp.pem }
- { name: Host2, user: root, host: 10.0.123.13, port: 22, password: 123456 }
```
