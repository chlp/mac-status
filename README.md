# mac-status
A lightweight Go application that runs an HTTP server and serves basic macOS system information (e.g., temperature, uptime) as JSON

## Features

- Reports:
  - **Maximum CPU temperature**
  - **System uptime**
  - **App uptime**

## Example output

```json
[
  {
    "name": "max_temperature",
    "value": 66,
    "unit": "celsius",
    "status": "ok"
  },
  {
    "name": "system_uptime",
    "value": 10485,
    "unit": "seconds",
    "status": "none"
  },
  {
    "name": "app_uptime",
    "value": 142,
    "unit": "seconds",
    "status": "none"
  }
]
```

## License

MIT License

## Author

[Aleksei Rytikov](https://github.com/chlp)