# RangeRadar
RangeRadar: A fast and versatile command-line tool to convert IP address ranges into CIDR notation, with support for multiple output formats.

## Summary
RangeRadar is a lightweight and efficient command-line tool for converting IP address ranges to CIDR notation. It is designed to be fast, easy to use, and provide output in multiple formats, making it an essential utility for network administrators, security analysts, and IT professionals.

## Features

    • Convert IP address ranges to CIDR notation
    • Supports both IPv4 and IPv6 address ranges
    • Parallel processing for faster performance
    • Output formats: JSON, CSV, and terminal display
    • Cross-platform compatibility (Linux, macOS, and Windows)
    
## Installation

1. Download or clone the repository to your local machine:

```shell
git clone https://github.com/ozfive/RangeRadar
```

2. Navigate to the repository directory and build the binary:

```shell
cd CIDRizer
go build
```

3. (Optional) Add the compiled binary to your system PATH for easier access.

## Usage

To use RangeRadar, run the following command in your terminal:

```shell
./RangeRadar -range="startIP-endIP" -parallel=true -concurrency=100 -output=json
```

Replace startIP and endIP with the desired IP address range, and customize the options as needed. For detailed usage instructions and available options, refer to the -help flag:

```shell
./RangeRadar -help
```

## Contributing

Contributions to RangeRadar are welcome! Please submit a pull request or open an issue with your suggestions, bug reports, or feature requests.

## License

RangeRagar is licensed under the MIT License. See the LICENSE file for more information.
