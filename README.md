boltbrowser
===========

A CLI Browser for BoltDB Files

![Image of About Screen](http://bullercodeworks.com/boltbrowser/ss2.png)

![Image of Main Browser](http://bullercodeworks.com/boltbrowser/ss1.png)

Installing
----------

Install in the standard way:

```sh
go get github.com/br0xen/boltbrowser
```

Then you'll have `boltbrowser` in your path.

Pre-built Binaries
------------------
Here are pre-built binaries:
* [Linux 64-bit](https://bullercodeworks.com/downloads/boltbrowser/boltbrowser.linux64)
* [Linux 32-bit](https://bullercodeworks.com/downloads/boltbrowser/boltbrowser.linux386)
* [Linux Arm](https://bullercodeworks.com/downloads/boltbrowser/boltbrowser.linuxarm)
* [Windows 64-bit](https://bullercodeworks.com/downloads/boltbrowser/boltbrowser.win64.exe)
* [Windows 32-bit](https://bullercodeworks.com/downloads/boltbrowser/boltbrowser.win386.exe)
* [Mac OS](https://bullercodeworks.com/downloads/boltbrowser/boltbrowser.darwin64)

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
