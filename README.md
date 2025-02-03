# Azure Blob Service Golang

## Overview

The **Azure Blob Service Golang** application provides an API-driven solution for managing file storage and retrieval using either **Azure Blob Storage** or local storage as a fallback. It is designed for cloud-native applications requiring a robust and scalable file storage system.

## Features

- **File Upload**: Upload files to Azure Blob Storage or local storage with optional overwrite functionality.
- **File Read**: Retrieve files stored in Azure Blob Storage or local storage.
- **File Deletion**: Delete files from the storage system.
- **Directory Operations**: Support for creating and deleting directories in local storage.
- **Event-Driven Architecture**: Kafka integration to process and log file events, such as uploads and deletions.
- **Elasticsearch Logging**: Log events and errors into Elasticsearch for observability and debugging.

## Key Components

### 1. **Storage Adapters**
   - **Azure Blob Storage**:
     - Connects to Azure Blob Storage using an account name, account key, and container name.
     - Supports file upload, read, delete, and list operations.
   - **Local Storage**:
     - Stores files locally on the server’s filesystem.
     - Provides an alternative when Azure credentials are not configured.

### 2. **Kafka Integration**
   - Kafka-based messaging for event-driven architecture.
   - Supports publishing and consuming events for file operations.

### 3. **Elasticsearch Logging**
   - Centralized logging for monitoring system events.
   - Sends logs to an Elasticsearch instance for analytics and debugging.

### 4. **REST API**
   - Built using **Gin Web Framework**.
   - Provides endpoints for file and directory operations.

## API Endpoints

### File Operations
- `POST /upload/:path`: Upload a file to the specified path.
- `GET /read/:path`: Retrieve a file from the specified path.
- `DELETE /delete/:path`: Delete a file from the specified path.

### Directory Operations
- `POST /directory/:path`: Create a directory at the specified path.
- `DELETE /directory/:path`: Delete a directory at the specified path.

### Event Operations
- `GET /events`: Fetch recent file operation events from Kafka.

## Configuration

The application uses a YAML-based configuration file (`config.yaml`) to manage settings, including:
- **Azure Storage**: `accountName`, `accountKey`, and `containerName`.
- **Kafka**: Brokers, consumer group, and topics.
- **Elasticsearch**: URL for logging.

### Steps
1. Clone the repository:
   ```bash
   git clone https://github.com/your-repo/AzureBlobServiceGolang.git
   cd AzureBlobServiceGolang

## Testing and Running the Application

### Prerequisites

Before running or testing the application, ensure the following are installed and configured:
- **Go** (version 1.17+ recommended)
- **Azure Storage Account** (if using Azure Blob Storage)
- **Kafka** broker (if running Kafka event-driven features)
- **Elasticsearch** (optional for logging)

### Running the Application

#### 1. Server
The server handles API requests for file and directory operations.

To start the server:
```bash
go run cmd/server/main.go