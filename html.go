package main

// type Style struct {
// 	Background string
// }

// var style = Style{
//   Background: "#0b141c",
// }

var authNotification = `<div class="notification %v" id="notty">
  <button class="delete" hx-get="/can" hx-trigger="click" hx-target="#notty" hx-swap="outerHTML"></button>
  %v
</div>`

var AdUserHTML = `<!DOCTYPE html>
<html>

<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>add user</title>
    <link rel="stylesheet" href="/static/bulma.min.css">
    <script src="/static/htmx.min.js"></script>
    <style>
        body {
            background-color: #0b141c;
        }
    </style>
</head>

<body>

    <div class="container">
        <div class="columns is-centered">
            <div class="column is-half">
                <div class="box has-background-black">
                    <h2 class="title is-2 has-text-info">add user</h2>

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
                                <input class="input is-outlined" type="password" name="password"
                                    placeholder="Enter password">
                            </div>
                        </div>

                        <button class="button is-info is-outlined" type="submit">add user</button>
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
  <title>login</title>
  <link rel="stylesheet" href="/static/bulma.min.css">
  <script src="/static/htmx.min.js"></script>
  <style>
    body {
      background-color: #0b141c;
    }
    .animate-spin {
  animation: spin 1s linear infinite;
}
@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);

  }
}
  </style>
</head>

<body>

  <div class="container">
    <div class="columns is-centered">
      <div class="column is-half">
        <div class="box has-background-black">
          <h2 class="title is-2 has-text-info">login</h2>

          <form hx-post="/login" hx-swap="outerHTML" class="has-background-black">
            <div class="field">
              <label class="label has-text-white">username</label>
              <div class="control">
                <input class="input is-outlined" type="text" name="username" placeholder="Enter your username">
              </div>
            </div>

            <div class="field">
              <label class="label has-text-white">password</label>
              <div class="control">
                <input class="input is-outlined" type="password" name="password" placeholder="Enter your password">
              </div>
            </div>
            <div>
            <button class="button is-info is-outlined" type="submit" hx-indicator="#progressMeter">login</button>
            <svg xmlns="http://www.w3.org/2000/svg" width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="#313e97" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="htmx-indicator animate-spin">
  <circle cx="12" cy="12" r="10"></circle>
</svg>
            </div>
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
  <title>screamery</title>
  <link rel="stylesheet" href="/static/bulma.min.css">
  <script src="/static/htmx.min.js"></script>
  <script src="https://unpkg.com/htmx.org/dist/ext/ws.js"></script>
  <script>
    document.addEventListener("htmx:wsAfterMessage", function (event) {
      var chatBox = document.getElementById("chat-box");
      chatBox.scrollTop = chatBox.scrollHeight;
    });
    
  </script>
  <style>
    .mydisplay {
      height: 500px;
      overflow-y: scroll;
    }

    body {
      background-color: #0b141c;
      /* Dark blue background*/
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
          <p class="menu-label">menu</p>
          <ul class="menu-list">
          <li><a href="/help" target="_blank" rel="noopener noreferrer" class="has-text-info">help</a></li>
            <li><a hx-post="/logout" class="has-text-info">logout</a></li>
            <li><a href="/add-room" class="has-text-info">add room</a></li>
            <li><a href="/profile" hx-boost="true" class="has-text-info">profile</a></li>
            <li><a hx-post="/history" class="has-text-info">history</a></li>
            <li><a hx-post="/rooms" class="has-text-info">rooms</a></li>
            <li><a href="/add-post" class="has-text-info">add post</a></li>
          </ul>
        </aside>
      </div>

      <div class="column">
        <div class="box has-background-black" hx-ext="ws" ws-connect="/ws/%v" hx-target="#chat-box" hx-swap="outerHTML">
          <div id="chat-box" hx-get="/messagehist/%v" hx-trigger="load" class="mydisplay"></div>
        </div>

        <form class="field" hx-post="/send-message" hx-trigger="submit" hx-swap="none">
          <div class="control is-expanded">
            <input class="input is-outlined" type="text" id="messageBox" name="message"
              placeholder="Type your message...">
          </div>
          <div class="control">
            <button class="button is-info is-outlined" type="submit">send</button>
            <button class="button is-info is-outlined" hx-target="#roomstats" hx-get="/roomstats" hx-trigger="click"
              hx-swap="outerHTML">room</button>
          </div>
          <input type="hidden" name="roomid" value="%v">
        </form>
      </div>
    </div>
    <div class="columns is-mobile">
      <div class="column is-is-one-quarter">
        <div class="content" id="roomstats">
          room: <strong>%v</strong>
        </div>

      </div>

    </div>
  </div>

</body>

</html>`

