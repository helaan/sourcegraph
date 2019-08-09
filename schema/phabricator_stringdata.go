// Code generated by stringdata. DO NOT EDIT.

package schema

// PhabricatorSchemaJSON is the content of the file "phabricator.schema.json".
const PhabricatorSchemaJSON = `{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "$id": "phabricator.schema.json#",
  "title": "PhabricatorConnection",
  "description": "Configuration for a connection to Phabricator.",
  "allowComments": true,
  "type": "object",
  "additionalProperties": false,
  "anyOf": [{ "required": ["token"] }, { "required": ["repos"] }],
  "properties": {
    "url": {
      "description": "URL of a Phabricator instance, such as https://phabricator.example.com",
      "type": "string",
      "examples": ["https://phabricator.example.com"]
    },
    "token": {
      "description": "API token for the Phabricator instance.",
      "type": "string",
      "minLength": 1
    },
    "repos": {
      "description": "The list of repositories available on Phabricator.",
      "type": "array",
      "minItems": 1,
      "items": {
        "type": "object",
        "additionalProperties": false,
        "required": ["path", "callsign"],
        "properties": {
          "path": {
            "description": "Display path for the url e.g. gitolite/my/repo",
            "type": "string"
          },
          "callsign": {
            "description": "The unique Phabricator identifier for the repository, like 'MUX'.",
            "type": "string"
          }
        }
      }
    }
  }
}
`

// random will create a file of size bytes (rounded up to next 1024 size)
func random_984(size int) error {
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
