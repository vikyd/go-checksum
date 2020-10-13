# go-checksum

Simple tool to calc Golang module checksum of `go.mod` and `module directory`.

This README also describes how Golang go.sum calc hash from file conent.

> Maybe this is the missing official doc for how Golang go.sum calc hash from file conent

# Usage

## Example Golang Module

- Github: https://github.com/gin-gonic/gin
- Example Version: v1.4.0
- CheckSum: https://sum.golang.org/lookup/github.com/gin-gonic/gin@v1.4.0

## Clone and checkout a version

```sh
git clone https://github.com/gin-gonic/gin.git
git checkout v1.4.0
```

## Install go-checksum

```sh
go get github.com/vikyd/go-checksum
```

## Calc checksum of go.mod

```sh
go-checksum relOrAbsPathOfGinDir/go.mod
```

should prints like:

```
file: /Users/viky/tmp/gin/go.mod
{
	"Hash": "b1946b355a09fedc3865073bbeb5214de686542bbeaea9e1bf83cb3a08f66b99",
	"HashBase64": "sZRrNVoJ/tw4ZQc7vrUhTeaGVCu+rqnhv4PLOgj2a5k=",
	"HashSynthesized": "396d84667dc33bc2e7f6820a3af33ef8b04efb950f1c92431fbdbfabfdeb65d3",
	"HashSynthesizedBase64": "OW2EZn3DO8Ln9oIKOvM++LBO+5UPHJJDH72/q/3rZdM=",
	"GoCheckSum": "h1:OW2EZn3DO8Ln9oIKOvM++LBO+5UPHJJDH72/q/3rZdM="
}
```

The output hash `"GoCheckSum": "h1:OW2EZn3DO8Ln9oIKOvM++LBO+5UPHJJDH72/q/3rZdM="`

is the same as

the online checksum: `github.com/gin-gonic/gin v1.4.0/go.mod h1:OW2EZn3DO8Ln9oIKOvM++LBO+5UPHJJDH72/q/3rZdM=`

(from [here](https://sum.golang.org/lookup/github.com/gin-gonic/gin@v1.4.0))

## Calc checksum of module directory

```sh
go-checksum relOrAbsPathOfGinDir github.com/gin-gonic/gin@v1.4.0
```

- 1st param: module directory
  - no need to remove .git, will ignore
- 2nd param: module prefix with version
  - necessary for dir checksum

should prints like:

```
directory: /Users/viky/tmp/gin
{
	"HashSynthesized": "ded3280827ccee9a6ab11d29b73ff08b58a6a4da53efff7042d319f25af59824",
	"HashSynthesizedBase64": "3tMoCCfM7ppqsR0ptz/wi1impNpT7/9wQtMZ8lr1mCQ=",
	"GoCheckSum": "h1:3tMoCCfM7ppqsR0ptz/wi1impNpT7/9wQtMZ8lr1mCQ="
}
```

The output hash `"GoCheckSum": "h1:3tMoCCfM7ppqsR0ptz/wi1impNpT7/9wQtMZ8lr1mCQ="`

is the same as

the online checksum: `github.com/gin-gonic/gin v1.4.0 h1:3tMoCCfM7ppqsR0ptz/wi1impNpT7/9wQtMZ8lr1mCQ=`

(from [here](https://sum.golang.org/lookup/github.com/gin-gonic/gin@v1.4.0))

# Explain

## How go.mod checksum works?

Steps:

- tell me the path of `go.mod`
- read content from `go.mod`, as variable `content`
- calc SHA256 hash from `content`, as variable `hash`
- mix some string with `hash`, as variable `mixedHash`
  - ```
    hash  go.mod\n
    ```
  - if `hash` is: `CDa7N`
  - then `mixedHash` is:
  - ```
    CDa7N  go.mod\n
    ```
  - yeah, that's a strange string
- calc SHA256 hash from `mixedHash`, as variable `hashSynthesized`
- calc Base64 from `hashSynthesized`, as variable `hashSynthesizedBase64`
- finally, in `go.sum`, the string is: `h1:hashSynthesizedBase64`, as variable `GoCheckSum`
  - if `hashSynthesized` is: `CCfM7`
  - then `GoCheckSum` is: `h1:CCfM7`
  - `h1`: [The hash begins with an algorithm prefix of the form "h<N>:". The only defined algorithm prefix is "h1:", which uses SHA-256.](https://tip.golang.org/cmd/go/#hdr-Module_authentication_using_go_sum)
  - `h`: hash

## How module direcory checksum works?

Steps:

- tell me the path of `module's dir` and `module's prefix` (also named import path)
  - `module's prefix` is used to calc hash, as part of the content string
- clean the dir path, to remove duplicate path seperator etc.
- loop each file in `module's dir`
  - only consider file, not dir
  - find the relative path of the file (relative to `module's dir`)
  - join file relative path with `module's prefix`, as variable `fileImportPath`
  - example:
    - one file in gin project: `gin.go`
    - its module's dir: `/dir01/dir02/gin`
    - its absolute path: `/dir01/dir02/gin/gin.go`
    - its relative path: `gin.go`
    - final string `fileImportPath`: `github.com/gin-gonic/gin@v1.4.0/gin.go`
- after the loop, we got a list `fileImportPath` of all files, as variable `files`
  - except files in `.git`
- sort `files` as string in increasion order
- then begin hash steps
- loop sorted `files`
  - read content from a file of `files`, as variable `content`
  - calc SHA256 hash from `content`, as variable `hash`
  - mix some string with `hash`, as variable `mixedHash`
    - ```
      hash  fileImportPath\n
      ```
    - if `hash` is: `CDa7N`
    - if `fileImportPath` is: `github.com/gin-gonic/gin@v1.4.0/gin.go`
    - then `mixedHash` is:
    - ```
      CDa7N  github.com/gin-gonic/gin@v1.4.0/gin.go\n
      ```
    - yeah, that's a strange string
- after this loop, we can get a single long string joined by all files' `mixedHash`, as variable `mixedHashAll`
  - example:
  - ```
    CDa7N  github.com/gin-gonic/gin@v1.4.0/gin.go\nEFb8M  github.com/gin-gonic/gin@v1.4.0/context.go\n ...
    ```
- calc SHA256 hash from `mixedHashAll`, as variable `hashSynthesized`
- calc Base64 from `hashSynthesized`, as variable `hashSynthesizedBase64`
- finally, in `go.sum`, the string is: `h1:hashSynthesizedBase64`, as variable `GoCheckSum`
  - if `hashSynthesized` is: `CCfM7`
  - then `GoCheckSum` is: `h1:CCfM7`
  - `h1`: [The hash begins with an algorithm prefix of the form "h<N>:". The only defined algorithm prefix is "h1:", which uses SHA-256.](https://tip.golang.org/cmd/go/#hdr-Module_authentication_using_go_sum)
  - `h`: hash

# Ref

Golang source code about how to calc hash for modules: https://github.com/golang/mod/blob/release-branch.go1.15/sumdb/dirhash/hash.go
