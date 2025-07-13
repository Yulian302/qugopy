<p align="center">
  <img src="https://github.com/user-attachments/assets/0dd167c7-99ba-4297-8305-378a58e083e4" width="200" alt="Description">
</p>

# QugoPy
[![Development Status](https://img.shields.io/badge/status-in_development-orange.svg)](https://github.com/yourusername/yourrepo)
![Go Version](https://img.shields.io/badge/go-1.23+-blue)
![Python](https://img.shields.io/badge/python-3.10%2B-darkgreen)
![License](https://img.shields.io/github/license/Yulian302/qugopy)
![Issues](https://img.shields.io/github/issues/Yulian302/qugopy)


A smart and highly-efficient task queue system that routes and processes priority-based jobs using multiple workers written in Go and Python.
Clients submit high-level tasks (e.g., email sending, report generation), and the system dispatches them to the appropriate runtime environment based on task type. Jobs are prioritized using a **custom min-heap priority queue** implemented in both Go and Python or are sent directly to **Redis**, depending on the chosen `mode`.

<p>&nbsp;</p>

# Features
- Fast task queue built in Go and Python
- CLI + REST API interfaces
- Custom min-heap priority queue
- Redis support for distributed task scheduling
- Autocomplete-powered interactive shell

<p>&nbsp;</p>

# Tasks
Here are a few predefined tasks that can be added by a task scheduler:
| Name       | Payload    |      | Description                                               |
|:----------------: |:------------------|----| ----------------------------------------------------|
| `send_email` | <code>{<br/>client_name,<br/>client_email,<br/>recipient_name,<br/> recipient_email,<br/>subject,<br/>html_content&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<br/>}</code> || Send email using Brevo SMTP to any recipient using your email address. |
| `download_file` | <code>{<br/>url,<br/>filename<br/>}</code> || Download a file from a given URL and store it locally in `/storage` directory as a `filename`. |

<p>&nbsp;</p>

# Requirements
Before building or running the project, ensure the following dependencies are installed:
| Requirement       | Version          | Description                                               |
|-------------------|------------------|-----------------------------------------------------------|
| [Go](https://golang.org/dl/) | 1.23 or higher   | Required to build and run the project          |
| [Python](https://www.python.org/downloads/) | 3.13.5 or higher   | Required to run Python workers|
| [Redis](https://redis.io/)    | Latest stable    | Required only if running in `--mode redis`    |
| Git               | Any modern version | Required to clone the repository                        |
| Unix-like Shell (optional) | — | Bash/Zsh recommended for script execution and CLI usage         |

<p>&nbsp;</p>

# Installation
Follow the steps **below** to set up and build the application.
## Go Project Setup
1. Navigate to your workspace
    ```bash
    cd your-directory
    ```

2. Clone the repository
    ```bash
    git clone https://github.com/Yulian302/qugopy.git
    cd qugopy
    ```

3. Build the project
    Make sure you have Go installed (Go 1.23+ recommended). Then build the binary:
    ```bash
    go build -o qugopy
    ```
    This will produce an executable named `qugopy` in the project root.

## Python Project Setup
Follow these steps to setup Python environment:
### Automated setup (Recommended)
<p align="left"><b>Linux/MacOS</b></p>

```bash
cd processing
chmod +x setup_python_env.sh # Ensure script is executable
./setup_python_env.sh
```
<p align="left"><b>Windows</b></p>

```powershell
cd processing
Set-ExecutionPolicy -Scope CurrentUser RemoteSigned # Allow execution
./setup_python_env.ps1
```

**OR**

### Manual Setup
<p align="left"><b>Linux/MacOS</b></p>

```bash
cd processing
python3 -m venv venv -upgrade-deps
source venv/bin/activate
pip install --upgrade pip setuptools wheel
pip install -r requirements.txt
```

<p align="left"><b>Windows</b></p>

```powershell
cd processing
python -m venv venv --upgrade-deps
.\venv\Scripts\Activate.ps1
python -m pip install --upgrade pip setuptools wheel
python -m pip install -r requirements.txt
```

### Verification
After setup, verify your environment:
```bash
python -c "import sys; print(sys.executable)"  # Should point to venv Python
pytest tests/  # Run basic tests if available
```



# Usage
Once built, you can run the server using the following command:
```bash
./qugopy start --mode <mode> --workers <workers>
```

## Arguments
| Flag       | Description                           |     | Default              |
| :-----     | :------------------------------------ |:--- | :-------------------:|
| `--mode`   | Queue mode to use: `redis` or `local` |     | `local`              |
| `--workers`| Number of workers to spawn            |     | `2`                  |

**Example:**
```bash
./qugopy start --mode redis --workers 5
```
This will start the server in Redis mode with 5 background workers. **The app will automatically start the Interactive Shell session.**

## Task Scheduling
Now it's time to schedule/queue some tasks. You can either use `CLI` or `REST API` mode and even **BOTH** of them.

## Interactive CLI
Comes with **autocomplete** and the **history** of commands.
- Press `TAB` in order to see available commands and options.
- Use the following command to enqueue a task:
```bash
add task --type <task_type> --payload <payload> --priority <n>
```

Example:

```bash
add task --type download_file --payload '{"url":"https://jsonplaceholder.typicode.com/todos/1","filename":"dummy.json"}' --priority 1
```

## REST API
You can also interact with the task scheduler programmatically via HTTP using the REST API.

|Method|Endpoint|Description|
|:------:|:--------:|-----------|
|`GET`|`/test`|Check if the REST API server is running and responsive|
|`POST`|`/tasks`|Enqueue a new task into the system|

The API accepts JSON-formatted task data in the request body.
**Default port: 5000**

Example using curl:
<br/>
Download a file:

```bash
curl -X POST http://localhost:5000/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "type": "download_file",
    "payload": {
      "url": "https://jsonplaceholder.typicode.com/todos/1",
      "filename": "dummy.json"
    },
    "priority": 1
  }'
```

Payload for sending email:

```json
{
  "type": "send_email",
  "payload": {
    "client_name": "Alice",
    "client_email": "alice@example.com",
    "recipient_name": "Bob",
    "recipient_email": "bob@example.com",
    "subject": "Welcome",
    "html_content": "<h1>Hello Bob!</h1>"
  },
  "priority": 2
}

```


<p>&nbsp;</p>

# Contributing
Contributions are welcomed and appreciated! If you'd like to improve this project, fix a bug, or add a new feature, please follow the steps below.

1. Fork the repository

    Click the "Fork" button on GitHub and clone your fork locally:
    ```bash
    git clone https://github.com/YOUR_USERNAME/qugopy.git
    cd qugopy
    ```
2. Create a new branch
    ```bash
    git checkout -b feature/your-feature-name
    ```
3. Make your changes
4. Run the tests
   ```bash
   go test ./...
   ```
    **and**
    ```bash
    pytest -v -s
    ```
    in **Python submodule**

5. Commit and Push
   ```bash
   git add .
   git commit -m "[mode]: Your feature description"
   git push origin feature/your-feature-name
   ```
6. Open a pull request
    - Go to GitHub and create a PR against the main branch.
    - Clearly describe what your PR does and why it’s needed.

# Development
Enter development mode using `air` command. You can also customize `.air.toml` config file located in a root directory for your needs.

# Miscellaneous
Here you can find some information about the project: system diagrams, charts, as well as execution flows.
## High-level architecture
![High-level arch](https://github.com/user-attachments/assets/cded4049-4cde-4027-b443-9ea2c856c567)
<p align="center">Picture 1. High-level system architecture</p>

# License
Distributed under the MIT License. See `LICENSE.txt` for more information.


# Contact
Yulian - [LinkedIn](https://www.linkedin.com/in/ybohomol/) - bohomolyulian3022003@gmail.com

Project Link: [https://github.com/Yulian302/qugopy](https://github.com/Yulian302/qugopy)
