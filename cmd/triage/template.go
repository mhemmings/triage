package main

// templateData defines the data availiable to a template.
type templateData struct {
	Repos []repo
}

// htmltemplate is the template to use when serving the HTML version of the issue triage.
var htmltemplate = `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <title>Issue Triage</title>
  <meta content="width=device-width" name="viewport">
  <link href="https://assets.ubuntu.com/v1/vanilla-framework-version-1.8.0.min.css" rel="stylesheet">
  <style>
    body {
      background-color: #f7f7f7;
    }
    .p-navigation__tagline {
      display: block;
      font-size: 1.5rem;
      padding: 0.85rem 0.5rem;
    }
    h2 {
      margin-top: 2rem;
      margin-bottom: 1rem;
    }
    .row {
      max-width: 1000px;
    }
    .p-card__title {
      max-width: 100%;
    }
    .refresh-form {
      margin-top: 0.3rem;
      right: 1.5rem;
      position: absolute;
    }
    .refresh-form i {
      margin-top: 0.7rem;
    }
  </style>
</head>
<body>
  <header class="p-navigation" id="navigation">
    <div class="p-navigation__banner">
      <div class="p-navigation__logo">
        <span class="p-navigation__tagline">Issue Triage</span>
        <form class="u-align--right refresh-form" action="/refresh" method="post" id="refresh" >
          <input class="p-button--neutral u-align--right" value="Refresh" type="submit">
        </form>
      </div>
  </header>

  {{ range $repoI, $repo := .Repos }}
  <div class="row">
    <div class="col-12">
      <h2><a href="{{$repo.IssuesLink}}" target="_blank">{{$repo.FullName}}</a></h2>
    </div>
    <div class="col-12">
      {{ range $issueI, $issue := $repo.Issues }}
      <div class="p-card--highlighted">
        <h4 class="p-card__title"><a href="{{$issue.Link}}" target="_blank">{{$issue.Title}}</a></h4>
        <p class="p-card__content">
          {{$issue.Created.Format "02 Jan 2006" }} |
          <i class="p-icon--user"></i> {{$issue.User}} |
          {{$issue.Comments}} comments
        </p>
      </div>
      {{end}}
      <hr>
    </div>
  </div>
  {{end}}

  <script>
    var refresh = document.getElementById('refresh');
    refresh.addEventListener('submit', function(){
      refresh.innerHTML = '<i class="p-icon--spinner u-animation--spin"></i>';
    });
  </script>

</body>
</html>
`
