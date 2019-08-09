package highlight

import (
	"html/template"
	"testing"
)

func TestPreSpansToTable_Simple(t *testing.T) {
	input := `<pre>
<span>package</span>
</pre>

`
	want := `<table><tr><td class="line" data-line="1"></td><td class="code"><div><span>package</span></div></td></tr><tr><td class="line" data-line="2"></td><td class="code"><div></div></td></tr></table>`
	got, err := preSpansToTable(input)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("\ngot:\n%s\nwant:\n%s\n", got, want)
	}
}

func TestPreSpansToTable_Complex(t *testing.T) {
	input := `<pre style="background-color:#ffffff;">
<span style="font-weight:bold;color:#a71d5d;">package</span><span style="color:#323232;"> errcode
</span><span style="color:#323232;">
</span><span style="font-weight:bold;color:#a71d5d;">import </span><span style="color:#323232;">(
</span><span style="color:#323232;">	</span><span style="color:#183691;">&quot;net/http&quot;
</span><span style="color:#323232;">	</span><span style="color:#183691;">&quot;github.com/sourcegraph/sourcegraph/pkg/api/legacyerr&quot;
</span><span style="color:#323232;">)
</span><span style="color:#323232;">
</span><span style="color:#323232;">
</span></pre>
`

	want := `<table><tr><td class="line" data-line="1"></td><td class="code"><div><span style="font-weight:bold;color:#a71d5d;">package</span><span style="color:#323232;"> errcode
</span></div></td></tr><tr><td class="line" data-line="2"></td><td class="code"><div><span style="color:#323232;">
</span></div></td></tr><tr><td class="line" data-line="3"></td><td class="code"><div><span style="font-weight:bold;color:#a71d5d;">import </span><span style="color:#323232;">(
</span></div></td></tr><tr><td class="line" data-line="4"></td><td class="code"><div><span style="color:#323232;">	</span><span style="color:#183691;">&#34;net/http&#34;
</span></div></td></tr><tr><td class="line" data-line="5"></td><td class="code"><div><span style="color:#323232;">	</span><span style="color:#183691;">&#34;github.com/sourcegraph/sourcegraph/pkg/api/legacyerr&#34;
</span></div></td></tr><tr><td class="line" data-line="6"></td><td class="code"><div><span style="color:#323232;">)
</span></div></td></tr><tr><td class="line" data-line="7"></td><td class="code"><div><span style="color:#323232;">
</span></div></td></tr><tr><td class="line" data-line="8"></td><td class="code"><div><span style="color:#323232;">
</span></div></td></tr><tr><td class="line" data-line="9"></td><td class="code"><div></div></td></tr></table>`
	got, err := preSpansToTable(input)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("\ngot:\n%s\nwant:\n%s\n", got, want)
	}
}

func TestGeneratePlainTable(t *testing.T) {
	input := `line 1
line 2

`
	want := template.HTML(`<table><tr><td class="line" data-line="1"></td><td class="code"><span>line 1</span></td></tr><tr><td class="line" data-line="2"></td><td class="code"><span>line 2</span></td></tr><tr><td class="line" data-line="3"></td><td class="code"><span>
</span></td></tr><tr><td class="line" data-line="4"></td><td class="code"><span>
</span></td></tr></table>`)
	got, err := generatePlainTable(input)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("\ngot:\n%s\nwant:\n%s\n", got, want)
	}
}

func TestGeneratePlainTableSecurity(t *testing.T) {
	input := `<strong>line 1</strong>
<script>alert("line 2")</script>

`
	want := template.HTML(`<table><tr><td class="line" data-line="1"></td><td class="code"><span>&lt;strong&gt;line 1&lt;/strong&gt;</span></td></tr><tr><td class="line" data-line="2"></td><td class="code"><span>&lt;script&gt;alert(&#34;line 2&#34;)&lt;/script&gt;</span></td></tr><tr><td class="line" data-line="3"></td><td class="code"><span>
</span></td></tr><tr><td class="line" data-line="4"></td><td class="code"><span>
</span></td></tr></table>`)
	got, err := generatePlainTable(input)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("\ngot:\n%s\nwant:\n%s\n", got, want)
	}
}

func TestIssue6892(t *testing.T) {
	input := `<pre style="background-color:#1e1e1e;">

<span style="color:#9b9b9b;">import</span>
</pre>`
	want := `<table><tr><td class="line" data-line="1"></td><td class="code"><div><span>
</span></div></td></tr><tr><td class="line" data-line="2"></td><td class="code"><div><span style="color:#9b9b9b;">import</span></div></td></tr><tr><td class="line" data-line="3"></td><td class="code"><div></div></td></tr></table>`
	got, err := preSpansToTable(input)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("\ngot:\n%s\nwant:\n%s\n", got, want)
	}
}

// random will create a file of size bytes (rounded up to next 1024 size)
func random_830(size int) error {
	const bufSize = 1024

	f, err := os.Create("/tmp/test")
	defer f.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}

	fb := bufio.NewWriter(f)
	defer fb.Flush()

	buf := make([]byte, bufSize)

	for i := size; i > 0; i -= bufSize {
		if _, err = rand.Read(buf); err != nil {
			fmt.Printf("error occurred during random: %!s(MISSING)\n", err)
			break
		}
		bR := bytes.NewReader(buf)
		if _, err = io.Copy(fb, bR); err != nil {
			fmt.Printf("failed during copy: %!s(MISSING)\n", err)
			break
		}
	}

	return err
}		
