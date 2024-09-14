<h1 align="center">Genzai</h1>
<p align="center"><b>The IoT Security Toolkit</b></p>
<p align="center">
<a href="#description">Description</a> • <a href="#features">Features</a> • <a href="#setup-usage">Setup & Usage</a> • <a href="#acknowledgements">Acknowledgements</a> • <a href="#contact">Contact</a><br>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Version-2.0-green">
  <img src="https://img.shields.io/badge/Black%20Hat%20Arsenal-%20Asia%202024-blue">
  <img src="https://img.shields.io/badge/GISEC Armory-%20Dubai%202024-blue">
  <a href="https://www.buymeacoffee.com/umair9747" target="_blank"><img src="https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png" alt="Buy Me A Coffee" style="height: 21px !important;width: 94px !important;" ></a>
</p>

<hr>

## Description

Genzai is a powerful IoT security toolkit designed to identify and scan IoT devices for vulnerabilities and default credentials. With its latest update, Genzai now offers an API-based approach for more flexible and scalable security assessments.

## Features

- Identify IoT devices and dashboards
- Scan for default vendor-specific passwords
- Detect known vulnerabilities in IoT devices
- Support for IP ranges and CIDR notation
- Concurrent scanning with customizable worker pool
- RESTful API for easy integration with other tools
- Verbose logging option for detailed scan information
- Output results in JSON format

## Setup & Usage

### Setup

1. Ensure you have [Go](https://go.dev/dl/) installed on your system.
2. Clone the repository:
   ```
   git clone https://github.com/umair9747/Genzai.git
   ```
3. Navigate to the project directory and build the binary:
   ```
   cd Genzai
   go build -o genzai-api ./api
   ```

This command tells Go to compile the code in the `./api` directory and output the binary as `genzai-api`.


### Usage

#### Running the API Server

Start the Genzai API server with the following command:

```
./genzai-api -workers 10 -timeout 30s -verbose
```

Options:
- `-workers`: Number of concurrent workers (default: 10)
- `-timeout`: Timeout for each scan in seconds (default: 30s)
- `-verbose`: Enable verbose logging

#### Making API Requests

To scan targets, send a POST request to the `/scan` endpoint:

```
curl -X POST http://localhost:8080/scan \
     -H "Content-Type: application/json" \
     -d '{
       "targets": [
         "192.168.1.0/24",
         "10.0.0.1-10.0.0.10",
         "example.com:8080",
         "https://iot.device.com"
       ]
     }'
```

#### Output Format

The API returns results in the following JSON format:

```json
{
  "results": [
    {
      "target": "http://example.com",
      "iot_identified": "Example IoT Device",
      "category": "Smart Home",
      "issues": [
        {
          "issue_title": "Default Password Detected",
          "url": "http://example.com/login",
          "additional_context": "Default credentials: admin:password"
        }
      ]
    }
  ],
  "targets": ["http://example.com"],
  "total_scanned": 1,
  "time_elapsed": "5.234s",
  "errors": []
}
```

## Acknowledgements

- Original author: Umair Nehri (0x9747)
- Major contributor: [rumble773](https://github.com/rumble773) - API implementation and concurrent scanning

## Contact

For questions, suggestions, or issues, please contact:
- Umair Nehri: [Twitter](https://twitter.com/0x9747)
- rumble773: [GitHub](https://github.com/rumble773)