package template

import (
	"bytes"
	"text/template"
)

const iosMailTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>HasRemovalPasscode</key>
	<false/>
	<key>PayloadContent</key>
	<array>
		<dict>
			<key>EmailAccountDescription</key>
			<string>{{.Email}}</string>
			<key>EmailAccountName</key>
			<string>{{.Email}}</string>
			<key>EmailAccountType</key>
			<string>EmailTypeIMAP</string>
			<key>EmailAddress</key>
			<string>{{.Email}}</string>
			<key>IncomingMailServerAuthentication</key>
			<string>EmailAuthPassword</string>
			<key>IncomingMailServerHostName</key>
			<string>{{.ImapHost}}</string>
			<key>IncomingMailServerPortNumber</key>
			<integer>{{.ImapPort}}</integer>
			<key>IncomingMailServerUseSSL</key>
			<true/>
			<key>IncomingMailServerUsername</key>
			<string>{{.Email}}</string>
			<key>OutgoingMailServerAuthentication</key>
			<string>EmailAuthPassword</string>
			<key>OutgoingMailServerHostName</key>
			<string>{{.SmtpHost}}</string>
			<key>OutgoingMailServerPortNumber</key>
			<integer>{{.SmtpPort}}</integer>
			<key>OutgoingMailServerUseSSL</key>
			<true/>
			<key>OutgoingMailServerUsername</key>
			<string>{{.Email}}</string>
			<key>OutgoingPasswordSameAsIncomingPassword</key>
			<true/>
			<key>PayloadDescription</key>
			<string>Configure Email Settings</string>
			<key>PayloadDisplayName</key>
			<string>{{.Email}}</string>
			<key>PayloadIdentifier</key>
			<string>cc.ubikom.autodiscover.com.apple.mail.managed.7A981A9E-D5D1-4EF8-87FE-39FD6A506FAC</string>
			<key>PayloadType</key>
			<string>com.apple.mail.managed</string>
			<key>PayloadUUID</key>
			<string>7A981A9E-D5D1-4EF8-87FE-39FD6A506FAC</string>
			<key>PayloadVersion</key>
			<real>1</real>
			<key>SMIMEEnablePerMessageSwitch</key>
			<false/>
			<key>SMIMEEnabled</key>
			<false/>
			<key>disableMailRecentsSyncing</key>
			<false/>
		</dict>
	</array>
	<key>PayloadDescription</key>
	<string>Configure Email Settings</string>
	<key>PayloadDisplayName</key>
	<string>{{.Email}}</string>
	<key>PayloadIdentifier</key>
	<string>cc.ubikom.autodiscover</string>
	<key>PayloadOrganization</key>
	<string>{{.Domain}}</string>
	<key>PayloadRemovalDisallowed</key>
	<false/>
	<key>PayloadType</key>
	<string>Configuration</string>
	<key>PayloadUUID</key>
	<string>48C88203-4DBA-49E8-B593-4831903605A0</string>
	<key>PayloadVersion</key>
	<integer>1</integer>
</dict>
</plist>
`

func IOSMail(args *Args) (string, error) {
	reportTmpl, err := template.New("report").Parse(iosMailTemplate)
	if err != nil {
		return "", err
	}
	var b bytes.Buffer
	err = reportTmpl.Execute(&b, args)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}
