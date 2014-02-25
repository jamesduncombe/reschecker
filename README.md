# reschecker

Program to help checking resolutions concurrently.

This expects JSON in the following form:

```json
{
  "config": {
    "callbackUrl": "http://url_callback_address",
    "token": "validation_token_for_user"
    "validation": {
      "width": 2048,
      "height": 2048,
    },
    "images": [
      "http://image_url",
    ]
  }
}
```

- `callbackUrl` a callback address that the results are posted to
- `token` a token used to attribute the images to a user for example
- `validation` is used to specify the minimum `width` and `height` of the images
- `images` is an array of images to validate (these can be on S3 etc)
