---
title: Simple Express.js Backend pt. 1
author: Garrett Davis
live: true
createdAt: "2022-05-23T15:23:00Z"
updatedAt: "2022-05-23T17:23:00Z"
description: In this article, I show the basics of using express.js to build a REST API
---

I've learned a lot of new things since my last post - sorry for the delay.
But, most notable is Node.js. I'm really growing to enjoy JavaScript, even
though before I said I like Python more. Python's interpreted nature makes
for a quick, engaging development experience, but JavaScript has the same
virtue! Additionally, with JavaScript because it runs alongside CSS and HTML,
it is much easier to make a UI.

I've also learned about Node.js, which is an awesome way to run JavaScript
outside of the browser. Express.js is a module for Node.js that allows you to
make a REST API for your web app using JavaScript.

Running JavaScript in Node.js is very similar to running it in the browser.
The main difference is the types of programs you're writing. Node.js allows
you to create a backend for your web apps using JavaScript.

## Express.js

Express.js is a great framework for creating a simple backend server, serving
a RESTful API. REST stands for representational state transfer, and API
stands for application programming interface. Together, a RESTful API
transfers a representation of the state of the resources needed for building
an application.

The main operations a RESTful API can do are GET, POST, PUT, and DELETE.
These 4 operations are done when a client "requests" them to be, using a HTML
Request.

- **GET** requests ask the backend for information.
- **POST** requests send new information to the backend, for example creating
  a new entry in an array. They are not indempotent, which means that the
  same requst multiple times will keep adding information. This is in
  contrast to:
- **PUT** requests, which send new information also, but are indempotent, so
  if you send the same PUT request repeatedly, only the first one has an
  effect. An example is updating an array's entry. If you set the value of an
  item at a specific index, if you reset it to the same value, you wouldn't
  expect it to change again.
- **DELETE** requests as you would expect are for deleting information.

## An Example

Using Express.js, you start an instance of the express server, by running
`const app = express();` Then for each of the requests, app can be assigned
a callback using a route method, for a specific path. That could look like
the following:

```js
// GET /data/id
app.get("/data/:id", (req, res) => {
  // get the id from the query parameters.
  const id = req.params.id;
  // look up data somehow. You want to make sure it exists.
  if (dataArray.findIndex(id) !== -1) {
    const data = dataArray[id];
    // send the data back to the client in the response.
    res.status(200).send(data);
  } else {
    // if the data doesn't exist, tell the client.
    res.status(404).send();
  }
});

// POST /data
app.post("/data", (req, res) => {
  // get the data from the request body.
  const data = req.body;
  // make sure the data is valid, then add the data entry.
  if (validData(data)) {
    dataArray.push(data);
    // tell the client that an entry was create, and what it was.
    res.status(201).send(data);
  } else {
    // if the data was invalid, tell the client it was rejected.
    res.status(403).send();
  }
});

// PUT /data/id
app.put("/data/:id", (req, res) => {
  // get the id and data.
  const id = req.params.id;
  const data = req.body;
  // check to see if the entry doesn't exist and tell the client.
  if (dataArray.findIndex(id) === -1) {
    res.status(404).send();
    // check to see if the data is invalid and tell the client.
  } else if (!validData(data)) {
    res.status(403).send();
    // data must be valid, so update it.
  } else {
    dataArray[id] = data;
    // tell the client that the entry was updated.
    res.status(204).send();
  }
});

// DELETE /data/id
app.delete("/data/:id", (req, res) => {
  // get the id.
  const id = req.params.id;
  // check to see if the entry exists, and delete it.
  if (dataArray.findIndex(id) !== -1) {
    dataArray.splice(id, 1);
    // tell the client the entry was deleted.
    res.status(204).send();
  } else {
    // if the entry didn't exist, tell the client.
    res.status(404).send();
  }
});
```

When a request is made, if it's type and path matches any of the route
methods, then that method's callback is executed. The path is specified using
the string passed to the route method. You can also pass an array of multiple
strings, or use regex to select match paths. In the example GET, PUT, and
DELETE all had the same path `'/data/:id'`. Because they all were different
request types though, each request will only correspond to one of them.

The `/:id` part is significant. This lets you match a parameter in the path.
If you were to make a request 'GET /data/87' then it would set `id` to 87.

The callback has parameters req which represents the request, and res which
represents the response. You can pass more parameters, like `(req, res, next)`
or `(err, req, res, next)` to add more functionality. err represents an error
thrown while handling a request. You should put a route method that handles
errors at the bottom. This is important, because Express.js can handle
multiple matching route methods, but it executes them in the order they are
declared in.

The other important callback parameter is next which represents the next
route method declared that matches. It is a function itself, so you can just
pass the error to it, as in `next(Error)`.

## Middleware `app.param()` and `app.use()`

Errors aren't the only use of next though. You can also use it to DRY your
code, or "Don't Repeat Yourself." You probably noticed, much of the route
methods from before were very similar, so in order to limit repetition, you
can use `app.param()` or `app.use()`.

The `app.param()` method lets you define parameters once, instead of for each
route method that uses that parameter. Here's an example of it's use:

```js
app.param("id", (req, res, next, id) => {
  // get the id and set it as an attribute of req.
  req.id = id;
  // call the next route handler.
  next();
});

app.get("/data/:id", (req, res, next) => {
  // look up data and make sure it exists.
  if (dataArray.findIndex(req.id) !== -1) {
    const data = dataArray[req.id];
    // send the data back to the client in the response.
    res.status(200).send(data);
  } else {
    // if the data doesn't exist, tell the client.
    res.status(404).send();
  } // call the next route handler.
  next();
});
```

Now for every route method that uses the id query parameter, you don't have
to redefine it! This helps avoid errors in your code, since there isn't
repeated blocks of code that all must be updated individually.

`app.use()` is used just like every other route handler, but it matches any
kind of request. If you wanted to use the example just given with param and
get, you might add an `app.use()` to perform logging for all requests at the
end of your file, but before error handling. Here's an example:

```js
app.use("/data/:id", (req, res, next) => {
  // log the request method and current time.
  console.log(`Recieved ${req.method} request at ${Date()}.`);
  // call the next route handler.
  next();
});
```
