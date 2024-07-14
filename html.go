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
      @media (max-width: 768px) { 
      .column.is-one-quarter {
        display: none; 
      }
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
          <li><a class="has-text-link">Settings</a></li>
          <li><a href="/profile" hx-boost="true" class="has-text-link">Profile</a></li>
          <li><a hx-post="/logout" hx-swap="outerHTML" class="has-text-link">Logout</a></li>
        </ul>
      </aside>
    </div>

    <div class="column">
      <div class="box has-background-black" hx-ext="ws" hx-trigger="load, every 1s" ws-connect="/ws/%v" hx-target="#chat-box" hx-swap="outerHTML">
      <div id="chat-box" hx-get="/messagehist/%v" hx-trigger="load"></div>
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
var profileView = `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Edit Profile</title>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bulma@0.9.4/css/bulma.min.css">
  <script src="https://unpkg.com/htmx.org@1.9.5"></script>
</head>
<body class="has-background-dark">
  <section class="section">
    <div class="container">
      <h1 class="title has-text-white">Edit Profile</h1>

      <form hx-post="/update-profile" class="box has-background-dark" hx-trigger="submit">
        <div class="field">
          <label class="label has-text-white">Email</label>
          <div class="control">
            <input class="input" type="email" name="email" placeholder="you@example.com" required>
          </div>
        </div>

        <div class="field">
          <label class="label has-text-white">First Name</label>
          <div class="control">
            <input class="input" type="text" name="first_name" placeholder="John">
          </div>
        </div>

        <div class="field">
          <label class="label has-text-white">Last Name</label>
          <div class="control">
            <input class="input" type="text" name="last_name" placeholder="Doe">
          </div>
        </div>

        <div class="field">
          <label class="label has-text-white">Password</label>
          <div class="control">
            <input class="input" type="password" name="password" placeholder="*********">
          </div>
        </div>

        <button type="submit" class="button is-primary is-link">Save Changes</button>
        <input type="hidden" name="userid" value="%v">
      </form>
    </div>
  </section>
</body>
</html>
`
