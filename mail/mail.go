package mail

import (
	"encoding/json"
	"log"
	"net/mail"
	"strings"
)

// Mail contains mail struct
type Mail struct {
	From             Email              `yaml:"from" mapstructure:"from" json:"from,omitempty" gorm:"column:from" bson:"from,omitempty" dynamodbav:"from,omitempty" firestore:"from,omitempty"`
	To               []Email            `json:"-,omitempty"`
	Cc               *[]Email           `json:"-,omitempty"`
	Bcc              *[]Email           `json:"-,omitempty"`
	Subject          string             `yaml:"subject" mapstructure:"subject" json:"subject,omitempty" json:"subject,omitempty" gorm:"column:subject" bson:"subject,omitempty" dynamodbav:"subject,omitempty" firestore:"subject,omitempty"`
	Personalizations []*Personalization `yaml:"personalizations" mapstructure:"personalizations" json:"personalizations,omitempty" gorm:"column:personalizations" bson:"personalizations,omitempty" dynamodbav:"personalizations,omitempty" firestore:"personalizations,omitempty"`
	Content          []Content          `yaml:"" mapstructure:"content" json:"content,omitempty" gorm:"column:content" bson:"content,omitempty" dynamodbav:"content,omitempty" firestore:"content,omitempty"`
	Attachments      []*Attachment      `yaml:"" mapstructure:"attachments" json:"attachments,omitempty" gorm:"column:attachments" bson:"attachments,omitempty" dynamodbav:"attachments,omitempty" firestore:"attachments,omitempty"`
	TemplateID       string             `json:"template_id,omitempty"`
	Sections         map[string]string  `yaml:"sections" mapstructure:"sections" json:"sections,omitempty"`
	Headers          map[string]string  `yaml:"headers" mapstructure:"headers" json:"headers,omitempty"`
	Categories       []string           `yaml:"categories" mapstructure:"categories" json:"categories,omitempty"`
	CustomArgs       map[string]string  `json:"custom_args,omitempty"`
	SendAt           int                `json:"send_at,omitempty"`
	BatchID          string             `json:"batch_id,omitempty"`
	Asm              *Asm               `yaml:"asm" mapstructure:"asm" json:"asm,omitempty"`
	IPPoolID         string             `json:"ip_pool_name,omitempty"`
	MailSettings     *MailSettings      `json:"mail_settings,omitempty"`
	TrackingSettings *TrackingSettings  `json:"tracking_settings,omitempty"`
	ReplyTo          *Email             `json:"reply_to,omitempty"`
}

// Personalization holds mail body struct
type Personalization struct {
	To                  []*Email               `yaml:"to" mapstructure:"to" json:"to,omitempty" `
	CC                  []*Email               `yaml:"cc" mapstructure:"cc" json:"cc,omitempty"`
	BCC                 []*Email               `yaml:"bcc" mapstructure:"bcc" json:"bcc,omitempty" `
	Subject             string                 `yaml:"subject" mapstructure:"subject" json:"subject,omitempty" gorm:"column:subject" bson:"subject,omitempty" dynamodbav:"subject,omitempty" firestore:"subject,omitempty"`
	Headers             map[string]string      `json:"headers,omitempty"`
	Substitutions       map[string]string      `json:"substitutions,omitempty"`
	CustomArgs          map[string]string      `json:"custom_args,omitempty"`
	DynamicTemplateData map[string]interface{} `json:"dynamic_template_data,omitempty"`
	Categories          []string               `json:"categories,omitempty"`
	SendAt              int                    `json:"send_at,omitempty"`
}

// Email holds email name and address info
type Email struct {
	Name    string `yaml:"name" mapstructure:"name" json:"name,omitempty" gorm:"column:name" bson:"name,omitempty" dynamodbav:"name,omitempty" firestore:"name,omitempty"`
	Address string `yaml:"address" mapstructure:"address" json:"email,omitempty" gorm:"column:address" bson:"address,omitempty" dynamodbav:"address,omitempty" firestore:"address,omitempty"`
}

