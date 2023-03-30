boltbrowser
===========

A CLI Browser for BoltDB Files

![Image of About Screen](https://git.bullercodeworks.com/brian/boltbrowser/raw/branch/master/build/aboutscreen.png)

![Image of Main Browser](https://git.bullercodeworks.com/brian/boltbrowser/raw/branch/master/build/mainscreen.png)

Installing
----------

Install in the standard way:

```sh
go get github.com/br0xen/boltbrowser
```

Then you'll have `boltbrowser` in your path.

Pre-built Binaries
------------------
Pre-build binaries are available on the [Releases Page](https://github.com/br0xen/boltbrowser/releases).

Usage
-----

Just provide a BoltDB filename to be opened as the first argument on the command line:

```sh
boltbrowser <filename>
```

To see all options that are available, run:

```
boltbrowser --help
```

Troubleshooting
---------------

If you're having trouble with garbled characters being displayed on your screen, you may try a different value for `TERM`.  
People tend to have the best luck with `xterm-256color` or something like that. Play around with it and see if it fixes your problems.
