# go-exec-communicator

Terraform has superpowers to run both remote and local commands with many different
authentication methods and networking constructs through provisioners.

Unfortunately that code is hidden inside `internal` packages and hard to use.

This repository aims to extract those execution communication capabilities out into
a simple, reusable package. Mostly for use in the [Steampipe exec plugin](https://github.com/turbot/steampipe-plugin-exec).
