# Cherry

``Cherry`` is a ``webchat engine`` wrote in [``Golang``](https://github.com/golang/go). It brings the main functionalities that you need in a webchat server.

With this application you are able to serve a bunch of rooms from your machine just editing a configuration file in a specific language.
So if you are interested you should read the [documentation](https://github.com/rafael-santiago/cherry/blob/master/doc/README.md) to learn how to master it.

![Sample](https://github.com/rafael-santiago/cherry/blob/master/etc/sample.gif)

Until now ``SSL connections`` are unsupported.

## How to build it?

You can use the standard ``go build`` or you can use [Hefesto](https://github.com/rafael-santiago/hefesto).

### Using hefesto

After following all steps to put Hefesto to work on your system just move to ``src`` subdirectory and invoke ``Hefesto`` from
inside. Something like:

```
doctor@TARDIS:~/web/git-hub/rafael-santiago/cherry/src# hefesto
```

If all worked a cherry binary was created under ``../bin/`` and

All done.

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