// Content defines content of the mail body
type Content struct {
	Type  string `yaml:"type" mapstructure:"type" json:"type,omitempty" gorm:"column:type" bson:"type,omitempty" dynamodbav:"type,omitempty" firestore:"type,omitempty"`
	Value string `yaml:"value" mapstructure:"value" json:"value,omitempty" gorm:"column:value" bson:"value,omitempty" dynamodbav:"value,omitempty" firestore:"value,omitempty"`
}

// Attachment holds attachement information
type Attachment struct {
	Content     string `json:"content,omitempty"`
	Type        string `json:"type,omitempty"`
	Name        string `json:"name,omitempty"`
	Filename    string `json:"filename,omitempty"`
	Disposition string `json:"disposition,omitempty"`
	ContentID   string `json:"content_id,omitempty"`
}

// Asm contains Grpip Id and int array of groups ID
type Asm struct {
	GroupID         int   `json:"group_id,omitempty"`
	GroupsToDisplay []int `json:"groups_to_display,omitempty"`
}

// MailSettings defines mail and spamCheck settings
type MailSettings struct {
	BCC                  *BccSetting       `json:"bcc,omitempty"`
	BypassListManagement *Setting          `json:"bypass_list_management,omitempty"`
	Footer               *FooterSetting    `json:"footer,omitempty"`
	SandboxMode          *Setting          `json:"sandbox_mode,omitempty"`
	SpamCheckSetting     *SpamCheckSetting `json:"spam_check,omitempty"`
}

// TrackingSettings holds tracking settings and mail settings
type TrackingSettings struct {
	ClickTracking        *ClickTrackingSetting        `json:"click_tracking,omitempty"`
	OpenTracking         *OpenTrackingSetting         `json:"open_tracking,omitempty"`
	SubscriptionTracking *SubscriptionTrackingSetting `json:"subscription_tracking,omitempty"`
	GoogleAnalytics      *GaSetting                   `json:"ganalytics,omitempty"`
	BCC                  *BccSetting                  `json:"bcc,omitempty"`
	BypassListManagement *Setting                     `json:"bypass_list_management,omitempty"`
	Footer               *FooterSetting               `json:"footer,omitempty"`
	SandboxMode          *SandboxModeSetting          `json:"sandbox_mode,omitempty"`
}

// BccSetting holds email bcc setings  to enable of disable
// default is false
type BccSetting struct {
	Enable *bool  `json:"enable,omitempty"`
	Email  string `json:"email,omitempty"`
}

// FooterSetting holds enaable/disable settings
// and the format of footer i.e HTML/Text
type FooterSetting struct {
	Enable *bool  `json:"enable,omitempty"`
	Text   string `json:"text,omitempty"`
	Html   string `json:"html,omitempty"`
}

// ClickTrackingSetting ...
type ClickTrackingSetting struct {
	Enable     *bool `json:"enable,omitempty"`
	EnableText *bool `json:"enable_text,omitempty"`
}

// OpenTrackingSetting ...
type OpenTrackingSetting struct {
	Enable          *bool  `json:"enable,omitempty"`
	SubstitutionTag string `json:"substitution_tag,omitempty"`
}

// SandboxModeSetting ...
type SandboxModeSetting struct {
	Enable      *bool             `json:"enable,omitempty"`
	ForwardSpam *bool             `json:"forward_spam,omitempty"`
	SpamCheck   *SpamCheckSetting `json:"spam_check,omitempty"`
}

// SpamCheckSetting holds spam settings and
// which can be enable or disable and
// contains spamThreshold value
type SpamCheckSetting struct {
	Enable        *bool  `json:"enable,omitempty"`
	SpamThreshold int    `json:"threshold,omitempty"`
	PostToURL     string `json:"post_to_url,omitempty"`
}

// SubscriptionTrackingSetting ...
type SubscriptionTrackingSetting struct {
	Enable          *bool  `json:"enable,omitempty"`
	Text            string `json:"text,omitempty"`
	Html            string `json:"html,omitempty"`
	SubstitutionTag string `json:"substitution_tag,omitempty"`
}

// GaSetting ...
type GaSetting struct {
	Enable          *bool  `json:"enable,omitempty"`
	CampaignSource  string `json:"utm_source,omitempty"`
	CampaignTerm    string `json:"utm_term,omitempty"`
	CampaignContent string `json:"utm_content,omitempty"`
	CampaignName    string `json:"utm_campaign,omitempty"`
	CampaignMedium  string `json:"utm_medium,omitempty"`
}

