<p align="center">
  <img src="https://github.com/user-attachments/assets/0dd167c7-99ba-4297-8305-378a58e083e4" width="200" alt="Description">
</p>

# QugoPy

![Go Version](https://img.shields.io/badge/go-1.23+-blue)
![Python](https://img.shields.io/badge/python-3.10%2B-darkgreen)
![License](https://img.shields.io/github/license/Yulian302/qugopy)
![Issues](https://img.shields.io/github/issues/Yulian302/qugopy)


A smart and highly-efficient task queue system that routes and processes priority-based jobs using multiple workers written in Go and Python.
Clients submit high-level tasks (e.g., email sending, report generation), and the system dispatches them to the appropriate runtime environment based on task type. Jobs are prioritized using a **custom min-heap priority queue** implemented in both Go and Python or are sent directly to **Redis**, depending on the chosen `mode`.

# Features
After starting the app via CLI, you can use either `CLI` or `REST API` to queue the tasks.

# Requirements ✅
Before building or running the project, ensure the following dependencies are installed:
| Requirement       | Version          | Description                                               |
|-------------------|------------------|-----------------------------------------------------------|
| [Go](https://golang.org/dl/) | 1.23 or higher   | Required to build and run the project          |
| [Python](https://www.python.org/downloads/) | 3.13.5 or higher   | Required to run Python workers|
| [Redis](https://redis.io/)    | Latest stable    | Required only if running in `--mode redis`    |
| Git               | Any modern version | Required to clone the repository                        |
| Unix-like Shell (optional) | — | Bash/Zsh recommended for script execution and CLI usage         |



# Installation
Follow the steps below to set up and build the project:
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
This will start the server in Redis mode with 5 background workers.

## Task Scheduling
Now it's time to schedule/queue some tasks. You can either use `CLI` or `REST API` mode and even **BOTH** of them.
### Using CLI
TODO

### Using REST API
TODO


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


# License
Distributed under the MIT License. See `LICENSE.txt` for more information.


# Contact
Yulian - [LinkedIn](https://www.linkedin.com/in/ybohomol/) - bohomolyulian3022003@gmail.com

Project Link: [https://github.com/Yulian302/qugopy](https://github.com/Yulian302/qugopy)
