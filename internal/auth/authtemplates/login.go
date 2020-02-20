package authtemplates

import (
	"html/template"
	"io"

	"github.com/place1/wg-access-server/internal/auth/authruntime"
)

type LoginPage struct {
	Providers []*authruntime.Provider
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
		background-color: #3899c9;
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
		margin: 0 0 35px;
		text-align: center;
		line-height: 1;
		color: black;
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
</style>

<section class="form">
  <h2>Sign In</h2>

	{{range $i, $p := .Providers}}
		<a href="/signin/{{$i}}">
			<button>{{$p.Type}}</button>
		</a>
	{{end}}

	<!--
	<form autocomplete="off">
		<input placeholder="Username" type="text" id="username"></input>
		<input placeholder="Password" type="password" id="password"></input>
		<button id="submit">Login</button>
	</form>
	<hr />
	-->

</section>
`