// Setting enables the mail settings
type Setting struct {
	Enable *bool `json:"enable,omitempty"`
}

// NewMail ...
func NewMail() *Mail {
	return &Mail{
		Personalizations: make([]*Personalization, 0),
		To:               make([]Email, 0),
		Content:          make([]Content, 0),
		Attachments:      make([]*Attachment, 0),
	}
}

// NewMailInit ...
func NewMailInit(from Email, subject string, to []Email, cc *[]Email, content ...Content) *Mail {
	m := new(Mail)
	m.SetFrom(from)
	m.Subject = subject
	p := NewPersonalization()
	for i := 0; i < len(to); i++ {
		p.AddTos(&to[i])
	}
	m.AddPersonalizations(p)
	m.AddContent(content...)
	m.To = to
	m.Cc = cc
	return m
}

// GetRequestBody ...
func GetRequestBody(m *Mail) []byte {
	b, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
	}
	return b
}

// AddPersonalizations ...
func (s *Mail) AddPersonalizations(p ...*Personalization) *Mail {
	s.Personalizations = append(s.Personalizations, p...)
	return s
}

// AddContent ...
func (s *Mail) AddContent(c ...Content) *Mail {
	s.Content = append(s.Content, c...)
	return s
}

// AddAttachment ...
func (s *Mail) AddAttachment(a ...*Attachment) *Mail {
	s.Attachments = append(s.Attachments, a...)
	return s
}

// SetFrom ...
func (s *Mail) SetFrom(e Email) *Mail {
	s.From = e
	return s
}

// SetReplyTo ...
func (s *Mail) SetReplyTo(e *Email) *Mail {
	s.ReplyTo = e
	return s
}

// SetTemplateID ...
func (s *Mail) SetTemplateID(templateID string) *Mail {
	s.TemplateID = templateID
	return s
}

// AddSection ...
func (s *Mail) AddSection(key string, value string) *Mail {
	if s.Sections == nil {
		s.Sections = make(map[string]string)
	}

	s.Sections[key] = value
	return s
}

// SetHeader ...
func (s *Mail) SetHeader(key string, value string) *Mail {
	if s.Headers == nil {
		s.Headers = make(map[string]string)
	}

	s.Headers[key] = value
	return s
}

// AddCategories ...
func (s *Mail) AddCategories(category ...string) *Mail {
	s.Categories = append(s.Categories, category...)
	return s
}

// SetCustomArg ...
func (s *Mail) SetCustomArg(key string, value string) *Mail {
	if s.CustomArgs == nil {
		s.CustomArgs = make(map[string]string)
	}

	s.CustomArgs[key] = value
	return s
}

// SetSendAt ...
func (s *Mail) SetSendAt(sendAt int) *Mail {
	s.SendAt = sendAt
	return s
}

// SetBatchID ...
func (s *Mail) SetBatchID(batchID string) *Mail {
	s.BatchID = batchID
	return s
}

// SetASM ...
func (s *Mail) SetASM(asm *Asm) *Mail {
	s.Asm = asm
	return s
}

// SetIPPoolID ...
func (s *Mail) SetIPPoolID(ipPoolID string) *Mail {
	s.IPPoolID = ipPoolID
	return s
}

// SetMailSettings ...
func (s *Mail) SetMailSettings(mailSettings *MailSettings) *Mail {
	s.MailSettings = mailSettings
	return s
}

// SetTrackingSettings ...
func (s *Mail) SetTrackingSettings(trackingSettings *TrackingSettings) *Mail {
	s.TrackingSettings = trackingSettings
	return s
}

// NewPersonalization ...
func NewPersonalization() *Personalization {
	return &Personalization{
		To:                  make([]*Email, 0),
		CC:                  make([]*Email, 0),
		BCC:                 make([]*Email, 0),
		Headers:             make(map[string]string),
		Substitutions:       make(map[string]string),
		CustomArgs:          make(map[string]string),
		DynamicTemplateData: make(map[string]interface{}),
		Categories:          make([]string, 0),
	}
}

