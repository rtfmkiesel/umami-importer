<h1 align="center">umami-importer</h1>

<p align="center">
  A simple tool to import web server logs into <a href="https://github.com/umami-software/umami">Umami</a>.<br>Benefit from statistics on JavaScript-free sites, or even if users block your tracker.
</p>

<p align="center">
  <a href="https://github.com/rtfmkiesel/umami-importer">
    <img src="https://img.shields.io/github/stars/rtfmkiesel/umami-importer" alt="Stars">
  </a>
  <a href="https://github.com/rtfmkiesel/umami-importer/blob/main/LICENSE">
    <img src="https://img.shields.io/github/license/rtfmkiesel/umami-importer" alt="LICENSE">
  </a>
  <a href="https://github.com/rtfmkiesel/umami-importer/releases)">
    <img src="https://img.shields.io/github/v/release/rtfmkiesel/umami-importer" alt="Releases">
  </a>
</p>

---

## How does this work?

`umami-importer` will parse web server logs and generate events, which are then sent to an Umami instance. A local [bbolt](https://github.com/etcd-io/bbolt) database is used to store [murmur3](https://github.com/twmb/murmur3) hashes of parsed files and log entries, so no log entry is imported twice.

> [!NOTE]  
> This is very much beta software. Expect all the bugs.

## Usage

+ Download the [prebuilt binaries](https://github.com/rtfmkiesel/umami-importer/releases) and unzip.
+ Copy `config.yaml.dist` and adjust the values to your needs.
  + Adjust `umami.collection_url` to point to your Umami instance.
  + Add a job under `imports[]`. Set the `id`, `base_url`, and specify log paths with `logs.paths`.
+ Run with `umami-importer -c config.yaml`

```
Usage of umami-importer:
  -c, --config string   config path (default "./config.yaml")
  -v, --verbose         enable verbose/debug output
```

### Log Formats

**Only the default log formats of Nginx and Apache2 are supported.**

However, using regex, any log format that has the needed information can be imported. Importing custom logs is done by supplying a regex with named groups.

+ In your config, adjust `imports[].logs.type` to be `custom`.
+ Set `imports[].logs.type_custom_regex` to your custom regex pattern, which defines the following named groups:
  + `remote_addr`
  + `timestamp`
  + `url`
  + `user_agent`
  + `referrer` (optional)
+ Set `imports[].logs.type_custom_timestamp` to the correct [Golang time format string](https://go.dev/src/time/format.go) for the timestamp in your logs.

## Contributing

Improvements in the form of PRs are welcome.

## Legal

This project is not affiliated with [Umami](https://github.com/umami-software/umami).
