# Monot

A simple, language agnostic Mono Repo manager.

__Project status: heavily WIP__

## Project Aims

1. Be simple, super simple. Having tried and failed to use Earthly and Bazel a few times. I realised nothing simple and language agnostic really existed (please let me know if there is...).
2. Be language agnostic. There's a few more tools which are heavily language specific. For example Learna, or Turborepo - these are amazing tools, if you're using lots of JavaScript. Monot should be language agnostic, focussing on managing project state, consistency, and executing commands across multiple services. That should be it. If this feels like a glorified Makefile, then that's exactly the aim.
3. Be modular/extendable. We should ensure this tool exposes some kind of way of exposing the functionality for developers to extend and build on-top of.

## Simple Example

Let's say you have a mono repo with two services:

```
services/service-a
services/service-b
```

In each service, include a service config file, named `monot-service.yaml`.

Which will contain something like:

```yaml
name: "service-a"
tasks:
  run-local:
    commands:
      - go run main.go
```

This creates a service called 'service-a' and defines a task called 'run-local'. A 'task' is a named action to run for each service. For example, run-local, build etc. These tasks should share a common name across all services.

In our project root, we need to create a 'manifest' config file, named `monot-manifest.yaml`.

It will look something like this:

```yaml
services:
  - "./example/service-a"
  - "./example/service-b"
```

This tells Monot where to look for each service.

## Run Monot

Now we need to initialise the repo to use Monot. We do this by running the following command:

```bash
$ monot init
```

This command will create a new sub-directory in the root of your project named `.monot` - this will be used to store some state/data. Including a `cache` file, which contains all of your parsed service config files etc in a single data structure. You shouldn't need to worry/care about this, just remember to add it to your `.gitignore`.

Now let's say we want to run `run-local` for all of our services. Here we go!

```bash
$ monot run run-local
```

Fin.

## Future Plans
- Better log output
- Some actual tests...
- Some sort of Docker networking integration so you can stitch services together locally
- Shared services - for example, running an instance of Postgres and Redis for all of the services, which could be defined as 'dependencies' in the `monot-manifest` file
- State management - as in, watch for services having been changed and storing that information - so that on CI/CD we can tell the pipeline which services have changed, and therefor need to be re-deployed
- Environment variable/config management