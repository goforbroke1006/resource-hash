# resource-hash

Parallel check links list, print md5 for each one.

### Requirements

* Go >=1.16
* [golangci-lint](https://golangci-lint.run/usage/install/)

### How to run

```shell
make
./resource-hash --concurrency=5 --filename=links.txt
```