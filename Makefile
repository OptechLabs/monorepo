# ------ conditional root directory variables ---------
# -----------------------------------------------------
%Gateway: ACTIVEDIR:=services/gateway
%Gateway: PRIMARYPORT:=8080
%Gateway: SECONDARYPORT:=8081
%Gateway: CONTAINERNAME:=gateway
%Gateway: DCKRFILE:=dockerfiles/local/servicegateway.Dockerfile
%Gateway: DCKRTAG:=monorepo.gateway
%Gateway: PSQLUSER:=servicegateway
%Gateway: PSQLDB:=monorepo_gateway_development
%Gateway: PSQLDBTEST:=monorepo_gateway_test
%Gateway: LOCALPSQLURL:='postgres://servicegateway:haveAtIt@localhost:5454/monorepo_gateway_development?sslmode=disable'

%Core: ACTIVEDIR:=services/core
%Core: PRIMARYPORT:=8080
%Core: SECONDARYPORT:=8081
%Core: CONTAINERNAME:=core
%Core: DCKRFILE:=dockerfiles/local/servicecore.Dockerfile
%Core: DCKRTAG:=servicecore
%Core: PSQLUSER:=servicecore
%Core: PSQLDB:=monorepo_core_development
%Core: PSQLDBTEST:=monorepo_core_test
%Core: LOCALPSQLURL:='postgres://servicecore:haveAtIt@localhost:5454/monorepo_core_development?sslmode=disable'

%Helpers: ACTIVEDIR:=helpers
%Foundation: ACTIVEDIR:=foundation
# -------------------------------------------------------

## This will restart the entire docker compose stack.
composeAllAndWatchLogs:
	docker compose -f compose.yml up --build

composeAllNoLogs:
	docker compose -f compose.yml up --build -d

# run this when you first clone the repo to get the databases up and running
# if you run into issues, run `make cleanDBs` and then `make initiateDBs` again.
initiateDBs:
	docker compose -f compose.yml up psql --build -d
	docker exec -it monorepo.psql psql -h localhost -U postgres -c "CREATE USER servicecore WITH PASSWORD 'haveAtIt';"
	docker exec -it monorepo.psql psql -h localhost -U postgres -c "CREATE USER servicegateway WITH PASSWORD 'haveAtIt';"
	docker exec -it monorepo.psql psql -h localhost -U postgres -c "CREATE DATABASE monorepo_core_development;" -c "CREATE DATABASE monorepo_core_test;"
	docker exec -it monorepo.psql psql -h localhost -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE monorepo_core_development TO servicecore; GRANT ALL PRIVILEGES ON DATABASE monorepo_core_test TO servicecore;"
	docker exec -it monorepo.psql psql -h localhost -U postgres -c "CREATE DATABASE monorepo_gateway_development;" -c "CREATE DATABASE monorepo_gateway_test;"
	docker exec -it monorepo.psql psql -h localhost -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE monorepo_gateway_development TO servicegateway; GRANT ALL PRIVILEGES ON DATABASE monorepo_gateway_test TO servicegateway;"

cleanDBs:
	docker compose -f compose.yml down psql


## These are the specific commands for each service should you want to run them individually. Largely, I use
## docker compose for this, but it's nice to have ptions.
tidyGateway: gotidy
venGateway: gotidy govendor
testGateway: gotest
covGateway: goCoverHTML
runGateway: gorun
composeGateway: composespecific
psqlGateway: psqlconnect
statusGateway: 
	curl http://localhost:3000/status

tidyCore: gotidy
venCore: gotidy govendor
testCore: gotest
covCore: goCoverHTML
runCore: gorun
composeCore: composespecific
psqlCore: psqlconnect
statusCore: 
	curl http://localhost:8080/status

tidyHelpers: gotidy
venHelpers: gotidy govendor
testHelpers: gotest
covHelpers: goCoverHTML

tidyFoundation: gotidy
venFoundation: gotidy govendor
testFoundation: gotest

# These are the sub commands that calling the above commands in turn call. 
# By separating out the two, the variables can be set automatically above and keep the commands DRY
# DON'T CALL THIS DIRECTLY
gotidy:
ifneq ($(strip $(ACTIVEDIR)),)
	@echo Use a project specific Make CMD, Please. DIR = $(ACTIVEDIR)
else
	@echo "running go mod tidy on $(ACTIVEDIR)"
	cd $(ACTIVEDIR) && go mod tidy
endif

govendor:
ifneq ($(strip $(ACTIVEDIR)),)
	@echo Use a project specific Make CMD, Please. DIR = $(ACTIVEDIR)
else
	@echo "running go mod tidy on $(ACTIVEDIR)"
	cd $(ACTIVEDIR) && go mod vendor
endif

gotest:
ifneq ($(strip $(ACTIVEDIR)),)
	@echo Use a project specific Make CMD, Please. DIR = $(ACTIVEDIR)
