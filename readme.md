[![Build Status](https://travis-ci.com/Antartical/gandalf.svg?branch=master)](https://travis-ci.com/Antartical/gandalf)
[![Coverage Status](https://coveralls.io/repos/github/Antartical/gandalf/badge.svg?branch=master)](https://coveralls.io/github/Antartical/gandalf?branch=master)
[![pre-commit](https://img.shields.io/badge/pre--commit-enabled-brightgreen?logo=pre-commit&logoColor=white)](https://github.com/pre-commit/pre-commit)


<p align="center">
  <img width="250" height="250" src="https://stickker.net/wp-content/uploads/2016/01/you-shall-not-pass.jpg.png">
</p>

## Gandalf in a nutshell

While frodo and sam were running away from Balrog, Gandalf said "You shall not pass" making them to be safe. This service works like an oauth2 server, it will logs user into the system by creating tokens.

## Local development

Gandalf is easy to develop in a local environment by using docker. just type in your terminal `make`
and everything you need will make up by itselt. Please copy the content of `build/env/.env.sample` to
your own *.env* in `build/env/.env`. You can do this by executting:
```cmd
cp ./build/env/.env.sample ./build/env/.env
```

Moreover you can perform the following operations:
 - **make**: setting up the containers
 - **make sh**: attach a console inside gandalf.
 - **make logs**: show gandalf logs
 - **make local.build**: recompiles gandalf image
 - **make tests**: launch tests
 - **make coverage_report**: open coverage report

## Migrations
In order to create new migrations, pleass follow the next steps:
 - **make**: if you have never setting up the containers
 - **make sh**: attach a shell to the `gandalf` container
 - **mgo create <migration_name> sql**: generates a new migration file in the `/migrations` directory

Write in the generated file the migration in SQL, Moreover make sure to run `mgo fix` when you had test your
migration in order to rename it automatically by adding a sequential number

## Configure pre-commit (Python3 required)
pre-commit is a useful tool which checks your files before any commit push preventings fails in early steps.

Install pre-commit is easy:
```
pip install pre-commit
python3 -m pre_commit install
```
