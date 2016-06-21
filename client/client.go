package freeling

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log"
	"fmt"
	"net"
	"sync"
)

const (
	flServerReady = "FL-SERVER-READY"
)

type Client struct {
	conn net.Conn
	br   *bufio.Reader
	Debug bool
	mu sync.Mutex
}

func New(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("net.Dial(_, %q) failed: %v", addr, err)
	}
	c := &Client{
		conn: conn,
		br:   bufio.NewReader(conn),
	}

	if _, err = c.Send("RESET_STATS"); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

type Token struct {
	CTag   string
	Case   string
	Form   string
	Gen    string
	ID     string
	Lemma  string
	Mood   string
	Num    string
	POS    string
	Person string
	Tag    string
	Tense  string
	Type   string
	WN string
	NEClass string
	NEC string
}

func (t Token) String() string {
	return fmt.Sprintf("%s %s %s", t.Form, t.Lemma, t.Tag)
}

// Phrase ("sentence")
type Phrase struct {
	ID     string
	Tokens []Token
}

func (p Phrase) String() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "Phrase %s:\n", p.ID)
	for _, t := range p.Tokens {
		fmt.Fprintf(&b, "  %s\n", t)
	}
	return b.String()
}

// Document ("paragraph")
type Document struct {
	Phrases []Phrase
}

func (d Document) String() string {
	var b bytes.Buffer
	for _, p := range d.Phrases {
		b.WriteString(p.String())
	}
	return b.String()
}

func (c *Client) Process(msg string) (*Document, error) {
	s, err := c.Send(msg)
	if err != nil {
		return nil, err
	}
	d := &Document{
		Phrases: make([]Phrase, 0),
	}
	if err := json.Unmarshal([]byte("["+s+"]"), &d.Phrases); err != nil {
		return nil, err
	}
	return d, nil
}

func (c *Client) Send(msg string) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.write(msg); err != nil {
		return "", err
	}

	ret := ""

	res, err := c.read()
	if err != nil {
		return "", err
	}
	if res != flServerReady {
		ret = res
	}

	if err := c.write("FLUSH_BUFFER"); err != nil {
		return "", err
	}

	res, err = c.read()
	if err != nil {
		return "", err
	}
	if res != flServerReady {
		ret += res
	}
	return ret, nil
}

func (c *Client) write(msg string) error {
	if c.Debug {
		log.Printf(">>> %q", msg)
	}
	_, err := fmt.Fprint(c.conn, msg+"\u0000")
	return err
}

func (c *Client) read() (string, error) {
	b, err := c.br.ReadBytes(0)
	if err != nil {
		return "", err
	}
	s := string(b[:len(b)-1])
	if c.Debug {
		log.Printf("<<< %q", s)
	}
	return s, err
}
