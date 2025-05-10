import smtplib
from email.mime.text import MIMEText
from email.mime.multipart import MIMEMultipart

# Create message
msg = MIMEMultipart()
msg['From'] = 'ai1@example.am'
msg['To'] = 'support@example.am'
msg['Subject'] = 'Test Email from Python'

body = 'This is a test email sent to Mailhog 111'
msg.attach(MIMEText(body, 'plain'))

# Send email
try:
    with smtplib.SMTP('localhost', 1025) as server:
        server.send_message(msg)
    print("Test email sent successfully")
except Exception as e:
    print(f"Failed to send email: {e}") 