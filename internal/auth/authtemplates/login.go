package authtemplates

import (
	"html/template"
	"io"

	"github.com/place1/wireguard-access-server/internal/auth/authconfig"
)

type LoginPage struct {
	Config *authconfig.AuthConfig
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
		font-family: monospace;
		font-size: 16px;
		-webkit-font-smoothing: antialiased;
		-moz-osx-font-smoothing: grayscale;
	}

	body {
		background-color: #44c4e7;
	}

	.form {
		position: absolute;
		top: 40%;
		left: 50%;
		background-color: #fff;
		width: 285px;
		margin: -140px 0 0 -182px;
		padding: 40px;
		box-shadow: 0 0 3px rgba(0, 0, 0, 0.3);
	}

	.form h2 {
		margin: 0 0 20px;
		line-height: 1;
		color: #44c4e7;
		font-size: 22px;
		font-weight: 400;
	}

	.form input {
		outline: none;
		display: block;
		width: 100%;
		padding: 10px 15px;
		border: 1px solid #ccc;
		color: #ccc;
		box-sizing: border-box;
		transition: 0.2s linear;
	}

	.form * {
		margin: 0 0 20px;
	}

	.form input:focus {
		color: #333;
		border: 1px solid #44c4e7;
	}

	.form button {
		cursor: pointer;
		background: #44c4e7;
		width: 100%;
		padding: 10px 15px;
		border: 0;
		color: #fff;
		text-transform: capitalize;
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

	.error, .valid{display:none;}
</style>



<section class="form animated flipInX">
  <h2>Login To Your Account</h2>
  <p class="valid">Valid. Please wait a moment.</p>
  <p class="error">Error. Please enter correct Username &amp; password.</p>
	<form class="loginbox" autocomplete="off">

		{{if .Config.Basic}}
			<input placeholder="Username" type="text" id="username"></input>
			<input placeholder="Password" type="password" id="password"></input>
			<button id="submit">Login</button>
		{{end}}

		<hr />

		{{if .Config.OIDC}}
			<button>{{.Config.OIDC.Name}}</button>
		{{end}}

		{{if .Config.Gitlab}}
			<button>{{.Config.Gitlab.Name}}</button>
		{{end}}

	</form>
</section>
`
