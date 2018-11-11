package cmd

import (
	"bufio"
	"bytes"
	"log"
	"os"
)

func ParseGoFile(fp string) (_tags []tag) {
	var (
		sq bool // inside of single-quote
		dq bool // inside of double-quote
		bt bool // inside of back-tick
		cl bool // inside of comment line
		cb bool // inside of comment block

		comment []byte
		n       int
		l       int
	)

	f, err := os.Open(fp)
	if err != nil {
		log.Fatal(err)
	}

	rd := bufio.NewReader(f)

	for {
		b, err := rd.ReadByte()
		if err != nil {
			break
		}
		n++

		if b == '\n' {
			l++
		}

		if !cb && !cl {
			// in single quote literal
			if b == '\'' {
				if sq {
					sq = false
				} else if !dq && !bt {
					sq = true
				}
				continue
			}

			// in double quote literal
			if b == '"' {
				if dq {
					dq = false
				} else if !sq && !bt {
					dq = true
				}
				continue
			}

			// in backtick literal
			if b == '`' {
				if bt {
					bt = false
				} else if !sq && !dq {
					bt = true
				}
				continue
			}
		}

		// continue if inside literal
		if sq || dq || bt {
			continue
		}

		// in comment line
		if cl {
			if b == '\n' {
				cl = false
				// log.Println("end comment line", n)
				if tt := parseGoComment(comment, fp, l); len(tt) > 0 {
					_tags = append(_tags, tt...)
					comment = []byte{}
				}
			} else {
				comment = append(comment, b)
			}
		}

		// in comment block
		if cb {
			// find end of comment block
			if b == '*' {
				_b, err := rd.Peek(1)
				if err != nil {
					break
				}

				if _b[0] == '/' {
					cb = false
					// log.Println("end comment block", n)
					if tt := parseGoComment(comment, fp, l); len(tt) > 0 {
						_tags = append(_tags, tt...)
						comment = []byte{}
					}
				} else {
					comment = append(comment, b)
				}
			} else {
				comment = append(comment, b)
			}
		}

		// find forward slash
		if b == '/' {
			// find second forward slash
			_b, err := rd.Peek(1)
			if err != nil {
				break
			}

			switch _b[0] {
			// find double slash comment
			case '/':
				cl = true
				// log.Println("start comment line", n)
			// find comment block
			case '*':
				cb = true
				// log.Println("start comment block", n)
			}

			if _b[0] == '*' {
				if !cb {
					cb = true
				}
			}
		}
	}

	return _tags
}

func parseGoComment(b []byte, fp string, l int) (_tags []tag) {
	lines := bytes.Split(b, []byte{'\n'})

	for _, line := range lines {
		if i := bytes.Index(line, []byte("TODO:")); i > -1 {
			_tags = append(_tags, tag{
				file:    fp[len(wd):],
				tag:     "TODO",
				line:    l,
				message: string(bytes.Trim(line[i+5:], " \n\t")),
			})
		}

		if i := bytes.Index(line, []byte("FIXME:")); i > -1 {
			_tags = append(_tags, tag{
				file:    fp[len(wd):],
				tag:     "FIXME",
				line:    l,
				message: string(bytes.Trim(line[i+6:], " \n\t")),
			})
		}

	}

	return _tags
}
