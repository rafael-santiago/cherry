# Cherry [![Go Report Card](https://goreportcard.com/badge/github.com/rafael-santiago/cherry)](https://goreportcard.com/report/github.com/rafael-santiago/cherry)

``Cherry`` is a ``webchat engine`` wrote in [``Golang``](https://github.com/golang/go). It brings the main functionalities that you need in a webchat server.

With this application you are able to serve a bunch of rooms from your machine just editing a configuration file in a specific language.
So if you are interested you should read the [documentation](https://github.com/rafael-santiago/cherry/blob/master/doc/README.md) to learn how to master it.

Until now ``SSL connections`` are unsupported.

## How to build it?

You can use the standard ``go build`` or you can use [Hefesto](https://github.com/rafael-santiago/hefesto).

### Using go build

You need to setup your ``GOPATH`` to the project root. Supposing that cherry repo was cloned under ``/home/doctor/web/git-hub/rafael-santiago/cherry``
just add this path to your ``GOPATH``.

Run ``go build`` from inside the ``src`` subdirectory.

### Using hefesto

After following all steps to put Hefesto to work on your system just move to ``src`` subdirectory and invoke ``Hefesto`` from
inside. Something like:

```
doctor@TARDIS:~/web/git-hub/rafael-santiago/cherry/src# hefesto
```

If all worked a cherry binary was created under ``../bin/`` and

All done.

Here you do not need to worry about ``GOPATH`` issues because Hefesto's script handles it for you on each build task that you invoke.

## How to run it?

This application works based on a configuration file (again: [documentation](https://github.com/rafael-santiago/cherry/blob/master/doc/README.md)).

You specify this configuration using the option ``--config``:

```
doctor@TARDIS:~/web/git-hub/rafael-santiago/cherry/bin# ./cherry --config=gallifrey-lounge.cherry

```

Supposing that ``TARDIS`` has the ``IP`` address ``192.30.70.3`` and ``Gallifrey lounge`` opens only one room at the port 1008.
Doctor should access the entrace form served at:

```
http://192.30.70.3:1008/join
```

That's all.
