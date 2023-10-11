# whawty.alerts

[![Go Report Card](https://goreportcard.com/badge/github.com/whawty/alerts)](https://goreportcard.com/report/github.com/whawty/alerts)

whawty.alerts is a simple daemon that handles notifications for monitoring alerts from the Prometheus Alertmanager. Notifications
can be sent via eMail as well as SMS. Unless most other solutions whawty.alerts tries to be of use without access to the internet.
This means sending eMails and SMS is not done using the API of some cloud provider but rather using local resources such as GSM modems
connected via USB.

## License

    3-clause BSD

    © 2023 whawty contributors (see AUTHORS file)
