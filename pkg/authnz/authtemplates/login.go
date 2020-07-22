package authtemplates

import (
	"html/template"
	"io"

	"github.com/place1/wg-access-server/pkg/authnz/authruntime"
	"github.com/place1/wg-access-server/pkg/authnz/authsession"
)

type LoginPage struct {
	Title     string
	Providers []*authruntime.Provider
	Banner    *authsession.Banner
}

func RenderLoginPage(w io.Writer, data LoginPage) error {
	tpl, err := template.New("login-page").Parse(loginPage)
	if err != nil {
		return err
	}
	return tpl.Execute(w, data)
}

const loginPage string = `
<style>
	* {
		font-family: 'Open Sans', -apple-system, BlinkMacSystemFont,
			"Segoe UI", Roboto, Helvetica, Arial, sans-serif,
			"Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol";
		font-size: 16px;
		-webkit-font-smoothing: antialiased;
		-moz-osx-font-smoothing: grayscale;
	}

	body {
		background: #1e6cc9;
		background: -webkit-linear-gradient(0deg, #3357cc 0%, #1e6cc9 100%);
	}

	.form {
		position: absolute;
		top: 50%;
		left: 50%;
		transform: translate(-50%, -50%);
		background-color: #fff;
		width: 285px;
		padding: 40px;
		box-shadow: 0 10px 20px rgba(0, 0, 0, 0.19), 0 6px 6px rgba(0, 0, 0, 0.23);
	}

	.form h2 {
		margin: 0 0 20px;
		text-align: center;
		line-height: 1;
		color: black;
		font-size: 22px;
		font-weight: 400;
	}

	.form label {
		font-size: 12px;
	}

	.form input {
		font-size: 14px;
		outline: none;
		display: block;
		width: 100%;
		padding: 8px 12px;
		border: 1px solid #ccc;
		border-radius: 3px;
		box-sizing: border-box;
		margin-bottom: 18px;
	}

	.form a {
		display: block;
	}

	.form > * {
		margin: 0 0 20px;
	}

	.form > *:last-child {
		margin-bottom: 0px;
	}

	.form input:focus {
		//color: #333;
		border: 1px solid #44c4e7;
	}

	.form button {
		position: relative;
		cursor: pointer;
		width: 100%;
		padding: 10px 34px;
		border-radius: 3px;
		border: 0;
		background: #44c4e7;
		color: white;
	}

	.form button img {
		position: absolute;
		left: 8px;
		top: 50%;
		transform: translateY(-50%);
		height: 18px;
	}

	.form button:hover {
		background: #369cb8;
	}

	.form hr {
		position: relative;
		width: 55%;
		margin-left: auto;
		margin-right: auto;
		overflow: visible;
		background-color: #4d4d4d;
	}

	.form hr:after {
		content: "";
		position: absolute;
		left: 50%;
		top: 50%;
		transform: translate(-50%, -50%);
		width: 4px;
		height: 4px;
		background-color: #4d4d4d;
		border-radius: 50%;
	}

	.banner {
		padding: 8px;
		font-size: 12px;
		margin-bottom: 16px;
	}

	.banner.danger {
		background: #d65463;
		color: white;
	}
</style>

<section class="form">
	<h2>{{.Title}}</h2>

	{{range $i, $p := .Providers}}
		{{if eq $p.Type "Basic"}}
			<form autocomplete="off" method="post" action="/signin/{{$i}}">
				{{if $.Banner}}
					<div
						class="banner {{$.Banner.Intent}}"
					>
						{{$.Banner.Text}}
					</div>
				{{end}}
				<div>
					<label for="username">username</label>
					<input id="username" name="username" placeholder="Username" type="text"></input>
				</div>
				<div>
					<label for="password">password</label>
					<input id="password" name="password" placeholder="Password" type="password"></input>
				</div>
				<button id="submit">Sign in</button>
			</form>

			{{ $length := len $.Providers }}
			{{ if gt $length 1 }}
				<hr />
			{{end}}

		{{else}}

			<a href="/signin/{{$i}}">
				<button
					style="
						{{if .Branding.Background}}
							background: {{.Branding.Background}}
						{{else}}
							background: #44c4e7;
						{{end}}
					"
				>
					{{if .Branding.Icon}}
						<img src="{{.Branding.Icon}}"></img>
					{{end}}
					<span
						style="
							{{if .Branding.Color}}
								color: {{.Branding.Color}}
							{{else}}
								color: white;
							{{end}}
						"
					>
						{{$p.Name}}
					</span>
				</button>
			</a>
		{{end}}
	{{end}}

</section>
`
