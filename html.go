package main

var AdUserHTML = `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Add User</title>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bulma@0.9.4/css/bulma.min.css">
    <script src="https://unpkg.com/htmx.org@1.9.5"></script> <style>
        body {
            background-color: #363636; 
        }
    </style>
</head>
<body>

<div class="container">
    <div class="columns is-centered">
        <div class="column is-half">
            <div class="box has-background-dark">
                <h2 class="title is-2 has-text-white">Add User</h2>

                <form hx-post="/adduser" hx-swap="outerHTML">  
                    <div class="field">
                        <label class="label has-text-white">Email</label>
                        <div class="control">
                            <input class="input is-outlined" type="email" name="email" placeholder="Enter email">
                        </div>
                    </div>

                    <div class="field">
                        <label class="label has-text-white">Password</label>
                        <div class="control">
                            <input class="input is-outlined" type="password" name="password" placeholder="Enter password">
                        </div>
                    </div>

                    <button class="button is-primary is-outlined" type="submit">Add User</button> 
                </form>
            </div>
        </div>
    </div>
</div>

</body>
</html>
`

var loginView = `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Login</title>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bulma@0.9.4/css/bulma.min.css">
  <script src="https://unpkg.com/htmx.org@1.9.5"></script>
  <style>
    body {
      background-color: #363636; 
    }
  </style>
</head>
<body>

<div class="container">
  <div class="columns is-centered">
    <div class="column is-half">
      <div class="box has-background-dark">
        <h2 class="title is-2 has-text-white">Login</h2>

        <form hx-post="/login" hx-swap="outerHTML">
          <div class="field">
            <label class="label has-text-white">Username</label>
            <div class="control">
              <input class="input is-outlined" type="text" name="username" placeholder="Enter your username">
            </div>
          </div>

          <div class="field">
            <label class="label has-text-white">Password</label>
            <div class="control">
              <input class="input is-outlined" type="password" name="password" placeholder="Enter your password">
            </div>
          </div>

          <button class="button is-primary is-outlined" type="submit">Login</button>
        </form>
      </div>
    </div>
  </div>
</div>

</body>
</html>
`

var chatView = `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Chat Interface</title>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bulma@0.9.4/css/bulma.min.css">
  <script src="https://unpkg.com/htmx.org@1.9.5"></script>
  <script src="https://unpkg.com/htmx.org/dist/ext/ws.js"></script>  </style>
  <style>
    body {
      background-color: #363636; /* Dark background */
    }
  </style>
</head>
<body>

<div class="container">
  <div class="columns is-mobile">
    <div class="column is-one-quarter">
      <aside class="menu">
        <p class="menu-label">Menu</p>
        <ul class="menu-list">
          <li><a>Settings</a></li>
          <li><a>Profile</a></li>
          <li><a hx-post="/logout" hx-swap="outerHTML">Logout</a></li>
        </ul>
      </aside>
    </div>

    <div class="column">
      <div class="box has-background-black" hx-ext="ws" hx-trigger="load, every 1s" ws-connect="/ws/%v" hx-target="#chat-box" hx-swap="outerHTML">
      <div id="chat-box"></div>
        </div>

      <form class="field" hx-post="/send-message" hx-trigger="submit" hx-swap="none">
        <div class="control is-expanded">
          <input class="input is-outlined" type="text" name="message" placeholder="Type your message...">
        </div>
        <div class="control">
          <button class="button is-info is-outlined" type="submit">Send</button>
        </div>
        <input type="hidden" name="roomid" value="%v">
      </form>
    </div>
  </div>
</div>

</body>
</html>
`
