# Event Trigger System Using Golang and Postgres

## Credits
Below are the resources utilized to create this server:

1. **YouTube Tutorials**: 
   - [Part 1](https://www.youtube.com/watch?v=dFQa06cK7xE&t=6s) 
   - [Part 2](https://www.youtube.com/watch?v=oslAAQdGTCk)
     These tutorials provided valuable insights into database events and listening to them, even though the `pg_Cron` solution didn't work in this case.
2. **ChatGPT**: For suggestions on debugging and refactoring code.
3. **Stack Overflow**: Helpful threads on issues related to events, deployment, etc., with many useful answers.
4. [**Deployment Tutorial**](https://www.youtube.com/watch?v=sSAWMr_-Co4&t=60s): This video guided the deployment of the server on AWS.

---

## About this Project

This project is a **Golang-based Trigger Scheduler App** that allows users to schedule and manage triggers. Triggers can be of two types:

1. **Scheduled Trigger**: Can be one-time or recurring.
2. **API Trigger**: Enables users to create a trigger to hit an API endpoint, which can be invoked later using another API.

### Logs
A key feature of the system is the detailed logs for executed triggers:
- Logs remain in the **ACTIVE** state for the first 2 hours after execution.
- After 2 hours, they are moved to the **ARCHIVE** state and are hidden by default.
- Archived logs can still be viewed and are deleted after 46 hours.

Timings for log states can be easily adjusted in the code.

### Additional Features
- Users can test **one-time scheduled triggers** and **API triggers** without storing them in the database by sending the same payload used for creating triggers. 
- Testing recurring triggers is still under development due to time constraints.
- User registration is available, but authentication is still in progress.

---

## Interacting With the Server

- **Base URL**: [http://18.205.239.53:8080/](http://18.205.239.53:8080/)
- **Swagger Documentation**: [http://18.205.239.53:8080/swagger/index.html#/](http://18.205.239.53:8080/swagger/index.html#/)

---

## Setting Up Locally

### Option 1: Using Docker
1. Clone this repository.
2. Open a terminal and execute the following command:

   ```bash
   docker-compose up --build
   ```

3. The server will be running at `localhost:8080`.

### Option 2: Without Docker
1. Clone this repository.
2. Update the `.env` file inside the `server` folder with your Postgres database details.
3. Navigate to the `server` folder in the terminal and run:

   ```bash
   go run main.go
   ```

4. The server will be running at `localhost:8080`.

---

## Technical Details

### 1. Database Details
We use **Postgres** as the database, which includes three tables:
- **Triggers**
- **Events**
- **Users**

Currently, there are no foreign key relationships to avoid complexity. The tables are straightforward.

### 2. Server Details
The system uses Golang to handle trigger scheduling, API requests, and log management. The architecture ensures flexibility and scalability.
