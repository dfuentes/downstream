### Downstream

`downstream` is a tool for managing and testing downstream node dependencies when you make changes upstream.

## Installing

Grab the latest release from [releases](https://github.com/dfuentes/downstream/releases/), decompress it, and place it somewhere in your path.

## Usage

All commands must be run from the root of your node module (where your package.json resides), and assumes that you have run npm install.

# list

List will search the parent directory of your node module for other node modules that are downstream dependencies of yours.  You can use the -p flag to only show non-dev dependencies.

```bash
$ downstream list
downstream-module@0.6.4 depends on ^0.3.0 of upstream
```

# build

Build will build the given downstream dependencies.  If none are specified, then it will build all.  This command can be very slow.  Built downstreams are placed in a ".downstream" directory in your modules directory.  Each downstream is given a copy of your module in its current state.

Running build again on a dependency will reinstall your module into it.

```bash
$ downstream build mod1 mod2
2014/11/26 10:13:44 Cloning into git@github.com:path/to/mod1...
2014/11/26 10:13:44 Npm installing...
2014/11/26 10:14:44 Cloning into git@github.com:path/to/mod2...
2014/11/26 10:14:44 Npm installing...
```

# test

Test will run "make test" for the downstreams specified.  If none are specified, it will test all.  Test takes a verbose flag (-v) which will show all stdout and stderr output, otherwise you will only get the exit code from your tests.

```bash
$ downstream test mod1
2014/11/26 10:20:15 Running tests for mod1...
2014/11/26 10:21:41 Tests failed: exit status 2
```
