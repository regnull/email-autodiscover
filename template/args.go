package template

type Args struct {
	EmailProvider string `yaml: "email-provider"`
	Domain        string `yaml: "domain"`
	ImapHost      string `yaml: "imap-host"`
	ImapPort      int    `yaml: "imap-port"`
	SmtpHost      string `yaml: "smtp-host"`
	SmtpPort      int    `yaml: "smtp-port`
	Email         string // This one is populated from the request.
}
