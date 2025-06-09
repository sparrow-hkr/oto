# Oto

**Oto** is an automated tool for extracting endpoints, sensitive data, and secrets from HTML and JavaScript source code.  
It is designed for security researchers, bug bounty hunters, and developers to discover hidden API endpoints, credentials, and other valuable information in web applications for reconnaissance and vulnerability assessment.

---

## Features

- 🚀 Automated discovery of endpoints and sensitive data in HTML/JS
- 🔍 Extracts API endpoints, paths, secrets, and sensitive keywords
- 🗂 Supports single domains or bulk domain lists
- ⚡ Fast, concurrent processing
- 📝 Outputs results in JSON (to file or stdout)
- 🎛 Customizable extraction types (`endpoint`, `path`, `info`, `critical`, `sensitive`)
- 🐧 Easy to use on Linux and in CI pipelines

---

## Installation

### 1. **Clone the repository**
```sh
git clone https://github.com/sparrow-hkr/oto.git
cd oto
```

### 2. **Build the project**
```sh
go build -o oto
```

---

## Usage

### **Basic Usage**
Scan a single domain:
```sh
./oto endpoint --domain example.com
```

Scan a list of domains from a file:
```sh
./oto endpoint --list domains.txt
```

### **Options**
- `-t, --result-types`: Types of results to extract (`endpoint`, `path`, `info`, `critical`, `sensitive`)
- `-c, --concurrency`: Number of concurrent threads (default: 5)
- `-T, --timeout`: Timeout for HTTP requests (default: 5s)
- `-v, --verbose`: Enable verbose output
- `-D, --debug`: Enable debug output
- `-o, --output`: Output file to save results (default: stdout)

### **Example**
```sh
./oto endpoint --domain example.com -t endpoint,path,info --output results.json
```

---

## Contributing

Pull requests are welcome! For major changes, please open an issue first to discuss what you would like to change.

---

## License

MIT License

---

## Author

Chandra (chandra@gmeil.com)