package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title string
	desc  string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	list list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

var (
	certPath string
	domain   string
	cert     *x509.Certificate
	header   string
)

func main() {
	flag.StringVar(&certPath, "file", "", "Path to the certificate file")
	flag.StringVar(&domain, "domain", "", "Domain to fetch the certificate from")
	flag.Parse()

	if certPath != "" {
		pemData, err := os.ReadFile(certPath)
		if err != nil {
			log.Fatal(err)
		}
		block, rest := pem.Decode([]byte(pemData))
		if block == nil || len(rest) > 0 {
			log.Fatal("Certificate decoding error")
		}
		cert, err = x509.ParseCertificate(block.Bytes)
		if err != nil {
			log.Fatal(err)
		}
		header = certPath
	} else if domain != "" {
		// Fetch certificate from domain
		conn, err := tls.Dial("tcp", domain+":443", &tls.Config{
			InsecureSkipVerify: true,
		})
		if err != nil {
			log.Fatalln("TLS connection failed: " + err.Error())
		}
		defer conn.Close()

		certs := conn.ConnectionState().PeerCertificates
		if len(certs) == 0 {
			log.Fatal("Error fetching certs from domain")
		}
		cert = certs[0]
		header = domain
	} else {
		fmt.Println("Please specify either a file path or a domain.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	certs := []list.Item{
		item{title: "Issued To", desc: cert.Subject.CommonName},
		item{title: "Issued By", desc: cert.Issuer.CommonName},
		item{title: "Issued On", desc: cert.NotBefore.String()},
		item{title: "Expires On", desc: cert.NotAfter.String()},
		item{title: "Public Key Algorithm", desc: cert.PublicKeyAlgorithm.String()},
		item{title: "Subject Alternative Names (DNS)", desc: strings.Join(cert.DNSNames, ",")},
	}

	m := model{list: list.New(certs, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = header

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
