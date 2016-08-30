# Gledger API

![circleci status](https://circleci.com/gh/gledger/api.svg?style=shield)

Gledger is a HTTP REST api for personal finance management. It provides transactions assigned to accounts and envelopes.

Gledger follows the [Envelope Budget](https://en.wikipedia.org/wiki/Envelope_system) system for managing your money.

It is under heavy developement, so you should not use this program for the main system of managing your money right now.

# Contributing

We welcome contribution! There's an actively maintained [issue list](https://github.com/gledger/api/issues)

Simply fork this repository and create a branch with the change and submit a pull request. Please ensure there are tests for any new features or bug fixes.

# Install

Gledger-api uses a postgresql database backend. There are [sql files](https://github.com/gledger/api-schema) that will create the schema. It was designed to be used with [flyway](https://flywaydb.org/) to manage your schema. Details on installing with flyway are out of scope of this document.

Once your database is migrated, you can run the api program with `make run`. You will need to set the `DATABASE_URL` environment variable pointed at your postgres database.
