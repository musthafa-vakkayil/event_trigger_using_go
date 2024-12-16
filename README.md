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
1. Clone this repository- (https://github.com/musthafa-vakkayil/event_trigger_using_go.git)
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
![Blank diagram - Page 1](https://github.com/user-attachments/assets/6b586fba-8862-4e53-a81c-72ae7dd63388)
### 1. Database Details
We use **Postgres** as the database, which includes three tables:
- **Triggers**
- **Events**
- **Users**

Currently, there are no foreign key relationships to avoid complexity. The tables are straightforward.

### 2. Server Details
The system uses Golang to handle trigger scheduling, API requests, and log management. The architecture ensures flexibility and scalability.


## Deployment Details
![image](https://github.com/user-attachments/assets/b68836dd-5068-4bd1-9c22-4c378f4a70e9)

For deployment, I used an **AWS EC2 instance** with my free-tier subscription. The deployment process was automated using **GitHub Actions**.


### Pricing
The approximate cost of running a **t3.micro EC2 instance** with a **16 GB EBS volume** for 30 days (24x7):

#### 1. EC2 Instance Costs:
- **Pricing**: $0.0104/hour in the US East (N. Virginia) region.
- **Monthly Usage**: 30 days × 24 hours = 720 hours.
- **Instance Cost**: 720 × $0.0104 = **$7.49**.

#### 2. EBS Volume Costs:
- **Volume Type**: General Purpose SSD (gp3).
- **Pricing**: $0.08/GB/month.
- **For a 16 GB volume**: 16 × $0.08 = **$1.28**.

### Total Cost:
**$7.49 + $1.28 = $8.77/month**.

## API's

### CRUD API's for Triggers

1. **Create Trigger**
   - **Endpoint**: `/v1/triggers/create` (POST)
   - Description: Use this API to create either an API trigger or a scheduled trigger.
     - For API trigger, refer below image:
       ![image](https://github.com/user-attachments/assets/e166c444-c35b-44f8-ac51-0e256eee32fb)
       
     **cURL Command:**
     ```bash
     curl --location 'http://18.205.239.53:8080/v1/triggers/create' \
      --header 'Content-Type: application/json' \
      --data '{
       "name": "API TRIGGER 1",
       "type": "API",
       "schedule_time": 0,
       "interval_seconds": 0,
       "api_method": "GET",
       "api_endpoint": "https://httpbin.org/get",
       "api_payload": {}
      }'
     ```
     - For scheduled trigger, refer below image:
       ![image](https://github.com/user-attachments/assets/0bc8c62a-008d-448d-b952-bf54f417a16c)

     **cURL Command:**
     ```bash
     curl --location 'http://18.205.239.53:8080/v1/triggers/create' \
      --header 'Content-Type: application/json' \
      --data '{
       "name": "SCHEDULED TRIGGER 1",
       "type": "SCHEDULED",
       "schedule_time": 120, 
       "interval_seconds": 0,
       "is_recurring": false,
       "api_method": "",
       "api_endpoint": "",
       "api_payload": {}
      }'
     ```


2. **Get Trigger Details**
   - **Endpoint**: `/v1/triggers/{trigger_id}` (GET)
   - Description: Use this API to retrieve the details of a particular trigger.

   **cURL Command:**
     ```bash
     curl --location 'http://18.205.239.53:8080/v1/triggers/b3acd2f9-6c0b-41c3-a27d-84eb24b4d9ea'
     ```

3. **Get All Triggers**
   - **Endpoint**: `/v1/triggers` (GET)
   - Description: Use this endpoint to get a list of all triggers in the system.
   - 
   **cURL Command:**
     ```bash
     curl --location 'http://18.205.239.53:8080/v1/triggers'
     ```

4. **Delete Trigger**
   - **Endpoint**: `/v1/triggers/{trigger_id}` (DELETE)
   - Description: Use this endpoint to delete a stored trigger.
   
   **cURL Command:**
     ```bash
     curl --location --request DELETE 'http://18.205.239.53:8080/v1/triggers/8637358e-9156-4e95-861a-f3fb7b81f23c'
     ```

5. **Update Trigger**
   - **Endpoint**: `/v1/triggers/{trigger_id}` (PUT)
   - Description: Use this endpoint to update the details of an existing trigger.

   **cURL Command:**
     ```bash
     curl --location --request PUT 'http://localhost/v1/triggers/b3acd2f9-6c0b-41c3-a27d-84eb24b4d9ea' \
      --header 'Content-Type: application/json' \
      --data '{
       "name": "Update Trigger Details",
       "type": "API",
       "api_endpoint": "https://httpbin.org/get",
       "api_payload": {
           "key": "value"
       },
       "api_method": "POST"
      }'
     ```

### API's for Events
1. **View Logs**
   - **Endpoint**: `/v1/events` (GET)
   - Description: Use this API to view the logs.
   - Additionaly this API supports two query params - includeArchived (to show only archived) and includeAll (to view all logs) with default value false
   ![image](https://github.com/user-attachments/assets/df69e6cb-3886-4f6f-8705-bc7cdcd768b6)

   **cURL Command:**
     ```bash
     curl --location 'http://18.205.239.53:8080/v1/events'
     ```
2. **View Logs**
   - **Endpoint**: `/v1/triggers/api/{trigger_id}` (GET)
   - Description: Use this API to execute the stored API trigger.

   ![image](https://github.com/user-attachments/assets/6ab370ee-f1e7-4082-b896-6dd387b88e24)

   **cURL Command:**
     ```bash
     curl --location 'http://18.205.239.53:8080/v1/triggers/api/b3acd2f9-6c0b-41c3-a27d-84eb24b4d9ea'
     ```

### API for Testing Event without Saving to DB
 - **Endpoint**: `/v1/test/trigger` (POST)
   - Description: Use this API to create either an API trigger or a scheduled trigger without saving it.
   - Currently this API supports API trigger and one time scheduled trigger

     **cURL Command:**
     ```bash
     curl --location 'http://18.205.239.53:8080/v1/test/trigger' \
      --header 'Content-Type: application/json' \
      --data '{
       "name": "API TRIGGER 1",
       "type": "API",
       "schedule_time": 0,
       "interval_seconds": 0,
       "api_method": "GET",
       "api_endpoint": "https://httpbin.org/get",
       "api_payload": {}
      }'
