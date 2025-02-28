# Electronic Library REST JSON API Documentation

## Table of Contents

1. [Overview](#overview)
2. [Prerequisites](#prerequisites)
3. [Installation](#installation)
4. [Running the App](#running-the-application)
5. [Debugging the App](#debugging-the-application)
6. [Concurrency Control in PostgreSQL](#concurrency-control-in-postgresql)
7. [Installing and Using Postman for API Testing](#installing-and-using-postman-for-api-testing)
8. [Ideas / Suggestions](#ideas--suggestions)

---

## Overview

This project is a simple REST API built with Go (Golang). The API enables users to perform the following actions:

- Search for a book by title using Postgres Full Text Search.
- Borrow a book.
- Extend the return date of a borrowed book.
- Return the book.

It supports CRUD (Create, Read, Update, Delete) operations on resources and is designed to handle HTTP requests, returning responses in JSON format.

## Features

### PostgreSQL Full Text Search

The API utilizes **PostgreSQL Full Text Search** to enable efficient searching of book titles. This feature allows users to search for books by title, using advanced text matching capabilities such as stemming and ranking of results. Full Text Search is optimized for fast and accurate search queries, even with large datasets.

For more details on how Full Text Search works in PostgreSQL, refer to the official PostgreSQL documentation:  
[PostgreSQL Full Text Search Documentation](https://www.postgresql.org/docs/current/textsearch.html)

#### How It Works

- The API uses the `tsvector` type in PostgreSQL to index the titles of the books.
- A search query is sent to the API where it is processed using PostgreSQL’s `to_tsquery` function to match the query against the indexed `tsvector`.
- The results are ranked based on relevance, providing the user with the most relevant books first.

Example SQL query used for the search:

```sql
SELECT title, author
FROM books
WHERE to_tsvector('english', title) @@ plainto_tsquery('english', 'your search query');
```

## Prerequisites

Before running the application, ensure the following dependencies are installed:

- Go (version 1.24+)
- Docker
- golang-migrate/migrate CLI
- Any IDE (e.g., VSCode, GoLand)

## Installation

### 1. Install Go

#### Linux

- Download the tarball from the official Go website:

    ```bash
    wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz
    ```

- Extract the archive to `/usr/local`:

    ```bash
    tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz
    ```

- Create the `go` directory in your home directory:

    ```bash
    mkdir ~/go
    ```

- Configure the `GOPATH` and `GOROOT` environment variables:

    ```bash
    export GOROOT=/usr/local/go
    export GOPATH=$HOME/go
    export PATH=$PATH:/usr/local/go/bin
    ```

#### macOS

To install Go on macOS, use Homebrew:

```bash
brew install go
```

### 2. Install docker

You can install docker for linux or Mac below:
[Get Docker](https://docs.docker.com/get-started/get-docker/)

### 3. Install golang-migrate/migrate CLI

#### With Go toolchain

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@4.1.0
```

#### With Homebrew

```bash
brew install golang-migrate
```

## Running the Application

### 1. Migrate the tables

Make sure the -database flag value is correct

```bash
migrate -database postgresql://electronic-library:electronic-library@localhost:5432/electronic-library?sslmode=disable -path internal/db/migrations -verbose up
```

### 2. Seed the database

```bash
go run cmd/seeder/main.go
```

### 3. Start your server

```bash
go run cmd/api/main.go
```

## Debugging the Application

To debug the application using **VSCode**, follow these steps:

1. **Install the Go extension** for VSCode:
   - Go to the [Go extension page on Visual Studio Marketplace](https://marketplace.visualstudio.com/items?itemName=golang.Go) and install the extension.

2. **Set up Delve for debugging**:
   - Delve is a powerful debugger for Go. In VSCode, you can configure it for the server.
   - Open your project in VSCode, then go to the **Run and Debug** panel.
   - Select **Run and Debug** to start the server with Delve.

3. **Run and Debug**:
   - Once everything is set up, use the **Run and Debug** option in VSCode to start the application.
   - You can now set breakpoints, step through code, and inspect variables to troubleshoot any issues in your Go application.

## Concurrency Control in PostgreSQL

PostgreSQL offers **Serializable Transactions** as the highest isolation level, preventing issues such as **dirty reads**, **nonrepeatable reads**, **phantom reads**, and **serialization anomalies**.

### Common Concurrency Issues and Scenarios in Your eBook API

1. **Dirty Read**  
   - **Scenario**: A user searches for a book and another user starts borrowing that same book but hasn’t committed the transaction yet.
   - **Problem**: Without proper isolation, the first user might read the data (i.e., the book is available) before the second user’s borrow transaction is committed, resulting in a **dirty read**.

   - **Example (Dirty Read)**:
     - **User 1** searches for a book titled "Go Programming".
     - **User 2** begins borrowing the same book, but the transaction is not yet committed.
     - **User 1** sees that the book is available and proceeds with borrowing it.
     - **User 1** may end up borrowing the book even though **User 2** has already borrowed it, leading to inconsistent data.

   - **Solution**: **Serializable Transactions** ensure that the transaction is isolated, meaning **User 1** will either see the committed state (i.e., the book is already borrowed) or not see the book at all until **User 2**'s transaction is completed. No dirty reads are allowed.

2. **Nonrepeatable Read**  
   - **Scenario**: A user reads the book’s availability (e.g., available for borrowing), but before they can act on it, another user modifies the data (e.g., borrows the book), and the original user re-reads the data, finding that the availability status has changed.
   - **Problem**: Without proper isolation, the first user could experience a **nonrepeatable read**, where data they read earlier has been modified by another transaction before they act on it.

   - **Example (Nonrepeatable Read)**:
     - **User 1** searches for a book titled "Advanced Go".
     - **User 2** borrows "Advanced Go", changing its availability status.
     - **User 1** reads the book and sees that it’s available.
     - **User 1** returns later to borrow the book and finds that it's no longer available, because **User 2**'s transaction modified the book's status.

   - **Solution**: **Serializable Transactions** would ensure that **User 1** sees the correct state of the book and is not affected by changes made by **User 2** during their transaction. They will either see the book as available or unavailable based on the final committed state, preventing nonrepeatable reads.

3. **Phantom Read**  
   - **Scenario**: A user queries for books available for borrowing, but while they are processing the query, another user adds or removes books, causing the result set to change.
   - **Problem**: Without proper isolation, the first user might get inconsistent results if the set of books changes due to another transaction before or during the query.

   - **Example (Phantom Read)**:
     - **User 1** queries for all available books titled "Programming".
     - **Admin** modifies the list of available books.
     - **User 1** re-runs the same query, but now the result set is different because **Admin**'s transaction removed a book from the list.

   - **Solution**: **Serializable Transactions** ensure that once **User 1** begins the query, no other transaction can modify the result set until **User 1**'s transaction completes, preventing phantom reads.

4. **Serialization Anomaly**  
   - **Scenario**: A group of transactions is executed concurrently, but the final state of the data is inconsistent with any possible sequence of those transactions.
   - **Problem**: Without proper isolation, the system could allow conflicting operations that result in inconsistent data or unexpected behaviors.

   - **Example (Serialization Anomaly)**:
     - **User 1** borrows a book and then returns it.
     - **User 2** tries to borrow the same book during **User 1**'s transaction.
     - If transactions are not properly isolated, **User 2** might be able to borrow the book before **User 1** returns it, violating the intended consistency of the system.

   - **Solution**: **Serializable Transactions** ensure that all transactions are executed in a serializable order, meaning no conflicting operations can occur, preventing serialization anomalies.

---

### Isolation Levels in PostgreSQL

Here’s how PostgreSQL handles these phenomena at different isolation levels:

| Isolation Level     | Dirty Read | Nonrepeatable Read | Phantom Read | Serialization Anomaly |
|---------------------|------------|---------------------|--------------|------------------------|
| Read Uncommitted    | Allowed (not in PG) | Possible | Possible | Possible |
| Read Committed      | Not possible | Possible | Possible | Possible |
| Repeatable Read     | Not possible | Not possible | Allowed (not in PG) | Possible |
| Serializable        | Not possible | Not possible | Not possible | Not possible |

For more details, refer to the [PostgreSQL Concurrency Control Documentation](https://www.postgresql.org/docs/current/mvcc.html).

## Installing and Using Postman for API Testing

### 1. Install Postman

To install the Postman Desktop App, follow the steps below:

  1. Visit the [Postman Download Page](https://www.postman.com/downloads/).
  2. Download the installer for Windows.
  3. Run the installer and follow the on-screen instructions to complete the installation.

### 2. Import API Collection to Postman

Once Postman is installed, you can import your eBook API request collection into the Postman Desktop App by following these steps:

1. **Copy the Collection Link**:  
   Copy the following URL to import the collection:  
   `https://www.postman.com/lively-spaceship-99649/workspace/public-applications/collection/16588736-ad74bcde-d340-48c2-af68-5327dd03c1e1?action=share&creator=16588736`

2. **Open Postman Desktop App**:
   - Launch the Postman Desktop App after installation.

3. **Import the Collection**:
   - Click on the **Import** button located in the top-left corner of the Postman app.
   - In the "Import" dialog, choose the **Link** tab.
   - Paste the copied URL into the input field.
   - Click **Continue** and then **Import** to bring the collection into your Postman.

4. **Access the API Requests**:
   - Once imported, you should see the collection listed under **Collections** on the left sidebar.
   - Click on the collection and select the specific request you wish to test.

Now you can start testing and interacting with your eBook API directly within Postman!

## Ideas / Suggestions

**Authentication Mechanisms**:  
There are several authentication methods that could be implemented to secure the API. These include:

- **Basic Authentication**: A simple method where the client sends a username and password in the HTTP request headers for validation.
  
- **Token-Based Authentication**: The client sends a token (usually a string) in the request headers to verify the user's identity, allowing access to protected resources.

- **JWT (JSON Web Tokens)**: A stateless authentication method where a server generates a signed token containing user information. The client sends this token in each request to prove their identity.

- **OAuth2**: A more complex, third-party authentication mechanism allowing users to authenticate via external services like Google, Facebook, or GitHub.

- **Session-Based Authentication**: The server keeps track of user sessions, and the client sends a session cookie with each request to maintain authentication.

**Reason for Not Implementing Yet**:  
Authentication has a broad specification, and the choice of which method to use depends on several factors, such as security needs, application requirements, and user base. For now, I haven’t implemented any of these methods to keep the focus on other functionalities.

---

**Pagination**:  
Pagination has been added to the book search functionality to manage large sets of results. This is particularly important for APIs that may deal with large amounts of data, ensuring better performance and a smoother user experience.

**Reasons for Adding Pagination**:  

- **Improved Performance**: Without pagination, the search query for books could return too many results, causing high memory usage and slowing down response times. Pagination allows the API to return only a subset of results at a time, improving speed.
  
- **User Experience**: Pagination makes it easier for users to browse through large datasets without overwhelming them with long lists of results. It also provides more control over the number of results displayed.
  
- **Reduced Server Load**: By limiting the number of results returned in each query, we reduce the load on the database and the server, as fewer records are processed and sent in each response.

---

**Unique Borrower per Book**:  
For the **borrow** and **return** endpoints, the idea is to ensure that each book is linked to a unique borrower at any given time, meaning only one borrower’s name can be associated with the book until it is returned.

**Reason for Not Implementing Yet**:  
Due to time constraints, this feature has not been implemented. While it is an important feature for data consistency and preventing multiple borrowers for a single book, it requires additional logic to check and enforce the uniqueness of the borrower per book.

To implement this properly, we would need to introduce a **User table** to manage borrowers and link each borrowing transaction to a unique user. This would require:

- Creating a **User** table to store user data (e.g., name, email, etc.).
- Updating the **Loan Detail** table to store the relationship between borrowed books and their borrowers.
- Modifying the **borrow** and **return** endpoints to ensure that a book can only be borrowed by one user at a time.

---

**Rate Limiting**:  
Rate limiting is a technique used to control the amount of incoming requests to an API in a given period. It helps prevent abuse, protects server resources, and ensures fair usage of the API.

**Reason for Not Implementing Yet**:  
Although rate limiting is an important feature to ensure the API remains secure and performant, it has not been implemented yet. The challenge lies in deciding where to apply the rate limit:  

- **On the AWS Side**: AWS offers services like **API Gateway** or **AWS WAF** that can be used to enforce rate limiting at the edge before requests hit the application server. This would reduce the load on the application server but may require configuring and managing these AWS services.  
- **On the Application Side**: Implementing rate limiting within the app itself could offer more granular control over the rate limiting logic. However, this would increase the load on the server and may complicate the implementation.

---

**API Versioning**:  
API versioning is a practice used to ensure backward compatibility when making changes to an API. By maintaining different versions of the API, we can ensure that existing clients continue to function as expected while new features are introduced.

**Reason for Not Implementing Yet**:  
API versioning is important for long-term maintenance, but it has not been implemented yet for the following reasons:

- **No Major Changes**: The API is still in its early stages, and no major breaking changes have been introduced that would necessitate versioning. Thus, it's not critical at this point.
- **Complexity**: Introducing versioning adds additional complexity to the routing and management of API endpoints. For instance, we'd need to ensure that the version is correctly handled by the server, and it could require changes to endpoint definitions, route handling, and documentation.

Possible strategies for API versioning include:

- **URI Versioning**: Including the version number in the URL, like `/api/v1/books`, `/api/v2/books`, etc.
- **Header Versioning**: Specifying the API version in the request header, such as `Accept: application/vnd.api.v1+json`.
- **Query Parameter Versioning**: Using a query parameter like `/api/books?version=1`.

---

**Serializable Transaction Error Retry Mechanism**
A serializable transaction error retry mechanism ensures that database transactions are handled in a way that they are executed serially without interference from concurrent transactions. When a transaction fails due to issues like deadlocks or temporary errors, this mechanism attempts to retry the operation a defined number of times, reducing the chances of system downtime or inconsistencies.

**Reason for Not Implementing Yet**
The retry mechanism has not been implemented due to its complexity in transaction management. Implementing a retry mechanism requires careful handling of transaction isolation levels to avoid issues like data inconsistencies or infinite loops. Additionally, retries can introduce delays, potentially affecting performance, and may require fine-tuning for efficiency.