// AddTos ...
func (p *Personalization) AddTos(to ...*Email) {
	p.To = append(p.To, to...)
}

// AddCCs ...
func (p *Personalization) AddCCs(cc ...*Email) {
	p.CC = append(p.CC, cc...)
}

// AddBCCs ...
func (p *Personalization) AddBCCs(bcc ...*Email) {
	p.BCC = append(p.BCC, bcc...)
}

// SetHeader ...
func (p *Personalization) SetHeader(key string, value string) {
	p.Headers[key] = value
}

// SetSubstitution ...
func (p *Personalization) SetSubstitution(key string, value string) {
	p.Substitutions[key] = value
}

// SetCustomArg ...
func (p *Personalization) SetCustomArg(key string, value string) {
	p.CustomArgs[key] = value
}

// SetDynamicTemplateData ...
func (p *Personalization) SetDynamicTemplateData(key string, value interface{}) {
	p.DynamicTemplateData[key] = value
}

// SetSendAt ...
func (p *Personalization) SetSendAt(sendAt int) {
	p.SendAt = sendAt
}

// NewAttachment ...
func NewAttachment() *Attachment {
	return &Attachment{}
}

// SetContent ...
func (a *Attachment) SetContent(content string) *Attachment {
	a.Content = content
	return a
}

// SetType ...
func (a *Attachment) SetType(contentType string) *Attachment {
	a.Type = contentType
	return a
}

// SetFilename ...
func (a *Attachment) SetFilename(filename string) *Attachment {
	a.Filename = filename
	return a
}

// SetDisposition ...
func (a *Attachment) SetDisposition(disposition string) *Attachment {
	a.Disposition = disposition
	return a
}

// SetContentID ...
func (a *Attachment) SetContentID(contentID string) *Attachment {
	a.ContentID = contentID
	return a
}

// NewASM ...
func NewASM() *Asm {
	return &Asm{}
}

// SetGroupID ...
func (a *Asm) SetGroupID(groupID int) *Asm {
	a.GroupID = groupID
	return a
}

// AddGroupsToDisplay ...
func (a *Asm) AddGroupsToDisplay(groupsToDisplay ...int) *Asm {
	a.GroupsToDisplay = append(a.GroupsToDisplay, groupsToDisplay...)
	return a
}

// NewMailSettings ...
func NewMailSettings() *MailSettings {
	return &MailSettings{}
}

// SetBCC ...
func (m *MailSettings) SetBCC(bcc *BccSetting) *MailSettings {
	m.BCC = bcc
	return m
}

// SetBypassListManagement ...
func (m *MailSettings) SetBypassListManagement(bypassListManagement *Setting) *MailSettings {
	m.BypassListManagement = bypassListManagement
	return m
}

// SetFooter ...
func (m *MailSettings) SetFooter(footerSetting *FooterSetting) *MailSettings {
	m.Footer = footerSetting
	return m
}

// SetSandboxMode ...
func (m *MailSettings) SetSandboxMode(sandboxMode *Setting) *MailSettings {
	m.SandboxMode = sandboxMode
	return m
}

// SetSpamCheckSettings ...
func (m *MailSettings) SetSpamCheckSettings(spamCheckSetting *SpamCheckSetting) *MailSettings {
	m.SpamCheckSetting = spamCheckSetting
	return m
}

// NewTrackingSettings ...
func NewTrackingSettings() *TrackingSettings {
	return &TrackingSettings{}
}

// SetClickTracking ...
func (t *TrackingSettings) SetClickTracking(clickTracking *ClickTrackingSetting) *TrackingSettings {
	t.ClickTracking = clickTracking
	return t

}

// SetOpenTracking ...
func (t *TrackingSettings) SetOpenTracking(openTracking *OpenTrackingSetting) *TrackingSettings {
	t.OpenTracking = openTracking
	return t
}

// SetSubscriptionTracking ...
func (t *TrackingSettings) SetSubscriptionTracking(subscriptionTracking *SubscriptionTrackingSetting) *TrackingSettings {
	t.SubscriptionTracking = subscriptionTracking
	return t
}

