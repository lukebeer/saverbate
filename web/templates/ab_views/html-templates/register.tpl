<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="RATING" content="RTA-5042-1996-1400-1577-RTA">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <meta name="description" content="Free webcam records from chaturbate performers">
    <title>Saverbate - free webcam records - registration</title>

    <meta name="theme-color" content="#2f3135">

    <link rel="stylesheet" href="/static/dist/main.css">
</head>
<body>
  <div class="container">
    <div class="row">
      <div class="col">
        <div class="site-header">
          <h1 class="site-logo"><span class="star">★</span>&nbsp;Saverbate&nbsp;<span class="star">★</span></h1>
          <h5 class="site-logo-subtitle">Don't miss the tastiest</h5>
        </div>


      </div>
    </div>

		<div class="row">
			<div class="col-sm">&nbsp;</div>
			<div class="col-sm-8">
				<form action="{{mountpathed "register"}}" method="post">
					<div class="form-group">
						{{with .errors}}{{with (index . "")}}{{range .}}<span>{{.}}</span><br />{{end}}{{end}}{{end -}}
						<label for="name">Name:</label>
						<input class="form-control" name="name" type="text" value="{{with .preserve}}{{with .name}}{{.}}{{end}}{{end}}" placeholder="Name" />
					</div>

					<div class="form-group">
						{{with .errors}}{{range .name}}<span>{{.}}</span><br />{{end}}{{end -}}
						<label for="email">E-mail:</label>
						<input class="form-control" name="email" type="text" value="{{with .preserve}}{{with .email}}{{.}}{{end}}{{end}}" placeholder="E-mail" />
					</div>

					<div class="form-group">
						{{with .errors}}{{range .email}}<span>{{.}}</span><br />{{end}}{{end -}}
						<label for="password">Password:</label>
						<input class="form-control" name="password" type="password" placeholder="Password" />
					</div>

					<div class="form-group">
						{{with .errors}}{{range .password}}<span>{{.}}</span><br />{{end}}{{end -}}
						<label for="confirm_password">Confirm Password:</label>
						<input class="form-control" name="confirm_password" type="password" placeholder="Confirm Password" />
					</div>

					{{with .errors}}{{range .confirm_password}}<span>{{.}}</span><br />{{end}}{{end -}}
					<button type="submit" class="btn btn-primary">Register</button>

					<a href="/">Cancel</a>

					{{with .csrf_token}}<input type="hidden" name="csrf_token" value="{{.}}" />{{end}}
				</form>
			</div>
			<div class="col-sm">&nbsp;</div>
		</div>
  </div>
	<footer>
    <div class="container">
      <div class="row">
        <div class="col">
          <small>
            Saverbate — The Chaturbate Archive.
            Chaturbate is an adult website providing live webcam
            performances by amateur camgirls, camboys and couples typically featuring nudity and sexual activity
            ranging from striptease and dirty talk to masturbation with sex toys.
            "Chaturbate" is a portmanteau of "chat" and "masturbate".
            Saverbate records your favorite live adult webcam broadcasts
            making by your lovely performers from Chaturbate.com to watch it later.
          </small>
        </div>
      </div>
    </div>
  </footer>
</body>
</html>
