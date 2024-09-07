# BlueNews Full RSS
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/ForAzens/bluesnews-full-rss)
![GitHub License](https://img.shields.io/github/license/ForAzens/bluesnews-full-rss)

A simple HTTP Server that exposes a RSS Feed of [BluesNews](https://www.bluesnews.com "BluesNews") by days instead of stories.

## Table of Contents
- [Motivation](#motivation)
- [Installation](#Installation)
	- [Docker](#Docker)
	- [Building from source](#building-from-source)
- [Usage](#Usage)
	- [Retrieving articles](#retrieving-articles)
		- [Docker usage](#docker-usage)
- [Roadmap](#roadmap)
- [License](#license)

## Motivation
There are two big motivations for building this project:
- **Learning Golang**: I've been toying with the language for a while, but never found out a problem where I could use it.
- **"Sane" RSS feed from BlueNews**: Probably this is only an issue for me, but I couldn't stand how the different stories polluted my RSS feed.

**A fair warning**: This project is a sort of a learning ground for me, so expect things not to work smoothly, at least until I start with the versioning and releasing.

## Installation
At the moment, the only way to use this project is by cloning this repository:
```bash
git clone https://github.com/ForAzens/bluesnews-full-rss.git
```
You can start the server using Docker or manually building the source.

### Docker
The simplest way of using it. You'll need to [install Docker](https://docs.docker.com/engine/install/) and, optionally, [Docker Compose](https://docs.docker.com/compose/install/).

Then, in a terminal of the newly cloned repository:
```bash
docker-compose up -d
```
Done! You should be able to see the RSS Feed in http://localhost:8080.

### Building from source
You will need to have [Go](https://go.dev) installed in your system.

Then, in a terminal open in the cloned repository:
```bash
go mod download # Download dependencies
go build -o bluesnews-rss # Build a binary executable named bluesnews-rss
```

Now you have an executable in the same folder as the repository. The server can be started with:
```bash
./bluenews-rss --mode=serve
```

## Usage
### Retrieving articles
You need to retrieve articles in some way, manually or with some sort of automation. There's a built-in command to do this:
```bash
./bluenews-rss --mode=fetch --lastDays=7
```
#### Docker Usage
If you use the Docker method of installation, it's already automated via a cronjob.
Either way, if you want to trigger the fetch manually:
```bash
docker-compose exec bluesnews-rss /app/bluesnews-full-rss --mode=fetch --lastDays=7
```

## Roadmap
- [ ] Implement some testing to fix issues.
- [ ] Publish the Docker image in some repository.
- [ ] Publish the binary with Github Releases.
...

**Warning**: This roadmap can be altered at any moment, as this is just a hobby project.

## License
[MIT](https://github.com/ForAzens/bluesnews-full-rss/blob/main/LICENSE.md)
