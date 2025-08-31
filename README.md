# Messaging-Go Core

A lightweight, extensible middleware framework for Go. Provides the building blocks to implement robust message consumers with backoff, retries, graceful shutdown, and other custom middleware.

This library focuses on core message processing mechanics, allowing developers to integrate it with any provider (Kafka, SQS, RabbitMQ, etc.) via provider-specific middleware.

## Features

- **Composable Middleware**: Chain multiple middlewares in any order.
- **Basic in-built middlewares**: Backoff, Graceful Shutdown.
- **Extensible**: Plug in different middleware to consume messages from different providers (kafka, sqs, rabbitmq, etc.)
- **Custom Middleware**: Implement your own middleware to handle custom message processing logic.
- **Provider-Agnostic**: Core doesn’t directly handle Kafka, SQS, or RabbitMQ — you provide that via middleware.

## Core Concepts

The core library revolves around three concepts:

1. Core Consumer Loop
   - Runs the message processing pipeline repeatedly.
   - Provides a single entry point for middleware chaining.
2. Middleware
   - Each middleware implements a Process function.
   - Middleware can:
     - Pre-process messages
     - Call next to pass control
     - Post-process or handle errors
3. Middleware Chain Flow
  - The first middleware in the chain is called first
  - The next middleware is called after the previous middleware has completed
  - The last middleware in the chain is called after all other middlewares have completed
  - And the message processing results are returned through the chain
4. Message
   - A message is a generic interface that can be used to represent any type of message.
   - Middleware can process messages of any type.
   - Some middleware may require a specific type of message.

## Extending the Core Library

If you plan to consume a message from a provider that isn't supported by the core library,
Here's the recommended approach:
1. Create an issue in this repository to request support for the provider and see if someone is willing to collaborate.
2. If you prefer to contribute, ask the maintainers to create a new repository for the provider with a blank template.
3. Implement the middleware for the provider.
4. Get the middleware reviewed, merged, and released.
5. Make sure things like producers, retries, deadletters are already taken care of.
