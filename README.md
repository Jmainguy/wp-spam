# wp-spam

This Go program serves as a webhook for a WordPress site to filter and handle incoming form submissions. The primary goal of this webhook is to detect and prevent spam by scanning form submissions for potential spam content (e.g., URLs or suspicious text). It also sends notification emails and Slack messages for legitimate form submissions.

## Features

- **Spam Detection**: Automatically flags form submissions as spam if they contain URLs or certain keywords.
- **Email Notifications**: Sends an email to a predefined list of recipients for valid form submissions.
- **Slack Notifications**: Posts a message to a Slack channel using a webhook for valid form submissions.
- **Customizable via Environment Variables**: Configure recipients, SMTP server, and Slack webhook via environment variables.


## Set Environment variables:

- **EMAIL_RECIPIENTS**: A comma-separated list of email addresses to notify.
- **SMTP_URL**: The URL of the SMTP server.
- **SMTP_PORT**: The SMTP server's port.
- **SMTP_USERNAME**: The username for SMTP authentication.
- **SMTP_PASSWORD**: The password for SMTP authentication.
- **SLACK_WEBHOOK**: The Slack webhook URL for sending messages.

## Usage

The webhook listens on /webhook and expects a JSON payload containing form data. Upon receiving the request:

* The program will first check for any spam-like content (such as links or .com mentions) in the form submission.
* If the submission is flagged as spam, it is ignored.
* If the submission is not spam, the following actions are taken:
    * An email is sent to the recipients specified in the **EMAIL_RECIPIENTS** environment variable.
    * A message is sent to the Slack webhook URL specified in the **SLACK_WEBHOOK** environment variable.

