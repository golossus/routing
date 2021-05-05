[![Go Report Card](https://goreportcard.com/badge/github.com/golossus/routing)](https://goreportcard.com/report/github.com/golossus/routing) 
[![Build Status](https://travis-ci.com/golossus/routing.svg?branch=master)](https://travis-ci.com/golossus/routing)
[![codecov](https://codecov.io/gh/golossus/routing/branch/master/graph/badge.svg?token=R4RDS0JM4X)](https://codecov.io/gh/golossus/routing)

<p align="center">
    <a href="https://www.golossus.com" target="_blank">
        <img height="100" src="https://avatars2.githubusercontent.com/u/58183018">
    </a>
</p>

[Golossus][1] is a set of reusable **Go modules** to facilitate the creation of 
web applications leveraging Go's standard packages, mainly **net/http**.

The routing module exposes a **Router** or commonly known as **Mux** to map urls to
specific handlers. It provides an enhanced set of features to improve default mux
capabilities:

* Binary tree search for static routes.
* Allows dynamic routes with parameters.
* Parameter constraints matching.
* Http verbs matching.
* Semantic interface.
* More to come...

Installation
------------

The routing package is just a common Go module. You can install it as any other Go module. 
To get more information just review the official [Go blog][2] regarding this topic.

Usage
-----

This is just a quick introduction, view the [GoDoc][3] for details.

Basic usage example:

```go
package main

import (
    "fmt"
    "net/http"
    "log"

    "github.com/golossus/routing"
)

func Index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request) {
    ps := routing.GetURLParameters(r)
    fmt.Fprintf(w, "hello, %s!\n", ps.GetByName("name"))
}

func main() {
    router := routing.NewRouter()
    router.Get("/", Index)
    router.Get("/hello/{name}", Hello)

    log.Fatal(http.ListenAndServe(":8080", router))
}
```

Documentation
-------------

[Official website][1] is still under construction and documentation is not yet finished. Stay
tunned to discover things to come, or [subscribe to our newsletter][4] to get direct notifications. 

Community
---------

* Join our [Slack][5] to meet the community and get support.
* Follow us on [GitHub][6].
* Read our [Code of Conduct][7].

Contributing
------------

Golossus is an Open Source project. The Golossus team wants to enable it to be community-driven 
and open to [contributors][8]. Take a look at [contributing documentation][9].

Security Issues
---------------

If you discover a security vulnerability within Golossus, please follow our
[disclosure procedure][10].

About Us
--------

Golossus development is led by the Golossus Team [Leaders][12] and supported by [contributors][8]. 
It started and supported as a **hackweek** project at [New Work SE][13], we can just say thank you!

[1]: https://www.golossus.com
[2]: https://blog.golang.org/using-go-modules
[3]: http://godoc.org/github.com/golossus/routing
[4]: mailto:subscribe@golossus.com
[5]: https://join.slack.com/t/golossus/shared_invite/zt-db4brnes-M8q1Lw2ouFT5X~gQg69NQQ
[6]: https://github.com/golossus
[7]: ./CODE_OF_CONDUCT.md
[8]: ./CONTRIBUTORS.md
[9]: ./CONTRIBUTING.md
[10]: ./CONTRIBUTING.md#reporting-a-security-issue
[12]: ./CONTRIBUTING.md#leaders
[13]: https://www.new-work.se/en/
