# is-time-odd
Determine if current time has an odd number of minutes

## Why?

This is a take-home assignment for my current job search.

#### Why Go?

Go has been on my learning wishlist for a while. This take-home would've been clunkier in an interpreted language like Javascript, so I took it as an opportunity to dive in and learn Go.

#### Why "is time odd"?

I wanted to keep the core of the app clear and non-distracting: any user can correlate the output against the log timestamps and intuitively know which version is supposed to be the latest.

I wanted a built-in way to demo the auto-update mechanism, without any extra user action required to trigger a new release.

## What is it?

The core app, *is_time_odd*, produces a steady stream of console messages declaring whether the current minute is odd or not. For example: 11:15am has an odd minute (15), and 19:54 has an even minute (54).

The app has a [Rube Goldberg machine](https://en.wikipedia.org/wiki/Rube_Goldberg_machine)-like architecture. Every version of the app has the "is odd" or "is even" message hardcoded.

In order to remain (mostly) accurate, the app continually checks in with an update server to make sure it has the latest version.

Every minute, the update server tags a different release as the _latest_, which then causes the `is_time_odd` app to download that new version and do an in-place upgrade of itself.

## Quick start

1. Find your operating system's file in the _dist_ folder, and unpack the .zip file to your computer.
1. Start the server app.
    - mac or linux: in your terminal, run `./server`
    - windows: double-click `server.exe`
1. Start the main app:
    - mac or linux: in a second terminal, run `./is_time_odd`
    - windows: double-click `is_time_odd.exe`
1. Watch the output for a minute or more. The time details are highlighted in blue. The rest of the output is verbose logging.

In the below example screenshot, you can see the app updating after the minute changes.

![example output](/docs/is_time_odd_demo.png)

## App usage

The `is_time_odd` app accepts some optional parameters:

```
  -server string
        base url for app updates API (default "http://localhost:3008")
  -test
        run a small self-test: print short output then exit.
```

Example:
```
    ./is_time_odd --server=http://localhost:5005
```


## Server usage

The server accepts some optional parameters:

```
  -app-data string
        data directory (default "./server_data")
  -port string
        port to listen on (default "3008")
```

Example:
```
    ./server --port=5005 --app-data=/Volumes/temp/is_time_odd/server/data
```

You can run the app and server on different computers.

> Caveat: to keep the files in `/dist` small, we ommitted downloads from other operating systems. If you want to for example run `is_time_odd` on a windows PC but run the server on linux, you should do a manual build to get the other windows files. (Or you copy them from the windows dist onto your linux server.)

# Developer docs

This main readme is geared towards users of the app. See [docs/developers.md](/docs/developers.md) for more about the internals, discussion of tradeoffs, and future improvements.