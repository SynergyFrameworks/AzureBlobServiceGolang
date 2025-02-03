Azure Blob Storage Service with Kafka Integration
Overview
This application provides a robust storage solution with seamless integration of Azure Blob Storage and Apache Kafka. It offers APIs for file and directory operations, event-driven processing, and supports both local and cloud storage. Designed for scalability, the app processes storage-related events in real-time using Kafka consumers.

Features
File Operations:

Upload, read, delete, and list files.
Supports overwriting files and retrieving metadata.
Directory Operations:

Create, delete, and list directories.
Storage Adapters:

Azure Blob Storage: Handles cloud-based storage.
Local Storage: Provides fallback for local file handling.
Event-Driven Processing:

Uses Kafka for storage event handling (e.g., file uploaded, deleted).
Publishes events to the storage-events topic.
Elasticsearch Logging:

Sends logs to Elasticsearch for centralized monitoring.
RESTful APIs:

Expose endpoints via a Gin-based HTTP server.
Architecture
Server:

Manages REST APIs and serves client requests.
Publishes file and directory events to Kafka.
Worker:

Listens to Kafka topics.
Processes events (e.g., file upload or deletion) in real-time.
Configurable Storage:

Uses Azure Blob Storage or falls back to local storage based on configuration.
Tech Stack
Go: Core language for implementation.
Azure SDK for Go: Handles cloud storage operations.
Kafka: Event-driven communication.
Elasticsearch: Centralized logging.
Gin: HTTP web framework for APIs.
gopkg.in/yaml.v3: For configuration management.
Usage
Setup Configuration:

Update config.yaml with Azure and Kafka credentials.
Run the Server:

bash
Copy
Edit
go run cmd/server/main.go
Run the Worker:

bash
Copy
Edit
go run cmd/worker/main.go
Testing:

Unit tests for storage operations:
bash
Copy
Edit
go test ./...
Example Endpoints
File Upload: POST /upload/:path
Delete File: DELETE /delete/:path
List Files: GET /files/:path
Fetch Events: GET /events
