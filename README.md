# updater

Simple app that runs `git checkout main && git fetch && git pull --ff-only` in all subdirectories of the current directory.
I could have used a shell script, but it was easier to write it in Go.

## Installation

```sh
go install github.com/SuddenGunter/updater@v0.1.0
```

## Usage

```sh
$ updater
2022/11/10 01:38:19 puller: 3, starting
2022/11/10 01:38:19 puller: 1, starting
2022/11/10 01:38:19 puller: 0, starting
2022/11/10 01:38:19 puller: 2, starting
2022/11/10 01:38:19 puller: 3, failed to fetch and pull repo: curltest, error: failed to run command: [git checkout main], error: exit status 128
2022/11/10 01:38:22 puller: 1, updated repo: argocd-manifests
2022/11/10 01:38:25 puller: 0, updated repo: mobile-api
2022/11/10 01:38:29 puller: 0, updated repo: spendshelf-backed
2022/11/10 01:38:29 puller: 2, updated repo: spendshelf-frontend
2022/11/10 01:38:32 puller: 0, finished
2022/11/10 01:38:32 puller: 3, finished
2022/11/10 01:38:32 puller: 1, finished
2022/11/10 01:38:32 puller: 2, finished

```
