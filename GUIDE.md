# Guide

## Extension

A the core of `ccsm` are extensions, in reality they are just yaml files containing a group of environment variables that you might import into multiple contexts. Take this example repo structure:

```
├── contexts
│   ├── extensions
│       ├── artifactory.yml
│   ├── context1.yml
│   ├── context2.yml
```

Where the content of `artifactory.yml` is:

```yml
skip_deploy: true

extends: []

secrets:
  secret1: secret1
```

And the content of `context1.yml` is:

```yml
context_id: <context_id>

extends:
- artifactory.yml

secrets:
  secret2: mysecret
```

And the content of `context2.yml` is:

```yml
context_id: <context_id>

extends:
- artifactory.yml

secrets:
  secret2: myothersecret
```

Would render the following contexts in circleci:

```txt
context1:
  environment_variables:
    secret1: secret1
    secret2: mysecret
context2:
  environment_variables:
    secret1: secret1
    secret2: myothersecret
```

As you can see, by utilising extensions you now have a single source of truth for `secret1` and a single update would update both contexts.

**Note**: Extensions are loaded in order, so one extension can override another, and secrets in a context can also override anything set in any extensions.

## Plan

`ccsm` contains functionality to show a plan of your changes for an indivual context.

For example:

```sh
ccsm plan --context context1.yml
```

This will run all extensions and give you a single view of what secrets are going to be configured.

## Dry run

For your safety and peice of mind `ccsm` also provides `dry-run` functionality. This will connext to circleci and display the actual changes that it is going to make.

**Note**: Any secrets that currently exist and are in your context yaml files will always show as `updated` this is because the circleci API returns a `***xzy` version of your secret that makes comparison impossible.

For example:

```sh
ccsm --api-token=<circleci_api_token> dry-run
```

Would output:

```sh
Context 'context1':
  New secrets:
    secret2
  Updating secrets:
    secret1

Context 'context2':
  New secrets:
    secret2
  Updating secrets:
    secret1
  Deleting secrets:
    secret3
```

This gives us a clear view for each context about what is new, what is still configured and what would be removed. If you have pipelines your context configuration then it would be recommended to run a `dry-run` before commiting to double check you are making the desired changes.

## Projects

Since creating this tool we have found some scenarios where managing environment variables on a context doesn't work. For example using them to access a private docker registery for an executor. In these scenarioes the best option is to configure you variables on a project.

Just because the variables on configured at a project scope doesn't mean we can't still follow good practice.

Example file:

```yml
extends:
- artifactory.yml

project_slug: github/myorg/myproject

skip_deploy: false

secrets:
  COMMON_SECRET: secret
```

As you can see it follows a similar stucture to we had with contexts, just replacing `context_id` with `project_slug`.

### Reference config structure

```
├── contexts
│   ├── extensions
│       ├── artifactory.yml
│   ├── context1.yml
│   ├── context2.yml
├── projects
│   ├── projects1.yml
```

As we would still consider contexts the "first class citizen" of the tool, we would recommend keeping your `extensions` together in one place.

It is for this reason that project commands get an additional optional flag.

```sh
ccsm dry-run-projects --extensions contexts/extensions
```

Using this flag we can specify that we use the extensions from the contexts folder even when configuring projects to make sure we don't duplicate config.

**Note**: If you specify the `extensions` flag you can still extend using files from within the `projects` repo, this means one project file can extend another as well as using the base extensions.
