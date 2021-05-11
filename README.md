
Upload files to local from anywhere.

## features

1. split big file into smaller pieces, to avoid request size limitation of 3rd reverse proxy (e.g., localtunnel)
2. use free reverse proxy for accessing server publicly

### TODO

1. data encryption
2. md5 hash check for file
3. use a stable free reverse proxy

## Usage

```
$ ./upload2local 
a tool to upload file to local from anywhere

Usage:
  upload2local [command]

Available Commands:
  help        Help about any command
  server      start server for receiving files
  testdata    generate binary data
  upload      upload local file to server
  version     binary version

Flags:
  -h, --help      help for upload2local
  -v, --verbose   show verbose log

Use "upload2local [command] --help" for more information about a command.
```

### server

```
$ ./upload2local server   
local access: http://127.0.0.1:1234
public access: https://quiet-fireant-51.loca.lt
```

### upload

```
# ll -h a.pcap 
-rw-r--r-- 1 root root 2.3M May 10 15:50 a.pcap

$ ./upload2local upload --host https://quiet-fireant-51.loca.lt -s 1000000 a.pcap
INFO[2021-05-11T09:05:34Z] uploading file part /tmp/a.pcap_0_part ...   
INFO[2021-05-11T09:05:39Z] uploading file part /tmp/a.pcap_1_part ...   
INFO[2021-05-11T09:05:42Z] uploading file part /tmp/a.pcap_2_part ...   
INFO[2021-05-11T09:05:45Z] file a.pcap uploaded                         

```

NOTE: localtunnel may be out of work after 10 http requests, try increase s, at most 1 MB.

## 3rd reverse proxy

### localtunnel

https://theboroer.github.io/localtunnel-www/

```
$ npm install -g localtunnel
$ lt --port 1234
your url is: https://smooth-dragon-53.loca.lt
```

go implementation: https://github.com/NoahShen/gotunnelme