// SetGoogleAnalytics ...
func (t *TrackingSettings) SetGoogleAnalytics(googleAnalytics *GaSetting) *TrackingSettings {
	t.GoogleAnalytics = googleAnalytics
	return t
}

// NewBCCSetting ...
func NewBCCSetting() *BccSetting {
	return &BccSetting{}
}

// SetEnable ...
func (b *BccSetting) SetEnable(enable bool) *BccSetting {
	setEnable := enable
	b.Enable = &setEnable
	return b
}

// SetEmail ...
func (b *BccSetting) SetEmail(email string) *BccSetting {
	b.Email = email
	return b
}

// NewFooterSetting ...
func NewFooterSetting() *FooterSetting {
	return &FooterSetting{}
}

// SetEnable ...
func (f *FooterSetting) SetEnable(enable bool) *FooterSetting {
	setEnable := enable
	f.Enable = &setEnable
	return f
}

// SetText ...
func (f *FooterSetting) SetText(text string) *FooterSetting {
	f.Text = text
	return f
}

// SetHTML ...
func (f *FooterSetting) SetHTML(html string) *FooterSetting {
	f.Html = html
	return f
}

// NewOpenTrackingSetting ...
func NewOpenTrackingSetting() *OpenTrackingSetting {
	return &OpenTrackingSetting{}
}

// SetEnable ...
func (o *OpenTrackingSetting) SetEnable(enable bool) *OpenTrackingSetting {
	setEnable := enable
	o.Enable = &setEnable
	return o
}

// SetSubstitutionTag ...
func (o *OpenTrackingSetting) SetSubstitutionTag(subTag string) *OpenTrackingSetting {
	o.SubstitutionTag = subTag
	return o
}

// NewSubscriptionTrackingSetting ...
func NewSubscriptionTrackingSetting() *SubscriptionTrackingSetting {
	return &SubscriptionTrackingSetting{}
}

// SetEnable ...
func (s *SubscriptionTrackingSetting) SetEnable(enable bool) *SubscriptionTrackingSetting {
	setEnable := enable
	s.Enable = &setEnable
	return s
}

// SetText ...
func (s *SubscriptionTrackingSetting) SetText(text string) *SubscriptionTrackingSetting {
	s.Text = text
	return s
}

// SetHTML ...
func (s *SubscriptionTrackingSetting) SetHTML(html string) *SubscriptionTrackingSetting {
	s.Html = html
	return s
}

// SetSubstitutionTag ...
func (s *SubscriptionTrackingSetting) SetSubstitutionTag(subTag string) *SubscriptionTrackingSetting {
	s.SubstitutionTag = subTag
	return s
}

// NewGaSetting ...
func NewGaSetting() *GaSetting {
	return &GaSetting{}
}

// SetEnable ...
func (g *GaSetting) SetEnable(enable bool) *GaSetting {
	setEnable := enable
	g.Enable = &setEnable
	return g
}

// SetCampaignSource ...
func (g *GaSetting) SetCampaignSource(campaignSource string) *GaSetting {
	g.CampaignSource = campaignSource
	return g
}

// SetCampaignContent ...
func (g *GaSetting) SetCampaignContent(campaignContent string) *GaSetting {
	g.CampaignContent = campaignContent
	return g
}

// SetCampaignTerm ...
func (g *GaSetting) SetCampaignTerm(campaignTerm string) *GaSetting {
	g.CampaignTerm = campaignTerm
	return g
}

// SetCampaignName ...
func (g *GaSetting) SetCampaignName(campaignName string) *GaSetting {
	g.CampaignName = campaignName
	return g
}

// SetCampaignMedium ...
func (g *GaSetting) SetCampaignMedium(campaignMedium string) *GaSetting {
	g.CampaignMedium = campaignMedium
	return g
}

// NewSetting ...
func NewSetting(enable bool) *Setting {
	setEnable := enable
	return &Setting{Enable: &setEnable}
}

