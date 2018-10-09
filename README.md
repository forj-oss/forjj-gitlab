# Introduction

This plugin implements gitlab configuration to build an SCM environment for GIT repositories upstream.

It has been implemented as REST API. See ...(TBD) for FORJJ REST API description.

Depending on tasks, the driver will concretely do several things described below:

## Create task

` + "`Create`" + ` will properly configure

The plugin will returns the list of source files managed by FORJJ gitlab plugin, generated in the local ` + "`infra`" + ` repo.

## Update task

Update mainly do update in the local ` + "`infra`" + ` repo and reports file updated to forjj. (The flow must be configured to push to the right place.)

## Maintain task
This action will ensure the SCM server side is properly configured and really update the server:

-  create repositories
