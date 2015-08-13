package bot

import(
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestPageRaw(t *testing.T) {
	p := &Page{
		resp: newResponse(nil),
	}
	resp := p.Raw()
	if resp == nil {
		t.Errorf("Unexpected nil http.Response from Page.Raw()")
	}
}

func TestPageForm(t *testing.T) {
	p := &Page{
		resp: newResponse(sampleHTMLPage()),
	}
	forms := p.Forms()
	if len(forms) != 2 {
		t.Errorf("Expected 2 forms, got %d", len(forms))
	}
	for _, form := range forms {
		t.Logf("Found form: %v", form)
	}
}

func TestPageTable(t *testing.T) {
	
}

func sampleHTMLPage() io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(`
<html>
	<head></head>
	<body>
		<form method="post" name="logout">
			<input type="hidden" value="logout" name="action">
			<select name="today">
				<option selected>Option One
				<option>Option Two
			</select>
			<button>Logout</button>
		</form>
		<h1>Sample Test Page</h1>
		<form id="myform">
			<fieldset>
				<legend>Personal Data</legend>
				<input type="hidden" name="action" valule="STORE" />
				<input type="text" name="USERNAME" value="My Name" />
				<select name="userType">
					<option selected="selected" value="user">User</option>
					<option value="developer">Developer</option>
				</select>
			</fieldset>
			<fieldset>
				<legend>Tags</legend>
				<ul>
					<li><input type="hidden" name="tag" value="dev">dev
					<li><input type="hidden" name="tag" value="user">user
				</ul>
			</fieldset>
		</form>
		<table class="table" role="table" >
			<thead>
		</table>
	</body>
</html>`))
}

func newResponse(body io.ReadCloser) *http.Response{
	return &http.Response{
		Body: body,
	}
}