// EscapeName adds quotes around the name to prevent errors from RFC5322 special
// characters:
//
//   ()<>[]:;@\,."
//
// To preserve backwards compatibility for people already quoting their name
// inputs, as well as for inputs which do not strictly require quoting, the
// name is returned unmodified if those conditions are met. Otherwise, existing
// intrastring backslashes and double quotes are escaped, and the entire input
// is surrounded with double quotes.
func EscapeName(name string) string {
	if len(name) > 1 && name[0] == '"' && name[len(name)-1] == '"' {
		return name
	}
	if strings.IndexAny(name, "()<>[]:;@\\,.\"") == -1 {
		return name
	}

	// This has to come first so we don't triple backslash after the next step
	name = strings.Replace(name, `\`, `\\`, -1)
	name = strings.Replace(name, `"`, `\"`, -1)
	name = `"` + name + `"`

	return name
}

// NewEmail ...
func NewEmail(name string, address string) *Email {
	name = EscapeName(name)
	return &Email{
		Name:    name,
		Address: address,
	}
}

// NewSingleEmail ...
func NewSingleEmail(from Email, subject string, to []Email, cc *[]Email, plainTextContent string, htmlContent string) *Mail {
	plainText := NewContent("text/plain", plainTextContent)
	html := NewContent("text/html", htmlContent)
	return NewMailInit(from, subject, to, cc, *plainText, *html)
}

func NewHtmlMail(mailFrom Email, subject string, mailTo []Email, cc *[]Email, htmlContent string) *Mail {
	//mailTo := []Email{{Name: toName, Address: toAddress}}
	//mailFrom := Email{Address: from}
	html := NewContent("text/html", htmlContent)
	mail := NewMailInit(mailFrom, subject, mailTo, cc, *html)
	// mail.Html = htmlContent
	mail.To = mailTo
	mail.Cc = cc
	return mail
}

func NewPlainTextMail(mailFrom Email, subject string, mailTo []Email, cc *[]Email, plainTextContent string) *Mail {
	//mailTo := []Email{{Name: toName, Address: toAddress}}
	//mailFrom := Email{Address: from}
	plainText := NewContent("text/plain", plainTextContent)
	mail := NewMailInit(mailFrom, subject, mailTo, cc, *plainText)
	// mail.Text = plainTextContent
	mail.To = mailTo
	return mail
}

// NewContent ...
func NewContent(contentType string, value string) *Content {
	return &Content{
		Type:  contentType,
		Value: value,
	}
}

// NewClickTrackingSetting ...
func NewClickTrackingSetting() *ClickTrackingSetting {
	return &ClickTrackingSetting{}
}

// SetEnable ...
func (c *ClickTrackingSetting) SetEnable(enable bool) *ClickTrackingSetting {
	setEnable := enable
	c.Enable = &setEnable
	return c
}

// SetEnableText ...
func (c *ClickTrackingSetting) SetEnableText(enableText bool) *ClickTrackingSetting {
	setEnable := enableText
	c.EnableText = &setEnable
	return c
}

// NewSpamCheckSetting ...
func NewSpamCheckSetting() *SpamCheckSetting {
	return &SpamCheckSetting{}
}

// SetEnable ...
func (s *SpamCheckSetting) SetEnable(enable bool) *SpamCheckSetting {
	setEnable := enable
	s.Enable = &setEnable
	return s
}

// SetSpamThreshold ...
func (s *SpamCheckSetting) SetSpamThreshold(spamThreshold int) *SpamCheckSetting {
	s.SpamThreshold = spamThreshold
	return s
}

// SetPostToURL ...
func (s *SpamCheckSetting) SetPostToURL(postToURL string) *SpamCheckSetting {
	s.PostToURL = postToURL
	return s
}

// NewSandboxModeSetting ...
func NewSandboxModeSetting(enable bool, forwardSpam bool, spamCheck *SpamCheckSetting) *SandboxModeSetting {
	setEnable := enable
	setForwardSpam := forwardSpam
	return &SandboxModeSetting{
		Enable:      &setEnable,
		ForwardSpam: &setForwardSpam,
		SpamCheck:   spamCheck,
	}
}

// ParseEmail parses a string that contains an rfc822 formatted email address
// and returns an instance of *Email.
func ParseEmail(emailInfo string) (*Email, error) {
	e, err := mail.ParseAddress(emailInfo)
	if err != nil {
		return nil, err
	}
	return NewEmail(e.Name, e.Address), nil
}
