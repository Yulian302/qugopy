# QugoPy
A smart and highly-efficient task queue system that routes and processes priority-based jobs using multiple workers written in Go and Python.
Clients submit high-level tasks (e.g., email sending, report generation), and the system dispatches them to the appropriate runtime environment based on task type. Jobs are prioritized using a **custom min-heap priority queue** implemented in both Go and Python or are sent directly to **Redis**, depending on the chosen `mode`.

# Getting Started
This is an example of how you can set up this project locally.

# License
Distributed under the MIT License. See `LICENSE.txt` for more information.

# Contact
Yulian - [LinkedIn](https://www.linkedin.com/in/ybohomol/) - bohomolyulian3022003@gmail.com

Project Link: [https://github.com/Yulian302/qugopy](https://github.com/Yulian302/qugopy)
