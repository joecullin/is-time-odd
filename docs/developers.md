### Caveat

In case you didn't notice my comment in the main readme, I want to be clear that I've never programmed in Go before this project.

I spent day 1 doing the [tour of Go](https://go.dev/tour/list), day 2 doing some scattered learning for the building blocks I knew I'd need. You can see that progress in my [go_sandbox](https://github.com/joecullin/go_sandbox) repo. On days 3 and 4 I put everything together in this app—you can see my progression in the git log.

I'm sure this code has lots of non-standard idioms and choices you'd frown upon if I were a seasoned Go programmer. I skimped on some error handling and I didn't include any automated tests. But I'm pretty proud of my learning curve, I had fun with the challenge, and I really like the brevity and power of the language so far. Thanks for giving me a good motivation to dive into Go!

## Map of the repo

Key directories and files:

- app - code for the main `is-time-odd` app
  - `app.go` - initialization and time-printing
  - `checkUpdates.go` - poll for updates, download new version and replace/restart.
  - `Makefile`, bin, buildtmp - used for creating builds & dists
- server - code and data for server
  - `server_data` dir
    - `appData.json` - metadata for all releases
    - one binary for each release
  - Makefile, bin - used for builds
- `README.md` - top-level user-facing doc
- `docs` - this doc, doc images
- `dist` - compact demo-ready zip file for each platform
- `utils` - small scripts used by Makefiles.

## Server API details

- Authentication and access controls - none yet.
- Error response statuses should be improved.
    - Most errors result in a 404 response. For example, if you pass an invalid version like "!!!!!", we're validating it and rejecting it. But the error handling flow is clunky and we wind up sending a 404 for that case when we should send a 400 with an explanation.
- Shortcuts I'd revisit for larger scale and/or longer maintenance:
    - Ditch appData.json in favor of a database or other store.
    - When serving binary files, don't read the whole file into a byte array. Buffer it in chunks.


### API routes

#### Get info about a release:

`GET /api/releases/<platform>/<version>/info`

Use `latest` to find the latest release, like `GET /api/releases/linux/latest/info`

Use version id to get that release, like `GET /api/releases/linux/1.2/info`

Example:
```
curl -v 'http://localhost:3008/api/releases/windows/1.2/info'


HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Wed, 11 Jun 2025 06:03:12 GMT
Content-Length: 127

{
    "File": "app_even.exe",
    "Md5": "542590e128ae93eeea07a3df0df4227a",
    "Platform": "windows",
    "Tags": [
        "archive",
        "even"
    ],
    "Version": "1.2"
}
```

#### Download a release:

`GET /api/releases/<platform>/<version>`

Path is the same as the "info" route, minus the `/info` at the end.

You can specify an exact version, or specify `latest` as the version.

Example:
```
curl -v 'http://localhost:3008/api/releases/windows/1.2' > /dev/null


HTTP/1.1 200 OK
Content-Length: 8962048
Content-Type: application/octet-stream
Date: Wed, 11 Jun 2025 06:51:21 GMT

...
```

### Makefiles

The one in the app directory is a little long, but hopefully simple to follow.

If this were a real work project, I'd put more thought into organizing it and using a smaller set of tools to make it more maintainable. The current state reflects me treating it as more of a peripheral task, just a tool to help me get the main task (repetitive builds for manual testing) done efficiently.

## TODO/Future

Misc. other things that I didn't touch on above.

- Validation of downloaded binaries.
    - I'm generating an md5 of each build and it's included in the info response, but it's ignored by the client so far.
    - Ideally, we should do as much verification as possible before taking the plunge to replace ourselves, including:
        - File size: reject if extremely small or large.
        - Run downloaded binary with `--test` param and validate the exitcode and output.
    - Symlinks - untested. I suspect if you created a symlink to the app, the replacing would break. I just didn't have time to dig into detecting and/or preventing that.

- Both apps could do more validation and checks. For example:
    - We're trusting the server urls and ports are valid and functioning. (The client does a decent job recovering from server outages, but it would be nice if it caught obvious problems like invalid an url early.)
    - It would be nice to do some self-checks on startup. For example, the server could examine its own data. (Currently that doesn't happen until it receives a request.)
    - There's no timeout (that I know of) on the client's API calls. I'd like to set an explicit limit, so we can kill the request and show a useful log message.