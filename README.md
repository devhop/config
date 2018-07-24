## Config Initialization Wrapper
Driven by configuration, vendor locked-in library using [viper](https://github.com/spf13/viper)

Right now it was enforced to use json format.

### Environment Variables

current environment retrieved from env variable named `APP_ENV` with possible value defined below:
```go
EnvDevelopment = "development"
EnvStaging     = "staging"
EnvUat         = "uat"
EnvProduction  = "production"
```
if none or not valid value provided, `production` is use. 

### Config Source

- In `development` environment, source could be retrieved from either one or both, local file or remote.

- In non-development environment, source **MUST** retrieved from remote, local file allowed but will be overridden by remote fetched config

- remote configuration, using consul which the host retrieved from env vars `APP_REMOTE`, if none given sensible default set to `consul:8500`, path is `/service/service_name/config.json`

- local configuration, attempt read from current execution directory, file name is `service_name.config.json`

### Log & Debug

Log level set by default based on environment,

- `development` - `DEBUG`
- `staging` - `INFO`
- `sanity` - `WARNING`
- `production` - `ERROR`

to override it into debug mode, through env vars `APP_DEBUG` set to `true` or `1`

### Usage

`Bootstrap()` to read config

```go
if err := config.Bootstrap("service_name"); err != nil {
    return err
}
```

further, reading config available directly from [viper repo](https://github.com/spf13/viper)

### Running The Test

TODO