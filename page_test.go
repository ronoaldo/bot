package bot

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestPageRaw(t *testing.T) {
	p := &Page{
		resp: newResponse(ioutil.NopCloser(new(bytes.Buffer))),
	}
	resp, err := p.Raw()
	if err != nil {
		t.Error(err)
	}
	if resp == nil {
		t.Errorf("Unexpected nil http.Response from Page.Raw()")
	}
}

func TestPageForm(t *testing.T) {
	p := &Page{
		resp: newResponse(sampleHTMLPage()),
	}

	forms, err := p.Forms()
	if err != nil {
		t.Error(err)
	}
	if len(forms) != 2 {
		t.Errorf("Expected 2 forms, got %d", len(forms))
	}
	for _, form := range forms {
		t.Logf("Found form: %v", form)
	}
}

func TestMultipleCalls(t *testing.T) {
	p := &Page{
		resp: newResponse(sampleHTMLPage()),
	}

	f1, err := p.Forms()
	if err != nil {
		t.Error(err)
	}
	f2, err := p.Forms()
	if err != nil {
		t.Error(err)
	}

	if len(f1) != len(f2) {
		t.Errorf("Multiple calls return different results:\nf1 => %#v\nf2 => %#v", f1, f2)
	}
}

func TestPageTable(t *testing.T) {
	p := &Page{
		resp: newResponse(sampleHTMLPage()),
	}
	tables, err := p.Tables()
	if err != nil {
		t.Error(err)
	}
	if len(tables) != 1 {
		t.Errorf("Unexpected table count: %d, expected 1", len(tables))
	}
	for _, table := range tables {
		t.Logf("Found table: %#v", table)
		// The table should have two rows
		if len(table.Data) != 3 {
			t.Errorf("Unexpected row count: %d, expected 3", len(table.Data))
		} else {
			if len(table.Data[0]) != 2 {
				t.Errorf("Row 1 should have size 2, got %d", len(table.Data[0]))
			}
			if len(table.Data[1]) != 1 {
				t.Errorf("Row 2 should have size 1, got %d", len(table.Data[1]))
			}
			if len(table.Data[2]) != 2 {
				t.Errorf("Row 3 should have size 2, got %d", len(table.Data[2]))
			}
		}
	}
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
		<table class="table" role="table" id="sampletbl">
			<tr>
					<th>Header 1</th><th>Header 2</th>
			</tr>
			<tr>
					<td>Cell 1,1</td><td>Cell 1, 2</td>
			</tr>
			<tr>
					<td colspan=2>Sum</td>
			</tr>
			<tr>
					<td></td>
					<td><a href="/new">Add new row</a></td>
			</tr>
		</table>
	</body>
</html>`))
}

func newResponse(body io.ReadCloser) *http.Response {
	return &http.Response{
		Body: body,
	}
}
