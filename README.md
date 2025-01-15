# llm-wrap

A basic LLM wrapping learning project aimed as serving as acai travel technical assesment task

## Build

- Go to source directory where main.go resides.
- Run `go install`
- Run `go build`
- Ready to run.

## Run

- Make sure to create config.vendor.json
- Fill it as the example shows
- It is ready to be run by `go run main.go`
- API will be listening from `localhost:11242`

## Example

API has one endpoint `/v1/api/command` which takes POST request with a json body where body schema is like:
<br/>

```json
{
  "prompt": "string"
}
```

It produces server side events named as `status` and `conversation`.
<br/>

`status` is a string and `conversation` is a json object.

```json
"conversation":{
    "context_id": "string",
	"payload":   "string",
	"response":  "string"
}
```