else
	@echo "running go test on $(ACTIVEDIR)"
	cd $(ACTIVEDIR) && go test ./... -coverprofile=cover.out
endif

goCoverHTML:
ifneq ($(strip $(ACTIVEDIR)),)
	@echo Use a project specific Make CMD, Please. DIR = $(ACTIVEDIR)
else
	@echo "opening test coverage on $(ACTIVEDIR)"
	cd $(ACTIVEDIR) && go test ./... -coverprofile=cover.out && go tool cover -html=cover.out
endif

gorun:
ifneq ($(strip $(ACTIVEDIR)),)
	@echo Use a project specific Make CMD, Please. DIR = $(ACTIVEDIR)
else
	@echo "running go mod run on $(ACTIVEDIR)"
	cd $(ACTIVEDIR) && go run main.go
endif

dockerbuild:
ifneq ($(strip $(ACTIVEDIR)),)
	@echo "Use a project specific Make CMD, Please. DIR = $(ACTIVEDIR)"
else
	docker build -t $(DCKRTAG) -f $(DCKRFILE) .
endif

dockerrunlatest:
ifneq ($(strip $(DCKRTAG)),)
	@echo "Use a project specific Make CMD, Please. DOCKER TAG = $(DCKRTAG)""
else
	docker stop $(CONTAINERNAME) || true
	docker rm $(CONTAINERNAME) || true
	@echo "running docker run on $(DCKRTAG)"
ifneq ($(strip $(SECONDARYPORT)),)
	docker run --net otos -p $(PRIMARYPORT):$(PRIMARYPORT) -p $(SECONDARYPORT):$(SECONDARYPORT) --name $(CONTAINERNAME) -d $$(docker images --filter reference=$(DCKRTAG):latest | awk '{print $$3}' | awk 'NR==2')
else
	docker run --net otos -p $(PRIMARYPORT):$(PRIMARYPORT) --name $(CONTAINERNAME) -d $$(docker images --filter reference=$(DCKRTAG):latest | awk '{print $$3}' | awk 'NR==2')
endif
endif

dockerrunjob:
ifneq ($(strip $(DCKRTAG)),)
	@echo "Use a project specific Make CMD, Please. DOCKER TAG = $(DCKRTAG)""
else
	docker stop $(CONTAINERNAME) || true
	docker rm $(CONTAINERNAME) || true
	@echo "running docker run on $(DCKRTAG)"
	docker run --net otos --name $(CONTAINERNAME) $$(docker images --filter reference=$(DCKRTAG):latest | awk '{print $$3}' | awk 'NR==2')
endif

dockerattachshell:
ifneq ($(strip $(DCKRTAG)),)
	@echo "Use a project specific Make CMD, Please. DOCKER TAG = $(DCKRTAG)""
else
	@echo "attaching to $(CONTAINERNAME)"
	docker attach $(CONTAINERNAME)
endif

#Docker Compose
composespecific:
ifneq ($(strip $(CONTAINERNAME)),)
	@echo "Use a project specific Make CMD, Please. CONTAINERNAME = $(CONTAINERNAME)"
else
	docker compose -f compose.yml up $(CONTAINERNAME) --build -d
endif

# Postgres Commands
psqlstart:
ifneq ($(strip $(CONTAINERNAME)),)
	@echo "Use a project specific Make CMD, Please. CONTAINERNAME = $(CONTAINERNAME)"
else
	docker compose -f compose.yml up psql-$(CONTAINERNAME) --build -d
endif
	
psqlstop:
ifneq ($(strip $(CONTAINERNAME)),)
	@echo "Use a project specific Make CMD, Please. CONTAINERNAME = $(CONTAINERNAME)"
else
	docker compose -f compose.yml down psql-$(CONTAINERNAME)
endif

psqlcreateuser:
ifneq ($(strip $(PSQLUSER)),)
	@echo "Use a project specific Make CMD, Please. PSQLUSER = $(PSQLUSER)"
else
	docker exec -it $(CONTAINERNAME).psql psql -U postgres -c "CREATE USER $(PSQLUSER) WITH PASSWORD 'haveAtIt';"
	docker exec -it $(CONTAINERNAME).psql psql -U $(PSQLUSER) -c "GRANT ALL PRIVILEGES ON DATABASE $(PSQLDB) TO $(PSQLUSER); GRANT ALL PRIVILEGES ON DATABASE $(PSQLDBTEST) TO $(PSQLUSER);"
endif

psqlcreatedb:
	docker exec -it $(CONTAINERNAME).psql psql -U $(PSQLUSER) -c "CREATE DATABASE $(PSQLDB);" -c "CREATE DATABASE $(PSQLDBTEST);"

psqldropdb:
	docker exec -it $(CONTAINERNAME).psql psql -U $(PSQLUSER) -c "DROP DATABASE IF EXISTS $(PSQLDB);" -c "DROP DATABASE IF EXISTS $(PSQLDBTEST);"

psqlconnect:
	psql $(LOCALPSQLURL)