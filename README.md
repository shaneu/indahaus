<h3 align="center">In Da Haus</h3>

---

<p align="center"> See if your favorite IP addresses are In Da Haus!!!
    <br> 
</p>

## 📝 Table of Contents

- [About](#about)
- [Getting Started](#getting_started)
- [A Note About The Data](#data)
- [Deployment](#deployment)
- [Usage](#usage)
- [k8s](#k8s)
- [Tests](#tests)
- [Built Using](#built_using)
- [TODO](#todo)
- [Regrets](#regrets)
- [Author](#author)
- [Acknowledgments](#acknowledgement)

## 🧐 About <a name = "about"></a>

indahaus is the exciting new graphql api that stores a users DNS blacklist queries for fast retrieval


## 🏁 Getting Started <a name = "getting_started"></a>

### Prerequisites

You'll need go 1.16+ installed to run and test locally.

For building and running the docker container you'll need [docker](https://docs.docker.com/get-docker/).

For running in a local kubernetes cluster you can use my favorite, [kind](https://kind.sigs.k8s.io/), install instructions [here](https://kind.sigs.k8s.io/docs/user/quick-start/#installation).

##### NOTE: not tested with minikube. I have scars and can't go back. I imagine it works?

Then you'll need [kubectl](https://kubernetes.io/docs/tasks/tools/).

Finally you'll need [helm](https://helm.sh/docs/intro/install/), whew, got all that?

To run locally:
```bash
make migrate
make run
```

To create and run a docker container:
```bash
make docker-build
make docker-run
```

For k8s instructions see below.

There are some other handy command to aid in local dev, like 

```bash
# tidy up the go.mod file
make tidy 

# reset the db
make resetdb

# run the tests 
make test
```


### Installing

Clone the repo:

```bash
git clone https://github.com/shaneu/indahaus.git
cd indahaus
```

Download the dependencies:

```bash
go mod download
```

Migrate the db
```bash
make migrate
```

Build (or run) the binary:

```bash
make build
# or
make run
```


and you're off the races.

```
> AUTH_PASSWORD=***** AUTH_USERNAME=***** PORT=8080 make run
go run cmd/api/main.go
API: 2021/05/22 16:02:02.985390 main.go:79: main: Application initializing: version "develop"
API: 2021/05/22 16:02:02.985514 main.go:84: main: Initializing database support
API: 2021/05/22 16:02:02.985711 main.go:106: main: Debug Listening  :4000/debug/vars
API: 2021/05/22 16:02:02.985979 main.go:127: main: Api listening on :8080
```


## 💾 A Note About The Data <a name = "data"></a>

A single IP address can have multiple results, for example the IP 103.35.191.44 has three results, 127.0.0.3, 127.0.0.4, 127.0.0.2.
I have decided to store all three as a comma separated list so the user can see any codes that may apply to the IP address they enqueued.
Conversely, when an address has no codes the user will receive `null`. 

Codes that represent an error from the spamhaus API, their equivalent of a 400, will not be stored. In other words, if the code received is 127.255.255.255 
meaning an excessive number of queries, that information is useful to us as the developers, but not the user so those codes won't be stored.

A possible future state of the app would be to have the response_code field return a slice of items that might contain the code and a
human readable message.

## 🔧 Running in k8s locally <a name = "k8s"></a>

If you have all the perquisites installed you can run:

```bash
make up
```

which will build the docker image, bring up the kind cluster, load the image into kind and install the helm chart. To interact locally run

```bash
# terminal 1
make port-forward-debug
```

```bash
# terminal 2
make port-forward-api
```

which will forward both the debug vars port 4000 and the api port 8080.

To view debug/metrics info you can use your favorite HTTP API tool such as postman, curl, or the good ole browser and visit http://localhost:4000/debug/vars

To interact with the graphql api you'll need pass a basic auth header.

To bring everything down, including the kind cluster, run:

```bash
make down
```

Some useful dev commands:

`make update-api` will rebuild container, load it into kind, and redeploy the pod. This is handy when you've made a change
and want to test it.


## 🔧 Running the tests <a name = "tests"></a>

```bash
make test
```

Runs the test suites

## 🎈 Design <a name="usage"></a>

This application takes a layered architecture approach.

At the bottom layer are the packages in the `pkg` dir. These are the kinds of packages that could be ripped out and moved
into any project that needs them. They contain no business logic or cross cutting concerns like logging and expose a simple api. 

The next layer is our business logic/data layer, in the `internal` dir. Packages in this layer can require packages from the `pkg` dir, but not visa versa. 

Finally we have our application layer where our graphql and our rest endpoints live. The `cmd` dir is where our binaries live, our our case
we have two, the main app and a thin admin app that does some useful things like migrating our database. Ideally the `graph` dir should be nested in
`cmd/api` but gqlgen seems happier when it isn't.


## ⛏️ Built Using <a name = "built_using"></a>

- [errors](https://github.com/pkg/errors) - Provides a nice way to capture errors along with their context
- [viper](https://github.com/spf13/viper) - Elegant project configuration that allows us to define sane default
parameters and override them with env vars
- [go-sqlite3](https://github.com/mattn/go-sqlite3) - The only sqlite driver that is included in and has passed the go driver compatibility test suite [doc](https://github.com/golang/go/wiki/SQLDrivers)
- [sqlx](https://github.com/jmoiron/sqlx) - provides some nice extensions beyond what is delivered by the built in sql library while still being fully interface compliant. Will be useful if we decide to change from sqlite to something like postgres
- [echo](https://github.com/labstack/echo) - A minimalist web framework we're using for it's routing and its auth and panic recover middlewares 
- [uuid](https://github.com/google/uuid) - Generates our trace IDs and db IDs
- [gqlgen](https://github.com/99designs/gqlgen) - Takes a lot of the boiler plate out of creating a graphql api all while providing a high level of type safety
- [go-cmp](https://github.com/google/go-cmp) - For doing easy comparisons between fields in our tests

## ✍️ TODO <a name = "todo"></a>

- Add support for tracing and metrics collection.
- Integration tests: go has amazing built in support for running integration tests using the httptest package
- Install a migration framework to allow us to update, and roll back, or database schema
- Improved error handling: We should create a subset of trusted errors or a custom error to respond to the user with without leaking too much information about our system


## 😩 Regrets <a name = "regrets"></a>

These are things I would have liked to have done/been able to do differently.

Right now requests are being handled by 2 different ServeHTTP methods - the echo framework and the gqlgen framework.
This means that my centralized error handling/reporting is split into two - I have to have one for the graphql requests and one for
the three REST routes. Now, the argument could be made that those 3 REST routes, /liveness, /readiness, and /debug/vars are not the business
deliverable here and so elegant error handling and the like isn't critical, and I could easily be persuaded. What if tomorrow though a client
comes and says I will give you 1 billion dollars if you write a rest endpoint for this service that I can use because Acme Corp doesn't 
use graphql? I mean, a billion dollars, that's a lotta green. I dug through the gqlgen source quite a bit and they don't expose a way to 
have it simply manage the graphql query execution and let the consumer choose the transport, which would have been nice.

## ✍️ Author <a name = "author"></a>

[@shaneu](https://github.com/shaneu)

## 🎉 Acknowledgements <a name = "acknowledgement"></a>

- My dog Chaos
- My cats Judy and Eris
- My wife Nina for being super supportive and taking over a bunch of the chores while I worked on this
