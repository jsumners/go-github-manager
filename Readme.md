# GitHub Manager

The intention of this tool is to provide a set of commands to query GitHub
repositories and organizations for interesting information or to perform
administrative tasks upon them. It is spawned from the fact that I kept
duplicating and tweaking [a simple script](https://gist.github.com/jsumners/e32f9c3e1e199184d0813720f67596fe)
for each such task that I want to accomplish. I decided to try and collect
all of these tasks into a single application that might be easier to use.

## Running

Currently, the project is in a bit of a hacked together state. I wanted to get
some information ready for the `fastify@5` release, and took the opportunity
to bootstrap this project. So, for now, clone the repo and:

```sh
$ export GHM_AUTH_TOKEN=$(pbpaste) # set GitHub token by pasteboard on macOS
$ go run ./... help
```

### Example

Let's say we want to generate a report of all contributors across an
organization. First, we want to generate a CSV of all the commit references
across the organization that we are interested in:

```sh
$ export GHM_AUTH_TOKEN=$(pbpaste)
$ go run ./... refs recent-releases -o fastify > git-refs.csv
```

We can then review that CSV, make any changes we desire, and:

```sh
$ go run ./... contributors list-all -o fastify -f git-refs.csv > contributors.csv
```
