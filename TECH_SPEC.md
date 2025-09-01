# Technical Specification

## GophKeeper Password Manager

### General Requirements

GophKeeper is a client–server system that allows users to securely and reliably store logins, passwords, binary data, and other private information.

The server must implement the following business logic:
- user registration, authentication, and authorization;
- storage of private data;
- synchronization of data between multiple authorized clients of the same owner;
- transfer of private data to the owner upon request.

The client must implement the following business logic:
- user authentication and authorization on the remote server;
- access to private data upon request.

Functions left to the discretion of the developer:
- creation, editing, and deletion of data on the server or client side;
- format of new user registration;
- choice of storage system and data storage format;
- ensuring security of data transmission and storage;
- client–server communication protocol;
- mechanisms for user authentication and access authorization.

Additional requirements:
- the client must be distributed as a CLI application runnable on Windows, Linux, and Mac OS;
- the client must provide the user with information about the version and build date of the client binary.

### Types of Stored Information

- login/password pairs;
- arbitrary text data;
- arbitrary binary data;
- bank card data.

For any type of data, it must be possible to store arbitrary text metadata (e.g., data ownership related to a website, person, or bank; lists of one-time activation codes; etc.).

### Abstract Interaction Flow

The following describes basic user interaction scenarios with the system. They are not exhaustive — handling of certain scenarios (such as resolving data conflicts on the server) is left to the developer’s discretion.

For a new user:
- The user downloads the client for their platform.
- The user completes the initial registration procedure.
- The user adds new data via the client.
- The client synchronizes data with the server.

For an existing user:
- The user downloads the client for their platform.
- The user goes through authentication.
- The client synchronizes data with the server.
- The user requests data.
- The client displays the data to the user.

### Testing and Documentation

The system code must be covered by unit tests at a level of no less than 80%. Every exported function, type, variable, and package of the system must contain complete documentation.

### Optional Features

The following features are optional but allow better assessment of the developer’s expertise. The developer may implement any number of them at their discretion:
- support for OTP (one-time password) data;
- support for a terminal user interface (TUI);
- use of a binary protocol;
- availability of functional and/or integration tests;
- description of the client–server communication protocol in Swagger format.
