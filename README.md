# Deployment event parser

Parses deployment events from files.

## Usage

```shell
depevent [--dir <path>] [--env <environment>] [--output-format <format>]
```

For example:

```shell
depevent --dir /path/to/deployment/events --output-format text

2023-05-13T11:47:52Z staging someorg/frontend-app 835048870362e4294023e823334a384e87f51e2a
2023-05-14T11:47:52Z build someorg/frontend-infra 835048870362e4294023e823334a384e87f51e2a
2023-05-15T11:47:52Z production someorg/frontend-ui fc7e86cb34b616ee89787ff52b68d054a75f6e77
```

Filtering by environment:

```shell
depevent --dir /path/to/deployment/events --env production --output-format text

2023-05-15T11:47:52Z production someorg/backend-api a80c1d6ee3c82a4ca4496ad10250547323a7cd85
2023-05-16T11:47:52Z production someorg/backend-api 830bac7de72d0c83bfaeb9285aea75a78c6f4f31
2023-05-17T11:47:52Z production someorg/backend-api 24deddf5cbd8b4cdf64218a1d15ea6c66786e578
```

The `dir` is searched recursively. Files are assumed to be JSON files.

## Implementation notes

### Input format

The parser expects the following format for the files:

```json
{
  "requestPayload":  {
    "time": "2023-05-15T11:47:52Z"
  },
  "responsePayload": {
    "message": {
      "deploymentStatus": "SUCCEEDED",
      "gitCommitSha":  "abcd1234",
      "gitRepository": "some-git-repo",
      "stackTemplateParameters": [
        {
          "ParameterKey": "Environment",
          "ParameterValue": "production"
        }
      ]
    }
  }
}
```

Other fields can be present and will be ignored.

> **Warning**
> Some files may contain multiple events. The parser will attempt to parse all events in a file.
> 
> See `main_test.go` for an example.

### Output format

By default, the parser outputs in JSON format.

The parser outputs the following format:

```json
{
  "eventTime": "2023-05-15T11:47:52Z",
  "environment": "production",
  "repoName": "someorg/frontend-app",
  "commit": "bc95cfc269ea482cef609748c3b788c68534bfe9"
}
```

The output format can be changed to text by using the `--output-format text` flag.

``` shell
$ depevent --output-format text

2023-05-15T11:47:52Z production someorg/frontend-app bc95cfc269ea482cef609748c3b788c68534bfe9
```
