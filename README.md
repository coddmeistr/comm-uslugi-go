
# Simple go web-server

This project was written for my course work at university. This project provides an API for running a utility website. Administrators have the ability to confirm/reject requests and assign workers to execution of these requests. While regular users can create requests

### NOTE: This project does not implement a RestfulAPI architecture, because it was written a long time ago

### Technology or what I learned
I was using Gin framework to handle endpoints\
I was using MySQL as my main database and I was using GORM to interact with this database\
I implemented JWT authorization to handle Login and Registration events




## API Reference

#### Sign up

```http
  GET /signup
```
Create a new user

#### Sign in (Login)

```http
  POST /login
```
Log in into existing account

#### Logout

```http
  DELETE /logout
```
Reset auth data

#### Validate

```http
  GET /validate
```
Get current authentification info

#### New Request

```http
  POST /newrequest
```
Create new request (admin only)

#### Approve Request

```http
  PUT /approverequest
```
Approve existing request (admin only)

#### Reject Request

```http
  PUT /rejectrequest
```
Reject existing request (admin only)

#### Done Request

```http
  PUT /donerequest
```
Mark existing request as dont, means that it is completed (admin only)

#### Get All Requests

```http
  GET /getallrequests
```
Get all requests (admin only)


#### Get All User's Requests

```http
  GET /getalluserrequests
```
Get all requests of some user id


#### New Worker

```http
  POST /newworker
```
Create a new worker (admin only)


#### Get All Workers

```http
  GET /getallworkers
```
Get all workers (admin only)


