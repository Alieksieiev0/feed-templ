# Project feed-templ

Feed Templ contains the frontend of the Feed project.
It uses HTMX, Templ and Go to display and dynamically replace data on the website.

## Getting Started

To run this service, just clone it, and start it 
using either [Make Run](#run-binary-after-executing-make-templ-and-make-tail) or [Make Watch](#live-reload). 
However, to run properly, it requires API to be online. 
In scope of the Feed Project multiple microservices were used, 
mapped to the appropriate routes using NGINX. Feed-templ was mapped to the "/" location in the NGINX config.
<br/>
<br/>
Microservices and NGINX locations:
- [auth-service](https://github.com/Alieksieiev0/auth-service) -- **/api/auth**
- [feed-service](https://github.com/Alieksieiev0/feed-service) -- **/api/feed**
- [notification-service](https://github.com/Alieksieiev0/notification-service) -- **/api/notify**


## Site Structure and Components

### Header
This element is displayed on every page on the website.
Header contains logo(at the moment random logo is used), links to the [Home](#home) and (Search)(#search).
<br/>
If user is logged in, then he will also see [Notification Menu](#notification-menu), 
Profile Icon that redirects to the [Profile](#profile) on click
and Sign out button.
<br/>
Otherwise, the only thing that user will see except for default items is Sign in button,
that redirects to the [Sign in Page](#sign-in) on click.

### Home
This page contains a feed, that will show to user posts sorted by created date, 
displaying most recent in the top. By default, only 10 are loaded, and user can load more
using "Load More" button (It won't be displayed if initially less than 10 posts were loaded).
User can see date of creation, title, body and owner name. It is possible to open owner [Profile](#profile)
by clicking either on owner name or owner icon.
<br/>
Signed in users will be able to see Post creation form, where they can specify title and body.

### Search
On a search page users can find other users by username, and subscribe/unsubscribe to/from them.

### Notification Menu
This component displays 10 most recent notifications for current day. Application supports
real time notification system. On the server side, websocket connection with [notification-service](https://github.com/Alieksieiev0/notification-service)
is created for each client, waiting for new notifications. As soon as notification is sent to the server, it passes it down
to the client using server-side events and is dynamically appended to the list of user notifications.
Currently, user will receive notifications in 2 scenarios:

1. Subscribed user created a new Post.
2. New user subscriber.

More details on how exactly notification-service was implemented can be found in corresponding repository.
<br />
<br />
It is possible to remove notifications from list, by clicking on close icon. In fact, these notifications
are not delete, but moved from "NEW" status into "REVIEWED". This approach was taken to implement notification page
in future

### Profile
This page contains all relevant information about user, like Subscribers and Posts numbers.
Also, you can find list of all user posts on this page.

### Sign in
Default sign in form, that provides user with ability to login into system using username and password

### Sign up
Default sign in form, that provides user with ability to sign up into system using username, email and password

## MakeFile

### Generate .go files from .templ files
```bash
make some
```

### Build tailwindccs
```bash
make tail
```

### Build binary after executing make templ and make tail
```bash
make build
```

### Run binary after executing make templ and make tail
```bash
make run
```

### Clean
```bash
make clean
```

### Live Reload
```bash
make watch
```

## Flags
This application supports startup flags, 
that can be passed to change server and API urls. 
However, be careful changing feed-templ server url 
if you are running it using docker-compose, because by default
only port 3003 is exposed

### Server
- Name: web-server
- Default: 3003

### API
- Name: api
- Default: http://localhost:80
