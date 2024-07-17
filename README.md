# Optech Labs Golang Monorepo Example 
### *by Christopher Hazlett, CTO/Co-founder @ Optech Group*

*I've written a Medium article to cover the whys and wherefores of this Monorepo example. Feel free to start there for greater context. However, this repo, the READMEs and comments should suffice unless you like a bit of a story to go along with your sample repos.*

---

This repo holds a skeleton for an easy to use and deploy, containerized multi-service monorepo built with Golang. It is remarkebly simple, employing common tools that a lot of developers will encounter in their day-to-day jobs. Even if you haven't used them extensively, in the past, I hope you can leverage this as a jumping off point for your own work. I'll continue to expand this repo to encompass not just the Software Development Lifecycle, but also gRPC, Event driven architecture utilizing GCP Pub/Sub, DB Migrations, Google Cloud Run Services & Jobs, Front-end for the Distributed Systems Engineer, Testing, etc. 

### So...why take the time to build this at all, you might ask?
And to that I'll say, "Great Question, dear reader." We started a project at the beginning of 2024 at Optech Group under our Optech Labs moniker and when I started, I was determined to implement some of the hard won lessons from starting fresh at my last company. So, like any good engineer, I started piecing it all together through weeks and weeks of research, trial and error, and asking smart people I know what they'd do. This repo is the result of those learnings, coalesced into one cohesive repository, in the hopes that you, whoever you may be, will be able to use this as a shortcut. Use it at your discretion and with all of my best wishes.

### The Technology in Play
1) *Mac OS* - I develop on a Macbook Pro, Apple M1 Max Chip. Somethings I take for granted because of that, so take that into consideration when using this repo.
2) *Make* - Originally I started the monorepo with Basil, but it had a steeper learning curve than I was willing to climb (and then manage), so I switched to the tried and true Makefile and Go Modules to make working easy and standard.
3) *Docker Compose* - This is for easy local development. I am probably not leveraging to its full capability, but boy does it make my life easy.
4) *Docker* - Put everything in a container and make deployments work for you, not against you. This really isn't a controversial take, but I started building software 20+ years ago, so it's still pretty awesome to me even though I've been doing it a really long time. 
5) *Github Actions* - This controls all of the CI/CD functions and interfaces with Google Cloud Platform for deployments, etc (which I prefer over AWS)
6) *Golang* - I build in Go, so a lot of the stuff in here will be particularly useful for the Go engineer (or the engineer who wants to get into Go). That said...there are a lot fungible ideas that anyone can leverage in any ecosystem.
7) *Gin* - For now, I'm using the Gin http server library, but will most likely update to the native server in the future.

### So How do I get started?
Allright, let's get you oriented to the folder structure.

---
***Root***
1) *Makefile* - All the shortcuts your heart desires. If I type a command against the repo more than once, I add it to the Makefile.
2) *compose.yml* - This is the Docker Compose file that utilizes the files in dockerfiles and interacts with your local Docker instance. The "services" defined in this file run on your Docker instance.

---
***./services***

The services directory contains each individual deployable service. For Golang services, these folders will also include the individual service's go.mod and go.sum files and are actually run from within those folders. You `cd` into those directories and work with them that way, but I find it's easier to put that all in the Makefile in `root` and shortcut your way to nirvana.

---

***./dockerfiles***

I keep a `local` and `remote` subdirectories under the dockerfiles directory. Locally, I work with `.env` or `.yml` config files that are not committed into the repo for deployment. Because of that, the local docker file does a little bit more work than than the remote version. For instance, it copies local config files and secrets and sets environment variables that the Google Cloud Vault and Cloud Run settings handle. This is by design.

---

***./github/workflows***

As you expect, this is where the Github Actions .yml files live. Each of them are heavily commented to explain what's happening. ***CAVEAT: I haven't put anything in place to stop the deployment of services should a library's tests fail. That's something I need to work on.***

---

***./foundation***

This is the main HTTP Server Library I use in Go services. It works with both gRPC and REST services, contains generic logging, and some helpful middleware. Feel free to use it if you're so inclined, but it's a use at your own risk type of deal. While I've used this or a version of it in production for a long time, I may be replacing the innards to the standard http library in the future. That said, it works quite well.


Then, make sure you have following installed on your machine. I use [Homebrew](https://brew.sh) on my Mac to install Git, Go, etc.
1) *Docker* - Running 4.32.0
2) *Git* - Running 2.32.0
3) *Golang* - Running 1.22.0

Once you've got those 3 items installed, you should be able to run `Make [FILL THIS IN]` to start up the (nearly) empty services locally. Obviously, they will not deploy to GCP or interface with Github Actions, but you can get it running locally fairly easily.