var editProfileView = `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>edit profile</title>
  <link rel="stylesheet" href="/static/bulma.min.css">
  <script src="/static/htmx.min.js"></script>
  <style>
  body {
      background-color: #0b141c; /* Dark background */
    }
  </style>
</head>
<body class="has-background-black">
  <section class="section">
    <div class="container">
      <h1 class="title has-text-info">edit profile</h1>
      <p class="has-text-info-light"> be careful not to leave any fields empty when saving.</p>

      <form hx-post="/update-profile" class="box has-background-black" hx-trigger="submit">
        <div class="field">
          <label class="label has-text-white">Email</label>
          <div class="control">
            <input class="input" type="email" name="email" placeholder="you@example.com" required value="%v">
          </div>
        </div>

        <div class="field">
          <label class="label has-text-white">First Name</label>
          <div class="control">
            <input class="input" type="text" name="first_name" placeholder="big" value="%v">
          </div>
        </div>

        <div class="field">
          <label class="label has-text-white">Last Name</label>
          <div class="control">
            <input class="input" type="text" name="last_name" placeholder="mac" value="%v">
          </div>
        </div>
        <div class="field">
          <label class="label has-text-white">about</label>
          <div class="control">
            <textarea maxlength="200" class="textarea" name="about" placeholder="about me...">%v</textarea>
          </div>
        </div>

        <button type="submit" class="button is-warning is-outlined">save changes</button>
        <button class="button is-primary is-outlined" href="#" type="button" onclick="history.back()">cancel</button>
        <input type="hidden" name="userid" value="%v">
      </form>
    </div>
  </section>
</body>
</html>
`
var clearAuthNotification = `<form hx-post="/login" hx-swap="outerHTML" class="has-background-black">
  <div class="field">
    <label class="label has-text-white">username</label>
    <div class="control">
      <input class="input is-outlined" type="text" name="username" placeholder="Enter your username">
    </div>
  </div>

  <div class="field">
    <label class="label has-text-white">password</label>
    <div class="control">
      <input class="input is-outlined" type="password" name="password" placeholder="Enter your password">
    </div>
  </div>

  <button class="button is-info is-outlined" type="submit">login</button>
  </form>
`

var addRoomView = `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>add room</title>
    <link rel="stylesheet" href="/static/bulma.min.css">
    <script src="/static/htmx.min.js"></script> <style>
        body {
            background-color: #0b141c; 
        }
    </style>
</head>
<body>

<div class="container">
    <div class="columns is-centered">
        <div class="column is-half">
            <div class="box has-background-black">
                <h2 class="title is-2 has-text-info">add user</h2>

                <form hx-post="/addroom" hx-swap="outerHTML">  
                    <div class="field">
                        <label class="label has-text-white">room name</label>
                        <div class="control">
                            <input class="input is-outlined" type="text" name="room" placeholder="room name">
                        </div>
                    </div>

                    <button class="button is-info is-outlined" type="submit">add room</button> 
                </form>
            </div>
        </div>
    </div>
</div>

</body>
</html>`

var addPostView = `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>add room</title>
    <link rel="stylesheet" href="/static/bulma.min.css">
    <script src="/static/htmx.min.js"></script> <style>
        body {
            background-color: #0b141c; 
        }
    </style>
</head>
<body>

<div class="container">
    <div class="columns is-centered">
        <div class="column is-half">
            <div class="box has-background-black">
                <h2 class="title is-2 has-text-info">add post</h2>

                <form hx-post="/addpost" hx-swap="outerHTML">  
                    <div class="field">
                        <label class="label has-text-white">add post</label>
                        <div class="control">
                            <textarea maxlength="200" class="textarea" name="post" id="post" placeholder="whats up, buttercup?"></textarea>
                        </div>
                    </div>

                    <button class="button is-info is-outlined" type="submit">post</button> 
                </form>
            </div>
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
  <title>user</title>
  <link rel="stylesheet" href="/static/bulma.min.css">
  <script src="/static/htmx.min.js"></script>
  <style>
    body {
      background-color: #0b141c;
    }
  </style>
</head>

<body>

  <div class="container">
    <div class="columns is-centered">
      <div class="column is-half">
        <div class="box has-background-black">
          %v
        </div>
        <div class="box has-background-black">
          %v
        </div>
      </div>
    </div>
  </div>

</body>

</html>
`

var splashView = `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>welcome!</title>
    <link rel="stylesheet" href="/static/bulma.min.css">
    <script src="/static/htmx.min.js"></script> <style>
        body {
            background-color: #0b141c; 
        }
    </style>
</head>
<body>

<div class="container">
    <div class="columns is-centered">
        <div class="column is-half">
            <div class="box has-background-black">
                <h2 class="title is-2 has-text-info">hi there</h2>
                <p class="subtitle has-text-info-light">welcome to the conversation. please do not be vulgar</p>
                <p class="subtitle has-text-info-light">this is a small chat application for a few trusted users. navigate to any room by going to https://url.com/<strong>room/literallywhatever</strong></p>
                <a href="/room/welcome" class="button is-info">enter the welcome chat</a>
                <a href="/help" class="button is-info">get help</a>
            </div>
        </div>
    </div>
</div>

</body>
</html>
`
