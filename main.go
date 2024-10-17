package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"gopkg.in/gomail.v2"
)

// WebhookPayload represents the structure of the incoming JSON payload
type WebhookPayload struct {
	Checkbox1             string `json:"checkbox_1"`
	Address1StreetAddress string `json:"address_1_street_address"`
	Address1AddressLine   string `json:"address_1_address_line"`
	Address1City          string `json:"address_1_city"`
	Address1State         string `json:"address_1_state"`
	Address1Zip           string `json:"address_1_zip"`
	Name1Prefix           string `json:"name_1_prefix"`
	Name1FirstName        string `json:"name_1_first_name"`
	Name1LastName         string `json:"name_1_last_name"`
	Email1                string `json:"email_1"`
	Phone1                string `json:"phone_1"`
	Textarea2             string `json:"textarea_2"`
	RefererURL            string `json:"referer_url"`
	WPHttpReferer         string `json:"_wp_http_referer"`
	PageID                string `json:"page_id"`
	FormType              string `json:"form_type"`
	CurrentURL            string `json:"current_url"`
	RenderID              string `json:"render_id"`
	Input7                string `json:"input_7"`
	Address1              string `json:"address_1"`
	Name1                 string `json:"name_1"`
	ForminatorUserIP      string `json:"_forminator_user_ip"`
	FormTitle             string `json:"form_title"`
	EntryTime             string `json:"entry_time"`
}

func sendEmail(to []string, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "jon@soh.re")
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	smtpURL := os.Getenv("SMTP_URL")
	smtpPortString := os.Getenv("SMTP_PORT")
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	smtpPort, err := strconv.Atoi(smtpPortString)
	if err != nil {
		panic(err)
	}

	d := gomail.NewDialer(smtpURL, smtpPort, smtpUsername, smtpPassword)

	err = d.DialAndSend(m)
	if err != nil {
		return err
	}
	return nil
}

func formatServicesHTML(services string) string {
	if services == "" {
		return "<p><em>None</em></p>"
	}
	serviceList := strings.Split(services, ",")
	var formattedServices strings.Builder
	formattedServices.WriteString("<ul>")
	for _, service := range serviceList {
		formattedServices.WriteString(fmt.Sprintf("<li>%s</li>", service))
	}
	formattedServices.WriteString("</ul>")
	return formattedServices.String()
}

func formatEmailMessage(payload WebhookPayload) string {
	return fmt.Sprintf(`
		<h2>New Quote Request, Please respond in less than 24 hours</h2>
		<p><strong>Services They are interested in:</strong> %s</p>
		<p><strong>Street Address:</strong> %s</p>
		<p><strong>Apartment, suite, etc:</strong> %s</p>
		<p><strong>City:</strong> %s</p>
		<p><strong>State/Province:</strong> %s</p>
		<p><strong>ZIP / Postal Code:</strong> %s</p>
		<p><strong>Prefix:</strong> %s</p>
		<p><strong>First Name:</strong> %s</p>
		<p><strong>Last Name:</strong> %s</p>
		<p><strong>Email Address:</strong> %s</p>
		<p><strong>Phone Number:</strong> %s</p>
		<p><strong>Additional notes:</strong> %s</p>
	`, formatServicesHTML(payload.Checkbox1), payload.Address1StreetAddress, payload.Address1AddressLine, payload.Address1City, payload.Address1State, payload.Address1Zip, payload.Name1Prefix, payload.Name1FirstName, payload.Name1LastName, payload.Email1, payload.Phone1, payload.Textarea2)
}

func formatServicesMarkdown(services string) string {
	if services == "" {
		return "_None_"
	}
	serviceList := strings.Split(services, ", ")
	var formattedServices strings.Builder
	for _, service := range serviceList {
		formattedServices.WriteString(fmt.Sprintf("• %s\n", service))
	}
	return formattedServices.String()
}

func sendSlackMessage(webhookURL string, payload WebhookPayload) error {
	message := map[string]interface{}{
		"attachments": []map[string]interface{}{
			{
				"color":   "#36a64f",
				"pretext": "New Quote Request",
				"fields": []map[string]interface{}{
					{"title": "Services they are interested in", "value": formatServicesMarkdown(payload.Checkbox1), "short": false},
					{"title": "Street Address", "value": payload.Address1StreetAddress, "short": true},
					{"title": "Apartment, suite, etc", "value": payload.Address1AddressLine, "short": true},
					{"title": "City", "value": payload.Address1City, "short": true},
					{"title": "State/Province", "value": payload.Address1State, "short": true},
					{"title": "ZIP / Postal Code", "value": payload.Address1Zip, "short": true},
					{"title": "Prefix", "value": payload.Name1Prefix, "short": true},
					{"title": "First Name", "value": payload.Name1FirstName, "short": true},
					{"title": "Last Name", "value": payload.Name1LastName, "short": true},
					{"title": "Email Address", "value": payload.Email1, "short": true},
					{"title": "Phone Number", "value": payload.Phone1, "short": true},
					{"title": "Additional notes", "value": payload.Textarea2, "short": false},
				},
			},
		},
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(messageBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to send message to Slack: %s", body)
	}

	return nil
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	var spam bool
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}

	var payload WebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "Unable to parse JSON", http.StatusBadRequest)
		return
	}

	if strings.Contains(payload.Textarea2, "http://") {
		spam = true
	} else if strings.Contains(payload.Textarea2, "https://") {
		spam = true
	} else if strings.Contains(payload.Textarea2, ".com") {
		spam = true
	}

	if !spam {
		fmt.Printf("Received webhook: %s\n", body)
		// Get email recipients from environment variable
		emailRecipients := os.Getenv("EMAIL_RECIPIENTS")
		to := strings.Split(emailRecipients, ",")

		// Send email
		subject := "New Quote Request"
		emailBody := formatEmailMessage(payload)
		if err := sendEmail(to, subject, emailBody); err != nil {
			log.Printf("Failed to send email: %v", err)
		}

		// Send Slack message

		webhookURL := os.Getenv("SLACK_WEBHOOK")
		if err := sendSlackMessage(webhookURL, payload); err != nil {
			log.Printf("Failed to send Slack message: %v", err)
		}
	}

	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/webhook", webhookHandler)

	port := "8080"
	fmt.Printf("Listening for webhooks on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
