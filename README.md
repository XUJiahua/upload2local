
Upload files to local from anywhere.

## features

1. split big file into smaller pieces, to avoid request size limitation of 3rd reverse proxy (e.g., localtunnel)
2. use free reverse proxy for accessing server publicly

### TODO

1. data encryption
2. md5 hash check for file

## Usage

```
```

## 3rd reverse proxy

### localtunnel

https://theboroer.github.io/localtunnel-www/

```
$ npm install -g localtunnel
$ lt --port 1234
your url is: https://smooth-dragon-53.loca.lt
```

go implementation: https://github.com/NoahShen/gotunnelme
