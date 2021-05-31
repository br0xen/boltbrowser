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
* [Linux 64-bit](https://git.bullercodeworks.com/attachments/29367198-79f9-4fb3-9a66-f71a0e605006)
* [Linux 32-bit](https://git.bullercodeworks.com/attachments/ba8b9116-a013-431d-b266-66dfa16f2a88)
* [Linux Arm](https://git.bullercodeworks.com/attachments/795108a6-79e3-4723-b9a8-83803bc27f20)
* [Windows 64-bit](https://git.bullercodeworks.com/attachments/649993d9-bf2c-46ea-98dd-1994f1c73020)
* [Windows 32-bit](https://git.bullercodeworks.com/attachments/c1662c27-524c-465a-8739-b021fb15066b)
* [Mac OS](https://git.bullercodeworks.com/attachments/10270b6f-9316-446d-8ab4-4022142323b3)

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
