# Go TCP Chat Server

This is a simple TCP chat server written in Go. It's designed to handle multiple client connections and facilitate real-time messaging, with built-in moderation features to prevent spam and abuse.
Features

    Concurrent Client Handling: Manages multiple client connections using goroutines and channels.

    Centralized Messaging: All messages (connection, disconnection, new message) are sent to a central server goroutine for processing.

    Rate Limiting: Prevents clients from spamming the chat by enforcing a message rate limit.

    Strike System: Clients who violate the rate limit or send invalid data receive strikes.

    Temporary IP Banning: Clients who reach the strike limit are temporarily banned, preventing them from reconnecting until the ban duration has expired.

    Safe Mode: An optional feature to redact sensitive information (like client IP addresses) from server logs.

How to Run the Server

    Ensure you have Go installed: If not, you can download it from the official Go website.

    Save the code: Place the provided code into a file named main.go.

    Run the server: Open your terminal and navigate to the directory where you saved the file. Then, run the following command:

    go run main.go

    The server will start listening for TCP connections on port 6969.

How to Connect

You can connect to the server using a simple TCP client. For a quick test, you can use the nc (netcat) command-line utility.

nc localhost 6969

Once connected, you can type your messages and press enter to send them.
Code Structure

    main(): The entry point of the application. It sets up the TCP listener, creates a message channel, and starts the server goroutine.

    server(): The core logic of the chat application. It manages connected clients, handles incoming messages from the channel, and enforces moderation rules.

    client(): A goroutine for each connected client that reads messages from the connection and sends them to the server goroutine.

    sensitive(): A utility function that redacts sensitive information if SafeMode is enabled.

Constants

The server's behavior is configured via constants at the top of the file:

    Port: The TCP port the server listens on.

    SafeMode: A boolean to enable or disable redaction of sensitive information.

    MessageRate: The minimum time interval (in seconds) between messages from a single client.

    BanLimit: The duration (in seconds) for a temporary ban.

    StrikeLimit: The number of strikes a client can receive before being banned.
