# legimi-go

Simple, alternative downloader of [Legimi](https://www.legimi.pl/) ebooks written in Go.

Basically, a rewrite of [previous downloader in Lua](https://github.com/tp86/legimi/).

It is completely unofficial, I am not affiliated with Legimi in any way.

> [!NOTE]
> This is still work in progress, however, it is already usable.

You can find more information about how it came about in [Background](#background).

## Installation

Simply download archive from Releases section, unpack and make it executable: `$ chmod +x legimi-go`.
You can add installation directory to your `PATH` variable to be able to run it from anywhere, of course.

Alternatively, if you have Go installed, you can install it using `go install` command:

```shell
$ go install github.com/tp86/legimi-go@<version>
```

`<version>` can be specific version tag from releases or `latest` to get code from `main` branch.
Note that `main` branch may contain unfinished features.
I'm doing my best to commit only working code, though.

## Usage

To view usage, invoke:
```shell
$ legimi-go --help
```

### Options

All command line switches are optional.

-   `--config path`

    Path to configuration file. Default value is `"$HOME/.config/legimi-go/config.ini"`.
    Configuration file contains your credentials and Kindle Id as assigned by Legimi service.
    It will be automatically created (with missing directories) on first command run, so generally you don't need to modify it by hand.
    If you don't want to store your login and password in file, you can provide credentials in command line (see `--login` and `--password` switches).

> [!TIP]
> You can create many configuration files so you can easily switch between multiple accounts.

-   `--login login`

    Your Legimi login.
    If you don't provide login from command line, it will be read from configuration file.
    If it is missing in configuration file as well, you will be asked to provide it during command execution.
    It will be then stored in configuration file, so you don't have to repeat it during future command runs.

-   `--password password`

    Your Legimi password.
    Same logic as for login applies.
    Note that login and password are stored in configuration file as plain text.

> [!NOTE]
> You can give switches with one (`-config`) or two dashes (`--config`).

### Commands

Available commands are:

-   `list`

    List books currently on your Legimi shelf.

-   `download <id> ...`

    Download book(s) given their id(s). Book id can be obtained by listing books (first value in book entry line).

Providing command is mandatory, there is no default command.

On the first command invocation, you will be prompted to provide credentials (if not given via command line switches, see [above](#options))
and Kindle Serial Number (Settings -> Device Options -> Device Info in Kindle).
Legimi Kindle Id will be automatically queried and stored in configuration file for future use.

### Basic usage scenario

1.  List books on your shelf

    ```shell
    $ legimi-go list
    ```

    If you're running script for the first time, or passing configuration file that does not have credentials yet, you will be prompted for login and password.

2.  Download selected book

    ```shell
    $ legimi-go download <book-id>
    ```

    Downloaded book will be saved to `<book-id>.mobi` file in current working directory.

> [!TIP]
> You can `cd` into mounted Kindle's `documents` directory before downloading to avoid copying files.

3.  Copy files to your Kindle's `documents` directory

    This is optional if you have downloaded books while in `documents` directory.

## Limitations

Obviously, only subset of functionality of official Legimi app is supported.

Most error responses are not recognized / handled yet. This should improve in the future.

Script does not provide a way to create account or register device unknown to the Legimi service. You should use official app for this.

There is no information about limit of book downloads in subscription period.
If you are trying to download more books than you limit, Legimi service will block downloads.

## Background

Official Legimi app does not support Linux.
I wanted to be able to download ebooks from Linux without the need to switch between OSes (or even have Windows installed).

[First version](https://github.com/tp86/legimi/) of downloader was created in Lua.
It worked fine, but proved to be hard to install sometimes (dependencies installed to different paths).
It is also hard to maintain and it causes issues with Legimi protocol updates.
Therefore, I decided to rewrite script in Go. I hope it will be easier to install, use and maintain.

I extracted the logic of downloading books based on traffic exchanged between official Legimi application and service.
As such, there are certainly missing pieces and features. Also, most error responses are not supported.
See [Limitations](#limitations) for more missing features.

## TODO

- [ ] handle more error responses from Legimi service
- [ ] refactors
- [ ] documentation update
- [ ] increase test coverage
- [ ] reuse session id if possible
- [ ] parallel downloader
- [ ] handle occassional EOF while getting list of books
- [ ] option to specify directory for downloaded books
