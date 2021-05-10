
a file upload server and ngrok like public ip/dns

## localtunnel

### setup
```
$ npm install -g localtunnel
$ lt --port 8090
your url is: https://smooth-dragon-53.loca.lt
```

### 413 Request Entity Too Large
do not support large file...

```
$ curl -F "data=@/Users/jiahua/Downloads/ngrok-stable-darwin-amd64.zip" https://smooth-dragon-53.loca.lt
<html>
<head><title>413 Request Entity Too Large</title></head>
<body>
<center><h1>413 Request Entity Too Large</h1></center>
<hr><center>nginx/1.17.9</center>
</body>
</html>
$ ll ~/Downloads/ngrok-stable-darwin-amd64.zip 
-rw-r--r--@ 1 jiahua  staff    13M May 11 00:56 /Users/jiahua/Downloads/ngrok-stable-darwin-amd64.zip
```