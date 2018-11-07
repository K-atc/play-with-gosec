Play with gosec
====

[gosec](https://github.com/securego/gosec) is security checker for golang code.


Install
----
### gosec
```shell
export GOPATH=$PWD
go get github.com/securego/gosec/cmd/gosec/...
```

Usage
----
### Try gosec
```shell
./bin/gosec -fmt=json -out gosec-md5.json ./src/md5
./bin/gosec -fmt=json -out gosec-ssrf.json ./src/ssrf
```

#### Example
```
K_atc% ./bin/gosec -fmt=json -out gosec-md5.json ./src/md5
... snipped ...

K_atc% cat gosec-md5.json 
{
    "Issues": [
        {
            "severity": "MEDIUM",
            "confidence": "HIGH",
            "rule_id": "G501",
            "details": "Blacklisted import crypto/md5: weak cryptographic primitive",
            "file": "/home/katc/project/gosec-jenkins/src/md5/main.go",
            "code": "\"crypto/md5\"",
            "line": "5"
        },
        {
            "severity": "MEDIUM",
            "confidence": "HIGH",
            "rule_id": "G401",
            "details": "Use of weak cryptographic primitive",
            "file": "/home/katc/project/gosec-jenkins/src/md5/main.go",
            "code": "md5.New()",
            "line": "14"
        }
    ],
    "Stats": {
        "files": 1,
        "lines": 17,
        "nosec": 0,
        "found": 2
    }
}%
```

### Reave issue detected by gosec on commit message
```shell
./bin/gosec -fmt=json -out gosec.json src/...
./comment_on_github.py gosec.json
```

see https://github.com/K-atc/play-with-gosec/commit/0dc89bcc22163cf3372e79dd5c2b185d3c68a9ff to check result


(Optional) Build sample golang apps
----
```shell
export GOPATH=$PWD
go build md5
go build ssrf
```