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
