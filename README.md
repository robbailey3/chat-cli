# ChatCLI

ChatCLI is a command line application written in Go that allows users to interact with OpenAI's GPT-4 language model. This application provides a way to communicate with OpenAI and offers an experience similar to that of [OpenAI Chat](https://chat.openai.com/chat) without using a browser.

## Features

- Send messages to OpenAI's GPT-4 language model
- View conversation in a chat-like format

## System Requirements

- An internet connection
- Go 1.16 or higher

## Installation

Clone the repository:

```bash
git clone https://github.com/robbailey3/chat-cli.git
```

Install required packages:

```bash
cd chat-cli
go mod tidy
```

Run the application:

```bash
go run main.go
```

Alternatively, the application can be compiled into a binary using the following command:

```bash
go build
```

**Please Note:** The application looks for a config file (`.chat-cli.yaml`) in the user's home directory. The contents of this file should include an API key for OpenAI.

Example:

```yaml
openAi: 
  apiKey: "sk-ABCDEFGHIJKLMNOP1234...."
```

## Usage

The chat-cli application provides a simple command line interface. After starting the application, users can send messages to the OpenAI's GPT-4 language model by typing them into the terminal and pressing enter.

Currently, there are no usage examples or tutorials available. However, I am working on creating these resources to help users get started with the application.

## Future Features

- Ability to choose the model
- Better loading 
- Chat histories
- Better settings
- Error handling