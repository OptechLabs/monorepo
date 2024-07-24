# Optech Labs Golang Monorepo Example 
### *by Christopher Hazlett, CTO/Co-founder @ Optech Group*

*I've written a [Medium Article](https://medium.com/@chrishazlett/how-to-build-a-friendly-monorepo-for-golang-services-92c8a62f3b9a) to cover the whys and wherefores of this Monorepo example. Feel free to start there for greater context. However, this repo, the READMEs and comments should suffice unless you like a bit of a story to go along with your sample repos.*

---

This repo holds a skeleton for an easy to use and deploy, containerized multi-service monorepo built with Golang. It is remarkebly simple, employing common tools that a lot of developers will encounter in their day-to-day jobs. Even if you haven't used them extensively, in the past, I hope you can leverage this as a jumping off point for your own work. I'll continue to expand this repo to encompass not just the Software Development Lifecycle, but also gRPC, Event driven architecture utilizing GCP Pub/Sub, DB Migrations, Google Cloud Run Services & Jobs, Front-end for the Distributed Systems Engineer, Testing, etc. 

### So...why take the time to build this at all, you might ask?
And to that I'll say, "Great Question, dear reader." We started a project at the beginning of 2024 at Optech Group under our Optech Labs moniker and when I started, I was determined to implement some of the hard won lessons from starting fresh at my last company. So, like any good engineer, I started piecing it all together through weeks and weeks of research, trial and error, and asking smart people I know what they'd do. This repo is the result of those learnings, coalesced into one cohesive repository, in the hopes that you, whoever you may be, will be able to use this as a shortcut. Use it at your discretion and with all of my best wishes.

### The Technology in Play
1) *Mac OS* - I develop on a Macbook Pro, Apple M1 Max Chip. I take some things for granted because of that, so take that into consideration when using this repo.
2) *Make* - Originally I started the monorepo with Basil, but it had a steeper learning curve than I was willing to climb (and then manage), so I switched to the tried and true Makefile with Go Modules to make working easy and standard.
3) *Docker Compose* - This is for improved local development. I am probably not leveraging to its full capability, but boy does it make my life easy.
4) *Docker* - Put everything in a container and make deployments work for you, not against you. This really isn't a controversial take, but I started building software 20+ years ago when this type of thing wasn't de rigeur, so it's still pretty awesome to me even though I've been doing it a really long time. 
5) *Github Actions* - This controls all of the CI/CD functions and interfaces with Google Cloud Platform for deployments, etc (which I prefer over AWS).
6) *Golang* - I build in Go, so a lot of the stuff in here will be particularly useful for the Go engineer (or the engineer who wants to get into Go). That said...there are a lot fungible ideas that anyone can leverage in any ecosystem and the infrastructure to support development and deployment will work for any type of language in the services directory, as long as compilation and containerization works.
7) *Gin* - For now, I'm using the Gin http server library, but will most likely update to the newer golang native server in the future.

### So How do I get started?
Allright, let's get you oriented to the folder structure.

---
***Root***
1) *Makefile* - All the shortcuts your heart desires. If I type a command against the repo more than once, I add it to the Makefile.
2) *compose.yml* - This is the Docker Compose file that utilizes the files in dockerfiles and interacts with your local Docker instance. The "services" defined in this file run on your Docker instance.

---
***./services***

The services directory contains each individual deployable service. For Golang services, these folders will also include the individual service's go.mod and go.sum files and are actually run from within those folders. You `cd` into those directories and work with them that way, but I find it's easier to put that all in the Makefile in `root` and shortcut your way to nirvana.

***./services/gateway*** - The `gateway` service acts as the outside world's interface to the `core` service containing the lion's share of logic for this very made up system. Because most of my history has been in retail and fulfillment, so the example I'll use is posting a simple ecommerce order to the API. If you've decided to get everything running, issue a post using CURL or Insomnia or Postman (whatever you fancy). You can also use the already made command `Make sendOrderToGateway`

***./services/core*** - The `core` service acts as the place where all the magic happens. Whereas the `gateway` service is meant to access the `core`, the `core` itself is not exposed to the outside world. The power of the monorepo is that you can develop a really effective client library in the `core` folder without having to worry about versions and publishing an all that. Sure, you will need to ensure that you are deploying things separately and not mingling your work between services, but it really just forces good behavior. Making your changes backward compatible, for instance.

*** SPECIAL NOTE ON MONOREPO vs MONOLITH *** - When you're first starting out on a project, the idea of adopting a service oriented architecture can seem a bit anathema to moving fast and reducing complexity, so while this is an example of a monorepo that supports multiple services, the truth is that you can start with just the `core` service and develop everything in one service directory. When it comes time to add something to your overall ecosystem that doesn't belong or should be decoupled from the service, simply add a new  in the service directory, docker files, etc and bob's your uncle, you've got a functioning monorepo with multiple services. Until that time, enjoy a monolith that can flex into services. I always like to say that a service should be as big as it needs to be, neither macro nor micro, but just right.

*** config.json *** - I made a few significant exceptions with this repository, the most noteworthy of which is the commiting of config files to a remote repository. I've also additionally set local db DEV passwords in the makefile, which I normally wouldn't do, instead using variables or uncommited local files. But that makes these types of repos all types of hard to follow, so do the right thing, exclude your `.env`, `config.json`, and passwords from your repos.

*ALSO...I never, not even for a teaching repo like this, commit a password or config that could jeapordize the security of a Github or GCP account. You'll note that the github/workflows directory has secret references in the important parts.*

---

***./dockerfiles***

I keep a `local` and `remote` subdirectory under the dockerfiles directory. Locally, I work with `.env` or `.yml` config files that are not committed into the repo for deployment. Because of that, the local docker file does a little bit more work than than the remote version. For instance, it copies local config files and secrets and sets environment variables that the Google Cloud Vault and Cloud Run settings handle. This is by design.

---

***./github/workflows***

As you expect, this is where the Github Actions .yml files live. Each of them are heavily commented to explain what's happening. ***CAVEAT: I haven't put anything in place to stop the deployment of services should a library's tests fail. That's something I need to work on.***

---

***./foundation***

This is the main HTTP Server Library I use in Go services. It works with both gRPC and REST services, contains generic logging, and some helpful middleware. Feel free to use it if you're so inclined, but it's a use at your own risk type of deal. While I've used this or a version of it in production for a long time, I may be replacing the innards to the standard http library in the future. That said, it works quite well.

---

***./helpers***

The helpers directory contains generic libraries useful throughout the stack. For example, it contains the `config` library which does some basic loading from a `json` file and sets reasonable defaults. It's simple, but I hate coding those types of things more than once.


Then, make sure you have following installed on your machine. I use [Homebrew](https://brew.sh) on my Mac to install Git, Go, etc.
1) *Docker* - Running 4.32.0
2) *Git* - Running 2.32.0
3) *Golang* - Running 1.22.0

Once you've got those 3 items installed, you should be able to run `make composeAllNoLogs` to start up the (nearly) empty services locally. Obviously, they will not deploy to GCP or interface with Github Actions, but you can get it running locally fairly easily.

# To Test it #
You can see that the servers are running by executing the following:
1) `make statusGateway` - Checks the status of the Gateway service.
2) `make statusCore` - Checks the status of the Core service.
3) `make psqlGateway` - Connects to the Gateway Dev DB.
4) `make psqlCore` - Connects to the Core Dev DB.

*Also, you can just look in the logs in the Docker Desktop application.*

