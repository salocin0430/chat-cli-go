# Chat CLI

A command-line chat application using NATS JetStream for message persistence and distribution.

## Features

- Real-time messaging between users
- Message persistence for 1 hour
- Automatic message history recovery
- Multiple chat channels support
- Simple CLI interface
- Hexagonal Architecture implementation

## Requirements

- Go 1.22 or higher
- Docker and Docker Compose

## Architecture

This project follows the Hexagonal Architecture (Ports and Adapters) pattern, also known as Ports and Adapters architecture. The main components are:

### Domain Layer
The core business logic and entities of the chat application:
- Message entity with its properties and behaviors
- Core business rules and logic

### Application Layer
Contains the use cases and orchestrates the flow of data:
- Chat service implementation
- Message handling and distribution
- Channel management

### Ports Layer (Interfaces)
Defines the contracts for incoming and outgoing interactions:
- Input Ports: Define how the application receives commands
- Output Ports: Define how the application interacts with external services in this case NATs

### Infrastructure Layer (Adapters)
Implements the interfaces defined by the ports:
- NATS Adapter: Handles message persistence and distribution
- CLI Adapter: Manages user interaction through command line

## Message Flow

### Message Publishing
When a user sends a message, it flows through the system:
1. User input is captured by the CLI adapter
2. Message is processed by the application service
3. NATS adapter persists and distributes the message

### Message Storage
Messages are stored using NATS JetStream with the following characteristics:
- One-hour retention period
- File-based storage for persistence (using volumn docker)
- Automatic message expiration
- Channel-based organization

### Message Reception
The system handles message reception in two ways:
1. Historical Messages:
   - Retrieved when joining a channel
   - Limited to the last hour
   - Chronologically ordered

2. Real-time Messages:
   - Immediate delivery to all channel subscribers
   - Handled through NATS pub/sub system


## Development

### Project Structure
The project follows a hexagonal and clear  architecture approach with clear separation of concerns:
- Command layer for application entry points
- Internal packages for core functionality
- Infrastructure implementations for external services
- Interface definitions for system boundaries


## Installation and Usage

### Quick Start
1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/chat-cli.git
   cd chat-cli
   ```

2. Start NATS server using Docker Compose:
   ```bash
   docker-compose up -d
   ```

3. Run chat clients in different terminals:
   There are two ways to run the chat client:

   A. Using local Go installation:
   The application accepts three optional command-line parameters:
   ```bash
   go run cmd/main.go [-nats <NATS URL>] [-name <username>] [-channel <channel name>]
   ```

   If not specified, the following defaults are used:
   - NATS URL: "nats://localhost:4222"
   - Username: "anonymous" 
   - Channel: "general"

   Examples:
   ```bash
   # Default values
   go run cmd/main.go

   # With some parameters
   go run cmd/main.go -name "Alice" -channel "tech"

   # With all parameters
   go run cmd/main.go -nats "nats://localhost:4222" -name "Bob" -channel "general"
   ```

   B. Using Docker (recommended):
   Use the provided script that handles building and running the Docker container:
   ```bash
   # Default values
   ./run-chat.sh

   # With some parameters
   ./run-chat.sh -name "Andres" -channel "tech"

   # With all parameters
   ./run-chat.sh -nats "nats://localhost:4222" -name "Mike" -channel "general"
   ```


### Message History
- When a new user joins a channel, they automatically receive messages from the last hour
- Messages older than 1 hour are automatically deleted
- Messages are persisted even if the NATS server restarts

### Chat Commands
Once connected:
1. Type your message and press Enter to send
2. Press Ctrl+C to exit the chat
